BIN_DIR := ./bin
GO := go

GOLANGCI_LINT := $(BIN_DIR)/golangci-lint
MOCKERY := $(BIN_DIR)/mockery
GOOSE := $(BIN_DIR)/goose

export PATH := $(BIN_DIR):$(PATH)

# ==========================
#    INSTALL TOOLS
# ==========================

deps: $(GOLANGCI_LINT) $(GOOSE) $(MOCKERY)

$(GOLANGCI_LINT):
	@echo "Installing golangci-lint..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN_DIR) v2.4.0

$(GOOSE):
	@echo "Installing goose..."
	@mkdir -p $(BIN_DIR)
	@$(GO) install github.com/pressly/goose/v3/cmd/goose@latest
	@cp $$GOPATH/bin/goose $(GOOSE)

$(MOCKERY):
	@echo "Installing mockery..."
	@mkdir -p $(BIN_DIR)
	@$(GO) install github.com/vektra/mockery/v2@latest
	@cp $$GOPATH/bin/mockery $(MOCKERY)

# ==========================
#    LINT
# ==========================

lint:
	@echo "Running golangci-lint..."
	@$(GOLANGCI_LINT) run ./...

# ==========================
#    TESTS
# ==========================

test:
	@echo "Running unit tests..."
	@$(GO) test ./... -count=1 -race

test-integration:
	@echo "Running integration tests..."
	@$(GO) test ./... -count=1 -race -tags=integration

test-e2e:
	@echo "Running e2e tests..."
	@$(GO) test ./... -count=1 -tags=e2e

ci: deps lint test test-integration test-e2e
	@echo "CI done."

# ==========================
#    DOCKER
# ==========================

up:
	@docker compose up --build

down:
	@docker compose down -v

logs:
	@docker compose logs -f

# ==========================
#    UTILITIES
# ==========================

mock:
	@$(MOCKERY) --all --keeptree

.PHONY: deps lint test test-integration test-e2e ci up down logs mock
