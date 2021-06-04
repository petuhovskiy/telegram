package apigen

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestExtractText(t *testing.T) {
	check := func(text string, src string) {
		n, err := html.Parse(strings.NewReader(src))
		assert.Nil(t, err)
		assert.Equal(t, text, extractText(n))
	}

	check("July 29, 2019", `<h4><a class="anchor" name="july-29-2019" href="#july-29-2019" id="july-29-2019"><i class="anchor-icon"></i></a>July 29, 2019</h4>`)
	check("Making requests", `<h3><a class="anchor" name="making-requests" href="#making-requests" id="making-requests"><i class="anchor-icon"></i></a>Making requests</h3>`)
	check("This object represents an incoming update.\nAt most one of the optional parameters can be present in any given update.", `<p>This <a href="#available-types">object</a> represents an incoming update.<br>At most <strong>one</strong> of the optional parameters can be present in any given update.</p>`)
	check("ParameterTypeRequiredDescription", `<thead><tr><th>Parameter</th><th>Type</th><th>Required</th><th>Description</th></tr></thead>`)
	check("FieldTypeDescription", `<thead><tr><th>Field</th><th>Type</th><th>Description</th></tr></thead>`)
	check("- InputMediaAnimation\n- InputMediaDocument\n- InputMediaAudio\n- InputMediaPhoto\n- InputMediaVideo", `<ul><li><a href="#inputmediaanimation">InputMediaAnimation</a></li><li><a href="#inputmediadocument">InputMediaDocument</a></li><li><a href="#inputmediaaudio">InputMediaAudio</a></li><li><a href="#inputmediaphoto">InputMediaPhoto</a></li><li><a href="#inputmediavideo">InputMediaVideo</a></li></ul>`)
}
