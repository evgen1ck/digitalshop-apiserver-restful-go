package api_v1

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"test-server-go/internal/models"
)

var (
	requestsProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "requests_processed_total",
		Help: "Total number of requests processed",
	})
	requestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "request_duration_seconds",
		Help:    "Duration of processed HTTP requests in seconds",
		Buckets: prometheus.DefBuckets,
	})
)

func SetupPrometheus(app models.Application) {
	prometheus.MustRegister(requestsProcessed)
	prometheus.MustRegister(requestDuration)

	app.Router.Use(PrometheusMiddleware)

	app.Router.Handle("/prometheus/metrics", promhttp.Handler())
}
