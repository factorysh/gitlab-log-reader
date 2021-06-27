package web

import (
	"fmt"
	"net/http"

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
		w.WriteHeader(400)
		return
	}
	if a.rg.Exists(r.Header.Get("x-project"), r.Header.Get("x-remote")) {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(403)
	}
}
