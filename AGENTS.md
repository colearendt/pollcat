# pollcat Agent Instructions

## Project Overview

`pollcat` is a lightweight Go CLI for polling network connectivity (TCP) and DNS resolution. It runs multiple concurrent pollers, displays a real-time TUI with a system clock, and can generate CLI-friendly reports and CSV exports.

## Architecture

- **cmd/**: Cobra commands (`poll`, `report`).
- **internal/poller/**: Pure polling logic (TCP connect, DNS lookup). Network calls are behind interfaces for testability.
- **internal/model/**: Shared data structures (`Result`, `PollConfig`).
- **internal/store/**: Thread-safe in-memory result storage.
- **internal/ui/**: Bubble Tea TUI model, update loop, and view.
- **internal/report/**: Report formatting (table, CSV).

## TDD Practices

1. **Red-Green-Refactor**: Write a failing test first, implement the minimal code to pass, then refactor.
2. **Interfaces for I/O**: All network I/O (TCP `net.Dial`, DNS `net.Resolver`) must be abstracted behind Go interfaces. Tests inject fakes or recordable mocks.
3. **Table-Driven Tests**: Use table-driven tests for all logic branches. Prefer `github.com/stretchr/testify/assert` for readable assertions.
4. **No Merging Without Green CI**: The GitHub Actions workflow runs tests, race detection (`go test -race`), linting (`golangci-lint`), and builds on every PR.
5. **Coverage Baseline**: Aim for >80% coverage in `internal/` packages. Coverage is informational, not a gate, but significant drops require justification.
6. **Test Naming**: `Test<Struct>_<Method>_<Scenario>` or `Test<Function>_<Scenario>`.
7. **Parallel Tests**: Use `t.Parallel()` where tests do not share mutable state.
8. **Golden Files**: For complex report output, use `testdata/` golden files and `testify` golden helpers.

## Coding Conventions

- Standard Go formatting (`gofumpt`).
- `golangci-lint` with default linters + `errcheck`, `goimports`, `govet`, `staticcheck`.
- Context-first: All long-running operations accept a `context.Context`.
- Errors: Wrap with `fmt.Errorf("...: %w", err)` where appropriate.

## Dependencies

- `github.com/spf13/cobra` — CLI framework.
- `github.com/charmbracelet/bubbletea` — TUI framework.
- `github.com/charmbracelet/lipgloss` — Terminal styling.
- `github.com/stretchr/testify` — Test assertions and suites.
