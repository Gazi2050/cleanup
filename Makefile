.PHONY: dev build install version test vet tidy clean help

BINARY  := cleanup
SRC     := ./cmd/cleanup
PREFIX  := /usr/local/bin
VERSION ?= dev
LDFLAGS := -s -w -X main.version=$(VERSION)

## dev: run the TUI directly without building a binary (fastest iteration)
dev:
	go run $(SRC)

## build: compile a stripped binary into ./$(BINARY)
build:
	go build -ldflags="$(LDFLAGS)" -o $(BINARY) $(SRC)

## install: build and move $(BINARY) into $(PREFIX) (needs sudo)
install: build
	sudo install -m 0755 $(BINARY) $(PREFIX)/$(BINARY)

## version: print the version a build would get
version:
	@echo $(VERSION)

## test: run all tests
test:
	go test ./...

## vet: run go vet
vet:
	go vet ./...

## tidy: sync go.mod/go.sum
tidy:
	go mod tidy

## clean: remove the local binary
clean:
	rm -f $(BINARY)

## help: list targets
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## //; s/: /\t/'

.DEFAULT_GOAL := help
