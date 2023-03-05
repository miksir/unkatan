package slack_commands

import (
	"errors"
	"fmt"
	"strings"

	"context"
	"github.com/slack-go/slack"
)

func (comm SlackCommands) KatanCommand(ctx context.Context, s slack.SlashCommand) (*slack.OutgoingMessage, error) {
	var fields = strings.Fields(s.Text)
	if len(fields) == 0 {
		return &slack.OutgoingMessage{Text: fmt.Sprintf("type `%s help` for help", s.Command)}, nil
	}

	var err error
	var res *slack.OutgoingMessage = nil

	switch fields[0] {
	case "help":
		res, err = comm.katanHelpCommand(s)
	case "ping":
		res, err = comm.katanPingCommand()
	case "deploy":
		if !comm.checkPermittedUsers(s) {
			return nil, errors.New("you are not allowed to use this action")
		} else {
			res, err = comm.katanDeployCommand(ctx, s)
		}
	case "grant":
		if !comm.checkPermittedUsers(s) {
			return nil, errors.New("you are not allowed to use this action")
		} else {
			res, err = comm.katanGrantCommand(ctx, s)
		}
	case "grantlist":
		res, err = comm.katanGrantListCommand(ctx, s)
	case "status":
		res, err = comm.katanStatusCommand(ctx, s)
	case "history":
		res, err = comm.katanHistoryCommand()
	default:
		err = errors.New(fmt.Sprintf("Unknown action %s, type `%s help` for help", fields[0], s.Command))
	}
	return res, err
}

func (comm SlackCommands) checkPermittedUsers(s slack.SlashCommand) bool {
	for _, permittedUser := range comm.cfg.GetStringSlice("permitted_users") {
		if s.UserID == permittedUser {
			return true
		}
		if s.UserName == permittedUser {
			return true
		}
	}

	return false
}
