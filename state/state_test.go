package state

import (
	"context"
	"testing"
	"time"

	"github.com/factorysh/gitlab-log-reader/metrics"
	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	s := NewState(ctx, 100*time.Millisecond, metrics.Collector)
	ok := s.SetWithTimestamp(Key{"b", "", ""}, time.Now().Add(-101*time.Millisecond), nil)
	assert.False(t, ok)

	s.Set(Key{"a", "", ""}, 42)
	v, ok := s.Get(Key{"a", "", ""})
	assert.True(t, ok)
	assert.Equal(t, 42, v)
	time.Sleep(101 * time.Millisecond)
	_, ok = s.Get(Key{"a", "", ""})
	assert.False(t, ok)
}
