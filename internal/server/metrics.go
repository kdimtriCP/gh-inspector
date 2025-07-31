package server

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/kdimtriCP/gh-inspector/internal/metrics"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gh_inspector_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "gh_inspector_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	repositoryAnalysisTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gh_inspector_repository_analysis_total",
			Help: "Total number of repository analyses",
		},
		[]string{"status"},
	)

	repositoryAnalysisDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "gh_inspector_repository_analysis_duration_seconds",
			Help:    "Repository analysis duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)

	cacheHits = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "gh_inspector_cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	cacheMisses = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "gh_inspector_cache_misses_total",
			Help: "Total number of cache misses",
		},
	)
)

type MetricsRecorder struct{}

func NewMetricsRecorder() metrics.Recorder {
	return &MetricsRecorder{}
}

func (m *MetricsRecorder) RecordHTTPRequest(method, endpoint, status string) {
	httpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
}

func (m *MetricsRecorder) RecordHTTPDuration(method, endpoint string, duration time.Duration) {
	httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

func (m *MetricsRecorder) RecordRepositoryAnalysis(status string, duration time.Duration) {
	repositoryAnalysisTotal.WithLabelValues(status).Inc()
	if status == "success" {
		repositoryAnalysisDuration.Observe(duration.Seconds())
	}
}

func (m *MetricsRecorder) RecordCacheHit() {
	cacheHits.Inc()
}

func (m *MetricsRecorder) RecordCacheMiss() {
	cacheMisses.Inc()
}

func recordHTTPRequest(method, endpoint, status string) {
	httpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
}

func recordRepositoryAnalysis(status string) {
	repositoryAnalysisTotal.WithLabelValues(status).Inc()
}
