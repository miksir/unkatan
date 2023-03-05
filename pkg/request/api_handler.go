package request

import (
	"context"
	"encoding/json"
	zlog "github.com/miksir/unkatan/pkg/log"
	"go.uber.org/zap"
	"net/http"
)

type RequestError struct {
	Status  int    `json:"-"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

func (r *RequestError) Error() string {
	return r.Message
}

type ApiHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, error)

type apiHandler struct {
	logger  zlog.Logger
	handler ApiHandler
	w       http.ResponseWriter
	ctx     context.Context
}

func (h apiHandler) WriteHeader(statusCode int) {
	h.w.WriteHeader(statusCode)
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.w = w
	ctx := r.Context()
	h.ctx = ctx

	data, err := h.handler(ctx, w, r)

	if ctx.Err() != nil {
		if ctx.Err() == context.Canceled {
			h.WriteHeader(http.StatusGatewayTimeout)
		} else {
			h.WriteHeader(http.StatusServiceUnavailable)
		}
		return
	}

	if err != nil {
		if v, ok := err.(*RequestError); ok {
			var dt []byte
			dt, err = json.Marshal(v)
			if err != nil {
				h.WriteHeader(http.StatusInternalServerError)
				return
			}
			h.WriteHeader(v.Status)
			_, _ = w.Write(dt)
			return
		}

		h.WriteHeader(http.StatusInternalServerError)
		return
	}

	if data == nil {
		h.WriteHeader(http.StatusNoContent)
		return
	}

	b, err := json.Marshal(data)
	if err != nil {
		h.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(ctx, "failed to marshal data", zap.Error(err))
		return
	}

	_, err = w.Write(b)
	if err != nil {
		h.logger.Warn(ctx, "failed to write data to response", zap.Error(err))
	}
}

func ApiResponseHandler(logger zlog.Logger, handler ApiHandler) func(http.ResponseWriter, *http.Request) {
	return apiHandler{
		logger:  logger,
		handler: handler,
	}.ServeHTTP
}
