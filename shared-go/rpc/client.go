package rpc

import (
	"google.golang.org/grpc"
)

type ClientCfg struct {
	CertPath string `envconfig:"CERT_PATH"`
	KeyPath  string `envconfig:"KEY_PATH"`
}

func NewClientConn(cfg ClientCfg, uri string, options ...ClientOption) (*grpc.ClientConn, error) {
	clientCfg := &clientConfig{}

	if cfg.KeyPath != "" && cfg.CertPath != "" {
		// TODO secure/mutual auth connections
	} else {
		clientCfg.AddOption(grpc.WithInsecure())
	}

	for _, option := range options {
		option(clientCfg)
	}

	return grpc.Dial(uri, clientCfg.ClientOptions()...)
}
