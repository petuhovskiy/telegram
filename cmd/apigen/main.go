package main

import (
	"github.com/petuhovskiy/telegram/tools/apigen"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	f, err := os.Open("api.html")
	if err != nil {
		log.WithError(err).Fatal("failed to read protocol html")
	}
	defer f.Close()

	p, err := apigen.Parse(f, apigen.DefaultParseOpts)
	if err != nil {
		log.WithError(err).Fatal("failed to parse protocol")
	}

	err = apigen.Codegen(p, &apigen.GenOpts{
		PackageName: "telegram",
		Dest:        "../../",
		TypeExceptions: []apigen.TypeException{
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
		MethodExceptions: []apigen.MethodException{
			{
				Method:       "setWebhook",
				OverrideType: "json.RawMessage",
			},
			{
				Method:       "getUpdates",
				OverrideType: "[]Update",
			},
		},
		StructExceptions: []apigen.StructException{
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
	if err != nil {
		log.WithError(err).Fatal("failed to generate code")
	}
}
