# hbactl - version and build
VERSION ?= 0.1.1

BINARY   = hbactl
MAIN     = .
LDFLAGS  = -s -w -X github.com/hrodrig/hbactl/cmd.Version=$(VERSION)

.PHONY: build test clean install release

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) $(MAIN)

test:
	go test ./...

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
