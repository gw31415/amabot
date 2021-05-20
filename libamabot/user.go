package libamabot

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	addHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
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
			switch m.Content {
			case ">>id":
				s.ChannelMessageSend(m.ChannelID, m.Author.ID)
			case ">>mfa":
				if m.Author.MFAEnabled {
					s.ChannelMessageSend(m.ChannelID, "true")
				} else {
					s.ChannelMessageSend(m.ChannelID, "false")
				}
			}
		})
}
