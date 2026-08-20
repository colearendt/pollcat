# ROADMAP

## v0.1.0 — MVP

- [x] Project scaffolding (`go.mod`, Cobra, Bubble Tea).
- [x] AGENTS.md and ROADMAP.md.
- [x] TCP poller with latency recording.
- [x] UDP poller with latency recording.
- [x] ICMP poller (ping) with latency recording.
- [x] DNS poller with latency and resolved-IP recording.
- [x] Thread-safe in-memory result store.
- [x] Bubble Tea TUI showing system clock and live poller status.
- [x] `poll` command (flags: `--tcp`, `--udp`, `--icmp`, `--dns`, `--interval`, `--duration`).
- [x] `report` command (formats: `table`, `csv`, `json`, `summary`).
- [x] Per-target filtering and summary in reports.
- [x] GitHub Actions CI (test, lint, build).
- [x] Unit tests with >80% coverage for `internal/`.

## v0.2.0 — Polish & Configuration

- [ ] Config file support (YAML/JSON) for complex multi-poller setups.
- [ ] Color-coded status in TUI (green = healthy, red = error, yellow = slow).
- [ ] Latency histogram / percentile summaries in reports.
- [ ] Export to JSON Lines (`.jsonl`) for streaming post-processing.
- [ ] Better error handling and retry logic (exponential backoff for DNS).

## v0.3.0 — Advanced Features

- [ ] HTTP/HTTPS polling mode (HEAD request latency + status code).
- [ ] Alerting thresholds (e.g., latency > 500 ms or DNS failure).
- [ ] WebSocket or SSE live dashboard alternative to TUI.
- [ ] Cross-compilation matrix in CI (Linux, macOS, Windows).
- [ ] Homebrew tap via GoReleaser.

## Future Ideas

- Prometheus metrics endpoint (for long-running daemon mode).
- SQLite persistence for historical analysis.
- Integration tests with `testcontainers-go`.
