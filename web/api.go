package web

import (
	"fmt"
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
	w.Header().Set("Server", "Log-Reader")
	fmt.Println("headers", r.Header)
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p := r.Header.Get("x-project")
	rip := r.Header.Get("x-remote")
	code := http.StatusForbidden
	if a.rg.Exists(p, rip) {
		code = http.StatusOK
	}
	w.WriteHeader(code)
	log.WithFields(log.Fields{"project": p, "remote ip": rip})
}
