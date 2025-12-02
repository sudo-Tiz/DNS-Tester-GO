# Monitoring

Prometheus metrics for performance and availability tracking.

---

## 📊 Metrics Endpoint

```bash
curl http://localhost:5000/metrics
```

---

## 📈 All Metrics

| Metric | Type | Description | Labels | Use Case |
|--------|------|-------------|--------|----------|
| `dns_lookup_total` | Counter | Total DNS lookups | `server`, `query_type`, `result` | Track query volume + success rate |
| `dns_lookup_duration_seconds` | Histogram | Lookup duration (all servers) | `server`, `query_type` | Measure latency, calculate P95/P99 |
| `dns_lookup_errors_total` | Counter | Total lookup errors | `server`, `error_type` | Identify problematic servers |
| `dns_tasks_total` | Counter | Total DNS tasks | `status` | Monitor async task processing |
| `dns_api_requests_total` | Counter | Total API requests | `endpoint` | Track API usage patterns |
| `dns_api_result_polls_total` | Counter | Result poll requests | - | Monitor polling frequency |
| `dns_response_time_seconds` | Histogram | DNS response time | `server` | Detailed server latency |
| `dns_total_queries` | Counter | Total queries per server | `server` | Query distribution |
| `dns_noerror_count` | Counter | Successful resolutions | `server` | Count successful queries |
| `dns_failure_count` | Counter | Failed queries | `server`, `rcode` | Track error types (NXDOMAIN, SERVFAIL) |
| `dns_avg_response_time_seconds` | Gauge | Average response time | `server` | Quick latency overview |
| `dns_query_types_count` | Counter | Queries by type | `qtype` | Distribution (A, AAAA, PTR, etc.) |

---

## 🔍 Common PromQL Queries

### Success Rate
```promql
# Per-server success rate (%)
sum by (server) (dns_lookup_total{result="success"}) /
sum by (server) (dns_lookup_total) * 100
```

### Latency Analysis
```promql
# P95 latency per server
histogram_quantile(0.95, 
  sum by (server, le) (rate(dns_lookup_duration_seconds_bucket[5m]))
)

# Average latency
rate(dns_lookup_duration_seconds_sum[5m]) /
rate(dns_lookup_duration_seconds_count[5m])
```

### Error Tracking
```promql
# Error rate per server
sum by (server) (rate(dns_lookup_errors_total[5m]))

# Top error types
topk(5, sum by (error_type) (rate(dns_lookup_errors_total[5m])))
```

### Query Distribution
```promql
# Queries per server
sum by (server) (rate(dns_total_queries[5m]))

# Query types distribution
sum by (qtype) (dns_query_types_count)
```

---

## 🎯 Quick Setup

### 1. Enable Metrics (API)
```yaml
# conf/config.yaml
api:
  enable_metrics: true
  metrics_port: 9090
```

### 2. Enable Metrics (Worker)
```bash
dnstestergo worker --enable-metrics --metrics-port 9091
```

### 3. Configure Prometheus
```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'dnstester-api'
    static_configs:
      - targets: ['localhost:9090']
  
  - job_name: 'dnstester-worker'
    static_configs:
      - targets: ['localhost:9091']
```

---

## 📊 Grafana Dashboard

**Key Panels:**

| Panel | Query | Visualization |
|-------|-------|---------------|
| Total Queries | `sum(dns_lookup_total)` | Single Stat |
| Success Rate | `sum(dns_lookup_total{result="success"})/sum(dns_lookup_total)*100` | Gauge |
| P95 Latency | `histogram_quantile(0.95, sum by (le) (rate(dns_lookup_duration_seconds_bucket[5m])))` | Graph |
| Errors/sec | `sum(rate(dns_lookup_errors_total[5m]))` | Graph |
| Query Types | `sum by (qtype) (dns_query_types_count)` | Pie Chart |
| Server Comparison | `sum by (server) (rate(dns_total_queries[5m]))` | Bar Chart |

---

## 🐛 Troubleshooting

| Problem | Check | Solution |
|---------|-------|----------|
| No metrics data | `curl localhost:9090/metrics` | Verify `enable_metrics: true` in config |
| Stale metrics | Perform test query | Metrics update on API requests |
| Missing worker metrics | Worker started with `--enable-metrics`? | Add flag and restart |
| Prometheus "down" | Check `http://localhost:9090/targets` | Verify ports and firewall |

---

## 📚 See Also

- [Configuration](05-configuration.md) - Enable metrics
- [Troubleshooting](07-troubleshooting.md) - Monitoring issues
- [API Reference](03-api.md) - API endpoints
