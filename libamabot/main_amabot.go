package libamabot

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

// Amabot instance
type Amabot struct {
	discord     *discordgo.Session
	isRunning   bool
	handlers_on []string
}

// Create new Amabot instance
func New(token string) (*Amabot, error) {
	discord_session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return &Amabot{
		discord:     discord_session,
		isRunning:   false,
		handlers_on: GetAllHandlersList(),
	}, nil
}

// Open session. If Running, restart Amabot
func (ama *Amabot) Run() error {
	if ama.isRunning {
		ama.Close()
		return ama.Run()
	}
	for _, id := range ama.handlers_on {
		if h := handlers_db[id]; h == nil {
			return errors.New("handler not found named: " + id)
		} else {
			for _, h := range h {
			ama.discord.AddHandler(h)
			}
		}
	}
	if err := ama.discord.Open(); err != nil {
		ama.isRunning = false
		return err
	} else {
		ama.isRunning = true
		return nil
	}
}

// Close session of Amabot
func (ama *Amabot) Close() error {
	if ama.isRunning {
		return nil
	}
	if ama.discord.Close() != nil {
		ama.isRunning = false
		return nil
	}
	return errors.New("unknown error on Closing.")
}

// Update the list of enabled Handlers
func (ama *Amabot) UpdateEnabledHandlers(ids ...string) error {
	runned := ama.isRunning
	if runned {
		if err := ama.Close(); err != nil {
			return err
		}
		if ama.isRunning {
			return errors.New("unreachable error.")
		}
	}
	ama.handlers_on = ids
	if runned {
		ama.Run()
	}
	return nil
}

// Get the list of id of enabled handlers
func (ama *Amabot) GetEnabledHandlers() []string {
	return ama.handlers_on
}
