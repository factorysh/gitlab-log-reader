package main

import (
	"context"
	"os"

	"github.com/factorysh/gitlab-log-reader/rg"
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
	r.Loop(ctx2)
}
