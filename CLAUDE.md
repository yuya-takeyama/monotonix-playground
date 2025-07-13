# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is a playground/demonstration repository for Monotonix, a CI/CD tool for monorepos. The repository contains two simple Go applications (`echo` and `hello-world`) that showcase Monotonix's capabilities for building, testing, and deploying applications in a monorepo structure.

## Repository Structure

```
monotonix-playground/
├── apps/
│   ├── monotonix-global.yaml   # Global configuration for AWS IAM roles and ECR registries
│   ├── echo/                   # Echo server application
│   │   ├── Dockerfile
│   │   ├── go.mod
│   │   ├── main.go
│   │   └── monotonix.yaml     # App-specific Monotonix configuration
│   └── hello-world/           # Hello World application
│       ├── Dockerfile
│       ├── go.mod
│       ├── main.go
│       └── monotonix.yaml     # App-specific Monotonix configuration
└── renovate.json              # Dependency update automation
```

## Common Development Commands

### Running Applications Locally

```bash
# Run the echo server
cd apps/echo
go run main.go

# Run the hello-world server
cd apps/hello-world
go run main.go
```

### Testing

```bash
# Run tests for a specific app (though no test files exist currently)
cd apps/echo
go test ./...

cd apps/hello-world
go test ./...
```

### Building Docker Images Locally

```bash
# Build echo app
cd apps/echo
docker build -t echo:local .

# Build hello-world app
cd apps/hello-world
docker build -t hello-world:local .
```

## Architecture & Key Concepts

### Monotonix Configuration

Each app has a `monotonix.yaml` file that defines:

- **Jobs**: Build and test workflows triggered by different events
- **Environments**: prd (production), dev_main, and dev_pr
- **Tagging strategies**:
  - `semver_datetime`: For production builds
  - `always_latest`: For dev_main builds
  - `pull_request`: For PR builds

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

## Adding a New Application

To add a new application to this monorepo:

1. Create a new directory under `apps/`
2. Add a `monotonix.yaml` configuration file (copy from existing apps)
3. Add a `Dockerfile` following the same multi-stage pattern
4. Add a `go.mod` file and your Go source code
5. Update the app name in `monotonix.yaml`
