package libamabot

import (
	"fmt"
	"log"
	"math/big"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/mathlava/bigc/math/rho"
)

func SetString(n *big.Int, str string) (nw *big.Int) {
	nw, ok := n.SetString(str, 10)
	if ok {
		return nw
	} else {
		panic("parse error.")
	}
}

func init() {
	handlers_db = append(handlers_db, handler{
		id:   "rho",
		help: "",
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
			if len(m.Content) < 5 {
				return
			}
			if m.Content[:5] == ">>rho" {
				num_str := strings.TrimSpace(m.Content[5:])
				//数値にする
				num := new(big.Int)
				num = SetString(num, num_str)
				//数値にならなければ
				list := rho.Primes(num)
				sort.Slice(list, func(i, j int) bool { return list[i].Cmp(list[j]) == -1 })
				result_bytes := make([]byte, 0)
				for i := 0; i < len(list)-1; i++ {
					result_bytes = append(result_bytes, list[i].String()...)
					result_bytes = append(result_bytes, "×"...)
				}
				result_bytes = append(result_bytes, list[len(list)-1].String()...)
				s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
					Title: num_str,
					Description: string(result_bytes),
					Color: 0x00e6ff,
				})
			}

		},
	})
}
