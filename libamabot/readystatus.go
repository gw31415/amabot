package libamabot

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

func init() {
	addHandler(&handler{
		help: "set ready status. no command is provided.",
		main: func(s *discordgo.Session, r *discordgo.Ready) {
			s.UpdateGameStatus(0, viper.GetString("game-status"))
			log.Println("Amabot is ready.")
		},
	})
}
