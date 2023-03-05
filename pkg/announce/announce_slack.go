package announce

import (
	"fmt"
	"github.com/miksir/unkatan/pkg/deploy"
	"github.com/miksir/unkatan/pkg/helpers"
	zlog "github.com/miksir/unkatan/pkg/log"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
	"time"
)

type AnnounceSlack struct {
	log              zlog.Logger
	slackAPI         *slack.Client
	announceChannels []string
}

func NewAnnounceSlack(announceChannels []string, slackAPI *slack.Client, log zlog.Logger) *AnnounceSlack {
	return &AnnounceSlack{log: log, slackAPI: slackAPI, announceChannels: announceChannels}
}

func (comm AnnounceSlack) DeployAnnounce(status bool, reason string, who deploy.ActionUser, nextDate *time.Time) {
	var err error

	if len(comm.announceChannels) == 0 {
		return
	}

	var fullReason, fullNext string

	if reason != "" {
		fullReason = ", причина: " + reason
	}
	if nextDate != nil {
		var relDate = helpers.RelativeDateText(*nextDate)
		fullNext = fmt.Sprintf(" до %s", relDate)
	}
	whoName := who.SlackName()
	if who.IsScheduler {
		whoName = "автоматически"
	}

	text := fmt.Sprintf(
		"Деплой *%s*%s! (%s%s)",
		helpers.DeployStatusRussianName(status),
		fullNext,
		whoName,
		fullReason,
	)

	go (func() {
		for _, channel := range comm.announceChannels {
			_, _, err = comm.slackAPI.PostMessage(channel, slack.MsgOptionText(text, false))
			if err != nil {
				comm.log.Error(
					nil,
					"[SLACK] failed to announce to channel",
					zap.Error(err),
					zap.String("channel", channel),
				)
			}
		}
	})()
}
