GO = go
GOFLAGS = -v

BUILD_DIR = build
CONFIG_DIR = $(PWD)/../production/config

BINARY_NAME = denon-to-mqtt-gw
DOCKER_IMAGE_NAME = pharndt/denon-to-mqtt-gw:latest

.PHONY: build
build:
	@echo "Building..."
	@$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

.PHONY: run
run: build
	@echo "Running the application..."
	@$(BUILD_DIR)/$(BINARY_NAME) $(CONFIG_DIR)/config.json

.PHONY: docker
docker:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE_NAME) .

.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

.PHONY: deps
deps:
	@echo "Downloading Go dependencies..."
	@$(GO) mod download
	@$(GO) mod tidy

.PHONY: lint
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

.PHONY: test
test:
	@echo "Running tests..."
	@$(GO) test ./...
