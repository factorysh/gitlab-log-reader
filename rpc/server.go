package rpc

import (
	"context"
	"encoding/gob"
	"fmt"
	"net"
)

type Server struct {
	rpc func(*Query) (*Answer, error)
}

func (s *Server) Listen(ctx context.Context, path string) error {
	ln, err := net.Listen("unix", path)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		ln.Close()
	}()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}
			go s.handles(conn)
		}
	}()

	return nil
}

func (s *Server) handles(conn net.Conn) {
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	var query Query
	err := dec.Decode(&query)
	if err != nil {
		fmt.Println(err)
		conn.Close()
	}
	answer, errA := s.rpc(&query)
	err = enc.Encode(answer)
	if err != nil {
		fmt.Println(err)
		conn.Close()
	}
	err = enc.Encode(Error{errA})
	if err != nil {
		fmt.Println(err)
		conn.Close()
	}
}
