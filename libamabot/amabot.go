package libamabot

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

// Amabot instance
type Amabot struct {
	discord   *discordgo.Session
	isRunning bool
}

// Create new Amabot instance
func New() (*Amabot, error) {
	discord_session, err := discordgo.New("Bot " + viper.GetString("token"))
	if err != nil {
		return nil, err
	}
	return &Amabot{
		discord:   discord_session,
		isRunning: false,
	}, nil
}

// Open session
func (ama *Amabot) Run() error {
	if ama.isRunning {
		return errors.New("this Amabot is already running.")
	}
	for _, h := range handlers_db {
		ama.discord.AddHandler(h.main)
	}
	if err := ama.discord.Open(); err != nil {
		ama.isRunning = false
		return err
	} else {
		ama.isRunning = true
		return nil
	}
}

// Get whether Amabot is running
func (ama *Amabot) IsRunning() bool {
	return ama.isRunning
}

// Restart the session
func (ama *Amabot) Restart() error {
	if !ama.isRunning {
		return errors.New("Amabot is not running")
	}
	if err := ama.Close(); err != nil {
		return err
	}
	return ama.Run()
}

// Close session of Amabot
func (ama *Amabot) Close() error {
	return ama.discord.Close()
}

// handler to catch request
type handler struct {
	/**
	  func(*discordgo.Session, *discordgo.MessageCreate)
	  |
	  func(*discordgo.Session, *discordgo.PresenceUpdate)
	  **/
	main interface{}
	// Identify name of the handler
	id   string
	// Help message of the handler
	help string
}

//the slice of handlers
var handlers_db []handler
