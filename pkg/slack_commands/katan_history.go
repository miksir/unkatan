package slack_commands

import (
	"fmt"
	"github.com/miksir/unkatan/pkg/helpers"
	"github.com/miksir/unkatan/pkg/katan"
	"github.com/slack-go/slack"
)

func (comm SlackCommands) katanHistoryCommand() (*slack.OutgoingMessage, error) {
	var deployText string = "Последние 20 изменений статуса деплоя\n"

	history := katan.DeployHistory()
	for _, item := range history {
		deployText += fmt.Sprintf(
			"%s *%s* (%s)",
			item.GetTime().Format("2006-01-02 15:04:05"),
			helpers.DeployStatusRussianName(item.IsDeployOn()),
			item.GetUser().SlackName(),
		)
		if item.GetReason() != "" {
			deployText += fmt.Sprintf(", причина: %s", item.GetReason())
		}
		deployText += "\n"
	}

	return &slack.OutgoingMessage{Text: deployText}, nil
}
