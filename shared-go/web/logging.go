package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ZapLogger(logger *zap.Logger) func(next http.Handler) http.Handler {
	logger = logger.WithOptions(
		// Drop useless caller info (would always be rest/logging.go)
		zap.WithCaller(false),
		// Only include a stacktrace for unhandled panics
		zap.AddStacktrace(zapcore.PanicLevel),
	).Named("WebLogger")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ww := middleware.NewWrapResponseWriter(writer, request.ProtoMajor)

			start := time.Now()

			defer func() {
				var logFn func(string, ...zap.Field)

				err := recover()

				switch {
				case err != nil:
					logFn = logger.
						WithOptions(
							zap.WithCaller(true),
							zap.AddCallerSkip(4), // 4 seems about right-ish
							zap.AddStacktrace(zapcore.ErrorLevel),
						).
						With(zap.String("Error", fmt.Sprint(err))).
						Error
					// ensure the status code reflects this panic
					if ww.Status() < 500 {
						ww.WriteHeader(http.StatusInternalServerError)
					}
				case ww.Status() < 400:
					logFn = logger.Info
				case ww.Status() < 500:
					logFn = logger.Warn
				default:
					logFn = logger.Error
				}
				logFn(
					fmt.Sprintf("[%d] %s %s", ww.Status(), request.Method, request.RequestURI),
					zap.String("RemoteAddr", request.RemoteAddr),
					zap.Int("ContentLength", ww.BytesWritten()),
					zap.Duration("ResponseTime", time.Since(start)),
				)
			}()

			next.ServeHTTP(ww, request)
		})
	}
}

func ZeroLogger(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ww := middleware.NewWrapResponseWriter(writer, request.ProtoMajor)

			start := time.Now()

			defer func() {
				var err error
				var logFn func() *zerolog.Event

				p := recover()

				switch {
				case p != nil:
					logFn = logger.Error().Stack
					// ensure the status code reflects this panic
					if ww.Status() < 500 {
						ww.WriteHeader(http.StatusInternalServerError)
					}
					err = errors.Errorf("%s", p)
				case ww.Status() < 400:
					logFn = logger.Info
				case ww.Status() < 500:
					logFn = logger.Warn
				default:
					logFn = logger.Error
				}
				log := logFn()
				if err != nil {
					log = log.Err(err)
				}
				log.Str("RemoteAddr", request.RemoteAddr).
					Int("ContextLength", ww.BytesWritten()).
					Dur("ResponseTime", time.Since(start)).
					Msgf("[%d] %s %s", ww.Status(), request.Method, request.RequestURI)
			}()

			next.ServeHTTP(ww, request)
		})
	}
}
