package announce

import (
	"github.com/miksir/unkatan/pkg/deploy"
	"time"
)

type announceNull struct {
}

func NewAnnounceNull() *announceNull {
	return &announceNull{}
}

func (n *announceNull) DeployAnnounce(_ bool, _ string, _ deploy.ActionUser, _ *time.Time) {
	return
}
