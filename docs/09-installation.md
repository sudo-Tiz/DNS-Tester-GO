# Installation

Choose your installation method.

---

## 📦 Installation Methods

| Method | Use Case | Command |
|--------|----------|---------|
| **Docker** | Production, testing | `docker compose --profile prod up -d` |
| **Go Install** | CLI or Development | `go install github.com/sudo-Tiz/dns-tester-go/cmd/dnstestergo@latest` |
| **Binary** | No Go/Docker | [Download from Releases](https://github.com/sudo-Tiz/dns-tester-go/releases) |
| **Source** | Development | `make build` |

---

## 🐳 Docker (Recommended)

```bash
git clone https://github.com/sudo-Tiz/dns-tester-go.git
cd dns-tester-go

# Production setup (API + Worker + Redis)
cp conf/config.example.yaml conf/config.yaml  # Create config file
# Edit conf/config.yaml with your DNS servers
docker compose --profile prod up -d

# OR Development (all-in-one, no Redis, uses config.example.yaml)
docker compose --profile dev up -d
```

---

## 🔧 Go Install

```bash
go install github.com/sudo-Tiz/dns-tester-go/cmd/dnstestergo@latest
```

---

## 🛠️ Build from Source

```bash
git clone https://github.com/sudo-Tiz/dns-tester-go.git
cd dns-tester-go
make build
```

**Binaries**: `bin/dnstestergo-query`, `bin/dnstestergo-server`, `bin/dnstestergo-worker`

---

## 🐛 Troubleshooting

[See Troubleshooting](07-troubleshooting.md)

---

## 📚 Next

- [Quick Start](01-quickstart.md) - Test the system
- [CLI Guide](04-cli.md) - All commands and flags
- [Configuration](05-configuration.md) - Configure DNS servers
