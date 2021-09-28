package web

import (
	"encoding/json"
	"flag"
	"net/http"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/factorysh/gitlab-log-reader/state"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// StateResult wraps all data about a state value
type StateResult struct {
	Project string        `json:"project"`
	IP      string        `json:"ip"`
	TTL     time.Duration `json:"ttl"`
	Access  time.Time     `json:"access"`
}

func toStateResults(entries state.StateValues) (allowed []StateResult) {
	allowed = make([]StateResult, 0)
	for k, v := range entries {
		access := v.Ts()
		allowed = append(allowed, StateResult{
			Project: k[0],
			IP:      k[1],
			TTL:     time.Now().Sub(access),
			Access:  access})
	}

	return allowed
}

// Admin is used to respond to requests on the admin endpoint
func Admin(a *API, w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		log.WithField("method", r.Method).Info("Admin handler received invalid method")
		return
	}

	// simple routing is okish for now
	switch strings.TrimRight(r.URL.Path, "/") {
	case "/allowlist":
		allowed := toStateResults(a.rg.State())
		data, err := json.Marshal(allowed)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	case "/sentry/hidden/test":
		if flag.Lookup("test.v") == nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			panic("Testing if panic is handled correctly by Sentry")
		}
	case "/metrics":
		h := promhttp.Handler()
		h.ServeHTTP(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}
