# Contributing to ozzo-validation

Thank you for your interest in contributing! This document covers how to build, test, and submit changes.

## Prerequisites

- **Go 1.21+** ([download](https://go.dev/dl/))
- **golangci-lint** (`go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest`)

## Building

```bash
go build ./...
```

## Running Tests

```bash
go test -race ./...
go vet ./...
```

## Code Style

- Run `go fmt ./...` before every commit. CI enforces this.
- Run `golangci-lint run --timeout=5m` and fix all issues.
- Follow standard Go naming conventions.
- Handle every error or explicitly ignore with `_ =` and a comment.
- Exported types and functions must have doc comments.

## Pull Request Workflow

1. Fork the repository and create a feature branch:
   ```bash
   git checkout -b feat/my-feature
   ```
2. Make your changes. Keep commits focused and well-described.
3. Verify locally:
   ```bash
   go fmt ./...
   go build ./...
   go test ./...
   golangci-lint run --timeout=5m
   ```
4. Push and open a pull request against `master`.
5. Wait for CI to pass. All checks must be green before merge.

Commit messages follow [Conventional Commits](https://www.conventionalcommits.org/):
```
feat: add UUID v7 validation rule
fix: require "+" prefix in E.164 phone validation
docs: update installation instructions
```

## Finding Work

Check the [issues](https://github.com/go-ozzo/ozzo-validation/issues) page for tasks labeled [`good first issue`](https://github.com/go-ozzo/ozzo-validation/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22).

## Reporting Bugs

Open an issue with a clear description, a minimal code example that reproduces the problem, and the Go version you're using.

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
