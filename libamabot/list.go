package libamabot

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	handlers_db = append(handlers_db, handler{
		id:   "list",
		help: "list all handler ids.",
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
				if m.Content[:6] == ">>list" {
					list := make([]string, 0)
					// send the list of handlers
					for _, h := range handlers_db {
						list = append(list, h.id)
					}
					sort.Strings(list)
					s.ChannelMessageSend(m.ChannelID, fmt.Sprint(list))
				}
			}
		},
	})
}
