package main

import (
	"bufio"
	"fmt"
	"github.com/miksir/unkatan/pkg/announce"
	"github.com/miksir/unkatan/pkg/katan"
	c "github.com/miksir/unkatan/pkg/lconfig"
	zlog "github.com/miksir/unkatan/pkg/log"
	"github.com/miksir/unkatan/pkg/schedule"
	"github.com/miksir/unkatan/pkg/slack_commands"
	"github.com/miksir/unkatan/pkg/storage"
	"github.com/slack-go/slack"
	"log"
	"os"
	"strings"
)

var version string

/*
*

	Enter slack commands via commandline for manual testing purposes
*/
func main() {
	reader := bufio.NewReader(os.Stdin)

	//zap, _ := zap2.NewProduction()
	//logger := zlog.NewWithLogger(zap)
	logger := zlog.NewNoopLogger()

	logger.Info(nil, fmt.Sprintf("VERSION %s", version))

	cfg, err := c.Init(logger)
	if err != nil {
		log.Fatal(err)
	}

	schd := schedule.NewSchedule(cfg.Sub("schedule"), logger)

	slackConfig := cfg.Sub("slack")

	slackAPI := slack.New("")

	slackCommands := slack_commands.NewSlackCommands(slackConfig, logger, slackAPI)

	katan.DeployStatusInit(
		storage.NewNullRegistry(logger),
		announce.NewAnnounceConsole(),
		*schd,
		logger,
	)

	fmt.Println("Type slack slash command or @user for switch slack user or . for exit")
	user := "test"

	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		text = strings.Replace(text, "\n", "", -1)

		if strings.Compare(".", text) == 0 {
			fmt.Println("Bye")
			os.Exit(0)
		}

		if strings.HasPrefix(text, "@") {
			user = text[1:]
			fmt.Println("Switched to:", user)
			continue
		}

		if strings.HasPrefix(text, "/") {
			strs := strings.SplitN(text, " ", 2)
			line := ""
			if len(strs) >= 2 {
				line = strs[1]
			}
			slackCmd := slack.SlashCommand{
				UserID:   user,
				UserName: user,
				Command:  strs[0],
				Text:     line,
			}
			resp := slackCommands.RunCommand(nil, slackCmd)
			fmt.Println(resp.Text)
			continue
		}

		fmt.Println("Type slack slash command or @user for switch slack user or . for exit")
	}
}
