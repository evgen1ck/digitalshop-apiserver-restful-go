package api_v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"test-server-go/internal/models"
)

type Resolver struct {
	App *models.Application
}

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

func (rs *Resolver) SetupPrometheus() *chi.Mux {
	r := rs.App.Router

	prometheus.MustRegister(requestsProcessed)
	prometheus.MustRegister(requestDuration)

	r.Use(prometheusMiddleware)

	r.Handle("/prometheus/metrics", promhttp.Handler())

	return r
}
