# libamabot

[![GoDoc](https://godoc.org/github.com/gw31415/amabot/libamabot?status.svg)](https://godoc.org/github.com/gw31415/amabot/libamabot)

An Discord bot Instance

## For Developers
### handler

```go
type handler struct {
	help  string
	main  interface{}
}
```

handler object.

`help` is help text to help users.
`main` is handler function. This is to be passed to discordgo.Session through function: [AddHandler](https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.AddHandler).

### handlers

When it comes to add new handler to Amabot-system, please add new instance of handler in an `init` function through `addHandler` function.

```go
func addHandler(h *handler)
```
