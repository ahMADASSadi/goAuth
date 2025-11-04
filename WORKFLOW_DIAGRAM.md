# Release Pipeline Workflow Diagram

## Overview

```
Developer                GitHub Actions              Docker Hub
    |                           |                         |
    | 1. git tag -a v1.0.0      |                         |
    |-------------------------->|                         |
    |                           |                         |
    |                    2. Checkout dev                  |
    |                           |                         |
    |                    3. Merge dev → main             |
    |                           |                         |
    |                    4. Build binaries               |
    |                      (5 platforms)                 |
    |                           |                         |
    |                    5. Create GitHub Release        |
    |                      + Upload binaries             |
    |                           |                         |
    |                    6. Build Docker image           |
    |                      (amd64 + arm64)               |
    |                           |                         |
    |                    7. Push to Docker Hub --------->|
    |                           |                         |
    |<------- 8. Release complete notification           |
    |                           |                         |
```

## Detailed Step-by-Step

### 1. Tag Creation (Developer)
```bash
git tag -a v1.0.0 -m "Release message"
git push origin v1.0.0
```

### 2. Workflow Triggered (GitHub Actions)
- Event: Push of tag matching `v*.*.*`
- Runner: Ubuntu latest

### 3. Merge to Main
- Checkout `dev` branch with full history
- Configure git user as `github-actions[bot]`
- Fetch and merge `dev` into `main`
- Push merged `main` branch

### 4. Build Binaries
Platforms built:
- `linux/amd64`
- `linux/arm64`
- `darwin/amd64` (macOS Intel)
- `darwin/arm64` (macOS Apple Silicon)
- `windows/amd64`

Build flags: `-ldflags="-s -w"` (strip debug info)

### 5. Create GitHub Release
- Generate changelog from git commits
- Create release on GitHub
- Upload all binaries + checksums
- Mark as pre-release if tag contains alpha/beta/rc

### 6. Build Docker Image
- Multi-stage build using Dockerfile
- Platforms: `linux/amd64`, `linux/arm64`
- Build cache enabled (GitHub Actions cache)

### 7. Push to Docker Hub
Image tags created:
- `<username>/goauth-backend-repo:v1.0.0` (specific version)
- `<username>/goauth-backend-repo:1.0` (minor version)
- `<username>/goauth-backend-repo:1` (major version)
- `<username>/goauth-backend-repo:latest` (always latest)

### 8. Summary
- Workflow summary posted to Actions tab
- Release available in GitHub Releases
- Docker image available on Docker Hub

## Required Secrets

| Secret | Purpose | How to Get |
|--------|---------|------------|
| `DOCKERHUB_USERNAME` | Docker Hub login | Your Docker Hub username |
| `DOCKERHUB_TOKEN` | Docker Hub auth | Generate at hub.docker.com → Security |
| `GITHUB_TOKEN` | GitHub API | Automatically provided by GitHub |

## Trigger Conditions

✅ **Will trigger:**
- `v1.0.0`
- `v2.5.3`
- `v0.1.0-alpha`
- `v1.0.0-beta.1`

❌ **Will NOT trigger:**
- `1.0.0` (missing 'v' prefix)
- `release-1.0.0` (wrong format)
- Push to dev/main branches
- Pull requests

## Timeline

| Step | Estimated Time |
|------|---------------|
| Merge dev → main | 30 seconds |
| Build binaries | 2-3 minutes |
| Create GitHub Release | 30 seconds |
| Build Docker images | 3-5 minutes |
| Push to Docker Hub | 1-2 minutes |
| **Total** | **~7-11 minutes** |

## Outputs

### GitHub Release
- URL: `https://github.com/ahMADASSadi/goAuth/releases/tag/v1.0.0`
- Contains:
  - Release notes (auto-generated)
  - Binaries for 5 platforms
  - `checksums.txt` file

### Docker Hub
- URL: `https://hub.docker.com/r/<username>/goauth-backend-repo`
- Tags: `latest`, `v1.0.0`, `1.0`, `1`
- Platforms: linux/amd64, linux/arm64

## Rollback Procedure

If you need to rollback a release:

1. **Delete the tag:**
   ```bash
   git tag -d v1.0.0
   git push --delete origin v1.0.0
   ```

2. **Delete GitHub Release manually:**
   - Go to Releases → Click on release → Delete release

3. **Delete Docker images:**
   - Go to Docker Hub → Tags → Delete tag

4. **Revert main branch if needed:**
   ```bash
   git checkout main
   git reset --hard HEAD~1  # Go back one commit
   git push origin main --force
   ```

## Success Indicators

✅ **Release successful when:**
- GitHub Actions workflow shows green checkmark
- GitHub Release exists with all binaries
- Docker Hub shows new tag
- `main` branch contains merge commit from `dev`

## Troubleshooting

### Workflow fails at merge step
**Cause:** Merge conflicts between dev and main
**Fix:** Manually merge and resolve conflicts, then re-push tag

### Docker push fails
**Cause:** Invalid Docker Hub credentials
**Fix:** Regenerate Docker Hub token, update GitHub secret

### Binaries missing
**Cause:** Build failure on specific platform
**Fix:** Check Actions logs for compilation errors

---

**Flow is fully automated after tag push!** 🎯
