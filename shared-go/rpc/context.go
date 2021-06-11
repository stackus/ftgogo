package rpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/stackus/edat/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Request tracking
const (
	requestIDKey     = "requestID"
	correlationIDKey = "correlationID"
	causationIDKey   = "causationID"
)

type clientStreamWrapper struct {
	grpc.ClientStream
}

func (s clientStreamWrapper) Context() context.Context {
	ctx := s.ClientStream.Context()

	md := metadata.New(map[string]string{
		requestIDKey:     core.GetRequestID(ctx),
		correlationIDKey: core.GetCorrelationID(ctx),
		causationIDKey:   core.GetCausationID(ctx),
	})

	return metadata.NewOutgoingContext(ctx, md)
}

type serverStreamWrapper struct {
	grpc.ServerStream
}

func (s serverStreamWrapper) Context() context.Context {
	ctx := s.ServerStream.Context()

	requestID := uuid.New().String()
	correlationID := requestID
	causationID := requestID

	md, _ := metadata.FromIncomingContext(ctx)
	vals := md.Get(requestIDKey)
	if len(vals) > 0 {
		requestID = vals[0]
	}

	vals = md.Get(correlationIDKey)
	if len(vals) > 0 {
		correlationID = vals[0]
	}

	vals = md.Get(causationIDKey)
	if len(vals) > 0 {
		causationID = vals[0]
	}

	return core.SetRequestContext(ctx, requestID, correlationID, causationID)
}

// Unary

func RequestContextUnaryServerInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	requestID := uuid.New().String()
	correlationID := requestID
	causationID := requestID

	md, _ := metadata.FromIncomingContext(ctx)
	vals := md.Get(requestIDKey)
	if len(vals) > 0 {
		requestID = vals[0]
	}

	vals = md.Get(correlationIDKey)
	if len(vals) > 0 {
		correlationID = vals[0]
	}

	vals = md.Get(causationIDKey)
	if len(vals) > 0 {
		causationID = vals[0]
	}

	return handler(core.SetRequestContext(ctx, requestID, correlationID, causationID), req)
}

func RequestContextUnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md := metadata.New(map[string]string{
		requestIDKey:     core.GetRequestID(ctx),
		correlationIDKey: core.GetCorrelationID(ctx),
		causationIDKey:   core.GetCausationID(ctx),
	})

	return invoker(metadata.NewOutgoingContext(ctx, md), method, req, reply, cc, opts...)
}

// Stream

func RequestContextStreamServerInterceptor(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return handler(srv, serverStreamWrapper{ss})
}

func RequestContextStreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	stream, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		return nil, err
	}
	return clientStreamWrapper{stream}, nil
}
