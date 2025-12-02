# Examples

Practical scenarios and use cases.

---

## 📋 Common Scenarios

| Use Case | Command | Purpose |
|----------|---------|---------|
| **Outage Response** | `dnstestergo query app.example.com udp://8.8.8.8:53 udp://internal.dns:53` | Compare public vs internal DNS |
| **Cache Validation** | `dnstestergo query cdn.example.com udp://cache1:53 udp://cache2:53` | Verify cache consistency |
| **Performance Test** | `dnstestergo query example.com udp://8.8.8.8:53 udp://1.1.1.1:53 --warn-threshold 0.05` | Compare server latency |
| **Protocol Comparison** | `dnstestergo query github.com udp://1.1.1.1:53 tls://one.one.one.one:853` | UDP vs DoT performance |
| **Migration Validation** | `dnstestergo query app.com udp://old-dns:53 udp://new-dns:53` | Before/after comparison |
| **Health Check** | `dnstestergo query critical.com udp://127.0.0.1:53 \|\| exit 1` | Load balancer endpoint |
| **Propagation Check** | Loop through NS servers with `dnstestergo query` | Verify DNS propagation |

---

## 🚨 Incident Response

### Quick Multi-Domain Check
```bash
for domain in api.example.com db.example.com auth.example.com; do
  dnstestergo query "$domain" udp://8.8.8.8:53 || echo "❌ $domain FAILED"
done
```

### Compare Internal vs External
```bash
dnstestergo query api.company.com \
  udp://8.8.8.8:53 \
  udp://10.0.0.1:53 \
  udp://10.0.0.2:53
```

**Analysis:**
- ✅ Public resolves, internal fails → Internal DNS issue
- ❌ All fail → Domain configuration problem  
- ⚠️ Different IPs → Cache inconsistency

---

## 🔄 Cache Testing

### Detect Stale Cache
```bash
# Compare authoritative vs cache
dnstestergo query example.com \
  udp://ns1.example.com:53 \
  udp://cache.corp:53 \
  --debug
```

**Check:** TTL values, IP addresses, response codes

---

## ⚡ Performance Testing

### Server Comparison
```bash
dnstestergo query example.com \
  udp://8.8.8.8:53 \
  udp://1.1.1.1:53 \
  udp://9.9.9.9:53 \
  --warn-threshold 0.05
```

### Protocol Benchmark
```bash
dnstestergo query github.com \
  udp://1.1.1.1:53 \
  tcp://1.1.1.1:53 \
  tls://one.one.one.one:853 \
  https://cloudflare-dns.com/dns-query
```

**Typical latency:**
- UDP: 10-20ms
- TCP: 15-30ms  
- DoT: 50-100ms
- DoH: 80-150ms

---

## 🔐 Multi-Protocol Testing

| Protocol | Example |
|----------|---------|
| **DoT** | `dnstestergo query example.com tls://dns.google:853` |
| **DoH** | `dnstestergo query example.com https://dns.google/dns-query` |
| **DoQ** | `dnstestergo query example.com quic://dns.adguard-dns.com:853` |

### Protocol Fallback Script
```bash
for proto in https://dns.google/dns-query tls://dns.google:853 udp://8.8.8.8:53; do
  dnstestergo query example.com "$proto" 2>/dev/null && break
done
```

---

## 📊 Monitoring & Alerting

### Continuous Monitoring (Cron)
```bash
#!/bin/bash
# Run every minute: * * * * * /opt/monitor-dns.sh

for server in udp://10.0.0.1:53 udp://10.0.0.2:53; do
  dnstestergo query api.example.com "$server" --warn-threshold 0.5 || \
    echo "ALERT: DNS failure on $server" | mail -s "DNS Alert" ops@example.com
done
```

### Health Check Endpoint
```bash
#!/bin/bash
# For load balancer health checks

for domain in api.example.com db.example.com; do
  dnstestergo query "$domain" udp://127.0.0.1:53 || exit 1
done
echo "OK"
```

---

## 🔄 DNS Migration

### Before/After Comparison
```bash
DOMAIN="example.com"
OLD_DNS="udp://10.0.0.1:53"
NEW_DNS="udp://10.1.0.1:53"

echo "=== Before ==="
dnstestergo query "$DOMAIN" "$OLD_DNS"

echo "=== After ==="
dnstestergo query "$DOMAIN" "$NEW_DNS"

echo "=== Both ==="
dnstestergo query "$DOMAIN" "$OLD_DNS" "$NEW_DNS"
```

---

## 🔌 API Integration

