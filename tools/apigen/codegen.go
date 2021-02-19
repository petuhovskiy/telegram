package apigen

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/dave/jennifer/jen"
	log "github.com/sirupsen/logrus"
)

type GenOpts struct {
	PackageName      string
	Dest             string
	TypeExceptions   []TypeException
	MethodExceptions []MethodException
	StructExceptions []StructException
}

type TypeException struct {
	Domain     string // if Domain matches type's location prefix
	TypeString string // and TypeString matches type
	GoType     string // use GoType instead
}

type MethodException struct {
	Method       string
	OverrideType string
}

type StructException struct {
	StructName string
	Skip       bool
}

func ChapterNameToFilename(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	return name + "_gen.go"
}

func fixURLSuffix(name string) string {
	if strings.HasSuffix(name, "Url") {
		name = strings.TrimSuffix(name, "Url")
		name += "URL"
	}
	return name
}

func TypeNameToGo(name string) (string, error) {
	if !unicode.IsUpper(rune(name[0])) {
		return "", fmt.Errorf("typename %s must be capitalized", name)
	}

	name = fixURLSuffix(name)
	return name, nil
}

func FieldToGo(name string) (string, error) {
	arr := strings.Split(name, "_")

	var res string
	for _, s := range arr {
		if strings.EqualFold(s, "id") {
			res += "ID"
			continue
		}

		if strings.EqualFold(s, "url") {
			res += "URL"
			continue
		}

		res += strings.Title(s)
	}

	return res, nil
}

func TypeToGo(t Type) (string, error) {
	if strings.HasPrefix(t.Name, "Array of ") {
		t.Name = strings.TrimPrefix(t.Name, "Array of ")

		res, err := TypeToGo(t)
		if err != nil {
			return "", err
		}

		res = strings.TrimPrefix(res, "*")
		res = "[]" + res
		return res, nil
	}

	switch t.Name {
	case "String":
		return "string", nil

	case "Float number":
		return "float64", nil

	case "Integer":
		return "int", nil

	case "Boolean":
		return "bool", nil

	case "Float":
		return "float64", nil

	case "True":
		return "bool", nil
	}

	if strings.Contains(t.Name, " ") {
		return "", fmt.Errorf("found whitespace in typename \"%s\"", t.Name)
	}

	var res string
	if t.HasLink {
		res += "*"
	}
	res += t.Name

	if !t.HasLink {
		log.Warn("Unrecognized simple type: ", t.Name)
	}

	res = fixURLSuffix(res)
	return res, nil
}

func FieldToCode(f Field, objectName string, opts *GenOpts) (jen.Code, error) {
	fieldName, err := FieldToGo(f.Name)
	if err != nil {
		return nil, err
	}

	var fieldType string
	var typeException bool
	for _, ex := range opts.TypeExceptions {
		if ex.TypeString != f.Type.Name {
			continue
		}

		domain := fmt.Sprintf("%s$%s", objectName, f.Name)
		if strings.HasPrefix(domain, ex.Domain) {
			typeException = true
			fieldType = ex.GoType
			break
		}
	}

	if !typeException {
		fieldType, err = TypeToGo(f.Type)
		if err != nil {
			return nil, err
		}
	}

	jsonTag := f.Name
	if f.IsOptional {
		jsonTag += ",omitempty"
	}

	return jen.Id(fieldName).Id(fieldType).Tag(map[string]string{"json": jsonTag}), nil
}

func FuncNameToGo(name string) (string, error) {
	return strings.Title(name), nil
}

func CodegenChapter(chap *Chapter, opts *GenOpts) (*jen.File, error) {
	f := jen.NewFile(opts.PackageName)

	for _, obj := range chap.Objects {
		if obj.IsType && obj.IsFunction {
			return nil, fmt.Errorf("impossible obj.IsType && obj.IsFunction. %#v", obj)
		}

		if obj.IsType {
			err := CodegenStruct(obj, f, opts)
			if err != nil {
				return nil, err
			}

			continue
		}

		if obj.IsFunction {
			err := CodegenFunc(obj, f, opts)
			if err != nil {
				return nil, err
			}

			continue
		}

		log.Warn("Unknown object is ignored ", obj)
	}

	return f, nil
}

