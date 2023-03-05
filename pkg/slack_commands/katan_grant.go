package slack_commands

import (
	"context"
	"fmt"
	"github.com/slack-go/slack"
	"strings"
)

func (comm SlackCommands) katanGrantCommand(ctx context.Context, s slack.SlashCommand) (*slack.OutgoingMessage, error) {
	var errMsg = &slack.OutgoingMessage{Text: fmt.Sprintf("%s grant {add|del} @user", s.Command)}

	tokens := strings.Fields(s.Text)

	if len(tokens) != 3 {
		return errMsg, nil
	}

	if string(tokens[2][0]) != "@" {
		return errMsg, nil
	}

	user := tokens[2][1:]

	found := -1
	changed := false
	for i, permittedUser := range comm.cfg.GetStringSlice("permitted_users") {
		if user == permittedUser {
			found = i
			break
		}
	}

	if tokens[1] == "del" {
		if found >= 0 {
			newlist := make([]string, 0, len(comm.cfg.GetStringSlice("permitted_users"))-1)
			newlist = append(newlist, comm.cfg.GetStringSlice("permitted_users")[0:found]...)
			newlist = append(newlist, comm.cfg.GetStringSlice("permitted_users")[found+1:]...)
			comm.cfg.Set("permitted_users", newlist)
			comm.cfg.Save()
			changed = true
		}
	}
	if tokens[1] == "add" {
		if found == -1 {
			newlist := comm.cfg.GetStringSlice("permitted_users")
			newlist = append(newlist, user)
			comm.cfg.Set("permitted_users", newlist)
			comm.cfg.Save()
			changed = true
		}
	}

	if changed {
		comm.log.Warn(ctx, fmt.Sprintf("<@%s> changed grant list: %s <@%s>", s.UserName, tokens[1], user))
	}

	return comm.katanGrantListCommand(ctx, s)
}
