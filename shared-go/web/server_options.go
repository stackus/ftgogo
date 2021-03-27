package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type ServerOption func(chi.Router)

func WithNotFoundHandler(fn http.HandlerFunc) ServerOption {
	return func(r chi.Router) {
		r.NotFound(fn)
	}
}

func WithMethodNotAllowed(fn http.HandlerFunc) ServerOption {
	return func(r chi.Router) {
		r.MethodNotAllowed(fn)
	}
}

func WithHeader(key, value string) ServerOption {
	return func(r chi.Router) {
		r.Use(middleware.SetHeader(key, value))
	}
}

func WithMiddleware(mws ...func(http.Handler) http.Handler) ServerOption {
	return func(r chi.Router) {
		r.Use(mws...)
	}
}

func WithHealthCheck(path string) ServerOption {
	return func(r chi.Router) {
		r.Use(middleware.Heartbeat(path))
	}
}

func WithCors(cfg CorsCfg) ServerOption {
	return func(r chi.Router) {
		r.Use(cors.New(cors.Options{
			AllowedOrigins:   cfg.Origins,
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: cfg.AllowCredentials,
			MaxAge:           cfg.MaxAge,
		}).Handler)
	}
}
func WithSecure() ServerOption {
	return func(r chi.Router) {
		r.Use(
			middleware.SetHeader("X-Content-Type-Options", "nosniff"),
			middleware.SetHeader("X-Frame-Options", "SAMEORIGIN"),
			middleware.SetHeader("X-XSS-Protection", "1; mode-block"),
			middleware.NoCache,
		)
	}
}
