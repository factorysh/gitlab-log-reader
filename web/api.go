package web

import (
	"net/http"

	"github.com/apex/log"
	"github.com/factorysh/gitlab-log-reader/rg"
)

type API struct {
	rg *rg.RG
}

func NewAPI(_rg *rg.RG) *API {
	return &API{
		rg: _rg,
	}
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.Header.Get("x-project")
	rip := r.Header.Get("x-remote")
	w.Header().Set("Server", "Log-Reader")
	if r.Method != "GET" {
		log.WithFields(log.Fields{
			"remote ip": rip,
			"method":    r.Method}).Info("Received invalid method")
		w.WriteHeader(http.StatusBadRequest)
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
