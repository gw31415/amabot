package libamabot

import (
	"github.com/bwmarrin/discordgo"
)

func init() {
	addHandler(&handler{
		help: "catch ping then reply pong, and vice versa",
		main: func(s *discordgo.Session, m *discordgo.MessageCreate) {
			if m.Author.ID == s.State.User.ID {
				return
			}
			if m.Author.Bot {
				return
			}
			if m.Content == "ping" {
				s.ChannelMessageSend(m.ChannelID, "Pong!")
			}

			if m.Content == "pong" {
				s.ChannelMessageSend(m.ChannelID, "Ping!")
			}
		},
	})
}
