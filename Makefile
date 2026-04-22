MODULE := github.com/3110Y/cap
BIN    := cap
BIN_DIR := bin
WIRE_DIR := internal/integration/di

GO   ?= go
WIRE ?= $(shell command -v wire 2>/dev/null)

.PHONY: init wire build build-all test vet fmt clean

init:
	$(GO) mod download
	$(GO) install github.com/google/wire/cmd/wire@latest

wire:
	@if [ -z "$(WIRE)" ]; then \
		echo "wire not found in PATH; run 'make init' first"; exit 1; \
	fi
	cd $(WIRE_DIR) && $(WIRE)

build:
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(BIN) ./cmd/$(BIN)

build-all:
	./scripts/build-all.sh

test:
	$(GO) test -race ./...

vet:
	$(GO) vet ./...

fmt:
	$(GO) fmt ./...

clean:
	rm -rf $(BIN_DIR)
