package state

import (
	"context"
	"sync"
	"time"
)

type State struct {
	lock   *sync.RWMutex
	values map[string]*Data
	maxAge time.Duration
}

type Data struct {
	ts    time.Time
	value interface{}
}

func NewState(ctx context.Context, maxAge time.Duration) *State {
	s := &State{
		lock:   &sync.RWMutex{},
		values: make(map[string]*Data),
		maxAge: maxAge,
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

func (s *State) Set(key string, value interface{}) {
	ok := s.SetWithTimestamp(key, time.Now(), value)
	if !ok {
		panic("We all gonna die")
	}
}

func (s *State) SetWithTimestamp(key string, ts time.Time, value interface{}) bool {
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
	return true
}

func (s *State) Get(key string) (interface{}, bool) {
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
	garbage := make([]string, 0)
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
		}
	}
}
