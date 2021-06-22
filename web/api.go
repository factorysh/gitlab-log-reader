package web

import (
	"net/http"

	"github.com/factorysh/gitlab-log-reader/rg"
)

type API struct {
	rg *rg.RG
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(400)
		return
	}
	q := r.URL.Query()
	if a.rg.Exists(q.Get("project"), q.Get("remote")) {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(403)
	}
}
