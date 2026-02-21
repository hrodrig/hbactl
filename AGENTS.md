# AGENTS.md

Context for AI coding agents working on **hbactl**: a Go CLI to manage PostgreSQL `pg_hba.conf` (list, add, remove, check, reload).

Ref: [agents.md](https://agents.md/)

---

## Project overview

- **Language:** Go; single binary, no CGO.
- **Layout:** `cmd/` (Cobra subcommands: list, add, remove, check, reload), `internal/hba` (parser, file, rules, sort), `internal/cli` (output), `internal/pg` (Postgres client).
- **Releases:** GoReleaser from `main` only. **VERSION** file (root) is the single source of truth; Makefile reads it, README badge must match (see `.cursor/rules/readme-badges-version.mdc`).

---

## Setup and build

- Install deps: `go mod download` (or `go build ./...`).
- Build: `make build` or `go build -o hbactl .` (version via `Makefile` LDFLAGS).
- Run: `./hbactl --help`; for commands that need a file, use `-f sample-pg_hba.conf` to avoid a live DB.

---

## Testing

- Run tests: `make test` or `go test ./...`.
- Unit tests do **not** require a running PostgreSQL server (parser, file, output tests are self-contained).
- Fix any failing test before finishing; add or update tests for changed behavior.

---

## Code style

- **Language:** All code (identifiers, comments, error messages, CLI output) must be in **English** (see `.cursor/rules/english-only.mdc`).
- Prefer standard library and existing patterns; Cobra for CLI, no extra style enforcer beyond `go vet` / `go test`.
- **Cobra/Viper CLI:** When adding or changing commands, flags, or CLI behavior, prefer following the [golang-cli-cobra-viper](https://skills.sh/bobmatnyc/claude-mpm-skills/golang-cli-cobra-viper) skill (e.g. use `RunE`/`PreRunE`, clear error messages, persistent vs local flags).

---

## Release and versioning

- Releases and tags are made **only from `main`**. Work on `develop` or feature branches, then merge to `main`.
- **All tests must pass before release:** run `make test` (or `go test ./...`) and fix any failure before tagging or running `make release`.
- Before release: bump the **VERSION** file (e.g. to `0.1.10`) and the version badge in `README.md` to match (e.g. `version-v0.1.10`).
- Release: `git tag v0.1.10 && make release` (requires goreleaser). Homebrew cask push requires `HOMEBREW_TAP_TOKEN` (PAT with `repo` scope) in the environment.

---

## Docs and assets

- User-facing docs: `README.md`. Sequence diagrams: `docs/README.md` and `docs/sequence-*.md`.
- Demo GIF: `docs/demo.gif` from `vhs docs/demo.tape` (run from repo root after `go build -o hbactl .`).

---

## Security and safety

- hbactl edits `pg_hba.conf` and may need elevated permissions; only use in trusted environments.
- Backup is created before file edits; no automated deployment or secrets in repo.
