# hbactl - version and build
# Version from VERSION file (single source of truth); override: make build VERSION=v0.2.0
VERSION   ?= $(shell v=$$(cat VERSION 2>/dev/null | tr -d '\n\r'); [ -n "$$v" ] && echo "v$$v" || echo "v0.1.0")
BINARY    = hbactl
MAIN      = .
LDFLAGS   = -s -w -X github.com/hrodrig/hbactl/cmd.Version=$(VERSION)

.PHONY: build test clean install release snapshot docker-build docker-scan lint lint-fix check

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) $(MAIN)

test:
	go test ./...

# Full pre-merge/release check: verify deps, build, test, lint, security scan (govulncheck + optional Grype on dir)
check:
	go mod verify && go build ./... && go test ./... && $(MAKE) lint && ./tools/scan.sh

# Lint: gofmt + gocyclo (run during development; CI runs this too)
lint:
	@echo "Checking gofmt -s..."
	@unformatted=$$(gofmt -s -l .); [ -z "$$unformatted" ] || { echo "Files not formatted (run make lint-fix):"; echo "$$unformatted"; exit 1; }
	@echo "Checking gocyclo (complexity <= 14)..."
	@command -v gocyclo >/dev/null 2>&1 || go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	@gocyclo -over 14 .

# Fix formatting only (gofmt -s -w); re-run make lint to verify gocyclo
lint-fix:
	gofmt -s -w .

clean:
	rm -f $(BINARY)

install: build
	go install -ldflags "$(LDFLAGS)" $(MAIN)

# Release: only from main. Merge develop → main, then: git tag v0.1.0 && make release
# Requires: brew install goreleaser
release:
	@branch=$$(git branch --show-current 2>/dev/null); \
	if [ "$$branch" != "main" ]; then \
		echo "Error: release only from main (current: $$branch). Merge and checkout main first."; \
		exit 1; \
	fi
	goreleaser release --clean

# Snapshot build (no tag required), outputs to dist/
snapshot:
	goreleaser release --snapshot --clean

# Docker image (VERSION from file; override: make docker-build VERSION=v0.1.10)
docker-build:
	docker build --build-arg VERSION=$(VERSION) -t hbactl .

# Build image as hbactl:scan and run Grype (--fail-on high: fails on high and critical). Requires: docker, grype on PATH.
docker-scan:
	@command -v grype >/dev/null 2>&1 || { echo "grype not found; install with: brew install grype or https://github.com/anchore/grype#installation"; exit 1; }
	docker build --build-arg VERSION=$(VERSION) -t hbactl:scan .
	grype hbactl:scan --fail-on high
