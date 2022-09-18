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
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/gw31415/amabot/libamabot"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "amabot",
		Short: "amabot is a Discord Bot developed by @Amadeus_vn for personal use.",
	}
	// Value of Discord API token
	token string
	// Path of the configuration file for this tool
	config = ""
)

func init() {
	rootCmd.Flags().StringVarP(&token, "token", "t", "", "Value of Discord API token")
	rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "", "Path of the configuration file for this tool")
	// Read configuration file when it exists.
	cobra.OnInitialize(func() {
		if config != "" {
			_, err := os.Stat(config)
			cobra.CheckErr(err)
			viper.SetConfigFile(config)
		} else {
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)

			viper.AddConfigPath(".")
			viper.AddConfigPath(home)
			viper.SetConfigType("yaml")
			viper.SetConfigName("amabot")
		}

		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err == nil {
			log.Println("Using config file:", viper.ConfigFileUsed())
			token = viper.GetString("token")
		}
	})
}

func main() {
	go func(){
		tiker := time.NewTicker(time.Second)
		for {
			select {
			case <- tiker.C:
				fmt.Println("Gorountine Count is: ", runtime.NumGoroutine())
			}
		}
	}()

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		amabot, e := libamabot.New(token)
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
