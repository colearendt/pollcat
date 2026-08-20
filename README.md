# pollcat

[![CI](https://github.com/colearendt/pollcat/actions/workflows/ci.yml/badge.svg)](https://github.com/colearendt/pollcat/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/colearendt/pollcat/branch/main/graph/badge.svg)](https://codecov.io/gh/colearendt/pollcat)
[![Go Reference](https://pkg.go.dev/badge/github.com/colearendt/pollcat.svg)](https://pkg.go.dev/github.com/colearendt/pollcat)

A lightweight CLI for polling network connectivity (TCP, UDP, ICMP) and DNS resolution. Features a real-time TUI with a live system clock, concurrent pollers, and exportable reports — perfect for observing how network behavior changes as you modify VPN, firewall, or Cloudflare Zero Trust configurations.

## Features

- **Multi-protocol polling** — TCP connect, UDP socket, ICMP ping, and DNS resolution.
- **Concurrent targets** — Run any mix of protocols against multiple endpoints in parallel.
- **Live TUI** — Bubble Tea interface with a running clock so you can correlate poll results with config changes.
- **Headless mode** — Auto-detects non-TTY environments (CI, SSH) and prints plain text.
- **Configurable timeout** — Per-poll timeout via `--timeout` (default: 5s).
- **Export & report** — Save raw JSON, then generate tables, CSV, JSON, or per-target summaries.
- **Small & fast** — Single static binary, zero runtime dependencies.

## Install

### Homebrew (macOS/Linux)

```bash
brew tap colearendt/tap
brew install pollcat
```

### Go

```bash
go install github.com/colearendt/pollcat@latest
```

### Pre-built binaries

Download from the [Releases](https://github.com/colearendt/pollcat/releases) page.

## Quick Start

Poll a TCP port and a DNS name simultaneously, watch the live TUI, and save results:

```bash
pollcat poll \
  --tcp 1.2.3.4:443 \
  --dns example.com \
  --interval 2s \
  --duration 60s \
  --out results.json
```

Press `q` to stop early. Then generate a summary report:

```bash
pollcat report -f results.json -t summary
```

## Usage

### `poll`

Run one or more pollers concurrently:

```bash
pollcat poll \
  --tcp 1.2.3.4:443 \
  --udp 1.2.3.4:53 \
  --icmp 1.2.3.4 \
  --dns example.com \
  --interval 2s \
  --duration 60s \
  --timeout 5s \
  --out results.json
```

| Flag | Description |
|------|-------------|
| `--tcp host:port` | TCP target(s). Repeatable. |
| `--udp host:port` | UDP target(s). Repeatable. |
| `--icmp host` | ICMP target(s). Repeatable. **Requires root/admin on most systems.** |
| `--dns hostname` | DNS target(s). Repeatable. |
| `-i, --interval` | Poll interval (default: `1s`). |
| `-d, --duration` | Total duration (`0` = run until interrupted). |
| `--timeout` | Per-poll timeout (default: `5s`). |
| `-o, --out` | Write raw results to JSON file. |
| `--no-tui` | Disable interactive UI (auto-detected in CI). |

### `report`

Generate reports from a JSON results file:

```bash
# All results as a human-readable table
pollcat report -f results.json -t table

# Filter to a specific target
pollcat report -f results.json -t table --target 1.2.3.4:443

# Per-target summary (count, min/max/avg latency, last result)
pollcat report -f results.json -t summary

# Export to CSV
pollcat report -f results.json -t csv > results.csv

# Pretty-printed JSON
pollcat report -f results.json -t json
```

| Flag | Description |
|------|-------------|
| `-f, --file` | Input JSON results file (**required**). |
| `-t, --format` | Output format: `table`, `csv`, `json`, `summary` (default: `table`). |
| `--target` | Filter results to specific target(s). Repeatable. |

## Development

### Requirements

- Go 1.23+
- `golangci-lint` (optional, for local linting)

### Quick start

```bash
make test      # run tests with race detection and coverage
make build     # build to bin/pollcat
make lint      # run linters
```

### Architecture

```
cmd/              # Cobra CLI commands (poll, report)
internal/
  poller/         # TCP, UDP, ICMP, DNS polling logic
  model/          # Shared data structures (Result, Target)
  store/          # Thread-safe in-memory result storage
  ui/             # Bubble Tea TUI model, update loop, and view
  report/         # Report formatting (table, csv, json, summary)
```

All network I/O is abstracted behind Go interfaces (`Dialer`, `Resolver`, `Pinger`) for testability.

### TDD

See [`AGENTS.md`](./AGENTS.md) for test-driven development practices, naming conventions, and coverage targets.

## Roadmap

See [`ROADMAP.md`](./ROADMAP.md) for planned features and milestones.

## Contributing

Contributions are welcome. Please open an issue to discuss significant changes before submitting a PR.

- Follow the existing Go style (`gofumpt`, `golangci-lint`).
- Write table-driven tests for new logic.
- Maintain >80% coverage in `internal/` packages.

## License

MIT
