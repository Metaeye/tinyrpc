// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a git
// license that can be found in the LICENSE file.

package tinyrpc

import (
	"io"
	"net/rpc"

	"tinyrpc/codec"
	"tinyrpc/compressor"
	"tinyrpc/serializer"
)

// Client rpc client based on net/rpc implementation
type Client struct {
	*rpc.Client
}

// Option provides options for rpc
type Option func(o *options)

type options struct {
	Format     compressor.Format
	serializer serializer.Serializer
}

// WithCompress set client compression format
func WithCompress(c compressor.Format) Option {
	return func(o *options) {
		o.Format = c
	}
}

// WithSerializer set client serializer
func WithSerializer(serializer serializer.Serializer) Option {
	return func(o *options) {
		o.serializer = serializer
	}
}

// NewClient Create a new rpc client
func NewClient(conn io.ReadWriteCloser, opts ...Option) *Client {
	options := options{
		Format:     compressor.Raw,
		serializer: serializer.Proto,
	}
	for _, option := range opts {
		option(&options)
	}
	return &Client{
		rpc.NewClientWithCodec(
			codec.NewClientCodec(conn, options.Format, options.serializer),
		),
	}
}

// Call synchronously calls the rpc function
func (c *Client) Call(serviceMethod string, args any, reply any) error {
	return c.Client.Call(serviceMethod, args, reply)
}

// AsyncCall asynchronously calls the rpc function and returns a channel of *rpc.Call
func (c *Client) AsyncCall(serviceMethod string, args any, reply any) chan *rpc.Call {
	return c.Go(serviceMethod, args, reply, nil).Done
}
