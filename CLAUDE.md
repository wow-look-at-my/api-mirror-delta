# api-mirror-delta — repo orientation for Claude

A development proxy for api-mirror: each request goes to the truth API and to every mirror. Each mirror answer is diffed against the truth's. `README.md` is the human front page. `docs/config.md` is the config, kind, JSONPath and report reference.

## File map

- `cmd/api-mirror-delta/` — cobra CLI, one command per file: `main.go` (root, shared flags, config loading), `serve.go`, `compare.go`.
- `cmd/fixtureserver/` — a test-only fake API that answers each path from one JSON file. The dats suite runs it.
- `internal/delta/config.go` — the YAML config, defaults, validation, ignore rules and their matching.
- `internal/delta/jsonpath.go` — concrete locations and the JSONPath pattern subset an ignore rule uses.
- `internal/delta/diff.go` — status, header and body comparison, and base-URL normalization.
- `internal/delta/delta.go` — the fan-out, the proxy handler, the `/_delta/` admin paths, and the drain.
- `internal/delta/report.go` — the log lines, the NDJSON file, the recent ring and the summary.
- `dats/delta.dats` — the built binaries end to end. `samples/github.yaml` — a starting config.

## Invariants

- **Only a request that reads fans out.** A write sent to the truth and to a mirror reaches the upstream twice. `compare` defaults to GET and HEAD. Anything else goes to the served target alone.
- **A rule names what it hides.** A rule with no path, header or kind is rejected, and the summary shows each rule's match count, zero included.
- **Nothing is dropped quietly.** A target that does not answer is an `error` difference. A full summary or ring counts and logs what it evicts. A cut value is marked `truncated`.
- **The comparison outlives the request.** Fetches run on a detached context, and `serve` drains them before it exits.
- **JSON is marshalled, never spliced.** Use `mustJSON`. Vet rejects a concatenated document.

## Commands

- `go-toolchain` — tidy, vet, test, build, then `dats/*.dats`. Never a bare `go` command, and never pipe its output.

## Conventions

- Tabs for indentation in Go and dats. Tests use testify.
- Comments state the invariant, in short sentences. No history.
