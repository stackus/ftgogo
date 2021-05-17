package rpc

import "google.golang.org/grpc"

type serverConfig struct {
	options []grpc.ServerOption
	unary   []grpc.UnaryServerInterceptor
	stream  []grpc.StreamServerInterceptor
}

func (c *serverConfig) AddOption(option grpc.ServerOption) {
	c.options = append(c.options, option)
}

func (c *serverConfig) AddUnaryInterceptor(interceptor grpc.UnaryServerInterceptor) {
	c.unary = append(c.unary, interceptor)
}

func (c *serverConfig) AddStreamInterceptor(interceptor grpc.StreamServerInterceptor) {
	c.stream = append(c.stream, interceptor)
}

func (c serverConfig) ServerOptions() []grpc.ServerOption {
	options := c.options

	options = append(options, grpc.ChainUnaryInterceptor(c.unary...))
	options = append(options, grpc.ChainStreamInterceptor(c.stream...))

	return options
}

type ServerOption func(config *serverConfig)
