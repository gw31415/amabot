package libamabot

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	handlers_db = append(handlers_db, handler{
		id:   "hash",
		help: "get hash.\n commands:\n`>>sha256sum [raw string]` `>>sha1sum [raw string]`",
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
				if m.Content[:11] == ">>sha256sum" {
					rawtext := strings.TrimSpace(m.Content[11:])
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```\n%x```", sha256.Sum256([]byte(rawtext))))
				} else if m.Content[:9] == ">>sha1sum" {
					rawtext := strings.TrimSpace(m.Content[9:])
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```\n%x```", sha1.New().Sum([]byte(rawtext))))
				}
			}
		},
	})
}
