package monitoring

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor returns a new unary server interceptor that tracks gRPC request duration.
func UnaryServerInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start).Seconds()

		st, _ := status.FromError(err)
		code := st.Code().String()

		GRPCRequestDuration.WithLabelValues(serviceName, info.FullMethod, code).Observe(duration)

		if err != nil {
			ErrorCounter.WithLabelValues(serviceName, "grpc_error").Inc()
		}

		return resp, err
	}
}

// UnaryClientInterceptor returns a new unary client interceptor that tracks gRPC request duration.
func UnaryClientInterceptor(serviceName string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		duration := time.Since(start).Seconds()

		st, _ := status.FromError(err)
		code := st.Code().String()

		GRPCRequestDuration.WithLabelValues(serviceName, method, code).Observe(duration)

		return err
	}
}
