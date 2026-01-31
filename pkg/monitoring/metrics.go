package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTPRequestDuration tracks the latency of HTTP requests.
	HTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests in seconds.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path", "status"})

	// GRPCRequestDuration tracks the latency of gRPC requests.
	GRPCRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "grpc_request_duration_seconds",
		Help:    "Duration of gRPC requests in seconds.",
		Buckets: prometheus.DefBuckets,
	}, []string{"service", "method", "code"})

	// ErrorCounter tracks the number of errors encountered.
	ErrorCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "service_errors_total",
		Help: "Total number of errors encountered by the service.",
	}, []string{"service", "error_type"})
)
