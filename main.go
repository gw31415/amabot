/*
Copyright © 2021 Amadeus_vn

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gw31415/amabot/libamabot"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	rootCmd = &cobra.Command{
		Use:   "amabot",
		Short: "amabot is a Discord Bot developed by @Amadeus_vn for personal use.",
	}
	// Path of the configuration file for this tool
	config_path = ""
)

func init() {
	// Flags of rootCmd
	rootCmd.PersistentFlags().StringVar(&config_path, "config-path", "", "Path of the configuration file for this tool")

	// Flags of Amabot
	rootCmd.PersistentFlags().StringP("token", "t", "", "Value of Discord API token")
	rootCmd.PersistentFlags().String("opts-prefix", ">>", "Prefix to fire MessageCmds")
	rootCmd.PersistentFlags().Duration("opts-timeout", 2*time.Second, "Set timeout duration")
	rootCmd.PersistentFlags().StringSlice("opts-guilds", make([]string, 0), "GuildIds to register ApplicationCommands")
	rootCmd.PersistentFlags().String("opts-sqlite", "amabot.sqlite3", "Sqlite3 Database-file to save data")

	// Read configuration file when it exists.
	cobra.OnInitialize(func() {
		if config_path != "" {
			_, err := os.Stat(config_path)
			cobra.CheckErr(err)
			viper.SetConfigFile(config_path)
		} else {
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)
			viper.AddConfigPath(".")
			viper.AddConfigPath(home)
			viper.SetConfigType("yaml")
			viper.SetConfigName("amabot")
		}
		if err := viper.ReadInConfig(); err == nil {
			log.Println("Using config file:", viper.ConfigFileUsed())
		}

		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
		viper.SetEnvPrefix("AMABOT")
		viper.AutomaticEnv()

		viper.BindPFlag("token", rootCmd.Flags().Lookup("token"))
		viper.BindPFlag("opts-prefix", rootCmd.Flags().Lookup("opts-prefix"))
		viper.BindPFlag("opts-timeout", rootCmd.Flags().Lookup("opts-timeout"))
		viper.BindPFlag("opts-guilds", rootCmd.Flags().Lookup("opts-guilds"))
		viper.BindPFlag("opts-sqlite", rootCmd.Flags().Lookup("opts-sqlite"))
	})
}

func main() {
	// go func() {
	// 	tiker := time.NewTicker(time.Second)
	// 	for {
	// 		select {
	// 		case <-tiker.C:
	// 			fmt.Println("Gorountine Count is: ", runtime.NumGoroutine())
	// 		}
	// 	}
	// }()
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		db, err := gorm.Open(sqlite.Open(viper.GetString("opts-sqlite")), &gorm.Config{})
		if err != nil {
			panic("failed to connect database.")
		}
		opts := libamabot.AmabotOptions{
			MessageCmdPrefix: viper.GetString("opts-prefix"),
			TimeoutDuration:  viper.GetDuration("opts-timeout"),
			AppCmdGuildIds:   viper.GetStringSlice("opts-guilds"),
			EnabledHandlers:  libamabot.GetAllHandlersList(),
			Db:               db,
		}
		token := viper.GetString("token")
		amabot, e := libamabot.New(token, opts)
		cobra.CheckErr(e)
		cobra.CheckErr(amabot.Run())
		defer amabot.Close()
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)
		log.Println("Press Ctrl+C to exit")
		<-stop
	}
	rootCmd.Execute()
}
