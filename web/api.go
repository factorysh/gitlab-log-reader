package web

import (
	"net/http"
	"time"

	"github.com/apex/log"
	"github.com/factorysh/gitlab-log-reader/metrics"
	"github.com/factorysh/gitlab-log-reader/rg"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

type httpHandler func(*API, http.ResponseWriter, *http.Request)

type API struct {
	rg            *rg.RG
	handler       httpHandler
	metrics       *metrics.Gatherer
	sentryHandler *sentryhttp.Handler
}

func NewAPI(_rg *rg.RG, _handler httpHandler, m *metrics.Gatherer) *API {
	return &API{
		rg:            _rg,
		handler:       _handler,
		metrics:       m,
		sentryHandler: sentryhttp.New(sentryhttp.Options{}),
	}
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.sentryHandler.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		a.handler(a, w, r)
	}).ServeHTTP(w, r)
}

// Auth responds to nginx auth requests
func Auth(a *API, w http.ResponseWriter, r *http.Request) {
	a.metrics.AuthRequestCounter.Inc()
	p := r.Header.Get("x-project")
	rip := r.Header.Get("x-remote")
	w.Header().Set("Server", "Log-Reader")
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"remote ip": rip,
			"method":    r.Method}).Info("Auth handler received invalid method")
		return
	}

	code := http.StatusForbidden
	if exp, valid := a.rg.Expires(p, rip); valid {
		exp.Truncate(time.Second)
		w.Header().Set("Expires", exp.String())
		code = http.StatusOK
		a.metrics.StatusOkRespCounter.Inc()
	} else {
		a.metrics.StatusForbiddenRespCounter.Inc()
	}
	w.WriteHeader(code)
	log.WithFields(log.Fields{
		"project":                  p,
		"remote ip":                rip,
		"status code":              code,
		"auth request remote addr": r.RemoteAddr,
		"auth request headers":     r.Header}).Info("Sending response")
}
