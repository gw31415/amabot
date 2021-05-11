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
package cmd

import (
	"log"
	"os"
	"os/signal"

	"github.com/gw31415/amabot/libamabot"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Amabot",
	Long:  `Run Amabot and start listening Discord Users.`,
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := recover(); err != nil {
				log.Fatalln(err)
			}
		}()
		amabot, err := libamabot.New()
		if err != nil {
			panic("Failed to instantiate Discord client")
		}

		err = amabot.Run()
		if err != nil {
			panic(err)
		}
		defer func() {
			log.Println("Closing Amabot....")
			amabot.Close()
			log.Println("done.")
		}()

		log.Println("Amabot is now running. Press CTRL-C to exit.")
		stop := make(chan os.Signal)
		signal.Notify(stop, os.Interrupt)
		<-stop
		log.Println("Keyboard Interrupt")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	runCmd.Flags().StringP("token", "t", "", "Discord Bot token")
	viper.BindPFlag("token", runCmd.Flags().Lookup("token"))
	viper.BindEnv("TOKEN")
	runCmd.Flags().StringP("game-status", "g", "Amabot", "The name of game at status of Discord")
	viper.BindPFlag("game-status", runCmd.Flags().Lookup("game-status"))
	viper.BindEnv("GAME_STATUS")
}
