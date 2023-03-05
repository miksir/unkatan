package slack_commands

import (
	"context"
	"fmt"
	"github.com/miksir/unkatan/pkg/deploy"
	"github.com/miksir/unkatan/pkg/katan"
	"github.com/slack-go/slack"
	"github.com/spf13/pflag"
	"strings"
)

func (comm SlackCommands) katanDeployCommand(ctx context.Context, s slack.SlashCommand) (*slack.OutgoingMessage, error) {
	_, err := processKatanDeployCommand(s)
	if err != nil {
		return nil, fmt.Errorf("неверный формат команды (%s), *%s help deploy* для инструкции", err, s.Command)
	}
	return comm.katanStatusCommand(ctx, s)
}

func processKatanDeployCommand(s slack.SlashCommand) (deploy.DeployCommandI, error) {
	var fields = strings.Fields(s.Text)

	if len(fields) < 2 {
		return nil, fmt.Errorf("не указано действие")
	}

	action := fields[1]
	var err error
	var usedCommand deploy.DeployCommandI
	var flags *pflag.FlagSet

	if action == "off" {
		command := deploy.DeployOffCommand{}
		flags = katanDeployOffFlagSet(&command)
		usedCommand = &command
	} else if action == "on" {
		command := deploy.DeployOnCommand{}
		flags = katanDeployOnFlagSet(&command)
		usedCommand = &command
	} else {
		return nil, fmt.Errorf("указано неизвестное действие")
	}

	err = katanDeployParseFlagSet(usedCommand, flags, fields[2:], s)
	if err != nil {
		return nil, fmt.Errorf("неверные параметры: %s", err)
	}

	return usedCommand, nil
}

func katanDeployParseFlagSet(command deploy.DeployCommandI, flags *pflag.FlagSet, tokens []string, s slack.SlashCommand) error {
	command.SetUser(deploy.ActionUser{SlackId: s.UserID, Name: s.UserName})
	err := flags.Parse(tokens)
	if err != nil {
		return err
	}
	if flags.NArg() > 1 {
		command.SetReason(strings.Join(flags.Args(), " "))
	}
	switch command.(type) {
	case *deploy.DeployOnCommand:
		names := command.(*deploy.DeployOnCommand).OnlyFor
		newNames := make([]string, 0, len(names))
		for _, name := range names {
			newName := strings.TrimLeft(name, "@")
			newNames = append(newNames, newName)
		}
		command.(*deploy.DeployOnCommand).OnlyFor = newNames
	}
	katan.DeployStatusCommand(command)
	return nil
}

func katanDeployOffFlagSet(options *deploy.DeployOffCommand) *pflag.FlagSet {
	flags := pflag.NewFlagSet("deploy", pflag.ContinueOnError)
	flags.BoolVarP(&options.Permanent, "permanent", "p", false,
		"отключить работу расписания (если не указывать, то деплой разблокируется по расписанию)")
	return flags
}

func katanDeployOnFlagSet(options *deploy.DeployOnCommand) *pflag.FlagSet {
	flags := pflag.NewFlagSet("deploy", pflag.ContinueOnError)
	flags.StringSliceVarP(&options.OnlyFor, "only-for", "u", []string{},
		"открыть только для указанных людей, можно через запятую без пробела или использовать опцию несколько раз; указывать логины или слак тег; новая команда заменяет предыдущий список людей; при запросе статуса деплой отображается заблокированым")
	return flags
}
