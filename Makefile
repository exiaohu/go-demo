# Go parameters
GOCMD=go
GOPATH=$(shell go env GOPATH)
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOVET=$(GOCMD) vet
GOLINT=$(GOPATH)/bin/golangci-lint

# Binary name
BINARY_NAME=go-demo
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_PATH=./cmd/main.go

# Build info
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME)"

.PHONY: all test clean coverage lint fmt install-deps swagger help

all: test build

# 帮助信息
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build        Build the application"
	@echo "  run          Run the application"
	@echo "  dev          Run the application with hot reload (Air)"
	@echo "  release      Release the application (Snapshot)"
	@echo "  test         Run tests"
	@echo "  coverage     Generate coverage report"
	@echo "  lint         Run linter"
	@echo "  fmt          Format code"
	@echo "  swagger      Generate Swagger docs"
	@echo "  clean        Clean build artifacts"
	@echo "  install-deps Install dependencies"

# 开发模式（使用 Air 热重载）
dev:
	$(GOPATH)/bin/air

# 发布（快照模式，不推送到 GitHub）
release:
	$(GOPATH)/bin/goreleaser release --snapshot --clean

# 编译项目
build:
	$(GOBUILD) $(LDFLAGS) -v -o $(BINARY_NAME) -v $(BINARY_PATH)

# 运行项目
run:
	$(GOCMD) run $(BINARY_PATH) server

# 生成 Swagger 文档
swagger:
	$(GOPATH)/bin/swag init -g internal/handler/handler.go -o docs

# 运行单元测试并生成代码覆盖率报告（包含竞态检测）
test:
	$(GOTEST) -v -race -coverprofile=coverage.out ./internal/... ./pkg/... ./cmd/...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

# 运行代码 lint 检查
lint:
	$(GOLINT) run ./internal/... ./pkg/... ./cmd/...

# 格式化代码
fmt:
	$(GOCMD) fmt ./internal/... ./pkg/... ./cmd/... ./config/...

# 检查代码中的常见问题
vet:
	$(GOVET) ./internal/... ./pkg/... ./cmd/... ./config/...

# 清理编译产物和临时文件
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f coverage.out
	rm -f coverage.html

# 安装项目依赖
install-deps:
	$(GOGET) -v ./internal/... ./pkg/... ./cmd/... ./config/...
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/air-verse/air@latest

# 交叉编译（Linux）
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_UNIX) -v $(BINARY_PATH)
