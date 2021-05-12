package libamabot

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	addHandler(&handler{
		help: "list all handler ids.\n**Command:** `>>list`",
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
					list := GetAllHandlersList() /// <- !
					// send the list of handlers
					sort.Strings(list)
					s.ChannelMessageSend(m.ChannelID, fmt.Sprint(list))
				}
			}
		},
	})
}
