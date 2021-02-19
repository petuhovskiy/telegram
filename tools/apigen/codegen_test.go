package apigen

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodegen(t *testing.T) {
	f, err := os.Open("api.html")
	assert.Nil(t, err)
	defer f.Close()

	p, err := Parse(f, DefaultParseOpts)
	assert.Nil(t, err)

	err = Codegen(p, &GenOpts{
		PackageName: "telegram",
		Dest:        "../../",
		TypeExceptions: []TypeException{
			{
				Domain:     "",
				TypeString: "Integer or String",
				GoType:     "string",
			},
			{
				Domain:     "",
				TypeString: "InputFile or String",
				GoType:     "Fileable",
			},
			{
				Domain:     "",
				TypeString: "InlineKeyboardMarkup or ReplyKeyboardMarkup or ReplyKeyboardRemove or ForceReply",
				GoType:     "AnyKeyboard",
			},
			{
				Domain:     "",
				TypeString: "InputMessageContent",
				GoType:     "InputMessageContent",
			},
		},
		MethodExceptions: []MethodException{
			{
				Method:       "setWebhook",
				OverrideType: "json.RawMessage",
			},
			{
				Method:       "getUpdates",
				OverrideType: "[]Update",
			},
		},
		StructExceptions: []StructException{
			{
				StructName: "InlineQueryResult",
				Skip:       true,
			},
			{
				StructName: "InputMessageContent",
				Skip:       true,
			},
		},
	})
	assert.Nil(t, err)
}
