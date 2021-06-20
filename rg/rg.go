package rg

import (
	"context"
	"fmt"
	"time"

	"github.com/factorysh/gitlab-log-reader/state"
	"github.com/influxdata/tail"
	"github.com/valyala/fastjson"
)

const timeFormat = "2006-01-02T15:04:05.000Z"

type RG struct {
	tail   *tail.Tail
	state  *state.State
	parser *fastjson.Parser
}

type Hit struct {
	User   string
	Remote string
	Ua     string
}

func New(ctx context.Context, path string) (*RG, error) {
	rg := &RG{
		state:  state.NewState(ctx, 3*time.Hour),
		parser: &fastjson.Parser{},
	}
	var err error
	rg.tail, err = tail.TailFile(path, tail.Config{
		ReOpen:    true,
		MustExist: false,
		Follow:    true,
	})
	if err != nil {
		return nil, err
	}
	return rg, nil
}

func (r *RG) processLine(line string) error {
	value, err := r.parser.Parse(line)
	if err != nil {
		return err
	}
	tsRaw := value.GetStringBytes("time")
	ts, err := time.Parse(timeFormat, string(tsRaw))
	if err != nil {
		return err
	}
	user := value.GetStringBytes("meta.user")
	if len(user) == 0 {
		return nil
	}
	remote := value.GetStringBytes("meta.remote_ip")
	ua := value.GetStringBytes("ua")
	r.state.SetWithTimestamp(string(remote), ts, Hit{
		User:   string(user),
		Remote: string(remote),
		Ua:     string(ua),
	})
	return nil
}

func (r *RG) Loop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			r.tail.Cleanup()
			return nil
		case line := <-r.tail.Lines:
			err := r.processLine(line.Text)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

func (r *RG) IPExists(ip string) bool {
	_, ok := r.state.Get(ip)
	return ok
}