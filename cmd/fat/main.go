package main

import (
	"context"
	"net/http"
	"os"

	"github.com/apex/log"
	"github.com/factorysh/gitlab-log-reader/rg"
	"github.com/factorysh/gitlab-log-reader/web"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r, err := rg.New(ctx, os.Args[1])
	if err != nil {
		panic(err)
	}
	log.WithField("file", os.Args[1]).Info("Reading log file")

	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	go r.Loop(ctx2)

	adm := &http.Server{
		Addr:    "0.0.0.0:8042",
		Handler: web.NewAPI(r, web.Admin),
	}
	log.WithField("addr", adm.Addr).Info("Admin endpoint ready for listen")

	s := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: web.NewAPI(r, web.Auth),
	}
	log.WithField("addr", s.Addr).Info("Auth endpoint ready for listen")
	s.ListenAndServe()
}
