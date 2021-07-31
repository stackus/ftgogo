package rpc

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stackus/edat/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func WithUnaryClientLogging(logger zerolog.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		startedAt := time.Now()

		requestID := core.GetRequestID(ctx)
		correlationID := core.GetCorrelationID(ctx)
		causationID := core.GetCausationID(ctx)

		defer func() {
			var logFn func() *zerolog.Event
			statusCode := codes.OK

			p := recover()

			switch {
			case p != nil:
				logFn = logger.Error().Stack
				err = errors.Errorf("%s", p)
			case err != nil:
				logFn = logger.Error
			default:
				logFn = logger.Info
			}
			log := logFn()
			if err != nil {
				log = log.Err(err)
				if s, ok := status.FromError(err); ok {
					statusCode = s.Code()
				} else {
					statusCode = codes.OK
				}
			}
			log = log.Str("Method", method).
				Dur("RoundTripTime", time.Since(startedAt))

			if requestID != "" {
				log = log.Str("RequestID", requestID).
					Str("CorrelationID", correlationID).
					Str("CausationID", causationID)
			}

			log.Msgf("[%s] Client", statusCode.String())
		}()

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func WithUnaryServerLogging(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		startedAt := time.Now()

		requestID := core.GetRequestID(ctx)
		correlationID := core.GetCorrelationID(ctx)
		causationID := core.GetCausationID(ctx)

		defer func() {
			var logFn func() *zerolog.Event
			statusCode := codes.OK

			p := recover()

			switch {
			case p != nil:
				logFn = logger.Error().Stack
				err = errors.Errorf("%s", p)
			case err != nil:
				logFn = logger.Error
			default:
				logFn = logger.Info
			}
			log := logFn()
			if err != nil {
				log = log.Err(err)
				if s, ok := status.FromError(err); ok {
					statusCode = s.Code()
				} else {
					statusCode = codes.OK
				}
			}

			log = log.Str("Method", info.FullMethod).
				Dur("ResponseTime", time.Since(startedAt))

			if requestID != "" {
				log = log.Str("RequestID", requestID).
					Str("CorrelationID", correlationID).
					Str("CausationID", causationID)
			}

			log.Msgf("[%s] Server", statusCode.String())
		}()

		return handler(ctx, req)
	}
}
