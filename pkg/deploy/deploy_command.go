package deploy

import "time"

type DeployOnCommand struct {
	OnlyFor           []string
	DeployBaseCommand `mapstructure:",squash"`
}

type DeployOffCommand struct {
	Permanent         bool
	DeployBaseCommand `mapstructure:",squash"`
}

type DeployBaseCommand struct {
	Reason string
	User   ActionUser
	Time   time.Time
}

type DeployCommandI interface {
	GetReason() string
	GetUser() ActionUser
	GetTime() time.Time
	SetReason(string)
	SetUser(ActionUser)
	SetTime(time.Time)
	IsDeployOn(...string) bool
}

func (c *DeployBaseCommand) GetReason() string {
	return c.Reason
}

func (c *DeployBaseCommand) GetUser() ActionUser {
	return c.User
}

func (c *DeployBaseCommand) GetTime() time.Time {
	return c.Time
}

func (c *DeployBaseCommand) SetReason(reason string) {
	c.Reason = reason
}

func (c *DeployBaseCommand) SetUser(user ActionUser) {
	c.User = user
}

func (c *DeployBaseCommand) SetTime(time time.Time) {
	c.Time = time
}

func (c *DeployBaseCommand) IsDeployOn(_ ...string) bool {
	panic("should be implemented in children")
}

func (c *DeployOffCommand) IsDeployOn(_ ...string) bool {
	return false
}

func (c *DeployOnCommand) IsDeployOn(forUsers ...string) bool {
	if len(c.OnlyFor) == 0 {
		return true
	}
	for _, user := range c.OnlyFor {
		for _, requestUser := range forUsers {
			if user == requestUser {
				return true
			}
		}
	}
	return false
}

func CommandName(command DeployCommandI) string {
	switch command.(type) {
	case *DeployOffCommand:
		return DeployActionOff
	case *DeployOnCommand:
		return DeployActionOn
	default:
		return ""
	}
}

const (
	DeployActionOn  string = "DeployOnCommand"
	DeployActionOff string = "DeployOffCommand"
)
