# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is a playground/demonstration repository for Monotonix, a CI/CD tool for monorepos. The repository contains two Go applications (`hello-world` and `web-app`) that showcase Monotonix's capabilities for building, testing, and deploying applications in a monorepo structure.

**Primary Purpose**: This repository is primarily used for testing and experimenting with Monotonix features. The services here have no business value and are purely for testing purposes.

## Repository Structure

```
monotonix-playground/
├── apps/
│   ├── monotonix-global.yaml   # Global configuration for AWS IAM roles and ECR registries
│   ├── hello-world/           # Simple Hello World application (no dependencies)
│   │   ├── Dockerfile
│   │   ├── go.mod
│   │   ├── main.go
│   │   └── monotonix.yaml     # App-specific Monotonix configuration
│   └── web-app/               # Web application with dependencies
│       ├── cmd/
│       │   ├── api-server/    # API server microservice
│       │   │   ├── Dockerfile
│       │   │   ├── main.go
│       │   │   └── monotonix.yaml
│       │   └── worker/        # Background worker microservice
│       │       ├── Dockerfile
│       │       ├── main.go
│       │       └── monotonix.yaml
│       ├── pkg/
│       │   └── common/        # Shared library package
│       │       ├── message.go
│       │       ├── message_test.go
│       │       └── monotonix.yaml
│       ├── go.mod
│       └── monotonix.yaml
└── renovate.json              # Dependency update automation
```

## Common Development Commands

### Running Applications Locally

```bash
# Run the hello-world server
cd apps/hello-world
go run main.go

# Run the web-app API server
cd apps/web-app/cmd/api-server
go run main.go

# Run the web-app worker
cd apps/web-app/cmd/worker
go run main.go
```

### Testing

```bash
# Run tests for hello-world
cd apps/hello-world
go test ./...

# Run tests for web-app (includes common package tests)
cd apps/web-app
go test ./...
```

### Building Docker Images Locally

```bash
# Build hello-world app
cd apps/hello-world
docker build -t hello-world:local .

# Build web-app API server
cd apps/web-app/cmd/api-server
docker build -t web-app-api-server:local .

# Build web-app worker
cd apps/web-app/cmd/worker
docker build -t web-app-worker:local .
```

## Architecture & Key Concepts

### Monotonix Configuration

Each app has a `monotonix.yaml` file that defines:

- **Jobs**: Build and test workflows triggered by different events
- **Environments**: prd (production), dev_main, and dev_pr
- **Dependencies**: Some apps depend on shared packages (e.g., `web-app/cmd/api-server` and `web-app/cmd/worker` both depend on `web-app/pkg`)
- **Tagging strategies**:
  - `semver_datetime`: For production builds
  - `always_latest`: For dev_main builds
  - `pull_request`: For PR builds

### Application Types

- **hello-world**: A simple standalone application with no dependencies, demonstrating basic Monotonix functionality
- **web-app**: A more complex application showcasing dependency management:
  - `web-app/pkg/common`: Shared library package
  - `web-app/cmd/api-server`: API server microservice that depends on the common package
  - `web-app/cmd/worker`: Background worker that depends on the common package

### Docker Build Process

All applications use multi-stage Docker builds:

1. **Builder stage**: Compiles Go binary using `golang:1.24.4`
2. **Runtime stage**: Uses distroless base image for minimal attack surface

### AWS ECR Repositories

- **Production**: `920373013500.dkr.ecr.ap-northeast-1.amazonaws.com/monotonix`
- **Dev Main**: `615299752259.dkr.ecr.ap-northeast-1.amazonaws.com/monotonix`
- **Dev PR**: `615299752259.dkr.ecr.ap-northeast-1.amazonaws.com/monotonix-pr`

## Working with Monotonix

When making changes to an app:

1. The Monotonix system automatically detects which apps have changed
2. Only the changed apps will be built and deployed
3. Jobs are triggered based on the configuration in each app's `monotonix.yaml`

## Testing Monotonix New Features

When testing new Monotonix features:

1. **Prioritize feature testing over stability** - This is a test environment, so breaking existing services is acceptable
2. **Be experimental** - Try edge cases and unusual configurations to thoroughly test new capabilities
3. **Document findings** - Note any unexpected behaviors or issues discovered during testing
4. **Use aggressive testing approaches** - Don't hesitate to make breaking changes if it helps validate Monotonix functionality

### Testing Monotonix PRs

This repository includes utility scripts to efficiently test Monotonix changes before they are merged. These scripts allow you to quickly switch the Monotonix action references in all workflow files to test pull requests or specific branches from the [monotonix repository](https://github.com/yuya-takeyama/monotonix).

**Available Scripts:**

1. **`scripts/switch-monotonix-ref.sh`** - Low-level ref switching utility
   ```bash
   # Switch to a specific branch/tag/commit
   ./scripts/switch-monotonix-ref.sh feat/metadata-support
   
   # Check current ref
   ./scripts/switch-monotonix-ref.sh --current
   
   # Reset to latest release
   ./scripts/switch-monotonix-ref.sh --reset
   ```

2. **`scripts/test-monotonix-pr.sh`** - High-level PR testing workflow
   ```bash
   # Test a specific PR from monotonix repository
   ./scripts/test-monotonix-pr.sh 151
   
   # After testing, cleanup and restore defaults
   ./scripts/test-monotonix-pr.sh --cleanup
   ```

**Typical Testing Workflow:**

When you need to test a Monotonix PR (e.g., [PR #151](https://github.com/yuya-takeyama/monotonix/pull/151)):

```bash
# 1. Switch workflows to use the PR branch
./scripts/test-monotonix-pr.sh 151

# 2. Push changes to trigger CI with the new Monotonix version
git push

# 3. Monitor CI results
gh pr checks --watch

# 4. Make test changes to apps/ to validate the new feature
# (e.g., add metadata fields to monotonix.yaml files)

# 5. After testing, cleanup and restore defaults
./scripts/test-monotonix-pr.sh --cleanup
git push
```

**Background:**

These scripts were created to streamline the testing process for Monotonix changes. Since this playground repository serves as the primary testing ground for Monotonix features, being able to quickly switch between different versions of Monotonix actions is essential for:

- Validating new features before merging PRs
- Testing breaking changes in isolation
- Debugging issues with specific commits
- Experimenting with unreleased features

The scripts handle all the tedious work of updating action references across multiple workflow files, creating proper test commits, and cleaning up afterward, allowing you to focus on testing the actual functionality.

## Adding a New Application

To add a new application to this monorepo:

1. Create a new directory under `apps/`
2. Add a `monotonix.yaml` configuration file (copy from existing apps)
3. Add a `Dockerfile` following the same multi-stage pattern
4. Add a `go.mod` file and your Go source code
5. Update the app name in `monotonix.yaml`
