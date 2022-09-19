package libamabot

import (
	"context"
	"math/big"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/mathlava/bigc/math/rho"
)

func init() {
	messageCmd(
		func(ctx context.Context, opts AmabotOptions, s *discordgo.Session, m *discordgo.MessageCreate) {
			s.ChannelTyping(m.ChannelID)
			prefix := opts.GetMessageCmdPrefix()
			cmd := ctx.Value("cmd").(string)
			num_str := m.Content[len(prefix)+len(cmd):]
			//数値にする
			num := new(big.Int)
			num, ok := num.SetString(strings.TrimSpace(num_str), 10)
			if !ok {
				//数値にならなければ
				panic("parse error")
			}
			childCtx, cancel := context.WithCancel(ctx)
			defer cancel()
			ch := rho.PrimesAsync(childCtx, num)
			defer close(ch)
			var list []*big.Int
			select {
			case <-ctx.Done():
				return
			case list = <-ch:
				sort.Slice(list, func(i, j int) bool { return list[i].Cmp(list[j]) == -1 })
				result_bytes := make([]byte, 0)
				for i := 0; i < len(list)-1; i++ {
					result_bytes = append(result_bytes, list[i].String()...)
					result_bytes = append(result_bytes, "×"...)
				}
				result_bytes = append(result_bytes, list[len(list)-1].String()...)
				_, e := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
					Title:       num_str,
					Description: string(result_bytes),
					Color:       0x00e6ff,
				})
				if e != nil {
					panic(e)
				} else {
					return
				}
			}
		})
}
