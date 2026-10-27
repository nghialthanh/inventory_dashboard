package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// CatalogCacheHits counts catalog cache hits by resource (makes|models).
	CatalogCacheHits = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "catalog_cache_hits_total",
		Help: "Catalog list responses served from Redis",
	}, []string{"resource"})

	// CatalogCacheMisses counts catalog cache misses by resource.
	CatalogCacheMisses = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "catalog_cache_misses_total",
		Help: "Catalog list responses that missed Redis and read Postgres",
	}, []string{"resource"})

	// CatalogCacheErrors counts Redis failures that fell open to Postgres.
	CatalogCacheErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "catalog_cache_errors_total",
		Help: "Redis errors on catalog cache get/set (fail-open)",
	}, []string{"resource", "op"})

	// DBQueries counts repository queries by role (primary|replica).
	DBQueries = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "db_query_total",
		Help: "Database queries by pool role",
	}, []string{"role", "op"})

	// HTTPRequests counts finished HTTP requests.
	HTTPRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "HTTP requests by method, route, and status class",
	}, []string{"method", "path", "status_class"})

	// HTTPDuration observes request latency.
	HTTPDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request latency in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})
)
