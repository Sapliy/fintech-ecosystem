package authutil

import (
	"context"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const InternalTokenHeader = "x-internal-token"

// UnaryInternalTokenServerInterceptor returns a new unary server interceptor that validates the internal service token.
func UnaryInternalTokenServerInterceptor() grpc.UnaryServerInterceptor {
	internalToken := os.Getenv("INTERNAL_SERVICE_TOKEN")
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// During development, if token is not set, we skip validation unless in production.
		if internalToken == "" {
			if os.Getenv("GO_ENV") == "production" {
				return nil, status.Error(codes.Unauthenticated, "INTERNAL_SERVICE_TOKEN must be set in production")
			}
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		values := md.Get(InternalTokenHeader)
		if len(values) == 0 || values[0] != internalToken {
			return nil, status.Error(codes.Unauthenticated, "invalid internal token")
		}

		return handler(ctx, req)
	}
}

// UnaryInternalTokenClientInterceptor returns a new unary client interceptor that injects the internal service token.
func UnaryInternalTokenClientInterceptor() grpc.UnaryClientInterceptor {
	internalToken := os.Getenv("INTERNAL_SERVICE_TOKEN")
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if internalToken != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, InternalTokenHeader, internalToken)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// ChainUnaryServer returns a single interceptor that executes the given interceptors in order.
func ChainUnaryServer(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ch := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			current := interceptors[i]
			next := ch
			ch = func(ctx context.Context, req interface{}) (interface{}, error) {
				return current(ctx, req, info, next)
			}
		}
		return ch(ctx, req)
	}
}
