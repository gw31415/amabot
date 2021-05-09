package libamabot

import (
	"fmt"
	"image/color"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/image/colornames"
)

func rgbaToInt(c color.RGBA) int {
	return 0x010000*int(c.R) + 0x000100*int(c.G) + 0x000001*int(c.B)
}

func init() {
	handlers_db = append(handlers_db, handler{
		id:   "color",
		help: "check color",
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
			if len(m.Content) < 7 {
				return
			}
			if m.Content[:7] == ">>color" {
				arg := strings.TrimSpace(m.Content[7:])
				for _, name := range colornames.Names {
					if arg == name {
						color := colornames.Map[name]
						s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
							Title:       name,
							Description: fmt.Sprintf("ColorName: %s\nColorCode: #%02x%02x%02x", name, color.R, color.G, color.B),
							Color:       rgbaToInt(color),
						})
						return
					}
				}
				panic("color not found: " + arg)
			}
		},
	})
}
