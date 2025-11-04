# Quick Reference - goAuth Development

## 📋 Common Commands

### Development
```bash
cd src

# Start with live reload
make watch

# Run tests
make test

# Run with coverage
make coverage

# Run linting
make lint

# Generate API docs
make doc
```

### Docker
```bash
cd src

# Build and run
make docker-build
make docker-run

# View logs
make docker-logs

# Stop everything
make docker-stop

# Rebuild from scratch
make docker-rebuild
```

### Git & Releases
```bash
# Create a release
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# Push to trigger CI
git push origin dev
```

## 🔗 Important URLs

- **Actions**: https://github.com/ahMADASSadi/goAuth/actions
- **Security**: https://github.com/ahMADASSadi/goAuth/security
- **Packages**: https://github.com/ahMADASSadi/goAuth/pkgs/container/goauth
- **Local Server**: http://localhost:8000

## 📦 Docker Images

```bash
# Pull from GitHub Container Registry
docker pull ghcr.io/ahmadassadi/goauth:latest
docker pull ghcr.io/ahmadassadi/goauth:dev
docker pull ghcr.io/ahmadassadi/goauth:v1.0.0
```

## 🔍 Troubleshooting Quick Fixes

```bash
# Fix go.mod
go mod tidy

# Clean and rebuild
make clean && make build

# Docker fresh start
docker-compose down -v && docker-compose up -d --build

# View full logs
docker-compose logs --tail=100 -f
```

## ✅ Pre-commit Checklist

- [ ] Run `make lint`
- [ ] Run `make test`
- [ ] Run `make doc` (if API changed)
- [ ] Update version in code (if releasing)
- [ ] Update CHANGELOG.md

## 🎯 CI/CD Status

Check these before merging:
- ✅ Lint Code
- ✅ Build and Test
- ✅ Security Scan
- ✅ Build Docker Image

## 📝 Environment Variables

Create `src/.env`:
```env
PORT=8000
HOST=0.0.0.0
DB_URL=./test.sqlite3
ACCESS_EXPIRY=24h
SECRET_KEY=your-secret-key-here
```
