package ulogger

import (
	zlog "github.com/miksir/unkatan/pkg/log"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type Slack struct {
	log      zlog.Logger
	slackAPI *slack.Client
	channels []string
}

func NewSlackLogger(channels []string, slackAPI *slack.Client, log zlog.Logger) *Slack {
	return &Slack{log: log, slackAPI: slackAPI, channels: channels}
}

func (s Slack) Message(msg string) {
	var err error
	if len(s.channels) == 0 {
		return
	}

	for _, channel := range s.channels {
		_, _, err = s.slackAPI.PostMessage(channel, slack.MsgOptionText(msg, false))
		if err != nil {
			s.log.Error(
				nil,
				"[LOG to SLACK] failed to write to channel",
				zap.Error(err),
				zap.String("channel", channel),
			)
		}
	}
}
