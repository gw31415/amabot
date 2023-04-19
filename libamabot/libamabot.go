package libamabot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

type AmabotOptions struct {
	MessageCmdPrefix string
	TimeoutDuration  time.Duration
	EnabledHandlers  []string
	AppCmdGuildIds   []string
	Db               *gorm.DB
}

// Amabot instance
type Amabot struct {
	discord                  *discordgo.Session
	isRunning                bool
	opts                     AmabotOptions
	registeredAppCmdsInGuild map[string][]*discordgo.ApplicationCommand
}

// Create new Amabot instance
func New(token string, opts AmabotOptions) (*Amabot, error) {
	discord_session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return &Amabot{
		discord:                  discord_session,
		isRunning:                false,
		opts:                     opts,
		registeredAppCmdsInGuild: make(map[string][]*discordgo.ApplicationCommand),
	}, nil
}

// Open session. If Running, restart Amabot
func (ama *Amabot) Run() error {
	if ama.isRunning {
		ama.Close()
		return ama.Run()
	}
	for _, id := range ama.opts.EnabledHandlers {
		if h := handlers_db[id]; h == nil {
			if appcmd_handlers_db[id] == nil {
				// If there is a handler that is not in the handlers_db or appcmd_handlers_db
				return errors.New("handler not found named: " + id)
			}
		} else {
			for _, h := range h {
				ama.discord.AddHandler(h(ama.opts))
			}
		}
	}
	ama.discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		o := ama.opts
		defer func() {
			if err := recover(); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags: discordgo.MessageFlagsEphemeral,
						Embeds: []*discordgo.MessageEmbed{
							{
								Title:       "Error",
								Color:       0xff0000,
								Description: "``` " + strings.ReplaceAll(fmt.Sprint(err), "```", " `` ") + " ```",
							},
						},
					},
				})
			}
		}()
		watch_res := make(chan struct{})
		watch_err := make(chan interface{})
		watch_timeout, cancel := context.WithTimeout(context.Background(), o.TimeoutDuration)
		defer cancel()
		go func() {
			defer func() {
				if err := recover(); err != nil {
					watch_err <- err
				}
			}()
			appcmd_handlers_db[i.ApplicationCommandData().Name](watch_timeout, o, s, i)
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
	if err := ama.discord.Open(); err != nil {
		ama.isRunning = false
		return err
	} else {
		// Register ApplicationCommands
		for _, id := range ama.opts.EnabledHandlers {
			if cmds := appcmd_db[id]; cmds != nil {
				for _, cmd := range cmds {
					if cmd.GuildID == "" {
						var guilds []string
						if cmd.GuildID == "" {
							guilds = ama.opts.AppCmdGuildIds
						} else {
							guilds = []string{cmd.GuildID}
						}
						for _, guild := range guilds {
							cmd, err := ama.discord.ApplicationCommandCreate(ama.discord.State.User.ID, guild, cmd)
							if err != nil {
								log.Println(err)
							} else {
								log.Println("Registered Appcmd:", cmd.Name, guild)
								ama.registeredAppCmdsInGuild[guild] = append(ama.registeredAppCmdsInGuild[guild], cmd)
							}
						}
					}
				}
			}
		}
		ama.isRunning = true
		return nil
	}
}

// Close session of Amabot
func (ama *Amabot) Close() error {
	for _, guild := range ama.opts.AppCmdGuildIds {
		for _, cmd := range ama.registeredAppCmdsInGuild[guild] {
			err := ama.discord.ApplicationCommandDelete(ama.discord.State.User.ID, guild, cmd.ID)
			if err != nil {
				log.Println("Error Cleaning Appcmd:", cmd.Name, guild, err)
			} else {
				log.Println("Cleaned Appcmd:", cmd.Name, guild)
			}
		}
	}
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
	key_map := make(map[string]struct{})
	for k := range handlers_db {
		key_map[k] = struct{}{}
	}
	for k := range appcmd_handlers_db {
		key_map[k] = struct{}{}
	}
	keys := make([]string, 0, len(key_map))
	for k := range key_map {
		keys = append(keys, k)
	}
	return keys
}

// The database that all handlers are in.
// DO NOT EDIT DIRECTRY
var handlers_db map[string][](func(AmabotOptions) interface{})

// The database that all ApplicationCommand are in.
// DO NOT EDIT DIRECTRY
var appcmd_db map[string][]*discordgo.ApplicationCommand

// The database that all ApplicationCommand-handlers are in.
// DO NOT EDIT DIRECTRY
var appcmd_handlers_db map[string]func(context.Context, AmabotOptions, *discordgo.Session, *discordgo.InteractionCreate)

// Implement slash commands
// appCmd : ApplicationCommand to be registered
// handler : Handler to be registered
func slashCmd(appCmd *discordgo.ApplicationCommand, handler func(ctx context.Context, opts AmabotOptions, s *discordgo.Session, i *discordgo.InteractionCreate)) {
	// Get the module filename
	_, name, _, _ := runtime.Caller(1)
	name = filepath.Base(name[:len(name)-len(filepath.Ext(name))])

	// Setup appcmd_db
	if appcmd_db == nil {
		appcmd_db = make(map[string][]*discordgo.ApplicationCommand, 0)
	}
	if appcmd_db[name] == nil {
		appcmd_db[name] = make([]*discordgo.ApplicationCommand, 0)
	}
	appcmd_db[name] = append(appcmd_db[name], appCmd)

	// Setup handlers_db
	if appcmd_handlers_db == nil {
		appcmd_handlers_db = make(map[string]func(context.Context, AmabotOptions, *discordgo.Session, *discordgo.InteractionCreate))
	}
	appcmd_handlers_db[appCmd.Name] = handler
}

// Get the module filename
func getFileName() string {
	_, name, _, _ := runtime.Caller(2)                             // Get the module filename
	return filepath.Base(name[:len(name)-len(filepath.Ext(name))]) // Without Ext
}

// Implement simple commands with raw message content intent
// handler : Handler to be registered
func rawMessageCmd(handler func(ctx context.Context, opts AmabotOptions, s *discordgo.Session, m *discordgo.MessageCreate), errorOutput bool) {
	if handlers_db == nil {
		handlers_db = make(map[string][](func(AmabotOptions) interface{}), 0)
	}
	// Get the module filename
	name := getFileName()
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
			defer func() {
				if err := recover(); errorOutput && err != nil {
					s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
						Title:       "Error",
						Color:       0xff0000,
						Description: "``` " + strings.ReplaceAll(fmt.Sprint(err), "```", " `` ") + " ```",
					})
				}
			}()
			watch_res := make(chan struct{})
			watch_err := make(chan interface{})
			watch_timeout, c1 := context.WithTimeout(context.Background(), o.TimeoutDuration)
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

// Implement simple commands with message content intent
// handler : Handler to be registered
func messageCmd(handler func(ctx context.Context, opts AmabotOptions, s *discordgo.Session, m *discordgo.MessageCreate)) {
	rawMessageCmd(func(ctx context.Context, opts AmabotOptions, s *discordgo.Session, m *discordgo.MessageCreate) {
		name := getFileName()
		// The written file name becomes a command.
		// ex rho.go; messageCmdPrefix == ">>" -> the command pattern is `>>rho`
		if len(m.Content) < len(name)+len(opts.MessageCmdPrefix) {
			return
		}
		if m.Content[:len(name)+len(opts.MessageCmdPrefix)] != opts.MessageCmdPrefix+name {
			return
		}
		handler(ctx, opts, s, m)
	}, true)
}
