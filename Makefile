# Go parameters
GOCMD=go
GOPATH=$(shell go env GOPATH)
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install
GOVET=$(GOCMD) vet
GOLINT=$(GOPATH)/bin/golangci-lint

# Binary name
BINARY_PATH=./cmd/main.go

# Build info
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME)"

.PHONY: build run test clean lint swagger docker-build deploy help

# Default target
all: build

# Help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

# Build the application
build: ## Build the binary
	$(GOBUILD) $(LDFLAGS) -o bin/server $(BINARY_PATH)

# Run the application
run: ## Run the server
	$(GOBUILD) $(LDFLAGS) -o bin/server $(BINARY_PATH)
	./bin/server server

# Run tests
test: ## Run tests
	$(GOTEST) $(LDFLAGS) -v -race -coverprofile=coverage.out ./internal/... ./pkg/... ./cmd/...

# Clean build artifacts
clean: ## Clean build artifacts
	rm -rf bin/ coverage.out

# Lint code
lint: ## Run linter
	$(GOLINT) run ./internal/... ./pkg/... ./cmd/...

# Generate Swagger docs
swagger: ## Generate Swagger docs
	$(GOPATH)/bin/swag init -d ./ -g cmd/main.go -o docs

# Build Docker image
docker-build: ## Build Docker image
	docker build -t go-demo:latest .

# Kubernetes 部署（示例）
deploy:
	kubectl apply -f deploy/k8s/go-demo.yaml

# 安装项目依赖
install-deps:
	$(GOGET) -v ./internal/... ./pkg/... ./cmd/... ./config/...
	$(GOINSTALL) -v github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOINSTALL) -v github.com/air-verse/air@latest

# 交叉编译（Linux）
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/server_linux -v $(BINARY_PATH)
