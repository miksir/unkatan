package announce

import (
	"github.com/miksir/unkatan/pkg/deploy"
	"time"
)

type Announce interface {
	DeployAnnounce(status bool, reason string, who deploy.ActionUser, nextDate *time.Time)
}
