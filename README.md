# cli-conn

A lightweight Go CLI for polling network connectivity (TCP, UDP, ICMP) and DNS resolution with a real-time TUI, system clock, and exportable reports.

## Features

- **TCP polling** — Connect to any `host:port` and record latency / success.
- **UDP polling** — Create UDP socket and record latency / success.
- **ICMP polling** — Send ICMP echo (ping) and record latency / reply.
- **DNS polling** — Resolve hostnames and capture resolved IPs and latency.
- **Concurrent pollers** — Run multiple probes in parallel.
- **Live TUI** — Watch status update in real time with a system clock (great for correlating with config changes).
- **Per-target reports** — Filter and summarize results by individual target.
- **Small & fast** — Single binary, zero runtime dependencies.

## Install

```bash
go install github.com/colearendt/cli-conn@latest
```

Or download a pre-built binary from the [Releases](https://github.com/colearendt/cli-conn/releases) page.

## Usage

### Poll

Run multiple pollers simultaneously:

```bash
cli-conn poll \
  --tcp 1.2.3.4:443 \
  --udp 1.2.3.4:53 \
  --icmp 1.2.3.4 \
  --dns example.com \
  --interval 2s \
  --duration 60s \
  --out results.json
```

Flags:
- `--tcp` — TCP target(s) (`host:port`). Repeatable.
- `--udp` — UDP target(s) (`host:port`). Repeatable.
- `--icmp` — ICMP target(s) (`host` or `ip`). Repeatable. **Requires root/admin on most systems.**
- `--dns` — DNS target(s) (`hostname`). Repeatable.
- `-i, --interval` — Poll interval (default `1s`).
- `-d, --duration` — Total polling duration (`0` = run until interrupted).
- `-o, --out` — Write raw results to a JSON file.
- `--no-tui` — Disable the interactive UI and print plain text (auto-detected in CI).

### Report

Generate a human-readable table or CSV from a saved JSON file:

```bash
# All results
cli-conn report -f results.json -t table

# Filter to a specific target
cli-conn report -f results.json -t table --target 1.2.3.4:443

# Per-target summary (stats per combo)
cli-conn report -f results.json -t summary

# Export to CSV
cli-conn report -f results.json -t csv > results.csv
```

Report formats:
- `table` — Human-readable aligned columns.
- `csv` — Machine-friendly CSV with headers.
- `json` — Pretty-printed JSON.
- `summary` — Aggregated stats per target (total, success/failure count, min/max/avg latency).

## Development

### Requirements

- Go 1.23+
- `golangci-lint` (optional, for local linting)

### Quick start

```bash
make test      # run tests with race detection
make build     # build to bin/cli-conn
make lint      # run linters
```

### Architecture

| Package | Purpose |
|---------|---------|
| `cmd/` | Cobra CLI commands (`poll`, `report`). |
| `internal/poller/` | TCP, UDP, ICMP and DNS polling logic behind testable interfaces. |
| `internal/model/` | Shared data structures (`Result`, `Target`). |
| `internal/store/` | Thread-safe in-memory result storage. |
| `internal/ui/` | Bubble Tea TUI model, update loop, and view. |
| `internal/report/` | Report formatting (table, CSV, JSON, summary) with target filtering. |

### TDD

See [`AGENTS.md`](./AGENTS.md) for our test-driven development practices and coding conventions.

## Roadmap

See [`ROADMAP.md`](./ROADMAP.md) for planned features and milestones.

## License

MIT
