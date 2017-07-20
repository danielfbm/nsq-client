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
	"sync"
	"time"

	nsq "github.com/nsqio/go-nsq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// consumeCmd represents the consume command
var consumeCmd = &cobra.Command{
	Use:   "consume",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("consume called", viper.GetString("TOPIC"), viper.GetString("NSQ"), viper.GetString("NSQ2"))

		consume4Ever(viper.GetString("LOOKUP"), viper.GetString("TOPIC"), viper.GetString("CHANNEL"))
		// consumeDirect(viper.GetString("TOPIC"), viper.GetString("NSQ"), viper.GetString("NSQ2"))
	},
}

func init() {
	RootCmd.AddCommand(consumeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// consumeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// consumeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func consume(lookupHost, topic, channel string) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	config := nsq.NewConfig()
	q, _ := nsq.NewConsumer(topic, channel, config)
	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Printf("Got a message: %v", message)
		wg.Done()
		return nil
	}))
	// err := q.ConnectToNSQD("127.0.0.1:4150")
	err := q.ConnectToNSQLookupd(lookupHost)
	if err != nil {
		log.Panic("Could not connect")
		wg.Done()
	}
	wg.Wait()
}

func consume4Ever(lookupHost, topic, channel string) {

	config := nsq.NewConfig()
	q, _ := nsq.NewConsumer(topic, channel, config)
	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Printf("Got a message: %s", message.Body)
		message.Finish()
		return nil
	}))

	err := q.ConnectToNSQLookupd(lookupHost)
	if err != nil {
		log.Panic("Could not connect")
	}
	for {
		time.Sleep(time.Minute)
	}

}

func consumeDirect(topic, channel string, nsqs ...string) {
	for i, nsqH := range nsqs {
		go func(index int, nsqHost string) {
			config := nsq.NewConfig()
			q, _ := nsq.NewConsumer(topic, channel, config)
			q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
				log.Printf("Got a message: %s", message.Body)
				return nil
			}))

			err := q.ConnectToNSQD(nsqHost)
			if err != nil {
				log.Panic("Could not connect")
			}
			for {
				time.Sleep(time.Minute)
			}
		}(i, nsqH)
	}
	for {
		time.Sleep(time.Minute)
	}

}
