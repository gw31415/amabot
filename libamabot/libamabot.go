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

type AmabotOptions struct {
	messageCmdPrefix string
	timeoutDuration  time.Duration
	enabledHandlers  []string
}

func NewAmabotOptions() AmabotOptions {
	return AmabotOptions{
		messageCmdPrefix: ">>",
		timeoutDuration:  2 * time.Second,
		enabledHandlers:  GetAllHandlersList(),
	}
}

func (opts *AmabotOptions) GetMessageCmdPrefix() string {
	return opts.messageCmdPrefix
}
func (opts *AmabotOptions) SetMessageCmdPrefix(prefix string) {
	opts.messageCmdPrefix = prefix
}
func (opts *AmabotOptions) GetTimeoutDuration() time.Duration {
	return opts.timeoutDuration
}
func (opts *AmabotOptions) SetTimeoutDuration(duration time.Duration) {
	opts.timeoutDuration = duration
}
func (opts *AmabotOptions) GetEnabledHandlers() []string {
	return opts.enabledHandlers
}
func (opts *AmabotOptions) SetEnabledHandlers(handlers []string) {
	opts.enabledHandlers = handlers
}

// Amabot instance
type Amabot struct {
	discord   *discordgo.Session
	isRunning bool
	opts      AmabotOptions
}

// Create new Amabot instance
func New(token string, opts AmabotOptions) (*Amabot, error) {
	discord_session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return &Amabot{
		discord:   discord_session,
		isRunning: false,
		opts:      opts,
	}, nil
}

// Open session. If Running, restart Amabot
func (ama *Amabot) Run() error {
	if ama.isRunning {
		ama.Close()
		return ama.Run()
	}
	for _, id := range ama.opts.GetEnabledHandlers() {
		if h := handlers_db[id]; h == nil {
			return errors.New("handler not found named: " + id)
		} else {
			for _, h := range h {
				ama.discord.AddHandler(h(ama.opts))
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

// Update AmabotOptions
func (ama *Amabot) UpdateOptions(opts AmabotOptions) error {
	runned := ama.isRunning
	if runned {
		if err := ama.Close(); err != nil {
			return err
		}
	}
	ama.opts = opts
	if runned {
		ama.Run()
	}
	return nil
}

// Get AmabotOptions
func (ama *Amabot) GetOptions() AmabotOptions {
	return ama.opts
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
var handlers_db map[string][](func(AmabotOptions) interface{})

// Implement simple commands with message content intent
// handler : Handler to be registered
func messageCmd(handler func(ctx context.Context, opts AmabotOptions, s *discordgo.Session, m *discordgo.MessageCreate)) {
	if handlers_db == nil {
		handlers_db = make(map[string][](func(AmabotOptions) interface{}), 0)
	}
	_, name, _, _ := runtime.Caller(1) // Get the module filename
	name = filepath.Base(name[:len(name)-len(filepath.Ext(name))])
	if handlers_db[name] == nil {
		handlers_db[name] = make([](func(AmabotOptions) interface{}), 0)
	}
	handlers_db[name] = append(handlers_db[name], func(o AmabotOptions) interface{} {
		return func(s *discordgo.Session, m *discordgo.MessageCreate) {
			if m.Author.ID == s.State.User.ID {
				return
			}
			if m.Author.Bot {
				return
			}
			// The written file name becomes a command.
			// ex rho.go; messageCmdPrefix == ">>" -> the command pattern is `>>rho`
			if len(m.Content) < len(name)+len(o.GetMessageCmdPrefix()) {
				return
			}
			if m.Content[:len(name)+len(o.GetMessageCmdPrefix())] != o.GetMessageCmdPrefix()+name {
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
			watch_timeout, c1 := context.WithTimeout(context.Background(), o.GetTimeoutDuration())
			defer c1()
			childCtx := context.WithValue(watch_timeout, "cmd", name)
			go func() {
				defer func() {
					if err := recover(); err != nil {
						watch_err <- err
					}
				}()
				handler(childCtx, o, s, m)
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
		}
	})
}
