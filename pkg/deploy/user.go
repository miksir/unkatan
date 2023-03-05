package deploy

import "fmt"

type ActionUser struct {
	SlackId     string
	Name        string
	Email       string
	IsScheduler bool
}

func (u ActionUser) PlainName() string {
	if u.Name == "" {
		return "unknown"
	}
	return u.Name
}

func (u ActionUser) SlackName() string {
	if u.SlackId == "" {
		return fmt.Sprintf("*%s*", u.PlainName())
	}
	return fmt.Sprintf("<@%s|%s>", u.SlackId, u.Name)
}
