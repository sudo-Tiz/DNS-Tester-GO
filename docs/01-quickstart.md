# Quick Start

Deploy and test dns-tester-go in 2 minutes.

> 📦 **Need installation help?** See [Installation Guide](09-installation.md)

---

## 🚀 Start Services

```bash
# First time: create config file
cp conf/config.example.yaml conf/config.yaml

# Start production stack
docker compose --profile prod up -d
```


---

## 🧪 Test

**CLI:**
```bash
dnstestergo query example.com udp://8.8.8.8:53
```

**API:**
```bash
# Submit
curl -X POST http://localhost:5000/dns-lookup \
  -H "Content-Type: application/json" \
  -d '{"domain":"example.com","qtype":"A"}'

# Get result (use task_id from response)
curl http://localhost:5000/tasks/TASK_ID
```

---

## 📚 Next

- [CLI Guide](04-cli.md) - All commands
- [API Reference](03-api.md) - REST API
- [Configuration](05-configuration.md) - DNS servers
- [Monitoring](06-monitoring.md) - Prometheus metrics
- [Architecture](02-architecture.md) - System design
