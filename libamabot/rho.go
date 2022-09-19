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
	slashCmd(&discordgo.ApplicationCommand{
		Name:        "rho",
		Description: "Factorize the given integer into prime factors",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "integer",
				Description: "Integer to prime-factorize",
				Required:    true,
			},
		},
	}, func(ctx context.Context, opts AmabotOptions, s *discordgo.Session, i *discordgo.InteractionCreate) {
		options := i.ApplicationCommandData().Options
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}
		num_str := optionMap["integer"].StringValue()

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
				result_bytes = append(result_bytes, " × "...)
			}
			result_bytes = append(result_bytes, list[len(list)-1].String()...)
			e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Prime-factorize: " + num_str,
							Description: string(result_bytes),
							Color:       0x00e6ff,
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

	})
}
