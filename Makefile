db-up:
	docker compose up -d 

db-clean:
	docker compose down
	sudo rm -rf db/data

db-reset: db-clean db-up

start-rcd:
	rclone rcd --rc-serve --rc-no-auth &

# VERSION is the release version embedded into the binary. Override on the
# command line (e.g. `make VERSION=v1.2.3 build`) or via the environment. When
# unset, fall back to `git describe` so local builds remain identifiable.
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
# Strip a leading `v` so the embedded Go version is bare semver, matching the
# convention used by published release artifacts.
GO_VERSION := $(VERSION:v%=%)
# `-s -w` strips the symbol table and DWARF debug info to keep release
# artifacts small. Use `go build` directly when a debuggable binary is needed.
LDFLAGS := -s -w -X main.version=$(GO_VERSION)

.PHONY: db-up db-clean db-reset start-rcd build install

build:
	go build -ldflags '$(LDFLAGS)' -o build/filesprawl ./cmd/filesprawl

install:
	go install -ldflags '$(LDFLAGS)' ./cmd/filesprawl
