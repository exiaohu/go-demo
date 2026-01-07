# Go Demo Playground

[![CI](https://github.com/exiaohu/go-demo/actions/workflows/ci.yml/badge.svg)](https://github.com/exiaohu/go-demo/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/exiaohu/go-demo)](https://goreportcard.com/report/github.com/exiaohu/go-demo)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

一个使用 Go 语言开发的现代化 Web 服务示例项目，遵循 Go 社区最佳实践，集成了完整的开发、测试、构建和部署流程。

本项目展示了如何构建一个生产级就绪的 Go 微服务，包含了**结构化日志**、**Prometheus 监控**、**Pprof 性能分析**、**Swagger 文档**、**Docker 容器化**以及**CLI 命令行工具**等核心功能。

## 🏗 项目结构

```text
├── cmd/                # 应用程序入口
│   ├── commands/       # Cobra 命令行命令定义 (root, server, version)
│   └── main.go         # 程序主入口
├── config/             # 配置定义与加载 (Viper)
├── deploy/             # 部署配置 (Kubernetes, Docker)
├── docs/               # Swagger 自动生成的文档
├── internal/           # 内部业务逻辑 (Clean Architecture)
│   ├── handler/        # HTTP 请求处理层
│   ├── math/           # 核心业务逻辑 (示例：数学运算)
│   ├── middleware/     # HTTP 中间件 (CORS, Gzip, RateLimit, etc.)
│   └── model/          # 数据模型定义
├── pkg/                # 通用工具包
│   ├── database/       # 数据库连接与工具
│   ├── errors/         # 自定义错误处理
│   ├── logger/         # 结构化日志 (Zap)
│   ├── response/       # 统一响应格式
│   └── util/           # 实用工具 (如 IP 获取)
├── go.mod              # 依赖管理
├── Makefile            # 自动化任务脚本
└── Dockerfile          # 多阶段构建 Dockerfile
```

## 特性

*   **Clean Architecture**: 清晰的代码结构，分离关注点。
*   **RESTful API**: 示例 API 实现（加减乘除）。
*   **SQLite Database**: 集成 GORM 和 **纯 Go SQLite 驱动** (无 CGO 依赖)，轻松跨平台编译。
*   **中间件**:
    *   Logger (Zap)
    *   Recovery
    *   Request ID
    *   Prometheus Metrics
    *   Gzip Compression
    *   Rate Limiting (Token Bucket)
    *   CORS
*   **配置管理**: 使用 Viper 加载配置。
*   **Swagger 文档**: 自动生成 API 文档。
*   **Docker 支持**: 基于 **Distroless** 的多阶段构建，生成极致轻量（~20MB）且安全的静态二进制镜像。
*   **优雅关闭**: 完善的信号处理与资源释放机制，确保异步任务（如日志记录）不丢失。
*   **Kubernetes 支持**: 完整的 Deployment 和 Service 配置。
*   **开发工具**:
    *   Makefile (常用命令)
    *   Air (热重载)
    *   Goreleaser (发布自动化)
    *   Pre-commit hooks (代码质量检查)

## 🚀 快速开始

### 本地运行

1. **安装依赖**
   ```bash
   make install-deps
   ```

2. **启动服务**
   ```bash
   # 使用 Makefile (推荐)
   make run

   # 或者使用 go run
   go run cmd/main.go server
   ```

3. **查看版本**
   ```bash
   go run cmd/main.go version
   ```

### Docker 运行

```bash
# 构建镜像
docker build -t go-demo .

# 启动容器
docker run -p 8080:8080 go-demo
```

### 访问服务

| 功能 | 地址 |
|------|------|
| **主页** | http://localhost:8080/ |
| **Swagger 文档** | http://localhost:8080/swagger/index.html |
| **健康检查** | http://localhost:8080/healthz |
| **监控指标** | http://localhost:8080/metrics |
| **性能分析** | http://localhost:8080/debug/pprof/ |
| **示例 API** | http://localhost:8080/add?a=1&b=2 |
| **计算历史** | http://localhost:8080/history |

## 🛠 开发指南

### 常用命令 (Makefile)

```bash
# 运行所有测试 (包含 Race Detector)
make test

# 执行静态代码检查
make lint

# 自动格式化代码
make fmt

# 开发模式（热重载）
make dev

# 发布快照（本地测试）
make release

# 更新 Swagger 文档
make swagger

# 编译二进制文件
make build

# 部署到 Kubernetes
make deploy
```

### 配置说明

配置文件位于 `config.yaml`，你也可以通过环境变量覆盖配置：

```yaml
app_name: "Playground"
version: "1.0.0"
port: 8080
debug: true
rate_limit:
  enabled: true
  rps: 10
  burst: 20
```

对应环境变量示例：`APP_PORT=9090`, `APP_DEBUG=false`

## 📦 技术栈

- **Web 框架**: 标准库 `net/http` + `ServeMux`
- **CLI**: [Cobra](https://github.com/spf13/cobra)
- **配置**: [Viper](https://github.com/spf13/viper)
- **日志**: [Zap](https://github.com/uber-go/zap)
- **ORM**: [GORM](https://gorm.io/) + [Pure Go SQLite](https://github.com/glebarez/sqlite)
- **测试**: [Testify](https://github.com/stretchr/testify)
- **热重载**: [Air](https://github.com/air-verse/air)
- **文档**: [Swag](https://github.com/swaggo/swag)
- **Lint**: [GolangCI-Lint](https://golangci-lint.run/)
- **Docker**: [Distroless](https://github.com/GoogleContainerTools/distroless)

## 📄 许可证

MIT
