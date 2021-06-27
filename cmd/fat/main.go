package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

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

	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	go r.Loop(ctx2)

	s := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: web.NewAPI(r),
	}
	fmt.Println("Listen", s.Addr)
	s.ListenAndServe()
}
