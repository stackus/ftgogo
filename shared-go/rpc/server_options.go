package rpc

import (
	"context"

	"github.com/stackus/errors"
	"google.golang.org/grpc"
)

type serverConfig struct {
	options []grpc.ServerOption
	unary   []grpc.UnaryServerInterceptor
	stream  []grpc.StreamServerInterceptor
}

func (c *serverConfig) AddOption(options ...grpc.ServerOption) {
	c.options = append(c.options, options...)
}

func (c *serverConfig) AddUnaryInterceptor(interceptors ...grpc.UnaryServerInterceptor) {
	c.unary = append(c.unary, interceptors...)
}

func (c *serverConfig) AddStreamInterceptor(interceptors ...grpc.StreamServerInterceptor) {
	c.stream = append(c.stream, interceptors...)
}

func (c serverConfig) ServerOptions() []grpc.ServerOption {
	options := c.options

	options = append(options, grpc.ChainUnaryInterceptor(c.unary...))
	options = append(options, grpc.ChainStreamInterceptor(c.stream...))

	return options
}

type ServerOption func(config *serverConfig)

func WithServerUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) ServerOption {
	return func(config *serverConfig) {
		config.AddUnaryInterceptor(interceptors...)
	}
}

func WithServerStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) ServerOption {
	return func(config *serverConfig) {
		config.AddStreamInterceptor(interceptors...)
	}
}

func WithServerUnaryEnsureStatus() ServerOption {
	return func(config *serverConfig) {
		config.AddUnaryInterceptor(
			func(
				ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
			) (resp interface{}, err error) {
				resp, err = handler(ctx, req)
				return resp, errors.SendGRPCError(err)
			},
		)
	}
}
