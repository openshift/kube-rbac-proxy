# AI Agent Instructions for kube-rbac-proxy

This document provides guidance for AI agents working with the kube-rbac-proxy codebase.

## Project Overview

kube-rbac-proxy is a small HTTP reverse proxy for a single upstream that performs RBAC authorization against the Kubernetes API using SubjectAccessReview. It's designed to protect endpoints (especially metrics endpoints) that would otherwise be accessible to any Pod in a cluster without NetworkPolicies.

**Key architectural principle**: This is a security-focused component. Changes must maintain or improve the security posture.

## Repository Structure

```text
kube-rbac-proxy/
├── cmd/kube-rbac-proxy/           # Main entry point and CLI
│   └── app/options/               # Flag definitions and configuration
├── pkg/                           # Core library code
│   ├── proxy/                     # HTTP proxy implementation
│   ├── filters/                   # Authentication and path filtering
│   ├── hardcodedauthorizer/       # Authorization logic
│   └── tls/                       # TLS certificate management
├── test/                          # Test utilities and e2e tests
│   ├── e2e/                       # End-to-end test scenarios
│   └── kubetest/                  # Kubernetes test helpers
├── examples/                      # Working examples of different configurations
└── scripts/                       # Build and generation scripts
```

## Important Context

### Security Considerations

1. **Token Security**: This proxy can receive bearer tokens. The receiving side can use these tokens to impersonate the client. Only use token authentication when the receiving side is already higher-privileged or the token is low-privileged. See [README - ServiceAccount token security](README.md#notes-on-serviceaccount-token-security).

2. **TLS First**: Prefer mTLS for authentication over bearer tokens for better security properties.

3. **Authorization Model**: The proxy performs a SubjectAccessReview against the Kubernetes API for every request. This is the core security mechanism - do not bypass or weaken it.

### Project Status

- **OpenShift fork**: This is the OpenShift downstream fork of kube-rbac-proxy, maintained for use in OpenShift clusters
- **Production use**: This proxy is widely used across OpenShift components to secure internal service endpoints

### Common Pitfalls

1. **Do not add flags for insecure operation**: The project is deprecating `--insecure-listen-address` and similar flags. New changes should not add insecure modes.

2. **License headers are required**: All Go files must have license headers. Run `make check-license` before submitting changes.

3. **Generated files**: The README includes embedded help text and examples. After changing flags or examples, run `make generate` and commit the changes.

4. **Kind cluster override**: `make test-local` creates/recreates the default kind cluster to run tests locally. This will destroy any existing default kind cluster.

## Making Code Changes

### Before You Start

1. Read [CONTRIBUTING.md](CONTRIBUTING.md) for process details
2. Check [ARCHITECTURE.md](ARCHITECTURE.md) for design context
3. Review recent commits to understand current work direction
4. Check existing issues/PRs to avoid duplicate work

### Development Workflow

```bash
# Build the binary
make build

# Run unit tests (fast, no k8s cluster needed)
make test-unit

# Run full local test suite (creates kind cluster, runs e2e tests)
make test-local

# Check generated files are current
make generate && git diff --exit-code
```

### Testing Requirements

- **Unit tests required** for all new functionality
- **E2E tests required** for changes to core proxy behavior
- Tests must pass with race detection (`-race` flag is used in CI)
- E2E tests run in a kind cluster
- Test locations:
  - Unit tests: `*_test.go` files alongside implementation
  - E2E tests: `test/e2e/` directory

### Code Quality Checks

All of these run in CI and must pass:

```bash
make check-license  # Verify license headers
make generate       # Update generated docs/examples
golangci-lint run   # Linting (configured in .golangci.yaml)
make test-unit      # Unit tests
make test-e2e       # E2E tests (requires kind cluster)
```

## Key Design Patterns

### Proxy Architecture

The proxy sits between a client and a backend service (the service being protected):

```text
Client -> kube-rbac-proxy -> Backend Service (e.g., Prometheus /metrics)
          |
          v
     Kubernetes API (SubjectAccessReview)
```

1. **Authentication**: Validate client identity via client TLS cert or bearer token (TokenReview against k8s API)
2. **Authorization**: Perform SubjectAccessReview against k8s API
3. **Proxy**: If authorized, forward request to the backend service

### Configuration Model

- **CLI flags**: Most configuration via command-line flags (see [`cmd/kube-rbac-proxy/app/options/options.go`](cmd/kube-rbac-proxy/app/options/options.go) for flag definitions)
- **Config file**: `--config-file` for SubjectAccessReview details and rewrites
  - Example minimal config:
    ```yaml
    authorization:
      resourceAttributes:
        namespace: default
        apiVersion: v1
        resource: services
        name: my-service
    ```
  - See [`examples/`](examples/) directory for complete working examples
- **Examples directory**: Contains working configurations for common scenarios

### Backend Service Communication

- Supports HTTP/1.1, HTTP/2, and h2c (HTTP/2 cleartext)
- Can use client certificates for backend authentication (`--upstream-client-cert-file`)
- Supports custom CA for backend TLS verification (`--upstream-ca-file`)

## Common Operations

### Adding a New Flag

1. Add flag definition in [`cmd/kube-rbac-proxy/app/options/options.go`](cmd/kube-rbac-proxy/app/options/options.go) within the `Flags()` method
   - Flag definition inherently includes help text, e.g.:
     ```go
     flagset.StringVar(&o.InsecureListenAddress, "insecure-listen-address", "", 
         "[DEPRECATED] The address the kube-rbac-proxy HTTP server should listen on.")
     ```
2. Run `make generate` to update embedded documentation in README
3. Add e2e tests in `test/e2e/` if the flag affects runtime behavior (minimal flag testing is the current convention)

### Modifying Authentication/Authorization

Changes to `pkg/filters/` (authentication filters) or `pkg/hardcodedauthorizer/` (authorization logic) handle the security-critical request validation flow:

- **`pkg/filters/`**: Performs client authentication (validating TLS certs, bearer tokens via TokenReview, OIDC tokens)
- **`pkg/hardcodedauthorizer/`**: Performs authorization checks (SubjectAccessReview against Kubernetes RBAC)

When modifying these components:

1. Understand how the current authentication and authorization flow works (see [ARCHITECTURE.md](ARCHITECTURE.md))
2. Write minimal unit tests covering the core behavior (current testing convention favors minimal over comprehensive test coverage)
3. Add e2e tests demonstrating the change works end-to-end
4. Document security implications in your commit message and/or PR description

### Adding a New Example

1. Create directory under `examples/` with descriptive name
2. Include working manifests and a README explaining the use case
3. Run `make generate` to update root README with example links

## Kubernetes API Integration

This project uses Kubernetes client libraries heavily:

- `k8s.io/client-go`: Kubernetes API client
- `k8s.io/apiserver`: Reuses apiserver components for authentication
- `k8s.io/component-base`: Version reporting and common utilities

When updating Kubernetes dependencies (`make update-go-deps`), test thoroughly - API changes can break authentication/authorization flows.

## References

- [Main README](README.md) - User-facing documentation
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contribution process
- [ARCHITECTURE.md](ARCHITECTURE.md) - Design decisions and tradeoffs
- [RELEASE.md](RELEASE.md) - Release process
- [examples/](examples/) - Working configuration examples
