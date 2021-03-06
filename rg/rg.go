package rg

import (
	"context"
	"fmt"
	"time"

	"github.com/factorysh/gitlab-log-reader/metrics"
	"github.com/factorysh/gitlab-log-reader/state"
	"github.com/influxdata/tail"
	"github.com/valyala/fastjson"
)

const TimeFormat = "2006-01-02T15:04:05.000Z"

type RG struct {
	tail   *tail.Tail
	state  *state.State
	parser *fastjson.Parser
}

func NewRG(tail *tail.Tail, state *state.State) *RG {
	return &RG{
		tail:   tail,
		state:  state,
		parser: &fastjson.Parser{},
	}
}

func New(ctx context.Context, path string, m *metrics.Gatherer) (*RG, error) {
	rg := &RG{
		state:  state.NewState(ctx, 3*time.Hour, m),
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

// State return current rg state
func (r *RG) State() map[state.Key]*state.Data {
	return r.state.Values()
}

func (r *RG) ProcessLine(line string) error {
	value, err := r.parser.Parse(line)
	if err != nil {
		return err
	}
	tsRaw := value.GetStringBytes("time")
	ts, err := time.Parse(TimeFormat, string(tsRaw))
	if err != nil {
		return err
	}
	project := value.GetStringBytes("meta.project")
	if len(project) == 0 {
		return nil
	}
	user := value.GetStringBytes("meta.user")
	if len(user) == 0 {
		return nil
	}
	remote := value.GetStringBytes("meta.remote_ip")
	//ua := value.GetStringBytes("ua")
	r.state.SetWithTimestamp(
		state.Key{
			string(project),
			string(remote),
			"",
			//Ua:      string(ua),
		},
		ts,
		nil,
	)
	return nil
}

func (r *RG) Loop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			r.tail.Cleanup()
			return nil
		case line := <-r.tail.Lines:
			if line == nil {
				continue
			}
			err := r.ProcessLine(line.Text)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

func (r *RG) Exists(project, remote string) bool {
	_, ok := r.state.Get(state.Key{project, remote, ""})
	return ok
}

func (r *RG) Expires(project, remote string) (time.Duration, bool) {
	data, ok := r.state.GetData(state.Key{project, remote, ""})
	if !ok {
		return 0, false
	}

	// data.Ts() ----> time.Now() ----> data.Ts+s.MaxAge
	// remaining duration is maxAge - delta(now, data.Ts())
	return r.state.MaxAge() - time.Since(data.Ts()), true
}
