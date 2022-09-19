package libamabot

import (
	"context"
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/image/colornames"
)

func rgbaToInt(c color.RGBA) int {
	return 0x010000*int(c.R) + 0x000100*int(c.G) + 0x000001*int(c.B)
}

func init() {
	messageCmd(func(ctx context.Context, opts AmabotOptions, s *discordgo.Session, m *discordgo.MessageCreate) {
		prefix := opts.GetMessageCmdPrefix()
		cmd := ctx.Value("cmd").(string)
		arg := strings.TrimSpace(m.Content[len(prefix)+len(cmd):])
		for _, name := range colornames.Names {
			if arg == name {
				color := colornames.Map[name]
				s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
					Title:       name,
					Description: fmt.Sprintf("ColorName: `%s`\nColorCode: `#%02x%02x%02x`", name, color.R, color.G, color.B),
					Color:       rgbaToInt(color),
				})
				return
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
			s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
				Title:       arg,
				Description: fmt.Sprintf("ColorName: `%s`\nColorCode: `%s`", name, arg),
				Color:       int(c_i),
			})
			return
		}
		panic("color not found: " + arg)
	})
}
