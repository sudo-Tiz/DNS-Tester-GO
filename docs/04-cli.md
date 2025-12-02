# CLI Guide

dns-tester-go provides three command-line tools for different use cases.

---

## `dnstestergo query` - DNS Query Tool

Perform DNS lookups with multi-protocol support (UDP, TCP, DoT, DoH, DoQ).

### Usage

```bash
dnstestergo query <domain> [dns_servers...] [flags]
```

### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-u, --api-url` | string | `http://localhost:5000` | API base URL |
| `-t, --qtype` | string | `A` | Query type (A, AAAA, MX, TXT, PTR, etc.) |
| `-i, --insecure` | bool | `false` | Skip TLS certificate verification |
| `-d, --debug` | bool | `false` | Show detailed error messages |
| `-p, --pretty` | bool | `false` | Enable emoji-enhanced output |
| `-w, --warn-threshold` | float | `1.0` | Response time warning threshold (seconds) |
| `-c, --config` | string | - | Path to config file |

### Examples

**Basic queries:**
```bash
# UDP query
dnstestergo query example.com udp://8.8.8.8:53

# Multiple protocols
dnstestergo query cloudflare.com \
  udp://1.1.1.1:53 \
  tls://one.one.one.one:853 \
  https://cloudflare-dns.com/dns-query

# IPv6 lookup
dnstestergo query google.com udp://8.8.8.8:53 -t AAAA

# DNS-over-QUIC
dnstestergo query example.com quic://dns.adguard-dns.com:853
```

**Advanced usage:**
```bash
# With config file (uses servers from config)
dnstestergo query example.com -c conf/config.yaml

# Debug mode
dnstestergo query example.com udp://8.8.8.8:53 -d

# Pretty output with emojis
dnstestergo query github.com udp://9.9.9.9:53 -p

# Custom response time threshold
dnstestergo query example.com udp://slow.dns:53 -w 0.05

# Different query types
dnstestergo query example.com udp://8.8.8.8:53 -t MX   # Mail servers
dnstestergo query example.com udp://8.8.8.8:53 -t TXT  # TXT records
dnstestergo query example.com udp://8.8.8.8:53 -t NS   # Nameservers
```

### Protocol Formats

| Protocol | Format | Example |
|----------|--------|---------|
| UDP | `udp://<ip>:<port>` | `udp://8.8.8.8:53` |
| TCP | `tcp://<ip>:<port>` | `tcp://1.1.1.1:53` |
| DoT | `tls://<host>:<port>` | `tls://dns.google:853` |
| DoH | `https://<host><path>` | `https://dns.google/dns-query` |
| DoQ | `quic://<host>:<port>` | `quic://dns.adguard-dns.com:853` |

### Output Format

| Status | Format |
|--------|--------|
| Success | `✅ <server> - <protocol> - <time>ms - TTL: <ttl>s - <answers>` |
| Warning | `⚠️ <server> - <protocol> - <time>ms` (slow or NXDOMAIN) |
| Error | `❌ <server> - connection issue` |

---

## `dnstestergo server` - API Server

Start the DNS-Tester API server with optional in-memory or distributed workers.

### Usage

```bash
dnstestergo server [flags]
```

### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-c, --config` | string | - | Path to config file |
| `-H, --host` | string | `0.0.0.0` | Server bind address |
| `-P, --port` | string | `5000` | Server port |
| `-r, --redis` | string | - | Redis URL (enables distributed workers) |
| `-w, --workers` | int | config/`4` | Max number of in-memory workers |
| `--dns-timeout` | int | config/`5` | DNS query timeout in seconds |
| `--max-servers` | int | config/`50` | Max DNS servers per request |
| `--max-concurrent` | int | config/`500` | Max concurrent DNS queries |
| `--max-retries` | int | config/`3` | Number of retries per query |
| `--rate-limit-rps` | int | config/`10` | Requests per second (`0` = disable) |
| `--rate-limit-burst` | int | config/`20` | Rate limit burst size |
| `--read-timeout` | int | config/`15` | HTTP read timeout in seconds |
| `--write-timeout` | int | config/`15` | HTTP write timeout in seconds |
| `--idle-timeout` | int | config/`60` | HTTP idle timeout in seconds |

### Examples

