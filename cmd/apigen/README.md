# apigen

Special command to generate client implementation. It parses html from telegram website, which describes the bot api, and generates code for api types and methods.

To update protocol:

```bash
curl https://core.telegram.org/bots/api/ > api.html
go run main.go
```