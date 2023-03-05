package main

import (
	"context"
	"fmt"
	"github.com/miksir/unkatan/pkg/announce"
	"github.com/miksir/unkatan/pkg/api"
	"github.com/miksir/unkatan/pkg/genie"
	html2 "github.com/miksir/unkatan/pkg/html"
	"github.com/miksir/unkatan/pkg/katan"
	c "github.com/miksir/unkatan/pkg/lconfig"
	zlog "github.com/miksir/unkatan/pkg/log"
	"github.com/miksir/unkatan/pkg/request"
	"github.com/miksir/unkatan/pkg/schedule"
	"github.com/miksir/unkatan/pkg/slack_commands"
	"github.com/miksir/unkatan/pkg/storage"
	"github.com/miksir/unkatan/pkg/ulogger"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
)
import "github.com/robfig/cron"

var version string

func main() {
	logger, err := zlog.NewLogger(zap.NewProductionConfig())
	if err != nil {
		log.Fatal(err)
	}

	logger.Info(nil, fmt.Sprintf("VERSION %s", version))

	cfg, err := c.Init(logger)
	if err != nil {
		log.Fatal(err)
	}

	schd := schedule.NewSchedule(cfg.Sub("schedule"), logger)

	slackConfig := cfg.Sub("slack")
	if slackConfig.GetString("bot_token") == "" {
		log.Fatal("slack config: bot_token required")
	}

	slackAPI := slack.New(slackConfig.GetString("bot_token"))

	slackLogger := ulogger.NewSlackLogger(
		cfg.GetStringSlice("katan.slack_warn_channels"),
		slackAPI,
		logger,
	)
	logger = logger.WithOptions(zap.Hooks(func(entry zapcore.Entry) error {
		if entry.Level == zapcore.WarnLevel {
			go slackLogger.Message(entry.Message)
		}
		return nil
	}))

	slackCommands := slack_commands.NewSlackCommands(slackConfig, logger, slackAPI)
	apiRequests := api.NewAPI(cfg.Sub("api"), logger, slackAPI)

	katan.DeployStatusInit(
		storage.NewRedisRegistry(cfg.Sub("redis"), logger),
		announce.NewAnnounceSlack(cfg.GetStringSlice("katan.slack_announce_channels"), slackAPI, logger),
		*schd,
		logger,
	)

	genie.InitCollection(cfg.Sub("genie"), logger)
	go genie.UpdateCollection()

	cronIns := cron.New()
	_ = cronIns.AddFunc("0 * * * *", func() {
		katan.CheckSchedule()
	})
	_ = cronIns.AddFunc("10 1 * * *", func() {
		genie.UpdateCollection()
	})
	cronIns.Start()

	html := html2.NewHtml()

	mux := http.NewServeMux()
	mux.HandleFunc("/", html.HandleMain)
	mux.HandleFunc("/slack/", request.ApiResponseHandler(logger, slackCommands.Commands))
	mux.HandleFunc("/api/", request.ApiResponseHandler(logger, apiRequests.APIHandler))
	mux.HandleFunc("/ready/", request.ApiResponseHandler(logger, apiRequests.APIHandler))
	mux.HandleFunc("/live/", request.ApiResponseHandler(logger, apiRequests.APIHandler))

	listen := cfg.GetString("http.listen")
	logger.Info(context.Background(), "Starting server", zap.String("port", listen))

	err = http.ListenAndServe(listen, request.AccessLogHandler(logger, mux))
	log.Fatal(err)
}
