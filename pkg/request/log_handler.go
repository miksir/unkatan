package request

import (
	zlog "github.com/miksir/unkatan/pkg/log"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	headerXRealIP = "X-Real-IP"
)

type logHandler struct {
	log zlog.Logger
	h   http.Handler
}

func (l logHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.h.ServeHTTP(w, r)
	passed := time.Since(start)

	l.log.Info(
		r.Context(),
		"",
		zap.Time("time", start),
		zap.Int64("duration", passed.Microseconds()),
		zap.String(headerXRealIP, r.Header.Get(headerXRealIP)),
		zap.String("URL", r.URL.RequestURI()),
		zap.String("method", r.Method),
	)
}

func AccessLogHandler(logger zlog.Logger, handler http.Handler) http.Handler {
	return logHandler{
		log: logger,
		h:   handler,
	}
}

func AccessLogResponseHandler(logger zlog.Logger, handler http.Handler) func(http.ResponseWriter, *http.Request) {
	return logHandler{
		log: logger,
		h:   handler,
	}.ServeHTTP
}
