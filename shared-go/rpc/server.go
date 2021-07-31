package rpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type ServerCfg struct {
	Network  string `envconfig:"NETWORK" default:"tcp"`
	Address  string `envconfig:"ADDRESS" default:":8000"`
	CertPath string `envconfig:"CERT_PATH"`
	KeyPath  string `envconfig:"KEY_PATH"`
}

type Server interface {
	grpc.ServiceRegistrar
	Start() error
	Shutdown(ctx context.Context) error
}

type server struct {
	s       *grpc.Server
	network string
	address string
}

var _ Server = (*server)(nil)

func NewServer(cfg ServerCfg, options ...ServerOption) Server {
	serverCfg := &serverConfig{}

	if cfg.KeyPath != "" && cfg.CertPath != "" {
		// setup TLS and listen for secure requests
		creds, err := credentials.NewServerTLSFromFile(cfg.CertPath, cfg.KeyPath)
		if err != nil {
			panic(err)
		}

		serverCfg.AddOption(grpc.Creds(creds))
	}

	for _, option := range options {
		option(serverCfg)
	}

	s := grpc.NewServer(serverCfg.ServerOptions()...)

	reflection.Register(s)

	return server{
		s:       s,
		network: cfg.Network,
		address: cfg.Address,
	}
}

func (s server) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	s.s.RegisterService(desc, impl)
}

func (s server) Start() error {
	listener, err := net.Listen(s.network, s.address)
	if err != nil {
		panic(err)
	}
	return s.s.Serve(listener)
}

func (s server) Shutdown(ctx context.Context) error {
	// GracefulStop does not accept a context so no timeout or cancel is possible
	stopped := make(chan struct{})
	go func() {
		s.s.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		// Force it to stop
		s.s.Stop()
		return fmt.Errorf("rpc server failed to stop gracefully")
	case <-stopped:
		return nil
	}
}
