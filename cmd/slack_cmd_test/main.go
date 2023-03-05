package main

import (
	"fmt"
	"github.com/miksir/unkatan/pkg/announce"
	"github.com/miksir/unkatan/pkg/api"
	"github.com/miksir/unkatan/pkg/helpers"
	"github.com/miksir/unkatan/pkg/katan"
	c "github.com/miksir/unkatan/pkg/lconfig"
	zlog "github.com/miksir/unkatan/pkg/log"
	"github.com/miksir/unkatan/pkg/schedule"
	"github.com/miksir/unkatan/pkg/slack_commands"
	"github.com/miksir/unkatan/pkg/storage"
	"github.com/slack-go/slack"
	zap2 "go.uber.org/zap"
	"log"
	"net/url"
	"time"
)

var version string

/*
*

	Functional test
*/
func main() {
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

	if !slackConfig.GetBool("testrun") {
		log.Fatalf("use unkatan.yml from test main.go dir")
	}

	slackAPI := slack.New("")

	slackCommands := slack_commands.NewSlackCommands(slackConfig, logger, slackAPI)

	katan.DeployStatusInit(
		storage.NewNullRegistry(logger),
		announce.NewAnnounceConsole(),
		*schd,
		logger,
	)

	apiRequests := api.NewAPI(cfg.Sub("api"), logger, slackAPI)

	httpCheckDeployStatus := func(expected bool, changedBy string, requestBy string) func() error {
		return func() error {
			rurl := "/api/deploy"
			form := make(url.Values)
			if requestBy != "" {
				rurl = rurl + "?user=" + requestBy
				form.Set("user", requestBy)
			}
			w, r := httpctx("GET", rurl, form)
			resp, err := apiRequests.APIHandler(r.Context(), w, r)
			if err != nil {
				return err
			}
			switch resp.(type) {
			case api.DeployResponse:
				if resp.(api.DeployResponse).Deploy != expected {
					return fmt.Errorf("deploy %t but shoud be %t: %+v", resp.(api.DeployResponse).Deploy, expected, resp)
				}
				if changedBy != "" && resp.(api.DeployResponse).User != changedBy {
					return fmt.Errorf("user %s but shoud be %s: %+v", resp.(api.DeployResponse).User, changedBy, resp)
				}
			default:
				return fmt.Errorf("unknown http response %+v", resp)
			}
			return nil
		}
	}

	slackCommand := func(user, text string, isError bool) func() error {
		return func() error {
			slackCmd := slack.SlashCommand{
				UserID:   user,
				UserName: user,
				Command:  "/katan",
				Text:     text,
			}
			message, err := slackCommands.KatanCommand(nil, slackCmd)
			logger.Info(nil, "slack output", zap2.Any("message", message), zap2.Error(err))
			if isError {
				if err != nil {
					return nil
				}
				return fmt.Errorf("no error received")
			} else {
				return err
			}
		}
	}

	checkSchedule := func(date time.Time, status bool, user string) func() error {
		return func() error {
			helpers.SetTime(date)
			katan.CheckSchedule()
			return httpCheckDeployStatus(status, user, "")()
		}
	}

	t("initial deploy status", httpCheckDeployStatus(false, "", ""))
	t("open deploy with unauthorized user failed", slackCommand("test", "deploy on", true))
	t("check deploy is still closed", httpCheckDeployStatus(false, "", ""))
	t("open deploy with slack command with authorized user", slackCommand("test2", "deploy on", false))
	t("check deploy is open", httpCheckDeployStatus(true, "test2", ""))
	t("check deploy is open for a.tarasov", httpCheckDeployStatus(true, "test2", "a.tarasov"))
	t("close deploy", slackCommand("test2", "deploy off", false))
	t("custom open for user i.petrov", slackCommand("test2", "deploy on -u @i.petrov", false))
	t("check deploy is closed", httpCheckDeployStatus(false, "test2", ""))
	t("check deploy is closed for a.tarasov", httpCheckDeployStatus(false, "test2", "a.tarasov"))
	t("check deploy is open for i.petorv", httpCheckDeployStatus(true, "test2", "i.petrov"))
	t("custom open for user s.ivanov", slackCommand("test2", "deploy on -u @s.ivanov", false))
	t("check deploy is closed", httpCheckDeployStatus(false, "test2", ""))
	t("check deploy is closed for i.petrov", httpCheckDeployStatus(false, "test2", "i.petrov"))
	t("check deploy is open for s.ivanov", httpCheckDeployStatus(true, "test2", "s.ivanov"))

	t("close deploy", slackCommand("test2", "deploy off", false))
	t("schedule open work day", checkSchedule(tDate(3, 22, 14, 50), true, "schedule"))
	t("schedule close for night", checkSchedule(tDate(3, 22, 22, 0), false, "schedule"))
	t("schedule exclude from special", checkSchedule(tDate(5, 01, 13, 0), true, "schedule"))
	t("schedule special day", checkSchedule(tDate(5, 9, 13, 0), false, "schedule"))

	t("schedule open work day", checkSchedule(tDate(3, 22, 14, 50), true, "schedule"))
	t("close deploy permanent", slackCommand("test2", "deploy off --permanent", false))
	t("check deploy is closed", httpCheckDeployStatus(false, "test2", ""))
	t("check deploy still closed at work day", checkSchedule(tDate(3, 22, 14, 50), false, "test2"))
	t("check deploy still closed at night", checkSchedule(tDate(3, 22, 22, 0), false, "test2"))
	t("check deploy still closed at work day", checkSchedule(tDate(3, 23, 14, 50), false, "test2"))
	t("close deploy but without permanent", slackCommand("test2", "deploy off", false))
	t("check deploy still closed at work day", checkSchedule(tDate(3, 23, 15, 30), false, "test2"))
	t("scheduler works at night", checkSchedule(tDate(3, 23, 22, 30), false, "schedule"))
	t("scheduler open work day", checkSchedule(tDate(3, 24, 12, 30), true, "schedule"))

	t("scheduler close next day", checkSchedule(tDate(3, 24, 21, 30), false, "schedule"))
	t("someone opened deploy at night", slackCommand("test2", "deploy on", false))
	t("deploy open after half hour", checkSchedule(tDate(3, 24, 22, 00), true, "test2"))
	t("deploy closed after hour", checkSchedule(tDate(3, 24, 22, 31), false, "schedule"))
}
