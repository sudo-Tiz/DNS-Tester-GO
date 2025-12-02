# Contributing

## Development Setup

```bash
git clone https://github.com/sudo-Tiz/dns-tester-go.git
cd dns-tester-go

# Install dependencies
go mod download

# Run tests
make test

# Build
make build
```

## Project Structure

```
dns-tester-go/
├── cmd/              # Entry points (api, worker, cli)
├── internal/         # Core logic
│   ├── api/         # HTTP handlers
│   ├── cli/         # CLI commands
│   ├── resolver/    # DNS resolution (dnsproxy wrapper)
│   ├── tasks/       # Asynq task handlers
│   ├── models/      # Data models
│   └── config/      # Configuration
├── docs/            # Documentation (Markdown)
└── website/         # Docusaurus site
```

## Code Style

**Go:**
- Use `gofmt` and `goimports`
- Follow [Effective Go](https://go.dev/doc/effective_go)
- Max line length: 120 chars

**Commit messages:**
```
type(scope): subject

body (optional)
```

**Types:** `feat`, `fix`, `docs`, `refactor`, `test`, `chore`

**Example:**
```
feat(api): add DoQ protocol support

Implements DNS-over-QUIC using AdGuard dnsproxy.
Closes #123
```

## Testing

**Run all tests:**
```bash
make test
```

**Run specific test:**
```bash
go test ./internal/resolver -v
```

**Coverage:**
```bash
make coverage
```

**Integration tests:**
```bash
make test-integration
```

## Pull Request Process

1. **Fork & branch:**
   ```bash
   git checkout -b feat/my-feature
   ```

2. **Make changes:**
   - Write tests
   - Update docs if needed
   - Follow code style

3. **Commit:**
   ```bash
   git commit -m "feat(api): add new endpoint"
   ```

4. **Push & PR:**
   ```bash
   git push origin feat/my-feature
   ```
   Open PR with description

5. **CI checks:**
   - Tests pass
   - Linting passes
   - No decrease in coverage

## Documentation

**Edit docs:**
```bash
cd docs/
# Edit .md files
```

**Preview website:**
```bash
cd website/
npm install
npm start
```

**Sync docs to website:**
```bash
cd website/
./sync-docs.sh
```

## Release Process

**Maintainers only:**

1. Update `VERSION` file
2. Update `CHANGELOG.md`
3. Tag release:
   ```bash
   git tag -a v1.2.3 -m "Release v1.2.3"
   git push origin v1.2.3
   ```
4. GitHub Actions builds & publishes

## Getting Help

| Channel | Link |
|---------|------|
| Issues | [GitHub Issues](https://github.com/sudo-Tiz/dns-tester-go/issues) |
| Discussions | [GitHub Discussions](https://github.com/sudo-Tiz/dns-tester-go/discussions) |
| Docs | [Documentation](https://sudo-Tiz.github.io/dns-tester-go) |

## Code of Conduct

Be respectful and constructive. Follow [Go Community Code of Conduct](https://go.dev/conduct).
