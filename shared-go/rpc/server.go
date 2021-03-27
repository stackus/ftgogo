package rpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ServerCfg struct {
	Port     string `envconfig:"PORT" default:":8000"`
	CertPath string `envconfig:"CERT_PATH"`
	KeyPath  string `envconfig:"KEY_PATH"`
}

type Server interface {
	grpc.ServiceRegistrar
	Start() error
	Shutdown(ctx context.Context) error
}

type server struct {
	s    *grpc.Server
	port string
}

var _ Server = (*server)(nil)

func NewServer(cfg ServerCfg) Server {
	options := make([]grpc.ServerOption, 0)

	if cfg.KeyPath != "" && cfg.CertPath != "" {
		// setup TLS and listen for secure requests
		creds, err := credentials.NewServerTLSFromFile(cfg.CertPath, cfg.KeyPath)
		if err != nil {
			panic(err)
		}

		options = append(options, grpc.Creds(creds))
	}

	return server{
		s:    grpc.NewServer(options...),
		port: cfg.Port,
	}
}

func (s server) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	s.s.RegisterService(desc, impl)
}

func (s server) Start() error {
	listener, err := net.Listen("tcp", s.port)
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
