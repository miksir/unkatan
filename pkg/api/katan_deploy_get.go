package api

import (
	"context"
	"github.com/miksir/unkatan/pkg/deploy"
	"github.com/miksir/unkatan/pkg/katan"
)

type DeployResponse struct {
	Deploy bool   `json:"deploy"`
	Reason string `json:"reason"`
	User   string `json:"user"`
}

func (api *API) KatanDeployGetApi(ctx context.Context) (DeployResponse, error) {
	users := []string{}
	by := ctx.HTTPRequest().FormValue("user")
	if by != "" {
		users = append(users, by)
	}
	status := katan.DeployStatus()
	active := status.IsDeployOn(users...)

	res := DeployResponse{
		Deploy: active,
		Reason: status.GetReason(),
		User:   status.GetUser().PlainName(),
	}

	return res, nil
}

func (api *API) KatanDeployStatusSwitchOff(ctx context.Context) (SuccessResponse, error) {
	by := ctx.HTTPRequest().PostFormValue("by")
	if by == "" {
		by = ctx.HTTPRequest().RemoteAddr
	}

	permanent := ctx.HTTPRequest().PostFormValue("permanent")
	permanentFlag := false
	if permanent == "true" || permanent == "yes" || permanent == "on" {
		permanentFlag = true
	}

	command := deploy.DeployOffCommand{}
	command.User = deploy.ActionUser{Name: by}
	command.Reason = ctx.HTTPRequest().PostFormValue("reason")
	command.Permanent = permanentFlag

	katan.DeployStatusCommand(&command)

	return SuccessResponse{Success: true}, nil
}

func (api *API) KatanDeployStatusSwitchOn(ctx context.Context) (SuccessResponse, error) {
	by := ctx.HTTPRequest().PostFormValue("by")
	if by == "" {
		by = ctx.HTTPRequest().RemoteAddr
	}

	command := deploy.DeployOnCommand{}
	command.User = deploy.ActionUser{Name: by}
	command.Reason = ctx.HTTPRequest().PostFormValue("reason")
	command.OnlyFor = make([]string, 0)

	katan.DeployStatusCommand(&command)

	return SuccessResponse{Success: true}, nil
}
