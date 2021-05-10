package libamabot

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	handlers_db = append(handlers_db, handler{
		id:   "uri",
		help: "encode and decode uri string.\n commands:\n`>>encodeuri [raw string]` `>>decodeuri [uri-parsed string]`",
		main: func(s *discordgo.Session, m *discordgo.MessageCreate) {
			defer func() {
				if err := recover(); err != nil {
					log.Println(err)
					s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
						Title:       "Error",
						Color:       0xff0000,
						Description: "``` " + strings.ReplaceAll(fmt.Sprint(err), "```", " `` ") + " ```",
					})
				}
			}()

			if m.Author.ID == s.State.User.ID {
				return
			}
			if m.Author.Bot {
				return
			}

			if len(m.Content) >= 11 {
				if m.Content[:11] == ">>encodeuri" {
					rawtext := strings.TrimSpace(m.Content[11:])
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```\n%s```", url.PathEscape(rawtext)))
				} else if m.Content[:11] == ">>decodeuri" {
					uritext := strings.TrimSpace(m.Content[11:])
					rawtext, err := url.PathUnescape(uritext)
					if err != nil {
						panic(err)
					}
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```\n%s```", rawtext))
				}
			}
		},
	})
}
