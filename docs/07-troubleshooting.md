# Troubleshooting

Quick solutions to common issues.

---

## 🔍 Quick Diagnostics

```bash
# Check services
docker compose ps

# View logs
docker compose logs -f

# Restart service
docker compose restart <service>

# Debug CLI
dnstestergo query example.com udp://8.8.8.8:53 --debug
```

---

## 🔴 Common Errors

| Error | Cause | Solution |
|-------|-------|----------|
| `connection refused` | API not running | `docker compose --profile prod up -d` |
| `503 Service Unavailable` | No workers | `docker compose restart dnstestergo-worker` |
| `429 Too Many Requests` | Rate limit hit | Increase `rate_limiting.requests_per_second` in config |
| `400 Bad Request: domain required` | Invalid JSON | Check request format with `jq` |
| `invalid server address format` | Missing protocol | Use `udp://8.8.8.8:53` not `8.8.8.8` |
| `task failed` | Worker error | `docker compose logs dnstestergo-worker` |
| `TLS verify failed` | Invalid cert | Use `--insecure` (test only) or fix certificate |
| `NXDOMAIN` | Domain doesn't exist | Verify domain with `whois` |
| `i/o timeout` | Network/firewall | Test with `nc -zvu 8.8.8.8 53` |
| `no workers available` | Worker not running | Check `docker compose ps dnstestergo-worker` |

---

## 🔧 Configuration Issues

| Problem | Fix |
|---------|-----|
| `do53/udp requires an IP address` | Add `ip` field for UDP/TCP services |
| `invalid IP address` | Use valid IPv4/IPv6 (e.g., `8.8.8.8`, `2001:4860:4860::8888`) |
| `failed to parse YAML` | Check indentation (2 spaces, no tabs). Use `yamllint` |
| `no servers configured` | Verify `servers:` section exists in config.yaml |

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
    services: ["do53/udp", "dot"]
```

---

## 🐳 Docker Issues

| Problem | Solution |
|---------|----------|
| Container exits immediately | `docker compose logs <service>` |
| Port already in use | Change port in `docker-compose.yml` or stop conflicting service |
| Redis connection failed | `docker compose restart redis` + check REDIS_URL |
| Permission denied (volumes) | `chmod 644 conf/config.yaml` |
| Worker not processing | `docker compose restart dnstestergo-worker` |

---

## 🌐 Network Issues

| Problem | Test | Solution |
|---------|------|----------|
| DoT/DoQ blocked (port 853) | `nc -zv dns.google 853` | Use DoH (port 443) instead |
| DNS timeout | `ping 8.8.8.8` | Try faster server or TCP |
| Firewall blocking UDP/53 | `nc -zvu 8.8.8.8 53` | Check firewall rules |
| High latency | Compare protocols with `--warn-threshold 0.1` | Use UDP for best performance |
| DoQ slow/unstable | Check UDP buffer sizes | Increase system UDP buffers (see below) |

### DoQ Performance: UDP Buffer Sizes

DNS-over-QUIC (DoQ) requires adequate UDP buffer sizes for optimal performance. If you experience slow queries or packet loss:

**Check current buffer sizes:**
```bash
# Linux
sysctl net.core.rmem_max
sysctl net.core.wmem_max

# macOS
sysctl kern.ipc.maxsockbuf
```

**Increase buffer sizes (Linux):**
```bash
# Temporary (until reboot)
sudo sysctl -w net.core.rmem_max=7500000
sudo sysctl -w net.core.wmem_max=7500000

# Permanent
echo "net.core.rmem_max=7500000" | sudo tee -a /etc/sysctl.conf
echo "net.core.wmem_max=7500000" | sudo tee -a /etc/sysctl.conf
sudo sysctl -p
```

**Increase buffer sizes (macOS):**
```bash
sudo sysctl -w kern.ipc.maxsockbuf=3014657
```

For more details, see [quic-go UDP Buffer Sizes guide](https://github.com/quic-go/quic-go/wiki/UDP-Buffer-Sizes).

---

## 📊 Monitoring Issues

| Problem | Check | Fix |
|---------|-------|-----|
| Metrics not updating | `curl localhost:5000/metrics` | Perform test query |
| Grafana "No data" | Prometheus targets: http://localhost:9090/targets | Verify data source config |
| Worker metrics missing | Worker started with `--enable-metrics`? | Add flag or check port 9091 |

---

## 🐛 Debug Checklist

Before reporting issues, collect:

```bash
# 1. Version
dnstestergo --version

# 2. Full error with debug
dnstestergo query example.com udp://8.8.8.8:53 --debug 2>&1 | tee error.log

# 3. Logs
docker compose logs > logs.txt

# 4. Config
cat conf/config.yaml

# 5. Services status
docker compose ps
```

---

## 🆘 Get Help

- **[GitHub Issues](https://github.com/sudo-Tiz/DNS-Tester-GO/issues)** - Bug reports
- **[GitHub Discussions](https://github.com/sudo-Tiz/DNS-Tester-GO/discussions)** - Questions

---

## 📚 See Also

- [Configuration](05-configuration.md) - Fix config errors
- [CLI Guide](04-cli.md) - Command reference
- [API Reference](03-api.md) - API documentation
