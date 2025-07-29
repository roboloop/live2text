MAKEFILE_DIR := $(patsubst %/,%,$(dir $(realpath $(lastword $(MAKEFILE_LIST)))))
BIN_DIR := $(MAKEFILE_DIR)/bin
export PATH := $(BIN_DIR):$(PATH)
export GOBIN := $(BIN_DIR)

GOLANGCI_LINT := $(BIN_DIR)/golangci-lint
MINIMOCK := $(BIN_DIR)/minimock
GO_COVER_TREEMAP := $(BIN_DIR)/go-cover-treemap

LIVE2TEXT_APP := ./cmd/live2text/main.go

GOLANGCI_LINT_VERSION := v2.1.6
MINIMOCK_VERSION := v3.4.5
GO_COVER_TREEMAP_VERSION := v1.5.0

.PHONY: all
all: install build

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

$(GOLANGCI_LINT): | $(BIN_DIR)
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

$(MINIMOCK): | $(BIN_DIR)
	go install github.com/gojuno/minimock/v3/cmd/minimock@$(MINIMOCK_VERSION)

$(GO_COVER_TREEMAP): | $(BIN_DIR)
	go install github.com/nikolaydubina/go-cover-treemap@$(GO_COVER_TREEMAP_VERSION)

.PHONY: install
install: check-portaudio $(GOLANGCI_LINT) $(MINIMOCK)

.PHONY: build
build:
	go build -trimpath -ldflags "-s -w" -o $(BIN_DIR)/live2text $(LIVE2TEXT_APP)

.PHONY: generate
generate: generate-mocks

.PHONY: generate-mocks
generate-mocks: $(MINIMOCK)
	find ./internal/ -iname '*_mock.go' -delete
	@echo "Using minimock from: $(MINIMOCK)"
	go generate --run "minimock" ./...

.PHONY: lint
lint: $(GOLANGCI_LINT)
	@echo "Using golangci-lint from: $(GOLANGCI_LINT)"
	$(GOLANGCI_LINT) run ./...

.PHONY: lint-fix
lint-fix: $(GOLANGCI_LINT)
	@echo "Using golangci-lint from: $(GOLANGCI_LINT)"
	$(GOLANGCI_LINT) run ./... --fix

.PHONY: format
format: $(GOLANGCI_LINT)
	@echo "Using golangci-lint from: $(GOLANGCI_LINT)"
	$(GOLANGCI_LINT) fmt ./...

.PHONY: test
test:
	@#go test -race ./... -v
	go test ./... -v

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)

.PHONY: check-portaudio
check-portaudio:
	@if test $$(uname) == "Darwin" && ! pkg-config -exists portaudio-2.0; then \
		echo "PortAudio not found via pkg-config."; \
		echo "Please install it: brew install portaudio"; \
	fi

.PHONY: show-coverage
show-coverage: $(GO_COVER_TREEMAP)
	@# hack that prevents multiple calls of `mktemp` command
	@$(MAKE) _show-coverage TMP_DIR=$$(mktemp -d)

.PHONY: _show-coverage
_show-coverage:
	go test -coverprofile $(TMP_DIR)/coverage.out ./...
	grep -vE '_(mock|string)\.go' $(TMP_DIR)/coverage.out > $(TMP_DIR)/filtered.out
	@echo "Using go-cover-treemap from: $(GO_COVER_TREEMAP)"
	$(GO_COVER_TREEMAP) -percent -h 1280 -w 2048 -coverprofile $(TMP_DIR)/filtered.out > $(TMP_DIR)/filtered.svg
	@echo "Check out the $(TMP_DIR)/filtered.svg file"
