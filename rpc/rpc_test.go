package rpc

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRpc(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "glr_")
	assert.NoError(t, err)
	tmp = fmt.Sprintf("%s/glr.sock", tmp)
	s := &Server{
		rpc: func(q *Query) (*Answer, error) {
			if q.IP == "127.0.0.1" {
				return &Answer{
					Ok: true,
				}, nil
			}
			if q.IP == "0.0.0.0" {
				return &Answer{}, errors.New("Bad IP")
			}
			return &Answer{
				Ok: false,
			}, nil
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = s.Listen(ctx, tmp)
	if err != nil {
		panic(err)
	}
	c, err := Dial(tmp)
	assert.NoError(t, err)
	a, err := c.Do(&Query{
		IP: "127.0.0.1",
	})
	assert.NoError(t, err)
	assert.True(t, a.Ok)

	a, err = c.Do(&Query{
		IP: "0.0.0.0",
	})
	assert.Error(t, err)
	assert.Nil(t, a)

	err = c.Close()
	assert.NoError(t, err)
}
