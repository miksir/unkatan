package slack_commands

import (
	"fmt"
	"github.com/miksir/unkatan/pkg/deploy"
	"github.com/slack-go/slack"
)

func (comm SlackCommands) katanHelpCommand(s slack.SlashCommand) (*slack.OutgoingMessage, error) {
	str := ""

	if s.Text == "help deploy" {
		flagSetOff := katanDeployOffFlagSet(&deploy.DeployOffCommand{})
		flagSetOn := katanDeployOnFlagSet(&deploy.DeployOnCommand{})

		str += fmt.Sprintf("*%s deploy off _options_ [причина]* - закрыть деплой, _options_:\n", s.Command)
		str += fmt.Sprintf("%s\n", flagSetOff.FlagUsages())
		str += fmt.Sprintf("*%s deploy on _options_ [причина]* - открыть деплой, _options_:\n", s.Command)
		str += fmt.Sprintf("%s\n", flagSetOn.FlagUsages())
		str += "*(!)* Если деплой открыт в нерабочее время, он будет закрыт автоматически через час\n"
		str += "Переключения из *закрыт* в *открыт только для пользователя* не анонсируются публично (так как для всех деплой считается закрытым), "
		str += "а из *открыт для всех* в *открыт только для пользователя* анонсируются как деплой закрыт.\n"
	} else {
		str += fmt.Sprintf("*%s deploy {on|off} _options_ [причина]* - переключение состояния деплоя\n", s.Command)
		str += fmt.Sprintf("*%s help deploy* - более подробно о переключении состояния деплоя\n", s.Command)
		str += fmt.Sprintf("*%s status* - статус деплоя\n", s.Command)
		str += fmt.Sprintf("*%s history* - история переключений деплоя\n", s.Command)
		str += fmt.Sprintf("*%s grant {add|del} @user* - права переключения деплоя\n", s.Command)
		str += fmt.Sprintf("*%s grantlist* - список прав переключения деплоя\n", s.Command)
	}

	return &slack.OutgoingMessage{Text: str}, nil
}
