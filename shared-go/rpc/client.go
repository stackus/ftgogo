package rpc

import "google.golang.org/grpc"

func NewClientConn(uri string, options ...ClientOption) (*grpc.ClientConn, error) {
	clientCfg := &clientConfig{}

	for _, option := range options {
		option(clientCfg)
	}

	return grpc.Dial(uri, clientCfg.ClientOptions()...)
}
