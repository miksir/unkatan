package slack_commands

import (
	"errors"
	"fmt"
	"github.com/miksir/unkatan/pkg/lconfig"
	"github.com/miksir/unkatan/pkg/request"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"context"
	zlog "github.com/miksir/unkatan/pkg/log"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type SlackCommands struct {
	cfg      lconfig.Reader
	log      zlog.Logger
	slackAPI *slack.Client
}

func NewSlackCommands(cfg lconfig.Reader, logger zlog.Logger, slackAPI *slack.Client) SlackCommands {
	if cfg.GetString("sign_secret") == "" {
		log.Fatal("slack config: sign_secret is required")
	}
	comm := SlackCommands{
		cfg:      cfg,
		log:      logger,
		slackAPI: slackAPI,
	}
	return comm
}

func (comm *SlackCommands) Commands(ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var err error
	var resp *slack.OutgoingMessage = nil
	var s *slack.SlashCommand

	w.Header().Set("Content-Type", "application/json")

	defer comm.log.Debug(ctx,
		"[SLACK] command request",
		zap.Any("headers", r.Header),
		zap.Any("command", s),
		zap.Any("response", resp),
	)

	verifier, err := slack.NewSecretsVerifier(r.Header, comm.cfg.GetString("sign_secret"))
	if err != nil {
		comm.log.Error(ctx, "[SLACK] Secret init error", zap.Error(err))
		return nil, &request.RequestError{
			Status: http.StatusBadRequest,
		}
	}

	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
	sl, err := slack.SlashCommandParse(r)
	s = &sl
	if err != nil {
		comm.log.Error(ctx, "[SLACK] failed to parse slash command", zap.Error(err))
		return nil, &request.RequestError{
			Status: http.StatusBadRequest,
		}
	}

	if err = verifier.Ensure(); err != nil {
		comm.log.Error(ctx, "[SLACK] failed to validate secret", zap.Error(err))
		return nil, &request.RequestError{
			Status: http.StatusUnauthorized,
		}
	}

	resp = comm.RunCommand(ctx, sl)

	return resp, nil
}

func (comm *SlackCommands) RunCommand(ctx context.Context, sl slack.SlashCommand) *slack.OutgoingMessage {
	var err error
	var resp *slack.OutgoingMessage = nil

	switch sl.Command {
	case "/katan", "/unkatan_test", "/unkatan_dev":
		resp, err = comm.KatanCommand(ctx, sl)
	default:
		err = errors.New(fmt.Sprintf("%s not supported by me", sl.Command))
	}

	if err != nil {
		comm.log.Error(
			ctx,
			"[SLACK] failed to execute command",
			zap.String("command", sl.Command),
			zap.Error(err),
		)
		resp = &slack.OutgoingMessage{Text: fmt.Sprintf("[ERROR] %s", err)}
	}

	return resp
}
