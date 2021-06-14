package state

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	s := NewState(ctx, 100*time.Millisecond)
	s.Set("a", 42)
	v, ok := s.Get("a")
	assert.True(t, ok)
	assert.Equal(t, 42, v)
	time.Sleep(101 * time.Millisecond)
	_, ok = s.Get("a")
	assert.False(t, ok)
}