### CI/CD Pipeline
```bash
#!/bin/bash
# Submit query via API
RESPONSE=$(curl -s -X POST http://dns-api:5000/dns-lookup \
  -H "Content-Type: application/json" \
  -d '{"domain":"staging.example.com","dns_servers":[{"target":"udp://10.0.0.1:53"}],"qtype":"A"}')

TASK_ID=$(echo "$RESPONSE" | jq -r .task_id)

# Poll for results
for i in {1..30}; do
  RESULT=$(curl -s "http://dns-api:5000/tasks/$TASK_ID")
  STATUS=$(echo "$RESULT" | jq -r .task_status)
  
  [ "$STATUS" = "SUCCESS" ] && echo "✅ Pass" && exit 0
  [ "$STATUS" = "FAILURE" ] && echo "❌ Fail" && exit 1
  
  sleep 1
done

echo "⏱️ Timeout"
exit 1
```

### Python Integration
```python
import requests, time

def check_dns(domain, servers):
    # Submit query
    r = requests.post('http://localhost:5000/dns-lookup', json={
        'domain': domain,
        'dns_servers': [{'target': s} for s in servers],
        'qtype': 'A'
    })
    task_id = r.json()['task_id']
    
    # Poll results
    for _ in range(30):
        result = requests.get(f'http://localhost:5000/tasks/{task_id}').json()
        if result['task_status'] == 'SUCCESS':
            return result['task_result']
        time.sleep(0.5)
    
    raise TimeoutError('DNS query timeout')

# Usage
results = check_dns('example.com', ['udp://8.8.8.8:53'])
```

### Slack Alerts
```bash
#!/bin/bash
SLACK_WEBHOOK="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"

dnstestergo query critical.example.com udp://8.8.8.8:53 || \
  curl -X POST "$SLACK_WEBHOOK" \
    -H "Content-Type: application/json" \
    -d '{"text":":warning: DNS Alert: critical.example.com failed"}'
```

---

## 🔍 Advanced Queries

### DNS Propagation Check
```bash
#!/bin/bash
DOMAIN="newsite.example.com"
NS_SERVERS=$(dig +short "$DOMAIN" NS | sed 's/\.$//')

for ns in $NS_SERVERS; do
  NS_IP=$(dig +short "$ns" A | head -1)
  echo "=== $ns ($NS_IP) ==="
  dnstestergo query "$DOMAIN" "udp://$NS_IP:53"
done
```

### Multiple Record Types
```bash
for type in A AAAA MX TXT NS SOA SRV CAA; do
  echo "=== $type ==="
  dnstestergo query example.com udp://8.8.8.8:53 --qtype "$type"
done
```

### Reverse DNS Bulk Check
```bash
for ip in 8.8.8.8 1.1.1.1 9.9.9.9; do
  echo "=== $ip ==="
  dnstestergo query "$ip" udp://8.8.8.8:53
done
```

---

## 🐳 Docker Compose Examples

### High Availability Setup
```yaml
version: '3.8'

services:
  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes

  api:
    image: sudo-tiz/dnstestergo:latest
    command: ["server", "--config", "/app/config.yaml"]
    ports:
      - "5000:5000"
    volumes:
      - ./conf/config.yaml:/app/config.yaml:ro
    environment:
      REDIS_URL: redis://redis:6379
    deploy:
      replicas: 3
    depends_on:
      - redis

  worker:
    image: sudo-tiz/dnstestergo:latest
    command: ["worker", "--config", "/app/config.yaml"]
    volumes:
      - ./conf/config.yaml:/app/config.yaml:ro
    environment:
      REDIS_URL: redis://redis:6379
    deploy:
      replicas: 4
    depends_on:
      - redis

volumes:
  redis_data:
```

### Nginx Load Balancer
```nginx
upstream dns_tester_api {
    least_conn;
    server api-1:5000;
    server api-2:5000;
    server api-3:5000;
}

server {
    listen 443 ssl http2;
    server_name dns-api.example.com;

    location / {
        proxy_pass http://dns_tester_api;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /metrics {
        deny all;  # Internal only
    }
}
```

---

## ☸️ Kubernetes Examples

### Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dns-tester-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: dns-tester-api
  template:
    metadata:
      labels:
        app: dns-tester-api
    spec:
      containers:
      - name: api
        image: sudo-tiz/dnstestergo:latest
        args: ["server", "--config", "/config/config.yaml"]
        ports:
        - containerPort: 5000
        env:
        - name: REDIS_URL
          value: redis://redis:6379
---
apiVersion: v1
kind: Service
metadata:
  name: dns-tester-api
spec:
  selector:
    app: dns-tester-api
  ports:
  - port: 5000
```

### CronJob Monitoring
```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: dns-monitor
spec:
  schedule: "*/5 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: monitor
            image: sudo-tiz/dnstestergo:latest
            command:
            - /bin/sh
            - -c
            - dnstestergo query critical.example.com udp://10.0.0.1:53 || echo "ALERT"
          restartPolicy: OnFailure
```

---

## 📚 See Also

- [CLI Guide](04-cli.md) - Command reference
- [API Reference](03-api.md) - API integration
- [Configuration](05-configuration.md) - Config options
- [Monitoring](06-monitoring.md) - Metrics setup