func CodegenStruct(obj *Object, f *jen.File, opts *GenOpts) error {
	name := obj.Name

	typeName, err := TypeNameToGo(name)
	if err != nil {
		return err
	}

	for _, ex := range opts.StructExceptions {
		if ex.StructName != typeName {
			continue
		}

		if ex.Skip {
			log.WithField("typeName", typeName).Info("skipping struct code generation")
			return nil
		}
	}

	var fields []jen.Code
	var skipType bool
	for i, f := range obj.Fields {
		if i != 0 {
			fields = append(fields, jen.Line())
		}

		if f.Description != "" {
			commentLines := processComments([]string{f.Description})
			for _, ln := range commentLines {
				fields = append(fields, jen.Comment(ln))
			}
		}

		field, err := FieldToCode(f, name, opts)
		if err != nil {
			log.Warn("Error while processing type. ", err)
			skipType = true
			continue
		}

		fields = append(fields, field)
	}

	if skipType {
		log.Warnf("Struct %s skipped!", typeName)
		return nil
	}

	commentLines := processComments(obj.Notes)
	for _, ln := range commentLines {
		f.Comment(ln)
	}

	f.Type().Id(typeName).Struct(fields...)
	f.Line()

	return nil
}

func CodegenFunc(obj *Object, f *jen.File, opts *GenOpts) error {
	name := obj.Name

	funcName, err := FuncNameToGo(name)
	if err != nil {
		return err
	}

	var fields []jen.Code
	var skipFunc bool
	for i, f := range obj.Fields {
		if i != 0 {
			fields = append(fields, jen.Line())
		}

		if f.Description != "" {
			descr := f.Description
			if f.IsOptional {
				descr = "Optional. " + descr
			}
			commentLines := processComments([]string{descr})
			for _, ln := range commentLines {
				fields = append(fields, jen.Comment(ln))
			}
		}

		field, err := FieldToCode(f, name, opts)
		if err != nil {
			log.Warn("Error while processing type. ", err)
			skipFunc = true
			continue
		}

		fields = append(fields, field)
	}

	if skipFunc {
		log.Warnf("Func %s skipped!", funcName)
		return nil
	}

	requestType := funcName + "Request"

	f.Type().Id(requestType).Struct(fields...)
	f.Line()

	commentLines := processComments(obj.Notes)
	for _, ln := range commentLines {
		f.Comment(ln)
	}

	tmp := f.Func().Params(
		jen.Id("b").Id("*Bot"),
	).Id(funcName).Params(
		jen.Id("req").Id("*" + requestType),
	)

	returnType := obj.ReturnType
	for _, exc := range opts.MethodExceptions {
		if exc.Method == name {
			if exc.OverrideType != "" {
				returnType = exc.OverrideType
			}
		}
	}

	if returnType == "json.RawMessage" {
		returnType = ""
	}

	switch returnType {
	case "":
		tmp.Params(
			jen.Qual("encoding/json", "RawMessage"),
			jen.Id("error"),
		).Block(
			jen.Return(jen.Id("b.makeRequest").Call(jen.Lit(name), jen.Id("req"))),
		)

	default:
		tmp.Params(
			jen.Id("*"+returnType),
			jen.Id("error"),
		).Block(
			jen.List(jen.Id("j"), jen.Id("err")).
				Op(":=").Id("b.makeRequest").
				Call(jen.Lit(name), jen.Id("req")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.Var().Id("resp").Id(returnType),
			jen.Err().
				Op("=").
				Qual("encoding/json", "Unmarshal").
				Call(jen.Id("j"), jen.Id("&resp")),
			jen.Return(jen.Id("&resp"), jen.Err()),
		)
	}

	return nil
}

func Codegen(api *ParsedAPI, opts *GenOpts) error {
	for _, chap := range api.Chapters {
		f, err := CodegenChapter(chap, opts)
		if err != nil {
			return err
		}

		filename := ChapterNameToFilename(chap.Name)
		file, err := os.Create(opts.Dest + filename)
		if err != nil {
			return err
		}

		fmt.Fprint(file, "// Code generated by telegram-apigen. DO NOT EDIT.\n\n")

		err = f.Render(file)
		if err != nil {
			return err
		}
	}

	return nil
}
