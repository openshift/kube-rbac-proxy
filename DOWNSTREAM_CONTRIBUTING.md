# Contributing to OpenShift kube-rbac-proxy

Thank you for your interest in contributing to the OpenShift fork of kube-rbac-proxy! This document provides guidelines specific to contributing to the OpenShift downstream version.

## OpenShift Fork Context

This repository is the OpenShift downstream fork of kube-rbac-proxy, maintained for use in OpenShift clusters. We sync changes from the upstream brancz/kube-rbac-proxy project while maintaining OpenShift-specific requirements and integrations.

**Upstream-first approach**: If a change can be made in the upstream project (brancz/kube-rbac-proxy), it should be submitted there first and then cherry-picked into this fork. Only downstream-specific work — for example, a customer request that provides no upstream benefit, or an OpenShift-specific integration — should be contributed directly to this repository. See [Syncing from Upstream](#syncing-from-upstream) for the full workflow.

## Getting Started

### Prerequisites

- Go 1.25.8 or later
- [kind](https://kind.sigs.k8s.io/) for local testing
- [golangci-lint](https://golangci-lint.run/) for linting
- Docker for building container images

### Development Environment Setup

1. Fork and clone the repository:
   ```bash
   git clone https://github.com/YOUR-USERNAME/kube-rbac-proxy.git
   cd kube-rbac-proxy
   ```

2. Add the remotes so you have the full chain — your fork → `upstream` (this OpenShift repository) → `original` (the brancz project this fork is based on):
   ```bash
   # upstream: the OpenShift downstream fork (this repository)
   git remote add upstream https://github.com/openshift/kube-rbac-proxy.git
   # original: the upstream brancz project this fork tracks
   git remote add original https://github.com/brancz/kube-rbac-proxy.git
   ```

3. Build the project:
   ```bash
   make build
   ```

4. Run tests locally:
   ```bash
   make test-local
   ```
   **Note**: This creates a kind cluster and will override your default kind cluster.

## Making Changes

### Coding Conventions

- Follow standard Go conventions and idioms
- Run `golangci-lint` before submitting - the CI will fail if linting errors exist
- All new code must include appropriate unit tests
- Public APIs and functions should include godoc comments
- Keep code changes focused and atomic - one logical change per PR

### License Headers

All Go files must include the appropriate license header. Run `make check-license` to verify:

```bash
make check-license
```

The CI will fail if license headers are missing or incorrect.

### Testing Requirements

#### Unit Tests

- Required for all new functionality
- Run unit tests: `make test-unit`
- Unit tests should be fast and not require external dependencies
- Tests must pass with race detection enabled (`-race` flag)

#### End-to-End Tests

- Required for changes affecting core proxy behavior
- Run e2e tests locally: `make test-local` (sets up kind cluster and runs full test suite)
- E2E tests run in a kind cluster and validate real Kubernetes integration

#### Test Coverage

- Maintain or improve existing test coverage
- Both unit and e2e tests run in CI and must pass before merge

## Submitting Changes

### Pull Request Process

#### For Red Hat / Control Plane Team Members

1. **Create a feature branch** from `master`:
   ```bash
   git checkout -b your-feature-name
   ```

2. **Make your changes** following the coding conventions above

3. **Run the full test suite locally**:
   ```bash
   make test-local
   ```

4. **Verify generated files are up to date**:
   ```bash
   make generate
   git diff --exit-code
   ```
   If there are differences, commit them.

5. **Submit your PR** with:
   - Clear description of what changed and why
   - Reference any related Jira issues (e.g., "Related: OCPBUGS-12345")
   - A clear, descriptive commit message; if the change is OpenShift-specific (not coming from upstream), say so in the commit body or PR description
   - Screenshots or example output for user-facing changes

6. **Request reviews**:
   - Tag appropriate team members from [OWNERS](OWNERS) file
   - For OpenShift-specific changes, ensure reviewers understand the downstream context
   - Internal Red Hat employees: You can also request reviews via the `#forum-kube-apiserver` Slack channel

#### For External Contributors

1. Follow the same process as above for creating your PR

2. **Request reviews**:
   - The team will automatically be notified via GitHub
   - Reviews from OWNERS are required for merge
   - You can comment on the PR requesting review if you don't hear back within a few business days

3. **Getting help**:
   - For questions: Comment on your PR or open a GitHub issue
   - For Red Hat employees: Use `#forum-kube-apiserver` Slack channel

### PR Review Standards

- PRs require approval from at least one OWNER (see [OWNERS](OWNERS) file)
- All CI checks must pass:
  - License header check
  - Go linting
  - Unit tests
  - E2E tests
  - Generate check (ensures generated docs/examples are current)
  - Build check

### Code Review Expectations

- **What reviewers look for**:
  - Correctness and security implications
  - Test coverage and quality
  - Code clarity and maintainability
  - Alignment with project goals and architecture
  - Breaking changes or deprecation path
  - Impact on OpenShift clusters

### Merge Requirements (Prow/Tide)

PRs must satisfy the following Prow/Tide requirements before merge:

- **`lgtm`** label: Requires `/lgtm` from an approver (review approval)
- **`approved`** label: Requires `/approve` from an OWNER
- All CI checks passing (e.g., `tide/merge-method-merge`)
- No unresolved review comments
- Commits should be clean and logical (squash if needed)

**Common Prow commands**:
- `/lgtm` - Adds lgtm label (reviewers)
- `/approve` - Adds approved label (approvers/owners)
- `/hold` - Prevents auto-merge
- `/hold cancel` - Removes hold
- `/retest` - Reruns failed tests

See [Prow Command Help](https://prow.ci.openshift.org/command-help) for full list.

## OpenShift-Specific Contributions

### Downstream-Only Changes

When making changes specific to OpenShift (not coming from upstream), write a clear, descriptive commit message and make the downstream-only nature explicit in the commit body and PR description. This helps reviewers and future maintainers identify downstream-specific patches when syncing with upstream.

```text
Add OpenShift-specific TLS cipher configuration

This is a downstream-only change: it adds support for OpenShift's
required cipher suite configuration that differs from upstream defaults
and is not intended to be contributed upstream.

Related: OCPBUGS-12345
```

Downstream syncs themselves land through dedicated merge branches (e.g.,
`merge-v0.22.1-downstream`) rather than via a commit-message prefix.

### Syncing from Upstream

The Control Plane team periodically syncs changes from upstream (brancz/kube-rbac-proxy). The process:

1. **Upstream tracking**: Monitor upstream releases and relevant PRs
2. **Cherry-pick or merge**: Bring in upstream changes, resolving conflicts with OpenShift-specific patches
3. **Testing**: Run full test suite to ensure upstream changes work with OpenShift
4. **Documentation**: Update CHANGELOG.md noting upstream version synced

**For contributors**: If you're implementing a feature that would benefit upstream, consider:
1. Submitting the change to upstream first (brancz/kube-rbac-proxy)
2. Once merged upstream, sync it to the OpenShift fork
3. Keep any OpenShift-specific modifications as separate, clearly described commits

## Updating Dependencies

To update Go dependencies:

```bash
make update-go-deps
```

This is typically done during release cycles. Include dependency updates in separate commits when possible.

## Building and Testing

### Local Build

```bash
make build          # Build for current OS/arch
make crossbuild     # Build for all supported platforms
```

Binary output: `_output/kube-rbac-proxy`

### Container Image

```bash
make container      # Build container image
make test-container # Build test container
```

### Running Examples

The [`examples/`](examples/) directory contains working examples:

- [non-resource-url](examples/non-resource-url)
- [resource-attributes](examples/resource-attributes)
- [oidc](examples/oidc)
- [rewrites](examples/rewrites)

## Project Status

This is the **OpenShift downstream fork** of kube-rbac-proxy, maintained for production use in OpenShift clusters. While we track upstream development, OpenShift-specific requirements take precedence.

## Release Process

Releases are handled by the OpenShift Control Plane team. See [RELEASE.md](RELEASE.md) for the release process documentation.

## Getting Help

### For Red Hat Employees

- **Internal Slack**: `#forum-kube-apiserver` for questions and discussions
- **Team contacts**: See [OWNERS](OWNERS) file for team members

### For External Contributors

- **Questions about contributing?** Open a GitHub issue or comment on your PR
- **Bug reports**: Open an issue with steps to reproduce
- **Feature requests**: Open an issue describing the use case and proposed solution

## Code of Conduct

This project follows the Kubernetes community [Code of Conduct](https://kubernetes.io/community/code-of-conduct/).
