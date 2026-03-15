BINARY  := gitgogit
PREFIX  ?= /usr/local

.PHONY: build install uninstall test clean

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
