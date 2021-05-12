package libamabot

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	addHandler(&handler{
		help: "show the help of each handler.\n**Command:**\n`>>help [handler name]`\n*cf:* if you want to list of handlers, please enter `>>list` command.",
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
			if len(m.Content) >= 6 {
				if m.Content[:6] == ">>help" {
					var handler_id string
					if strings.TrimSpace(m.Content) == ">>help" {
						handler_id = "help"
					} else {
						handler_id = strings.TrimSpace(m.Content[6:])
					}
					// search and send
					if help, err := GetHelp(handler_id); err == nil {
						s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
							Title:       fmt.Sprintf("Help of handler: `%s`", handler_id),
							Description: help,
						})
					} else {
						panic(err)
					}
				}
			}
		},
	})
}
