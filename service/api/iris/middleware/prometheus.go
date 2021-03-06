package middleware

import (
	"strconv"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	reqsName    = "http_requests_total"
	latencyName = "http_request_duration_seconds"
)

type metrics struct {
	reqs    *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

func newMetrics(name string) *metrics {
	m := metrics{}

	m.reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        reqsName,
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
			ConstLabels: prometheus.Labels{"service": name},
		},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.reqs)

	m.latency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        latencyName,
			Help:        "How long it took to process the request, partitioned by status code, method and HTTP path.",
			ConstLabels: prometheus.Labels{"service": name},
		},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.latency)

	return &m
}

func NewPrometheus(serviceName string) iris.Handler {
	m := newMetrics(serviceName)
	return func(ctx context.Context) {
		start := time.Now()
		ctx.Next()
		elapsed := float64(time.Since(start).Nanoseconds() / (int64(time.Millisecond) / int64(time.Nanosecond)))

		r := ctx.Request()
		statusCode := strconv.Itoa(ctx.GetStatusCode())

		m.reqs.WithLabelValues(statusCode, r.Method, r.URL.Path).Inc()
		m.latency.WithLabelValues(statusCode, r.Method, r.URL.Path).Observe(elapsed)
	}
}
