BINARY  := gitgogit
PREFIX  ?= /usr/local
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
DIST    := dist

# Strip debug info (symbol table and DWARF) and embed version
GO_FLAGS += -ldflags="-s -w -X 'main.version=$(VERSION)'"

# Avoid embedding build path in executable
GO_FLAGS += -trimpath

.PHONY: build install system-install system-uninstall test clean platform-all platform-unixlike platform-darwin-arm64 platform-darwin-amd64 platform-linux-amd64 platform-linux-arm64 release

build:
	go build -o $(BINARY) .

## install puts the binary on PATH via GOPATH/bin (no sudo required).
install:
	go install .

## system-install copies the binary to $(PREFIX)/bin (may require sudo).
system-install: build
	install -d $(PREFIX)/bin
	install -m 755 $(BINARY) $(PREFIX)/bin/$(BINARY)

## system-uninstall removes the binary from $(PREFIX)/bin.
system-uninstall:
	rm -f $(PREFIX)/bin/$(BINARY)

test:
	go test ./...

clean:
	rm -f $(BINARY)
	rm -rf $(DIST)

  ##################
 # Platform Build #
##################

platform-all:
	@echo "Building all platform binaries..."
	@mkdir -p $(DIST)
	@$(MAKE) --no-print-directory -j4 \
		platform-darwin-arm64 \
		platform-darwin-amd64 \
		platform-linux-amd64 \
		platform-linux-arm64

platform-unixlike:
	@test -n "$(TGOOS)"   || (echo "GOOS must be set"  && false)
	@test -n "$(TGOARCH)" || (echo "GOARCH must be set" && false)
	CGO_ENABLED=0 GOOS="$(TGOOS)" GOARCH="$(TGOARCH)" \
		go build $(GO_FLAGS) -o "$(DIST)/$(BINARY)-$(TGOOS)-$(TGOARCH)"

platform-darwin-arm64:
	@$(MAKE) --no-print-directory TGOOS=darwin TGOARCH=arm64 platform-unixlike

platform-darwin-amd64:
	@$(MAKE) --no-print-directory TGOOS=darwin TGOARCH=amd64 platform-unixlike

platform-linux-amd64:
	@$(MAKE) --no-print-directory TGOOS=linux TGOARCH=amd64 platform-unixlike

platform-linux-arm64:
	@$(MAKE) --no-print-directory TGOOS=linux TGOARCH=arm64 platform-unixlike

  ###########
 # Release #
###########

release: platform-all
	@echo "Checking for gh CLI..." && gh --version > /dev/null
	@echo "Checking for a version tag..." && \
		git describe --tags --exact-match HEAD > /dev/null 2>&1 || \
		(echo "HEAD is not tagged — run: git tag vX.Y.Z" && false)
	@echo "Checking for uncommitted changes..." && \
		test -z "`git status --porcelain`" || \
		(echo "Cannot release with uncommitted changes:" && git status --porcelain && false)
	@echo "Checking for main branch..." && \
		test main = "`git rev-parse --abbrev-ref HEAD`" || \
		(echo "Cannot release from non-main branch" && false)
	@echo "Checking for unpushed commits..." && git fetch && \
		test "" = "`git cherry`" || (echo "Cannot release with unpushed commits" && false)
	gh release create "$(VERSION)" $(DIST)/$(BINARY)-* \
		--title "$(VERSION)" \
		--generate-notes