```bash
# Start with default settings (in-memory workers)
dnstestergo server

# Start with custom host/port
dnstestergo server --host 0.0.0.0 --port 8080

# Start with Redis backend (distributed workers)
dnstestergo server --redis redis://localhost:6379/0

# Start with custom config
dnstestergo server --config /path/to/config.yaml

# Override DNS settings
dnstestergo server --dns-timeout 10 --max-retries 5

# Disable rate limiting
dnstestergo server --rate-limit-rps 0

# Allow 200 servers per request
dnstestergo server --max-servers 200

# Combine config file with CLI overrides
dnstestergo server --config prod.yaml --port 8080 --max-servers 100

# Docker
docker run -p 5000:5000 sudo-tiz/dnstestergo:latest server
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DNS_TESTER_HOST` | `0.0.0.0` | API bind address |
| `DNS_TESTER_PORT` | `5000` | API port |
| `MAX_WORKERS` | `4` | Worker concurrency (in-memory mode) |
| `REDIS_URL` | - | Redis connection URL |
| `RATE_LIMIT_IP_SOURCE` | `RemoteAddr` | IP source for rate limiting |

**Priority**: CLI flags > Environment variables > Config file > Defaults

---

## `dnstestergo worker` - Standalone Worker

Start a standalone worker that processes tasks from Redis queue.

### Usage

```bash
dnstestergo worker [flags]
```

### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-r, --redis` | string | **required** | Redis URL (e.g., `redis://localhost:6379/0`) |
| `-n, --concurrency` | int | `4` | Number of concurrent workers |
| `-c, --config` | string | - | Path to config file |
| `-M, --enable-metrics` | bool | `false` | Enable Prometheus metrics endpoint |
| `-m, --metrics-port` | int | `9091` | Metrics port (if enabled) |
| `--dns-timeout` | int | config/`5` | DNS query timeout in seconds |
| `--max-concurrent` | int | config/`500` | Max concurrent DNS queries |
| `--max-retries` | int | config/`3` | Number of retries per query |

### Examples

```bash
# Start worker with default concurrency
dnstestergo worker --redis redis://localhost:6379/0

# Start worker with custom concurrency
dnstestergo worker --redis redis://localhost:6379/0 --concurrency 8

# Start worker with metrics (useful for single worker)
dnstestergo worker \
  --redis redis://localhost:6379/0 \
  --enable-metrics \
  --metrics-port 9091

# Start worker with custom config
dnstestergo worker \
  --config /path/to/config.yaml \
  --redis redis://localhost:6379/0

# Override DNS settings
dnstestergo worker \
  --redis redis://localhost:6379/0 \
  --dns-timeout 10 \
  --max-retries 5 \
  --max-concurrent 1000

# Combine config file with CLI overrides
dnstestergo worker \
  --config prod.yaml \
  --redis redis://localhost:6379/0 \
  --concurrency 16 \
  --dns-timeout 15

# Docker
docker run sudo-tiz/dnstestergo:latest worker \
  --redis redis://host.docker.internal:6379/0
```

### Notes

- **Redis is required** for standalone workers
- Use `--enable-metrics` carefully to avoid port conflicts when running multiple workers
- Workers automatically register with Redis and process tasks from the queue
- CLI flags override config file settings
- See [Configuration](05-configuration.md) for DNS server settings

---

## Docker Usage

### Query Tool
```bash
docker run --rm sudo-tiz/dnstestergo:latest \
  query example.com udp://8.8.8.8:53
```

### Server + Workers (docker-compose)
```bash
docker compose up -d
```

### With Custom Config
```bash
docker run --rm \
  -v $(pwd)/conf:/conf \
  sudo-tiz/dnstestergo:latest \
  query example.com -c /conf/config.yaml
```

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Invalid server format | Use `protocol://host:port` format |
| Connection refused | Ensure API server is running: `docker compose up -d` |
| Task failed | Check worker logs: `docker compose logs dnstestergo-worker` |
| TLS certificate error | Use `--insecure` flag (testing only) |
| Redis connection error | Verify Redis URL and that Redis is running |
| Port already in use | Change port with `--port` or `DNS_TESTER_PORT` |

---

## See Also

- [API Reference](03-api.md) - REST API documentation
- [Configuration](05-configuration.md) - Config file format and DNS server settings
- [Quick Start](01-quickstart.md) - Deployment guide
