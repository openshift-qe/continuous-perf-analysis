package Notify

import (
	"fmt"

	"github.com/slack-go/slack"
)

type SlackConfig struct {
	UserID     string `json:"userID"`
	ChannelID  string `json:"channelID"`
	SlackToken string `json:"slackToken"`
}

// TODO Add function to Read Slack Config

func (s SlackConfig) SlackNotify(message string) {
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
