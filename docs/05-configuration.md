# Configuration

YAML configuration reference for dns-tester-go.

---

## 📂 Config File Location

**Search order** (highest priority first):
1. `--config` flag → `dnstestergo server --config /path/to/config.yaml`
2. `./config.yaml` (current directory)
3. `conf/config.yaml` (default)

---

## 🔄 Configuration Precedence

**Priority** (highest → lowest):
1. **CLI flags** (`--host`, `--port`, `--dns-timeout`, `--max-servers`, etc.)
2. **Environment variables** (`DNS_TESTER_HOST`, `DNS_TESTER_PORT`)
3. **Config file** (YAML)
4. **Built-in defaults**

---

## ⚡ Quick Start Example

```yaml
servers:
  - ip: "8.8.8.8"
    hostname: "dns.google"
    services: ["do53/udp", "dot", "doh"]
    tags: ["GOOGLE"]
  - ip: "1.1.1.1"
    hostname: "one.one.one.one"
    services: ["do53/udp", "dot", "doh"]
    tags: ["CLOUDFLARE"]

rate_limiting:
  requests_per_second: 10
  burst_size: 20

server:
  host: "0.0.0.0"
  port: "5000"

worker:
  max_workers: 4
  cleanup_interval: 10

dns:
  timeout: 5
  max_servers_per_req: 50
```

---

## 📋 Configuration Reference

### Servers (Required)

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `ip` | string | ✅* | - | IPv4/IPv6 address |
| `hostname` | string | ✅* | - | Hostname for TLS |
| `port` | int | ❌ | Protocol default | Custom port |
| `services` | array | ✅ | - | Protocol list |
| `tags` | array | ❌ | `[]` | Identification tags |

**\* Required:** `ip` for UDP/TCP | `hostname` for DoT/DoH/DoQ

**Services:**

| Service | Protocol | Default Port | Requires |
|---------|----------|--------------|----------|
| `do53/udp` | UDP | 53 | `ip` |
| `do53/tcp` | TCP | 53 | `ip` |
| `dot` | TLS | 853 | `hostname` |
| `doh` | HTTPS | 443 | `hostname` |
| `doq` | QUIC | 853 | `hostname` |

### Rate Limiting (Optional)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `requests_per_second` | int | `10` | Max req/s per IP |
| `burst_size` | int | `20` | Burst capacity |

### Server (Optional)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `host` | string | `"0.0.0.0"` | Listen address |
| `port` | string | `"5000"` | Listen port |

### Worker (Optional)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `max_workers` | int | `4` | Concurrent workers |
| `cleanup_interval` | int | `10` | Task cleanup (minutes) |

### DNS (Optional)

Controls DNS query behavior and limits.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `timeout` | int | `5` | Query timeout in seconds |
| `max_servers_per_req` | int | `50` | Max DNS servers per API request |
| `max_concurrent_queries` | int | `500` | Max servers queried in parallel (per request) |
| `max_retries` | int | `3` | Number of retry attempts per query |

**Notes:**
- `max_servers_per_req`: Limits total number of servers a client can request
- `max_concurrent_queries`: Controls internal parallelism (goroutines) when querying multiple servers
- `max_retries`: Applied per server, not globally

**Example:**
```yaml
dns:
  timeout: 10                  # Wait max 10s per query
  max_servers_per_req: 100     # Allow 100 servers per request
  max_concurrent_queries: 50   # Query 50 servers in parallel
  max_retries: 5               # Retry failed queries 5 times
```

---

## 🌐 Public DNS Servers

| Provider | IP | Hostname | Services |
|----------|----|----|----------|
| Google | `8.8.8.8` | `dns.google` | UDP, TCP, DoT, DoH |
| Cloudflare | `1.1.1.1` | `one.one.one.one` | UDP, TCP, DoT, DoH |
| Quad9 | `9.9.9.9` | `dns.quad9.net` | UDP, TCP, DoT, DoH |
| AdGuard | `94.140.14.14` | `dns.adguard-dns.com` | UDP, TCP, DoT, DoH, DoQ |

