// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"

	"time"

	randomdata "github.com/Pallinder/go-randomdata"
	nsq "github.com/nsqio/go-nsq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// produceCmd represents the produce command
var produceCmd = &cobra.Command{
	Use:   "produce",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("produce called", viper.GetString("NSQ"), viper.GetString("TOPIC"))

		// produce(viper.GetString("NSQ"), viper.GetString("TOPIC"))

		produce(viper.GetString("TOPIC"), viper.GetString("NSQ"), viper.GetString("NSQ2"))
	},
}

func init() {
	RootCmd.AddCommand(produceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// produceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// produceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func produce(topic string, nsqs ...string) {

	for i, nsqH := range nsqs {
		go func(index int, nsqHost string) {
			config := nsq.NewConfig()
			w, _ := nsq.NewProducer(nsqHost, config)

			for {
				err := w.Publish(topic, []byte(randomdata.Adjective()+randomdata.FirstName(randomdata.Male)))
				if err != nil {
					log.Panic("Could not connect")
				}
			}
		}(i, nsqH)
	}
	for {
		time.Sleep(time.Minute)
	}
}
