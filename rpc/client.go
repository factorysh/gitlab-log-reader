package rpc

import (
	"encoding/gob"
	"net"
	"os"
	"time"
)

type Client struct {
	conn net.Conn
	enc  *gob.Encoder
	dec  *gob.Decoder
}

func Dial(path string) (*Client, error) {
	for i := 0; i < 10; i++ {
		_, err := os.Stat(path)
		if err != nil {
			break
		}
		if os.IsNotExist(err) {
			time.Sleep(100 * time.Millisecond)
		}
	}
	conn, err := net.Dial("unix", path)
	if err != nil {
		return nil, err
	}

	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)
	return &Client{
		conn: conn,
		enc:  enc,
		dec:  dec,
	}, nil
}

func (c *Client) Do(query *Query) (*Answer, error) {
	err := c.enc.Encode(query)
	if err != nil {
		return nil, err
	}
	var answer Answer
	err = c.dec.Decode(&answer)
	if err != nil {
		return nil, err
	}
	var errA Error
	err = c.dec.Decode(&errA)
	if err != nil {
		return nil, err
	}
	return &answer, errA.Error
}

func (c *Client) Close() error {
	return c.conn.Close()
}
