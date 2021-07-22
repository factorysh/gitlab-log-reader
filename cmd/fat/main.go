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

	s := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: web.NewAPI(r, web.Auth),
	}
	log.WithField("addr", s.Addr).Info("Ready for listen")
	s.ListenAndServe()
}