> ⚠️ **Note**: Some providers block ads/malware domains. Test before production use.

---

## 🔧 Configuration Examples

### Custom Port
```yaml
- ip: "192.168.1.1"
  port: 5353
  services: ["do53/udp"]
```

### DoH Only
```yaml
- hostname: "doh.opendns.com"
  services: ["doh"]
```

### IPv6
```yaml
- ip: "2001:4860:4860::8888"
  hostname: "dns.google"
  services: ["do53/udp", "dot"]
```

**CLI format:**
- IPv4: `udp://8.8.8.8:53`
- IPv6: `udp://[2001:4860:4860::8888]:53`
- TLS: `tls://dns.google:853`

### Tags and Organization
```yaml
servers:
  # Public
  - ip: "8.8.8.8"
    hostname: "dns.google"
    services: ["do53/udp", "dot"]
    tags: ["PUBLIC", "GOOGLE", "PRIMARY"]
  
  # Internal
  - ip: "10.0.0.1"
    hostname: "dns.internal.corp"
    services: ["do53/udp", "do53/tcp"]
    tags: ["INTERNAL", "CORP", "PRIMARY"]
  
  # Test
  - ip: "127.0.0.1"
    port: 5353
    services: ["do53/udp"]
    tags: ["TEST", "LOCAL"]
```

---

## 🔐 Environment Variables

| Variable | Type | Default | Overrides | Description |
|----------|------|---------|-----------|-------------|
| `DNS_TESTER_HOST` | string | `0.0.0.0` | `server.host` | API bind address |
| `DNS_TESTER_PORT` | string | `5000` | `server.port` | API bind port |
| `MAX_WORKERS` | int | `4` | `worker.max_workers` | Worker pool size |
| `REDIS_URL` | string | - | - | Redis backend (e.g., `redis://localhost:6379/0`) |
| `RATE_LIMIT_IP_SOURCE` | string | `RemoteAddr` | - | IP source for rate limiting |

**Rate Limit IP Source** (for proxies/load balancers):
- `RemoteAddr` (default) - Direct connection IP
- `X-Forwarded-For` - X-Forwarded-For header
- `X-Real-IP` - X-Real-IP header

**Example (behind nginx):**
```bash
export RATE_LIMIT_IP_SOURCE=X-Real-IP
dnstestergo server
```

---

## ✅ Validation

### Test Configuration
```bash
dnstestergo server --config conf/config.yaml
# Check logs for errors
```

### Common Errors

| Error | Fix |
|-------|-----|
| `do53/udp requires an IP address` | Add `ip` field for UDP/TCP services |
| `invalid IP address` | Use valid IPv4/IPv6 format |
| `services must not be empty` | Add at least one service |
| `no servers configured` | Add `servers:` section with at least one entry |
| `config file not found` | Check path: `ls -la conf/config.yaml` |

**Example fix:**
```yaml
# ❌ Wrong
servers:
  - hostname: "dns.google"
    services: ["do53/udp"]

# ✅ Correct
servers:
  - ip: "8.8.8.8"
    hostname: "dns.google"
    services: ["do53/udp"]
```

---

## 🐳 Docker Configuration

**Mount config file:**
```yaml
services:
  api:
    volumes:
      - ./conf/config.yaml:/app/config.yaml
    command: ["--config", "/app/config.yaml"]
```

---

## 💡 Best Practices

✅ **Do:**
- Use tags to organize servers
- Configure multiple providers for redundancy
- Use hostname for DoT/DoH (proper TLS validation)
- Set realistic timeouts based on network
- Test configuration before deployment

❌ **Don't:**
- Use `do53/udp` with only hostname (requires IP)
- Set rate limits too low in production
- Use default cleanup_interval for high-volume systems

---

## 📚 See Also

- [API Reference](03-api.md) - REST API
- [CLI Guide](04-cli.md) - Command-line usage
- [Monitoring](06-monitoring.md) - Metrics setup
