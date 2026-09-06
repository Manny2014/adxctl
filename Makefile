BINARY_NAME ?= adxctl
BIN_DIR     ?= bin
DIST_DIR    ?= dist

# Compiler and toolchain detection: prefer go1.27.1, fallback to system go
GO ?= $(shell GOPATH=$$(go env GOPATH 2>/dev/null); \
	if [ -x "$$GOPATH/bin/go1.27.1" ]; then echo "$$GOPATH/bin/go1.27.1"; \
	elif command -v go1.27.1 >/dev/null 2>&1; then command -v go1.27.1; \
	else command -v go; fi)

# Metadata for build stamping
VERSION    ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "0.1.0-dev")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
BUILD_DATE ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
MODULE     ?= $(shell $(GO) list -m 2>/dev/null || echo "adxctl")

# Production-grade compiler flags:
# -trimpath: removes host file paths from debug traces and stack traces for privacy and reproducible builds
# -s: omit symbol table
# -w: omit DWARF debugging information (shrinks binary size significantly)
# CGO_ENABLED=0: produces static, portable binaries with no runtime C library dependencies
LDFLAGS = -s -w \
	-X '$(MODULE)/cmd.Version=$(VERSION)' \
	-X '$(MODULE)/cmd.GitCommit=$(GIT_COMMIT)' \
	-X '$(MODULE)/cmd.BuildDate=$(BUILD_DATE)'

GO_BUILD = CGO_ENABLED=0 $(GO) build -trimpath -ldflags="$(LDFLAGS)"

# Supported cross-compilation architectures
PLATFORMS := \
	darwin/amd64 \
	darwin/arm64 \
	linux/amd64 \
	linux/arm64

.PHONY: all build build-all clean test vet tidy help $(PLATFORMS) nats-stream-info nats-stream-view nats-kv-state nats-kv-list

all: clean build

## help: Display available make targets
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build                Build optimized binary for current host (saved to $(BIN_DIR)/$(BINARY_NAME))"
	@echo "  build-all            Cross-compile optimized binaries for Darwin and Linux (saved to $(DIST_DIR)/)"
	@echo "  build-darwin-amd64   Build for macOS Intel (amd64)"
	@echo "  build-darwin-arm64   Build for macOS Apple Silicon (arm64)"
	@echo "  build-linux-amd64    Build for Linux x86_64 (amd64)"
	@echo "  build-linux-arm64    Build for Linux ARM64 (arm64)"
	@echo "  clean                Remove build and distribution artifacts"
	@echo "  vet                  Run go vet static analysis"
	@echo "  test                 Run tests"
	@echo "  tidy                 Run go mod tidy"
	@echo "  nats-stream-info     View details and statistics of the JIRA_EVENTS NATS stream"
	@echo "  nats-stream-view     View JIRA_EVENTS stream messages (override default with SUBJECT=...)"
	@echo "  nats-kv-state        View current poll state watermark in NATS KV store"
	@echo "  nats-kv-list         List all registered Key-Value buckets in NATS"

## build: Build optimized binary for current OS/ARCH
build:
	@echo "==> Building $(BINARY_NAME) for host with production optimizations..."
	@mkdir -p $(BIN_DIR)
	$(GO_BUILD) -o $(BIN_DIR)/$(BINARY_NAME) .
	@echo "==> Binary created at $(BIN_DIR)/$(BINARY_NAME)"

## build-all: Cross-compile for all Darwin and Linux architectures
build-all: build-darwin-amd64 build-darwin-arm64 build-linux-amd64 build-linux-arm64
	@echo "==> All Darwin and Linux binaries built in $(DIST_DIR)/"
	@ls -lh $(DIST_DIR)

## Individual platform build targets
build-darwin-amd64:
	@echo "==> Building for darwin/amd64..."
	@mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 $(GO_BUILD) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 .

build-darwin-arm64:
	@echo "==> Building for darwin/arm64..."
	@mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=arm64 $(GO_BUILD) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 .

build-linux-amd64:
	@echo "==> Building for linux/amd64..."
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 $(GO_BUILD) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 .

build-linux-arm64:
	@echo "==> Building for linux/arm64..."
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=arm64 $(GO_BUILD) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 .

## clean: Remove build artifacts
clean:
	@echo "==> Cleaning $(BIN_DIR) and $(DIST_DIR)..."
	@rm -rf $(BIN_DIR) $(DIST_DIR)

## vet: Run go vet
vet:
	@echo "==> Running go vet..."
	$(GO) vet ./...

## test: Run tests
test:
	@echo "==> Running tests..."
	$(GO) test -v ./...

## tidy: Run go mod tidy
tidy:
	@echo "==> Tidying dependencies..."
	$(GO) mod tidy

nats-start:
	./bin/adxctl nats start --runner docker --config config.example.yaml   

nats-stop:
	./bin/adxctl nats stop --runner docker --config config.example.yaml   

## nats-stream-info: Display info of JIRA_EVENTS stream
nats-stream-info:
	nats stream info JIRA_EVENTS

## nats-stream-view: View stream messages (optional: SUBJECT="jira.*.KAN")
SUBJECT ?= "jira.*.*"
nats-stream-view:
	nats stream view JIRA_EVENTS --subject $(SUBJECT) --raw

## nats-kv-state: View current watermark state in NATS KV
nats-kv-state:
	nats kv view jira_poller_state last_poll_timestamp

## nats-kv-list: List KV buckets in NATS
nats-kv-list:
	nats kv ls

