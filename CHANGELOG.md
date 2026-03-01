# Changelog

All notable changes to this project are documented in this file.

Format based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added

- (Nothing yet.)

---

## [0.2.2] - 2026-03-01

### Added

- **Makefile:** `make check` (go mod verify, build, test, lint, scan.sh) for pre-merge/release; `make docker-scan` (build hbactl:scan + Grype --fail-on high). Documented in tools/README.md, AGENTS.md, release-tests.mdc, git-flow.mdc.
- **Docker:** Dockerfile.release for GoReleaser dockers_v2 (multi-arch image to ghcr.io/hrodrig/hbactl); layout linux/<arch>/hbactl for snapshot/release. .gitignore Dockerfile.release.

### Changed

- **Go:** go.mod and README badge to Go 1.26; indirect deps aligned with pgwd (golang.org/x/sync v0.19.0, x/text v0.34.0).
- **Grype:** `--fail-on high` (single severity; high and critical). Makefile, CI, tools/README updated (was invalid `high,critical`).
- **README:** Published image ghcr.io/hrodrig/hbactl, docker run examples; Go 1.26 badge.

---

## [0.2.0] - 2026-02-28

### Added

- **list:** `--no-index` to omit the rule index column; output is copy-paste friendly. Header and separator lines are prefixed with `# ` so they become pg_hba.conf comments when pasted (aligned so TYPE lines up with data). Works with `--group-by user` (group separators `# === user: X ===` also commented).
- **Docker & security:** Dockerfile (multi-stage Go/Alpine, non-root user), `make docker-build`, `.dockerignore`. Security workflow (`.github/workflows/security.yml`): lint (gofmt + gocyclo), govulncheck, Grype on built image only (aligned with pgwd). [tools/README.md](tools/README.md) and `tools/scan.sh` for pre-merge/release scanning.
- **Community & docs:** [CONTRIBUTING.md](CONTRIBUTING.md), [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) (Contributor Covenant 2.1), this CHANGELOG. Cursor rules: changelog, diagrams-mermaid, git-flow, gitignore-whitelist, language-english; release-tests and english-only updated.

### Changed

- **README:** Badges and structure aligned with pgwd (Version, Release, Go, License, pkg.go.dev, Go Report Card; Repo/Releases line; Documentation line; Docker section; License footer with CONTRIBUTING, CODE_OF_CONDUCT, CHANGELOG).
- **Makefile:** Added `lint` (gofmt -s + gocyclo -over 14), `lint-fix`, `docker-build`; `.PHONY` updated.
- **AGENTS.md:** Layout mentions tools/; release section mentions security scan; docs mention CONTRIBUTING, CODE_OF_CONDUCT, CHANGELOG, Docker.
- **cmd/add.go:** Refactored into normalizeAddFlags, resolveAddPath, printAddDryRun, doAdd, runAdd (clearer flow, testable).
- **cmd/remove.go:** Refactored into validateRemoveFlags (removeMode), resolveRemovePath, findRulesToRemove, removeCriteriaString, doRemove, runRemove (clearer flow, testable).
- **internal/cli/output.go:** WriteRulesTableNoIndex, WriteRulesTableGroupedByUserNoIndex, writeRulesTableNoIndexTo (copy-paste output with `# TYPE` header alignment).
- **internal/hba/parser.go:** Whitespace alignment in hostTypes map.
- **internal/hba/parser_test.go:** Extracted checkRule0–checkRule3 helpers for readability.
- **.gitignore:** Whitelist for CONTRIBUTING.md, CODE_OF_CONDUCT.md, CHANGELOG.md, .github/, tools/, Dockerfile, .dockerignore; .DS_Store.
- **Demo:** docs/demo.tape and docs/demo.gif updated (minor).

---

## [0.1.10] - 2026-02-28

### Added

- **Docker:** Dockerfile (multi-stage Go/Alpine), non-root user, `make docker-build`; README Docker section.
- **tools/:** Security scanning before merge/release: `tools/scan.sh` (govulncheck + optional Grype), [tools/README.md](tools/README.md). CI: `.github/workflows/security.yml` (govulncheck + Grype on dir and on built image). Release rule and AGENTS/README updated.
- **Community:** [CONTRIBUTING.md](CONTRIBUTING.md), [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) (Contributor Covenant 2.1), this CHANGELOG.

---

## [0.1.9] - 2026-02-16

### Added

- **Homebrew:** Cask for tap `hrodrig/hbactl` (install via `brew install hrodrig/hbactl/hbactl`).

---

## [0.1.8] - 2026-02-16

### Added

- **remove:** Criteria-based removal: `--user` (optional `--db`), `--addr`; remove by index unchanged.

---

## [0.1.4] - 2026-02-14

### Added

- **list:** `--group-by user` (separators between users); **add:** `--after-user <name>` (insert after last rule for that user).
- **add/remove:** Dry-run output now describes the action (e.g. "would insert after last rule for user …").
- **Docs:** Sequence diagrams (docs/), demo GIF (VHS), anonymized sample config.

---

## [0.1.1] - 2026-02-14

### Added

- **Packaging:** .deb and .rpm via GoReleaser nfpm.

### Changed

- **Version:** Consistent `v` prefix in Makefile, README badge, and release asset names.

---

## [0.1.0] - 2026-02-14

### Added

- **CLI:** Commands list, add, remove (by index), check, reload. Auto-discovery of `pg_hba.conf` via Postgres or `--file`.
- **Safety:** Backup before edits; validation with `hbactl check` (pg_hba_file_rules).
- **Releases:** GoReleaser (binaries, checksums); README, Cursor rules (version badge, english-only).

---

[Unreleased]: https://github.com/hrodrig/hbactl/compare/v0.2.2...HEAD
[0.2.2]: https://github.com/hrodrig/hbactl/compare/v0.2.0...v0.2.2
[0.2.0]: https://github.com/hrodrig/hbactl/compare/v0.1.10...v0.2.0
[0.1.10]: https://github.com/hrodrig/hbactl/compare/v0.1.9...v0.1.10
[0.1.9]: https://github.com/hrodrig/hbactl/compare/v0.1.8...v0.1.9
[0.1.8]: https://github.com/hrodrig/hbactl/compare/v0.1.4...v0.1.8
[0.1.4]: https://github.com/hrodrig/hbactl/compare/v0.1.1...v0.1.4
[0.1.1]: https://github.com/hrodrig/hbactl/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/hrodrig/hbactl/releases/tag/v0.1.0
