# Tag-Based Release Pipeline - Quick Start

## What Changed

✅ **Release workflow** (`.github/workflows/release.yaml`) now:
- Triggers **ONLY on tag creation** (e.g., `v1.0.0`)
- Merges `dev` → `main` automatically
- Creates GitHub release with binaries
- Builds multi-arch Docker images (amd64 & arm64)
- Pushes to **Docker Hub** as `<username>/goauth-backend-repo`

✅ **Build workflow disabled** (`build.yaml.disabled`)
- No automatic builds on push/PR
- Only manual tag-based releases

✅ **Dockerfile fixed**:
- Added `wget` for healthcheck
- Removed invalid `.env` copy from builder

✅ **docker-compose.yml** validated:
- Environment variables explicit
- Healthcheck aligned with Dockerfile

## Setup (One-Time)

### 1. Add Docker Hub Secrets to GitHub

**Required secrets:**
- `DOCKERHUB_USERNAME` - Your Docker Hub username
- `DOCKERHUB_TOKEN` - Docker Hub access token (not password!)

**Steps:**
1. Go to [Docker Hub](https://hub.docker.com/) → Account Settings → Security
2. Create new access token with **Read, Write, Delete** permissions
3. Copy the token
4. In GitHub: Settings → Secrets and variables → Actions → New repository secret
5. Add both secrets

### 2. Create Your First Release

```bash
# On dev branch, ensure everything is committed
git checkout dev
git status

# Create and push tag
git tag -a v1.0.0 -m "Release v1.0.0 - Initial release"
git push origin v1.0.0
```

### 3. Monitor Progress

- Go to **GitHub Actions** tab
- Watch "Release" workflow run (~5-10 minutes)
- Check GitHub Releases page for binaries
- Check Docker Hub for images

## Pull Your Docker Image

```bash
docker pull <your-dockerhub-username>/goauth-backend-repo:latest
docker pull <your-dockerhub-username>/goauth-backend-repo:v1.0.0
```

## Documentation

- **Full guide**: `RELEASE_GUIDE.md` - Complete instructions, troubleshooting, best practices
- **Quick reference**: `QUICK_REFERENCE.md` - Common commands
- **Setup summary**: `SETUP_SUMMARY.md` - Original CI/CD documentation

---

**Ready to release!** Just create a tag starting with `v` and push it. 🚀
