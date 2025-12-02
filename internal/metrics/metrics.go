// Package metrics exports Prometheus metrics for DNS lookups.
// Delegates metric collection to prometheus/client_golang.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// DNSLookupTotal tracks the total number of DNS lookups by server, query type, and result
	DNSLookupTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dns_lookup_total",
			Help: "Total number of DNS lookups",
		},
		[]string{"server", "query_type", "result"},
	)

	// DNSLookupDuration tracks DNS lookup duration in seconds
	DNSLookupDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "dns_lookup_duration_seconds",
			Help:    "DNS lookup duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"server", "query_type"},
	)

	// DNSLookupErrors tracks DNS lookup errors by server and error type
	DNSLookupErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dns_lookup_errors_total",
			Help: "Total number of DNS lookup errors",
		},
		[]string{"server", "error_type"},
	)

	// TasksTotal tracks the total number of DNS tasks by status
	TasksTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dns_tasks_total",
			Help: "Total number of DNS tasks",
		},
		[]string{"status"},
	)

	// APIRequestsTotal tracks API requests to submit DNS lookups
	APIRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dns_api_requests_total",
			Help: "Total number of API requests to submit DNS lookups",
		},
		[]string{"endpoint"},
	)

	// APIResultPollsTotal tracks result polling for coherence monitoring
	APIResultPollsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "dns_api_result_polls_total",
			Help: "Total number of times clients polled for task results (for coherence tracking)",
		},
	)

	// DNSResponseTime tracks DNS resolution time (Python dnstester compat).
	DNSResponseTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "dns_response_time_seconds",
			Help:    "Time taken for DNS resolution",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"server"},
	)

	// DNSTotalQueries tracks total DNS queries (Python dnstester compat).
	DNSTotalQueries = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dns_total_queries",
			Help: "Total number of DNS queries",
		},
		[]string{"server"},
	)

	// DNSNoErrorCount tracks successful DNS resolutions (Python dnstester compat).
	DNSNoErrorCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dns_noerror_count",
			Help: "Count of successful DNS resolutions (NoError)",
		},
		[]string{"server"},
	)

	// DNSFailureCount tracks failed DNS queries (Python dnstester compat).
	DNSFailureCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dns_failure_count",
			Help: "Total number of failed DNS queries",
		},
		[]string{"server", "rcode"},
	)

	// DNSAvgResponseTime tracks average DNS response time (Python dnstester compat).
	DNSAvgResponseTime = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dns_avg_response_time_seconds",
			Help: "Average DNS response time",
		},
		[]string{"server"},
	)

	// DNSQueryTypesCount tracks queries per query type (Python dnstester compat).
	DNSQueryTypesCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dns_query_types_count",
			Help: "Total number of DNS queries per query type",
		},
		[]string{"qtype"},
	)
)

// RecordQueryMetrics updates legacy metrics for Python dnstester dashboard compat.
func RecordQueryMetrics(server string, responseTimeSec float64, rcode, qtype string) {
	DNSTotalQueries.WithLabelValues(server).Inc()
	DNSResponseTime.WithLabelValues(server).Observe(responseTimeSec)
	DNSAvgResponseTime.WithLabelValues(server).Set(responseTimeSec)
	DNSQueryTypesCount.WithLabelValues(qtype).Inc()

	if rcode == "NOERROR" {
		DNSNoErrorCount.WithLabelValues(server).Inc()
	} else {
		DNSFailureCount.WithLabelValues(server, rcode).Inc()
	}
}
