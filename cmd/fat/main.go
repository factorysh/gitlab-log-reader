package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/factorysh/gitlab-log-reader/metrics"
	"github.com/factorysh/gitlab-log-reader/rg"
	"github.com/factorysh/gitlab-log-reader/web"
	"github.com/getsentry/sentry-go"
)

var sentryTimeout = 2 * time.Second

func main() {
	if os.Getenv("SENTRY_DSN") == "" {
		log.Warn("No SENTRY_DSN environment variable specified")
	}

	err := sentry.Init(sentry.ClientOptions{})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer sentry.Flush(sentryTimeout)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r, err := rg.New(ctx, os.Args[1], metrics.Collector)
	if err != nil {
		panic(err)
	}
	log.WithField("file", os.Args[1]).Info("Reading log file")

	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	go r.Loop(ctx2)

	adm := &http.Server{
		Addr:    "0.0.0.0:8042",
		Handler: web.NewAPI(r, web.Admin, nil),
	}
	log.WithField("addr", adm.Addr).Info("Admin endpoint ready for listen")
	go adm.ListenAndServe()

	s := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: web.NewAPI(r, web.Auth, metrics.Collector),
	}
	log.WithField("addr", s.Addr).Info("Auth endpoint ready for listen")
	s.ListenAndServe()
}
