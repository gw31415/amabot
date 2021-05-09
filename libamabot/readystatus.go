package libamabot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func init() {
	handlers_db = append(handlers_db, handler{
		id:   "readystatus",
		help: "set ready status",
		main: func(s *discordgo.Session, r *discordgo.Ready) {
			s.UpdateGameStatus(0, "The Answer to Life, the Universe, Everything")
			log.Println("Amabot is ready.")
		},
	})
}
