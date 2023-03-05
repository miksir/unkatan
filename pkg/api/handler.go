package api

import (
	"github.com/miksir/unkatan/pkg/lconfig"
	"github.com/miksir/unkatan/pkg/request"
	"net/http"
	"net/url"
	"strings"

	"context"
	zlog "github.com/miksir/unkatan/pkg/log"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type API struct {
	cfg      lconfig.Reader
	log      zlog.Logger
	slackAPI *slack.Client
}

func NewAPI(cfg lconfig.Reader, log zlog.Logger, slackAPI *slack.Client) *API {
	comm := &API{
		cfg:      cfg,
		log:      log,
		slackAPI: slackAPI,
	}
	return comm
}

func (api *API) APIHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var err error
	w.Header().Set("Content-Type", "application/json")

	u, err := url.Parse(r.RequestURI)
	if err != nil {
		api.log.Error(
			ctx,
			"Malformed request URL",
			zap.String("url", r.RequestURI),
			zap.Error(err),
		)
		return nil, &request.RequestError{
			Status: http.StatusBadRequest,
		}
	}

	u.Path = strings.TrimRight(u.Path, "/")
	methodAndPath := r.Method + " " + u.Path

	switch methodAndPath {
	case "GET /api/deploy":
		return api.KatanDeployGetApi(ctx)
	case "GET /ready":
		return &SuccessResponse{Success: true}, nil
	case "GET /live":
		return &SuccessResponse{Success: true}, nil
	case "POST /api/deploy/off":
		if tokenValid := api.checkToken(ctx, r); tokenValid != nil {
			return nil, tokenValid
		}
		return api.KatanDeployStatusSwitchOff(ctx)
	case "POST /api/deploy/on":
		if tokenValid := api.checkToken(ctx, r); tokenValid != nil {
			return nil, tokenValid
		}
		return api.KatanDeployStatusSwitchOn(ctx)
	}

	return nil, &request.RequestError{
		Status: http.StatusNotFound,
	}
}

func (api *API) checkToken(ctx context.Context, r *http.Request) error {
	token := r.Header.Get("Authorization")
	match := false
	if token != "" {
		for _, apiToken := range api.cfg.GetStringSlice("auth_tokens") {
			if token == "Bearer "+apiToken {
				match = true
			}
		}
	}
	if !match {
		api.log.Error(
			ctx,
			"Authorization error",
			zap.String("url", r.RequestURI),
		)
		status := http.StatusForbidden
		if token == "" {
			status = http.StatusUnauthorized
		}
		return &request.RequestError{
			Status: status,
		}
	}
	return nil
}
