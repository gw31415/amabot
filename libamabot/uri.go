package libamabot

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	messageCmd(func(ctx context.Context, opts AmabotOptions, s *discordgo.Session, m *discordgo.MessageCreate) {
		prefix := opts.MessageCmdPrefix
		cmd := ctx.Value("cmd").(string)
		cmd_len := len(prefix) + len(cmd)
		if m.Content[cmd_len:cmd_len+6] == "encode" {
			rawtext := strings.TrimSpace(m.Content[cmd_len+6:])
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```\n%s```", url.PathEscape(rawtext)))
		} else if m.Content[cmd_len:cmd_len+6] == "decode" {
			uritext := strings.TrimSpace(m.Content[cmd_len+6:])
			rawtext, err := url.PathUnescape(uritext)
			if err != nil {
				panic(err)
			}
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```\n%s```", rawtext))
		}
	})
}
