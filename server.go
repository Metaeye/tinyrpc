// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tinyrpc

import (
	"log"
	"net"
	"net/rpc"

	"tinyrpc/codec"
	"tinyrpc/serializer"
)

// Server rpc server based on net/rpc implementation
type Server struct {
	*rpc.Server
	serializer.Serializer
}

// NewServer Create a new rpc server
func NewServer(opts ...Option) *Server {
	options := options{
		serializer: serializer.Proto,
	}
	for _, opt := range opts {
		opt(&options)
	}

	return &Server{&rpc.Server{}, options.serializer}
}

// Register register rpc function
func (s *Server) Register(rcvr any) error {
	return s.Server.Register(rcvr)
}

// RegisterName register the rpc function with the specified name
func (s *Server) RegisterName(name string, rcvr any) error {
	return s.Server.RegisterName(name, rcvr)
}

// Serve start service
func (s *Server) Serve(lis net.Listener) {
	log.Printf("tinyrpc started on: %s", lis.Addr().String())
	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}
		go s.Server.ServeCodec(codec.NewServerCodec(conn, s.Serializer))
	}
}
