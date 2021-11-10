package Notify

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/slack-go/slack"
	"gopkg.in/yaml.v2"
)

const configPath = "./config/"

type slackConfig struct {
	UserID     string `json:"userid"`
	ChannelID  string `json:"channelid"`
	SlackToken string `json:"slacktoken"`
}

func (c *slackConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

func ReadslackConfig() (config slackConfig, err error) {
	data, err := ioutil.ReadFile(configPath + "slack.yaml")
	msg := fmt.Sprintf("Cound't read %sslack.yaml", configPath)
	if err != nil {
		return config, fmt.Errorf(msg)
	}
	if err := config.Parse(data); err != nil {
		log.Fatal(err)
		return config, err
	}
	return config, nil
}

func (s slackConfig) SlackNotify(message, thread_ts string) string {
	// api := slack.New(s.SlackToken, slack.OptionDebug(true)) // To debug api requests, you uncomment this line and comment the one below
	api := slack.New(s.SlackToken)
	msgText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Hi <@%s>, %s", s.UserID, message), false, false)
	msgSection := slack.NewSectionBlock(msgText, nil, nil)
	msgBlock := slack.MsgOptionBlocks(
		msgSection,
	)
	var err error
	if thread_ts != "" {
		// if we have thread_ts - use it to send new messages on the thread
		_, _, _, err = api.SendMessage(s.ChannelID, msgBlock, slack.MsgOptionTS(thread_ts))
	} else {
		// if thread_ts was empty, assume that this is first message we are sending and retrieve thread_ts and return it for subsequent requests
		_, thread_ts, _, err = api.SendMessage(s.ChannelID, msgBlock)
	}
	if err != nil {
		log.Fatal(err)
	}
	return thread_ts
}

func (s slackConfig) Notify(c chan string, thread_ts string) {
	waitChars := []string{"/", "-", "\\", "|"}
	for {
		select {
		case msg := <-c:
			msgFmt := fmt.Sprintf(`
%s
Received following on the channel: %s
%[1]s
			`, strings.Repeat("~", 80), msg)
			fmt.Println(msgFmt)
			s.SlackNotify("Following query failed:"+msg, thread_ts)
		default:
			fmt.Printf("\r%s Please Wait. No new message received on the channel....", waitChars[rand.Intn(4)])
			time.Sleep(time.Millisecond * 500)
		}
	}

}
