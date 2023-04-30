package server

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strings"
	"test-server-go/internal/models"
	"time"
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

func setupPrometheus(app models.Application, path, promPort string) {
	prometheus.MustRegister(requestsProcessed)
	prometheus.MustRegister(requestDuration)

	app.Router.Use(prometheusMiddleware)
	app.Router.Use(blockPrometheusMetricsMiddleware(path, promPort))

	app.Router.Handle(path, promhttp.Handler())
}

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)

		requestDuration.Observe(elapsed.Seconds())
		requestsProcessed.Inc()
	})
}

func blockPrometheusMetricsMiddleware(path, promPort string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, path) && !strings.HasSuffix(r.Host, ":"+promPort) {
				http.Error(w, "Forbidden. Prometheus metrics is not available here", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
