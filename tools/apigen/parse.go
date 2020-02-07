package apigen

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

var DefaultParseOpts = &ParseOpts{
	IgnoreH3: []string{
		"Recent changes",
		"Authorizing your bot",
		"Making requests",
	},
	IgnoreH4: []string{
		"Available types$Sending files",
		"Available types$Inline mode objects",
		"Available methods$Formatting options",
	},
}

type ParseOpts struct {
	IgnoreH3 []string
	IgnoreH4 []string
}

func (o *ParseOpts) skipH3(h3 string) bool {
	for _, v := range o.IgnoreH3 {
		if v == h3 {
			return true
		}
	}

	return false
}

func (o *ParseOpts) skipH4(h3 string, h4 string) bool {
	for _, v := range o.IgnoreH4 {
		if v == h3+"$"+h4 {
			return true
		}
	}

	return false
}

type ParsedAPI struct {
	Chapters map[string]*Chapter
}

func (p *ParsedAPI) GetChapter(key string) *Chapter {
	c, ok := p.Chapters[key]
	if !ok {
		c = &Chapter{
			Name: key,
		}
		p.Chapters[key] = c
	}

	return c
}

type Chapter struct {
	Name    string
	Objects []*Object
}

func (p *Chapter) GetObject(name string) *Object {
	for _, obj := range p.Objects {
		if obj.Name == name {
			return obj
		}
	}

	obj := &Object{
		Name: name,
	}

	p.Objects = append(p.Objects, obj)
	return obj
}

type Object struct {
	IsType     bool
	IsFunction bool
	Name       string
	Notes      []string
	Fields     []Field

	// Function-specific fields
	ReturnType string
}

type Field struct {
	Name        string
	Type        Type
	Description string
	IsOptional  bool
	IsRequired  bool
}

type Type struct {
	Name    string
	HasLink bool // indicated reference to object
}

type sharedContext struct {
	err error
}

type parseContext struct {
	opts   *ParseOpts
	api    *ParsedAPI
	shared *sharedContext
}

func hasAttr(n *html.Node, tag string, val string) bool {
	for _, attr := range n.Attr {
		if attr.Key == tag && attr.Val == val {
			return true
		}
	}

	return false
}

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	if n.Type == html.ElementNode && n.Data == "br" {
		return "\n"
	}

	res := ""

	if n.Type == html.ElementNode && n.Data == "li" {
		res = "- "
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		child := extractText(c)

		switch {
		case n.Type == html.ElementNode && (n.Data == "ul" || n.Data == "ol"):
			if res != "" {
				res += "\n"
			}
			res += child

		default:
			res += child
		}
	}

	return res
}

func checkTag(n *html.Node, tag string) bool {
	return n != nil && n.Type == html.ElementNode && n.Data == tag
}

// parseRecursive1 searches for dev_page_content, which contains all we possibly need
func parseRecursive1(ctx parseContext, n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "div" && hasAttr(n, "id", "dev_page_content") {
		parsePageContent(ctx, n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseRecursive1(ctx, c)
	}
}

func parsePageContent(ctx parseContext, n *html.Node) {
	h3 := ""
	h4 := ""

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Data == "h3" {
			h3 = extractText(c)
			h4 = ""

			if ctx.opts.skipH3(h3) {
				h3 = ""
			}

			continue
		}

		if h3 == "" {
			// skip
			continue
		}

		if c.Data == "h4" {
			h4 = extractText(c)

			if ctx.opts.skipH4(h3, h4) {
				h4 = ""
			}

			continue
		}

		if h4 == "" {
			// skip
			continue
		}

		obj := ctx.api.GetChapter(h3).GetObject(h4)

		parseObjectContext(ctx, obj, c)
	}
}

