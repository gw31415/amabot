package libamabot

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Amabot instance
type Amabot struct {
	discord          *discordgo.Session
	isRunning        bool
	enabled_handlers []string
}

// Create new Amabot instance
func New(token string) (*Amabot, error) {
	discord_session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return &Amabot{
		discord:          discord_session,
		isRunning:        false,
		enabled_handlers: GetAllHandlersList(),
	}, nil
}

// Open session. If Running, restart Amabot
func (ama *Amabot) Run() error {
	if ama.isRunning {
		ama.Close()
		return ama.Run()
	}
	for _, id := range ama.enabled_handlers {
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
	return errors.New("unknown error on Closing")
}

// Update the list of enabled Handlers
func (ama *Amabot) UpdateEnabledHandlers(ids ...string) error {
	runned := ama.isRunning
	if runned {
		if err := ama.Close(); err != nil {
			return err
		}
	}
	ama.enabled_handlers = ids
	if runned {
		ama.Run()
	}
	return nil
}

// Get the list of id of enabled handlers
func (ama *Amabot) GetEnabledHandlers() []string {
	return ama.enabled_handlers
}

// Get the list of all handlers
func GetAllHandlersList() []string {
	keys := make([]string, 0, len(handlers_db))
	for k := range handlers_db {
		keys = append(keys, k)
	}
	return keys
}

// The database that all handlers are in.
// DO NOT EDIT DIRECTRY
var handlers_db map[string][]interface{}

// Implement simple commands with message content intent
// handler : Handler to be registered
func messageCmd(handler func(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate)) {
	if handlers_db == nil {
		handlers_db = make(map[string][]interface{}, 0)
	}
	_, name, _, _ := runtime.Caller(1) // Get the module filename
	name = filepath.Base(name[:len(name)-len(filepath.Ext(name))])
	if handlers_db[name] == nil {
		handlers_db[name] = make([]interface{}, 0)
	}
	handlers_db[name] = append(handlers_db[name], func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		if m.Author.Bot {
			return
		}
		// The written file name becomes a command.
		// ex rho.go -> the command pattern is `>>rho`
		if len(m.Content) < len(name)+2 {
			return
		}
		if m.Content[:len(name)+2] != ">>"+name {
			return
		}
		defer func() {
			if err := recover(); err != nil {
				s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
					Title:       "Error",
					Color:       0xff0000,
					Description: "``` " + strings.ReplaceAll(fmt.Sprint(err), "```", " `` ") + " ```",
				})
			}
		}()
		watch_res := make(chan struct{})
		watch_err := make(chan interface{})
		childCtx, cancel := context.WithTimeout(context.Background(), 0xfff*time.Millisecond)
		defer cancel()
		watch_timeout, cancel2 := context.WithCancel(childCtx)
		defer cancel2()
		go func() {
			defer func() {
				if err := recover(); err != nil {
					watch_err <- err
				}
			}()
			handler(childCtx, s, m)
			close(watch_res)
		}()
		select {
		case err := <-watch_err:
			panic(err)
		case <-watch_res:
			return
		case <-watch_timeout.Done():
			panic("timeout")
		}
	})
}
