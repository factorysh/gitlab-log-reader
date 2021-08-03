package web

import (
	"net/http"

	"github.com/apex/log"
	"github.com/factorysh/gitlab-log-reader/rg"
)

type httpHandler func(*API, http.ResponseWriter, *http.Request)

type API struct {
	rg      *rg.RG
	handler httpHandler
}

func NewAPI(_rg *rg.RG, _handler httpHandler) *API {
	return &API{
		rg:      _rg,
		handler: _handler,
	}
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler(a, w, r)
}

// Auth responds to nginx auth requests
func Auth(a *API, w http.ResponseWriter, r *http.Request) {
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
	if a.rg.Exists(p, rip) {
		code = http.StatusOK
	}
	w.WriteHeader(code)
	log.WithFields(log.Fields{
		"project":                  p,
		"remote ip":                rip,
		"status code":              code,
		"auth request remote addr": r.RemoteAddr,
		"auth request headers":     r.Header}).Info("Sending response")
}