func parseObjectContext(ctx parseContext, obj *Object, c *html.Node) {
	var addNote *html.Node

	if c.Type == html.ElementNode {
		switch {
		case c.Data == "p":
			addNote = c

		case c.Data == "blockquote":
			addNote = c

		case c.Data == "table":
			parseTable(ctx, obj, c)

		case c.Data == "div" && hasAttr(c, "class", "blog_image_wrap"):
			return

		case c.Data == "ul" || c.Data == "ol":
			addNote = c

		case c.Data == "hr":
			return

		default:
			ctx.shared.err = fmt.Errorf("failed to recognize <%s> tag", c.Data)
		}
	}

	if addNote != nil {
		noteText := extractText(addNote)
		obj.Notes = append(obj.Notes, noteText)

		if obj.ReturnType == "" {
			for a := addNote.FirstChild; a != nil; a = a.NextSibling {
				if checkTag(a, "a") {
					possibleReturnType := extractText(a)
					if strings.Title(possibleReturnType) != possibleReturnType {
						continue
					}

					if strings.Contains(possibleReturnType, " ") {
						continue
					}

					obj.ReturnType = possibleReturnType
				}
			}
		}

		firstSentence := strings.ToLower(noteText)
		if pos := strings.Index(firstSentence, "."); pos != -1 {
			firstSentence = firstSentence[:pos]
		}

		if strings.Contains(firstSentence, "this method") {
			obj.IsFunction = true
		}

		if strings.Contains(firstSentence, "a simple method") {
			obj.IsFunction = true
		}

		if strings.Contains(firstSentence, "this object") {
			obj.IsType = true
		}
	}
}

func parseTable(ctx parseContext, obj *Object, table *html.Node) {
	var (
		thead *html.Node
		tbody *html.Node
	)

	for c := table.FirstChild; c != nil; c = c.NextSibling {
		if c.Data == "thead" {
			thead = c
		}
		if c.Data == "tbody" {
			tbody = c
		}
	}

	if !checkTag(thead, "thead") {
		ctx.shared.err = fmt.Errorf("failed to parse table[0] <%s>", thead.Data)
		return
	}

	if !checkTag(tbody, "tbody") {
		ctx.shared.err = fmt.Errorf("failed to parse table[1] <%s>", tbody.Data)
		return
	}

	theadText := extractText(thead)
	theadText = strings.ReplaceAll(theadText, " ", "")
	theadText = strings.ReplaceAll(theadText, "\n", "")

	iterateTR := func(n *html.Node, cb func(node *html.Node)) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if checkTag(c, "tr") {
				cb(c)
			}
		}
	}

	extractTD := func(n *html.Node) []*html.Node {
		var arr []*html.Node

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if checkTag(c, "td") {
				arr = append(arr, c)
			}
		}

		return arr
	}

	if theadText == "ParameterTypeRequiredDescription" {
		// must be a function
		obj.IsFunction = true

		iterateTR(tbody, func(f *html.Node) {
			td := extractTD(f)
			tableLen := 4

			if len(td) != tableLen {
				ctx.shared.err = fmt.Errorf("unknown function table structure, %v", td)
			}

			param := extractText(td[0])
			tp := parseType(td[1])
			required := extractText(td[2])
			descr := extractText(td[3])

			field := Field{
				Name:        param,
				Type:        tp,
				Description: descr,
				IsOptional:  required == "Optional",
				IsRequired:  required == "Yes",
			}

			if !field.IsOptional && !field.IsRequired {
				ctx.shared.err = fmt.Errorf("unknown required info, %v", required)
			}

			obj.Fields = append(obj.Fields, field)
		})

		return
	}

	if theadText == "FieldTypeDescription" {
		// must be an object
		obj.IsType = true

		iterateTR(tbody, func(f *html.Node) {
			td := extractTD(f)
			tableLen := 3

			if len(td) != tableLen {
				ctx.shared.err = fmt.Errorf("unknown type table structure, %v", td)
			}

			param := extractText(td[0])
			tp := parseType(td[1])
			descr := extractText(td[2])
			isOptional := strings.HasPrefix(descr, "Optional.")

			field := Field{
				Name:        param,
				Type:        tp,
				Description: descr,
				IsOptional:  isOptional,
				IsRequired:  !isOptional,
			}

			obj.Fields = append(obj.Fields, field)
		})

		return
	}

	ctx.shared.err = fmt.Errorf("unknown table structure %s", theadText)
}

func parseType(n *html.Node) Type {
	t := Type{
		Name:    extractText(n),
		HasLink: false,
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if checkTag(c, "a") {
			t.HasLink = true
		}
	}

	return t
}

func Parse(r io.Reader, opts *ParseOpts) (*ParsedAPI, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	rootCtx := parseContext{
		opts: opts,
		api: &ParsedAPI{
			Chapters: make(map[string]*Chapter),
		},
		shared: &sharedContext{},
	}

	parseRecursive1(rootCtx, doc)

	return rootCtx.api, rootCtx.shared.err
}
