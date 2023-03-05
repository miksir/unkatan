package slack_commands

import (
	"context"
	"fmt"
	"github.com/slack-go/slack"
)

func (comm SlackCommands) katanGrantListCommand(_ context.Context, _ slack.SlashCommand) (*slack.OutgoingMessage, error) {

	str := "Могут переключать деплой:\n"
	for _, permittedUser := range comm.cfg.GetStringSlice("permitted_users") {
		str += fmt.Sprintf("<@%s>\n", permittedUser)
	}
	return &slack.OutgoingMessage{Text: str}, nil
}
