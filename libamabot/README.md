# libamabot

[![GoDoc](https://godoc.org/github.com/gw31415/amabot/libamabot?status.svg)](https://godoc.org/github.com/gw31415/amabot/libamabot)

An Discord bot Instance

## For Developers
### handler

```go
type handler struct {
	id    string
	help  string
	main  interface{}
}
```

handler object.

`id` is string to identificate each handler.
`help` is help text to help users.
`main` is handler function. This is to be passed to discordgo.Session through function: [AddHandler](https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.AddHandler).

### handlers

```go
var handlers_db []handlers
```

All database of handlers. Add handlers here.
