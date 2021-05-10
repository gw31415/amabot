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
		id:   "encodeuri",
		help: "encode string data into uri.",
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
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```\n%s```",url.PathEscape(rawtext)))
				}
			}
		},
	})
}
