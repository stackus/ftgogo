package rpc

import (
	"context"

	"github.com/stackus/errors"
	"google.golang.org/grpc"
)

type clientConfig struct {
	options []grpc.DialOption
	unary   []grpc.UnaryClientInterceptor
	stream  []grpc.StreamClientInterceptor
}

func (c *clientConfig) AddOption(options ...grpc.DialOption) {
	c.options = append(c.options, options...)
}

func (c *clientConfig) AddUnaryInterceptor(interceptors ...grpc.UnaryClientInterceptor) {
	c.unary = append(c.unary, interceptors...)
}

func (c *clientConfig) AddStreamInterceptor(interceptors ...grpc.StreamClientInterceptor) {
	c.stream = append(c.stream, interceptors...)
}

func (c clientConfig) ClientOptions() []grpc.DialOption {
	options := c.options

	options = append(options, grpc.WithChainUnaryInterceptor(c.unary...))
	options = append(options, grpc.WithChainStreamInterceptor(c.stream...))

	return options
}

type ClientOption func(config *clientConfig)

func WithClientUnaryConvertStatus() ClientOption {
	return func(config *clientConfig) {
		config.AddUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return errors.ReceiveGRPCError(invoker(ctx, method, req, reply, cc, opts...))
		})
	}
}
