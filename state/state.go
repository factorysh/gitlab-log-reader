package state

import (
	"context"
	"sync"
	"time"

	"github.com/factorysh/gitlab-log-reader/metrics"
)

type Key [3]string

type StateValues map[Key]*Data

type State struct {
	lock    *sync.RWMutex
	values  StateValues
	maxAge  time.Duration
	metrics *metrics.Gatherer
}

type Data struct {
	ts    time.Time
	value interface{}
}

func (d *Data) Ts() time.Time {
	return d.ts
}

func NewState(ctx context.Context, maxAge time.Duration, m *metrics.Gatherer) *State {
	s := &State{
		lock:    &sync.RWMutex{},
		values:  make(StateValues),
		maxAge:  maxAge,
		metrics: m,
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(maxAge):
				s.gc()
			}
		}
	}()
	return s
}

func (s *State) Set(key [3]string, value interface{}) {
	ok := s.SetWithTimestamp(key, time.Now(), value)
	if !ok {
		panic("We all gonna die")
	}
}

func (s *State) Values() StateValues {
	return s.values
}

func (s *State) SetWithTimestamp(key [3]string, ts time.Time, value interface{}) bool {
	delta := time.Since(ts)
	if delta > s.maxAge { // false
		return false
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	s.values[key] = &Data{
		ts:    ts,
		value: value,
	}
	s.metrics.AllowListSize.Inc()
	return true
}

func (s *State) Get(key [3]string) (interface{}, bool) {
	s.lock.RLock()
	v, ok := s.values[key]
	if !ok {
		s.lock.RUnlock()
		return nil, false
	}
	if -time.Until(v.ts) > s.maxAge { // value is rotten, lets delete it
		s.lock.RUnlock()
		s.lock.Lock()
		delete(s.values, key)
		s.lock.Unlock()
		return nil, false
	}
	s.lock.RUnlock()
	return v.value, true
}

func (s *State) gc() {
	garbage := make([]Key, 0)
	s.lock.RLock()
	for k, v := range s.values {
		if time.Until(v.ts) > s.maxAge { // value is rotten, lets delete it
			garbage = append(garbage, k)
		}
	}
	s.lock.RUnlock()
	if len(garbage) > 0 {
		s.lock.Lock()
		defer s.lock.Unlock()
		for _, k := range garbage {
			delete(s.values, k)
			s.metrics.AllowListSize.Dec()
		}
	}
}
