package announce

import (
	"fmt"
	"github.com/miksir/unkatan/pkg/deploy"
	"github.com/miksir/unkatan/pkg/helpers"
	"time"
)

type announceConsole struct {
}

func NewAnnounceConsole() *announceConsole {
	return &announceConsole{}
}

func (n *announceConsole) DeployAnnounce(status bool, reason string, who deploy.ActionUser, nextDate *time.Time) {
	var fullReason, fullNext string
	if reason != "" {
		fullReason = ", причина: " + reason
	}
	if nextDate != nil {
		var relDate = helpers.RelativeDateText(*nextDate)
		fullNext = fmt.Sprintf(" до %s", relDate)
	}
	fmt.Printf("Деплой %s by %s%s%s\n",
		helpers.DeployStatusRussianName(status),
		who.PlainName(),
		fullNext,
		fullReason)
}
