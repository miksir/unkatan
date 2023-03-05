package deploy

import (
	"github.com/miksir/unkatan/pkg/helpers"
	"strings"
)

func DeployCommandFullRussianName(command DeployCommandI) string {
	answer := ""
	switch command.(type) {
	case *DeployOffCommand:
		answer = "закрыт в " + helpers.RelativeDateText(command.GetTime()) + " "
		if command.(*DeployOffCommand).Permanent {
			answer += "permanent "
		}
		if command.GetUser().IsScheduler {
			answer += "(" + "автоматически"
		} else {
			answer += "(" + command.GetUser().PlainName()
		}
		if command.GetReason() != "" {
			answer += "; " + command.GetReason()
		}
		answer += ")"
	case *DeployOnCommand:
		answer = "открыт в " + helpers.RelativeDateText(command.GetTime()) + " "
		if len(command.(*DeployOnCommand).OnlyFor) > 0 {
			answer += "только для " + strings.Join(command.(*DeployOnCommand).OnlyFor, ",") + " "
		}
		if command.GetUser().IsScheduler {
			answer += "(" + "автоматически"
		} else {
			answer += "(" + command.GetUser().PlainName()
		}
		if command.GetReason() != "" {
			answer += "; " + command.GetReason()
		}
		answer += ")"
	}
	return answer
}
