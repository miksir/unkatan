package slack_commands

import (
	"github.com/slack-go/slack"
)

func (comm SlackCommands) katanPingCommand() (*slack.OutgoingMessage, error) {

	return &slack.OutgoingMessage{Text: "pong"}, nil
}
