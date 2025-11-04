# GitHub Actions CI/CD Setup

This document describes the GitHub Actions workflows configured for the goAuth project.

## Workflows

### 1. CI/CD Pipeline (`build.yaml`)

**Triggers:**
- Push to `main`, `dev`, or `feature/**` branches
- Pull requests to `main` or `dev`
- Manual workflow dispatch

**Jobs:**

#### Lint
- Runs `golangci-lint` with comprehensive linting rules
- Checks if `go.mod` and `go.sum` are tidy
- Ensures code quality and formatting standards

#### Build and Test
- Builds the application for multiple platforms
- Generates Swagger documentation
- Runs all tests with race detection
- Generates code coverage report
- Uploads coverage to Codecov (requires `CODECOV_TOKEN` secret)
- Uploads build artifacts

#### Security
- Runs Gosec security scanner for Go code vulnerabilities
- Runs Trivy vulnerability scanner for dependencies
- Uploads results to GitHub Security tab

#### Docker Build
- Builds multi-architecture Docker images (amd64, arm64)
- Pushes images to GitHub Container Registry
- Tags images appropriately based on branch/tag
- Scans Docker images for vulnerabilities
- Uses build cache for faster builds

#### Integration Test
- Starts services using Docker Compose
- Runs basic integration tests
- Verifies service health

#### Deploy
- Deploys to production (only on main branch)
- **Note:** Add your specific deployment steps in this job

#### Summary
- Provides a comprehensive pipeline summary
- Shows status of all jobs

### 2. Release Workflow (`release.yaml`)

**Triggers:**
- Push of version tags (e.g., `v1.0.0`)
- Manual workflow dispatch

**Features:**
- Builds binaries for multiple platforms:
  - Linux (amd64, arm64)
  - macOS (amd64, arm64)
  - Windows (amd64)
- Generates SHA256 checksums
- Creates GitHub releases with changelog
- Builds and pushes versioned Docker images
- Supports semantic versioning

## Setup Instructions

### 1. Required Secrets

Add these secrets in GitHub: Settings → Secrets and variables → Actions

#### Optional Secrets:
- `CODECOV_TOKEN`: For uploading code coverage (get from codecov.io)

#### For Deployment (if using):
- Add any deployment-specific secrets based on your deployment target

### 2. GitHub Container Registry

The workflow automatically publishes Docker images to GitHub Container Registry (ghcr.io).

**To pull images:**
```bash
# Authenticate with GitHub
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin

# Pull the image
docker pull ghcr.io/ahmadassadi/goauth:latest
```

### 3. Permissions

Ensure your repository has the following permissions enabled:
- Settings → Actions → General → Workflow permissions: "Read and write permissions"
- Settings → Code security and analysis → Enable Dependabot alerts

### 4. Branch Protection

Recommended branch protection rules for `main`:

- Require pull request reviews before merging
- Require status checks to pass before merging:
  - `Lint Code`
  - `Build and Test`
  - `Security Scan`
  - `Build Docker Image`
- Require branches to be up to date before merging
- Include administrators

## Usage

### Running Workflows Manually

1. Go to Actions tab in GitHub
2. Select the workflow
3. Click "Run workflow"
4. Select branch and input parameters (if any)

### Creating a Release

**Option 1: Using Git Tags**
```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

**Option 2: Manual Dispatch**
1. Go to Actions → Release workflow
2. Click "Run workflow"
3. Enter version (e.g., `v1.0.0`)

### Viewing Build Artifacts

1. Go to Actions tab
2. Select a workflow run
3. Scroll to "Artifacts" section
4. Download the artifacts

### Checking Security Scan Results

1. Go to Security tab
2. Click "Code scanning alerts"
3. Review any findings

## Docker Compose

The updated `docker-compose.yml` includes:
- Health checks for the application
- Persistent volumes for database
- Network isolation
- Environment variable management
- Easy PostgreSQL integration (commented out)

**Usage:**
```bash
# Start services
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

## Dockerfile

The updated `Dockerfile` includes:
- Multi-stage build for smaller images
- Health checks
- Proper volume management
- Fixed environment variable syntax
- Security best practices

## Continuous Improvement

### Adding Tests

Add your tests in the appropriate `*_test.go` files. The CI pipeline will automatically run them.

### Adding Deployment

Update the `deploy` job in `build.yaml`:

```yaml
- name: Deploy to production
  run: |
    # Add your deployment commands here
    # Examples:
    # - kubectl apply -f k8s/
    # - aws ecs update-service ...
    # - docker-compose -f production.yml up -d
```

### Environment-Specific Deployments

Create additional workflow files for different environments:
- `.github/workflows/deploy-staging.yaml`
- `.github/workflows/deploy-production.yaml`

## Monitoring

- **Build Status**: Visible in README with badge
- **Code Coverage**: Track in Codecov
- **Security Alerts**: GitHub Security tab
- **Dependency Updates**: Dependabot PRs

## Badge for README

Add this to your `README.md`:

```markdown
[![CI/CD Pipeline](https://github.com/ahMADASSadi/goAuth/actions/workflows/build.yaml/badge.svg)](https://github.com/ahMADASSadi/goAuth/actions/workflows/build.yaml)
[![Release](https://github.com/ahMADASSadi/goAuth/actions/workflows/release.yaml/badge.svg)](https://github.com/ahMADASSadi/goAuth/actions/workflows/release.yaml)
```

## Troubleshooting

### Build Failures

1. Check the logs in Actions tab
2. Run locally: `cd src && go build ./cmd/main.go`
3. Check lint errors: `golangci-lint run`

### Docker Build Issues

1. Test locally: `docker build -t goauth:test .`
2. Check Docker logs: `docker-compose logs`
3. Verify .env file exists in `src/` directory

### Test Failures

1. Run tests locally: `cd src && go test ./...`
2. Run with verbose output: `go test -v ./...`
3. Check race conditions: `go test -race ./...`

## Support

For issues or questions:
1. Check the Actions logs
2. Review this documentation
3. Open an issue in the repository
