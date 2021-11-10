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
	fmt.Println(config)
	return config, nil
}
func (s slackConfig) SlackNotify(message string) {
	api := slack.New(s.SlackToken, slack.OptionDebug(true))
	msgText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Hi <@%s>, following query failed:%s", s.UserID, message), false, false)
	msgSection := slack.NewSectionBlock(msgText, nil, nil)
	msgBlock := slack.MsgOptionBlocks(
		msgSection,
	)
	_, _, _, err := api.SendMessage(s.ChannelID, msgBlock)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

}

func (s slackConfig) Notify(c chan string) {
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
			s.SlackNotify(msg)
		default:
			fmt.Printf("\r%s Please Wait. No new message received on the channel....", waitChars[rand.Intn(4)])
			time.Sleep(time.Millisecond * 500)
		}
	}

}
