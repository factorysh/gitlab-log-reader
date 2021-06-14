package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/factorysh/gitlab-log-reader/state"
	"github.com/influxdata/tail"
	"github.com/valyala/fastjson"
)

func main() {
	t, err := tail.TailFile(os.Args[1], tail.Config{
		ReOpen:    true,
		MustExist: false,
		Follow:    true,
	})
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := state.NewState(ctx, 3*time.Hour)
	parser := &fastjson.Parser{}
	for {
		line := <-t.Lines
		value, err := parser.Parse(line.Text)
		if err != nil {
			fmt.Println(err)
			continue
		}
		tsRaw := value.GetStringBytes("time")
		ts, err := time.Parse(time.RFC1123, string(tsRaw))
		if err != nil {
			fmt.Println(err)
			continue
		}
		user := value.GetStringBytes("meta.user")
		if len(user) == 0 {
			continue
		}
		remote := value.GetStringBytes("meta.remote_ip")
		ua := value.GetStringBytes("ua")
		fmt.Printf("%v %s %s %s\n", ts, string(user), string(remote), string(ua))
		s.SetWithTimestamp(string(remote), ts, nil)
	}
}
