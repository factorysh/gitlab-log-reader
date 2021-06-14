package main

import (
	"fmt"
	"os"

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
	parser := &fastjson.Parser{}
	for {
		line := <-t.Lines
		value, err := parser.Parse(line.Text)
		if err != nil {
			fmt.Println(err)
			continue
		}
		ts := value.GetStringBytes("time")
		user := value.GetStringBytes("meta.user")
		if len(user) == 0 {
			continue
		}
		remote := value.GetStringBytes("meta.remote_ip")
		ua := value.GetStringBytes("ua")
		fmt.Printf("%s %s %s %s\n", string(ts), string(user), string(remote), string(ua))
	}
}
