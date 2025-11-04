# Release Pipeline - Setup & Usage

## Overview

The release pipeline **only runs when you create a version tag** (e.g., `v1.0.0`). It performs the following steps automatically:

1. **Merges `dev` into `main`** branch
2. **Creates a GitHub Release** with binaries for multiple platforms
3. **Builds a multi-architecture Docker image** (linux/amd64, linux/arm64)
4. **Pushes the image to Docker Hub** as `<username>/goauth-backend-repo`

## Required GitHub Secrets

You **must** add these secrets to your GitHub repository before creating a tag:

### Steps to Add Secrets:

1. Go to your GitHub repository
2. Click **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Add the following secrets:

| Secret Name | Description | Example Value |
|------------|-------------|---------------|
| `DOCKERHUB_USERNAME` | Your Docker Hub username | `ahmadassadi` |
| `DOCKERHUB_TOKEN` | Your Docker Hub access token (not password!) | `dckr_pat_...` |

### How to Get Your Docker Hub Token:

1. Log in to [Docker Hub](https://hub.docker.com/)
2. Click your username (top-right) → **Account Settings**
3. Go to **Security** → **New Access Token**
4. Give it a name (e.g., "GitHub Actions goAuth")
5. Select permissions: **Read, Write, Delete**
6. Click **Generate**
7. **Copy the token immediately** (you won't be able to see it again!)
8. Add it to GitHub Secrets as `DOCKERHUB_TOKEN`

## How to Trigger a Release

### Step 1: Ensure Your Changes Are in `dev` Branch

```bash
# Make sure you're on dev and everything is committed
git checkout dev
git status
```

### Step 2: Create and Push a Version Tag

```bash
# Create a version tag (must start with 'v')
git tag -a v1.0.0 -m "Release v1.0.0 - Initial release"

# Push the tag to GitHub (this triggers the release workflow)
git push origin v1.0.0
```

### Tag Naming Convention:

- **Production release**: `v1.0.0`, `v1.2.5`, `v2.0.0`
- **Pre-release (optional)**: `v1.0.0-alpha`, `v1.0.0-beta`, `v1.0.0-rc1`

The workflow will automatically:
- Detect alpha/beta/rc versions and mark them as "pre-release" on GitHub
- Mark regular versions as full releases

## What Happens After You Push a Tag

1. **GitHub Actions starts** (check the Actions tab)
2. **Merges dev → main** with a commit message
3. **Builds binaries** for Linux, macOS, Windows (amd64 & arm64)
4. **Creates GitHub Release** with:
   - Release notes (auto-generated changelog)
   - Binary downloads
   - Checksums file
5. **Builds Docker images** for linux/amd64 and linux/arm64
6. **Pushes to Docker Hub** with tags:
   - `<username>/goauth-backend-repo:v1.0.0`
   - `<username>/goauth-backend-repo:1.0`
   - `<username>/goauth-backend-repo:1`
   - `<username>/goauth-backend-repo:latest`

## Monitoring the Release

### Check Workflow Progress:

1. Go to your GitHub repository
2. Click **Actions** tab
3. You'll see "Release" workflow running
4. Click on it to see detailed logs

### Expected Duration:

- Total time: ~5-10 minutes
- Merge: ~30 seconds
- Build binaries: ~2-3 minutes
- Docker build & push: ~3-5 minutes

## After Release Completes

### Your Docker image will be available at:

```bash
# Pull the specific version
docker pull <your-dockerhub-username>/goauth-backend-repo:v1.0.0

# Pull the latest
docker pull <your-dockerhub-username>/goauth-backend-repo:latest

# Run it
docker run -d -p 8000:8000 \
  -e SECRET_KEY="your-secret-key" \
  -e ACCESS_EXPIRY="24h" \
  <your-dockerhub-username>/goauth-backend-repo:latest
```

### GitHub Release will include:

- **Binaries** for download (Linux, macOS, Windows)
- **Checksums** for verification
- **Changelog** (commits since last release)
- **Docker pull instructions**

## Troubleshooting

### If the workflow fails:

1. **Check secrets are set correctly**
   - Go to Settings → Secrets → Actions
   - Verify `DOCKERHUB_USERNAME` and `DOCKERHUB_TOKEN` exist

2. **Check Docker Hub token has correct permissions**
   - Token must have Read, Write, Delete permissions
   - Regenerate if needed

3. **Check merge conflicts**
   - If dev can't merge into main, you'll need to resolve conflicts manually:
   ```bash
   git checkout main
   git pull origin main
   git merge dev
   # Resolve conflicts
   git push origin main
   # Then re-push the tag
   git push origin v1.0.0 --force
   ```

4. **View detailed logs**
   - Go to Actions tab → Click the failed workflow
   - Expand each step to see error messages

### If you need to delete a tag and retry:

```bash
# Delete tag locally
git tag -d v1.0.0

# Delete tag remotely
git push --delete origin v1.0.0

# Create a new tag and push again
git tag -a v1.0.0 -m "Release v1.0.0 - Fixed"
git push origin v1.0.0
```

## Testing the Docker Image

After the release completes, test your image:

```bash
# Create a test directory
mkdir -p test-goauth
cd test-goauth

# Create .env file
cat > .env << 'EOF'
PORT=8000
HOST=0.0.0.0
DB_URL=/app/data/test.sqlite3
ACCESS_EXPIRY=24h
SECRET_KEY=test-secret-key-for-testing
EOF

# Run the container
docker run -d \
  --name goauth-test \
  -p 8000:8000 \
  --env-file .env \
  -v $(pwd)/data:/app/data \
  <your-dockerhub-username>/goauth-backend-repo:latest

# Check logs
docker logs -f goauth-test

# Test health endpoint
curl http://localhost:8000/health

# Stop and remove
docker stop goauth-test && docker rm goauth-test
```

## Best Practices

1. **Always test on dev first** before tagging
2. **Use semantic versioning**: Major.Minor.Patch (e.g., v1.2.3)
3. **Write descriptive tag messages**: `git tag -a v1.0.0 -m "Release v1.0.0 - Added user authentication"`
4. **Review the changelog** in the GitHub Release after it's created
5. **Test the Docker image** before announcing the release

## Version Numbering Guide

- **Major version** (v**2**.0.0): Breaking changes, major new features
- **Minor version** (v1.**2**.0): New features, backwards compatible
- **Patch version** (v1.0.**3**): Bug fixes, minor improvements

## Next Release Checklist

Before creating your next release tag:

- [ ] All changes committed to `dev` branch
- [ ] Tests pass locally (`cd src && make test`)
- [ ] Linting passes (`cd src && make lint`)
- [ ] Documentation updated (if needed)
- [ ] Secrets are configured in GitHub
- [ ] Docker Hub repository exists or will be auto-created
- [ ] Choose appropriate version number

Then:

```bash
git tag -a v1.0.0 -m "Release v1.0.0 - Description"
git push origin v1.0.0
```

---

**That's it!** The pipeline will handle everything else automatically. 🚀
