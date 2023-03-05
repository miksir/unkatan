package slack_commands

import (
	"context"
	"fmt"
	"github.com/miksir/unkatan/pkg/deploy"
	"github.com/miksir/unkatan/pkg/helpers"
	"github.com/miksir/unkatan/pkg/katan"
	"github.com/slack-go/slack"
)

func (comm SlackCommands) katanStatusCommand(_ context.Context, s slack.SlashCommand) (*slack.OutgoingMessage, error) {
	var deployText string
	status := katan.DeployStatus()
	deployOn := status.IsDeployOn()
	deployText = fmt.Sprintf(
		"*%s* (%s)",
		helpers.DeployStatusRussianName(deployOn),
		status.GetUser().SlackName(),
	)
	if status.GetReason() != "" {
		deployText += fmt.Sprintf(", причина: %s", status.GetReason())
	}
	switch status.(type) {
	case *deploy.DeployOffCommand:
		if status.(*deploy.DeployOffCommand).Permanent {
			deployText += ", расписание выключено"
		}
	}

	deployOnMe := status.IsDeployOn(s.UserName)
	if deployOnMe != deployOn {
		deployText += fmt.Sprintf("\nДля вас деплой *%s*", helpers.DeployStatusRussianName(deployOnMe))
	}

	switch status.(type) {
	case *deploy.DeployOnCommand:
		if len(status.(*deploy.DeployOnCommand).OnlyFor) != 0 && comm.checkPermittedUsers(s) {
			deployText += "\nДеплой открыт для"
			for _, openFor := range status.(*deploy.DeployOnCommand).OnlyFor {
				deployText += fmt.Sprintf(" <@%s>", openFor)
			}
		}
	}

	return &slack.OutgoingMessage{Text: fmt.Sprintf("Статус деплоя: %s", deployText)}, nil
}
