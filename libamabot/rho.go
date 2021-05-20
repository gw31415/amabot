package libamabot

import (
	"fmt"
	"log"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mathlava/bigc/math/rho"
)

func init() {
	setString := func(n *big.Int, str string) (nw *big.Int) {
		nw, ok := n.SetString(strings.TrimSpace(str), 10)
		if ok {
			return nw
		} else {
			panic("parse error.")
		}
	}

	addHandler(
		func(s *discordgo.Session, m *discordgo.MessageCreate) {
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
				watch_time := make(chan *discordgo.MessageEmbed, 1)
				s.ChannelTyping(m.ChannelID)
				go func() {
					num_str := m.Content[5:]
					//数値にする
					num := new(big.Int)
					num = setString(num, num_str)
					//数値にならなければ
					list := rho.Primes(num)
					sort.Slice(list, func(i, j int) bool { return list[i].Cmp(list[j]) == -1 })
					result_bytes := make([]byte, 0)
					for i := 0; i < len(list)-1; i++ {
						result_bytes = append(result_bytes, list[i].String()...)
						result_bytes = append(result_bytes, "×"...)
					}
					result_bytes = append(result_bytes, list[len(list)-1].String()...)
					watch_time <- &discordgo.MessageEmbed{
						Title:       num_str,
						Description: string(result_bytes),
						Color:       0x00e6ff,
					}
				}()
				select {
				case messageEmbed := <-watch_time:
					s.ChannelMessageSendEmbed(m.ChannelID, messageEmbed)
				case <-time.After(0xfff * time.Millisecond):
					panic("timeout.")
				}
			}
		})
}
