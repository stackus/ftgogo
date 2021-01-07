package web

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type ServerCfg struct {
	Port              string        `envconfig:"PORT" default:":80"`
	CertPath          string        `envconfig:"CERT_PATH"`
	KeyPath           string        `envconfig:"KEY_PATH"`
	ReadTimeout       time.Duration `envconfig:"READ_TIMEOUT" default:"1s"`
	WriteTimeout      time.Duration `envconfig:"WRITE_TIMEOUT" default:"1s"`
	IdleTimeout       time.Duration `envconfig:"IDLE_TIMEOUT" default:"30s"`
	ReadHeaderTimeout time.Duration `envconfig:"READ_HEADER_TIMEOUT" default:"2s"`
	RequestTimeout    time.Duration `envconfig:"REQUEST_TIMEOUT" default:"60s"`
}

type CorsCfg struct {
	Origins          []string `envconfig:"ORIGINS" default:"*"`
	AllowCredentials bool     `envconfig:"ALLOW_CREDENTIALS" default:"true"`
	MaxAge           int      `envconfig:"MAX_AGE" default:"300"`
}

type Server interface {
	Start() error
	Shutdown(ctx context.Context) error
	Mount(path string, routeFn func(router chi.Router) http.Handler, options ...ServerOption)
}

type server struct {
	s    *http.Server
	root *chi.Mux
}

type tlsServer struct {
	server
}

func NewServer(cfg ServerCfg, options ...ServerOption) Server {
	root := chi.NewRouter()

	root.Use(
		middleware.Recoverer,
		middleware.Compress(5),
		middleware.Timeout(cfg.RequestTimeout),
	)

	for _, option := range options {
		option(root)
	}

	httpServer := &http.Server{
		Addr:              cfg.Port,
		Handler:           root,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	// setup TLS and listen for secure requests
	if cfg.CertPath != "" && cfg.KeyPath != "" {
		cert, err := tls.LoadX509KeyPair(cfg.CertPath, cfg.KeyPath)
		if err != nil {
			panic(err)
		}

		httpServer.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		root.Use(middleware.SetHeader("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload"))

		return &tlsServer{server: server{s: httpServer, root: root}}
	}

	return &server{s: httpServer, root: root}
}

// Start starts listening to non-TLS requests
func (s server) Start() error {
	return s.s.ListenAndServe()
}

// Shutdown executes a graceful shutdown of the http server
func (s server) Shutdown(ctx context.Context) error {
	return s.s.Shutdown(ctx)
}

func (s *server) Mount(path string, routeFn func(router chi.Router) http.Handler, options ...ServerOption) {
	router := chi.NewRouter()

	for _, option := range options {
		option(router)
	}

	s.root.Mount(path, routeFn(router))
}

// Start starts listening to TLS requests
func (s tlsServer) Start() error {
	return s.s.ListenAndServeTLS("", "")
}
