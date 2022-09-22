package libamabot

import (
	"context"
	"fmt"
	"image/color"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/image/colornames"
)

func rgbaToInt(c color.RGBA) int {
	return 0x010000*int(c.R) + 0x000100*int(c.G) + 0x000001*int(c.B)
}

func init() {
	slashCmd(&discordgo.ApplicationCommand{
		Name:        "color",
		Description: "Get color information",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "color",
				Description: "Color name or code",
				Required:    true,
			},
		},
	}, func(ctx context.Context, opts AmabotOptions, s *discordgo.Session, i *discordgo.InteractionCreate) {
		options := i.ApplicationCommandData().Options
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}
		arg := optionMap["color"].StringValue()
		for _, name := range colornames.Names {
			if arg == name {
				color := colornames.Map[name]
				e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Title:       name,
								Description: fmt.Sprintf("ColorName: `%s`\nColorCode: `#%02x%02x%02x`", name, color.R, color.G, color.B),
								Color:       rgbaToInt(color),
							},
						},
					},
				})
				if e != nil {
					panic(e)
				} else {
					return
				}
			}
		}
		if len(arg) == 7 && arg[0] == '#' {
			c_s := arg[1:]
			c_i, err := strconv.ParseInt(c_s, 16, 0)
			if err != nil {
				panic(err)
			}
			if 0 > c_i || c_i > 0xffffff {
				panic("overflow")
			}
			name := "(none)"
			for _, color_index := range colornames.Names {
				if rgbaToInt(colornames.Map[color_index]) == int(c_i) {
					name = color_index
					break
				}
			}
			e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       arg,
							Description: fmt.Sprintf("ColorName: `%s`\nColorCode: `%s`", name, arg),
							Color:       int(c_i),
						},
					},
				},
			})
			if e != nil {
				panic(e)
			} else {
				return
			}
		}
		panic("color not found: " + arg)
	})
}
