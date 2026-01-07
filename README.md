# Go Demo Playground

一个使用 Go 语言开发的现代化 Web 服务示例项目，遵循 Go 社区最佳实践，集成了完整的开发、测试、构建和部署流程。

本项目展示了如何构建一个生产级就绪的 Go 微服务，包含了**结构化日志**、**Prometheus 监控**、**Pprof 性能分析**、**Swagger 文档**、**Docker 容器化**以及**CLI 命令行工具**等核心功能。

## 🏗 项目结构

```
├── cmd/                # 应用程序入口
│   ├── commands/       # Cobra 命令行命令定义 (root, server, version)
│   └── main.go         # 程序主入口
├── config/             # 配置定义与加载 (Viper)
├── docs/               # Swagger 自动生成的文档
├── internal/           # 内部业务逻辑 (Clean Architecture)
│   ├── handler/        # HTTP 请求处理层
│   ├── math/           # 核心业务逻辑 (示例：数学运算)
│   └── middleware/     # HTTP 中间件 (CORS, Gzip, RateLimit, etc.)
├── pkg/                # 通用工具包
│   ├── errors/         # 自定义错误处理
│   ├── logger/         # 结构化日志 (Zap)
│   └── response/       # 统一响应格式
├── go.mod              # 依赖管理
├── Makefile            # 自动化任务脚本
└── Dockerfile          # 多阶段构建 Dockerfile
```

## ✨ 功能特性

- **RESTful API**: 清晰的路由设计与统一的 JSON 响应格式
- **命令行工具**: 基于 Cobra 构建，支持子命令 (`server`, `version`)
- **中间件生态**:
  - **Request Logger**: 携带 Request ID 的全链路日志追踪
  - **Recovery**: 恐慌捕获与优雅恢复
  - **Prometheus Metrics**: 标准的 `/metrics` 监控端点
  - **Rate Limiting**: 基于 IP 的令牌桶限流
  - **Gzip**: 智能响应压缩
  - **CORS**: 灵活的跨域资源共享配置
- **可观测性**:
  - **Swagger UI**: 自动生成交互式 API 文档 (`/swagger/index.html`)
  - **Health Check**: K8s 就绪/存活探针 (`/healthz`)
  - **Pprof**: 实时性能剖析 (`/debug/pprof/`)
- **工程化实践**:
  - **配置管理**: 支持 YAML 文件与环境变量 (`APP_` 前缀)
  - **优雅停机**: 捕捉信号，等待请求处理完成
  - **依赖注入**: 模块化设计，易于测试
- **质量保证**:
  - **单元测试**: 基于 Testify，支持并发竞态检测 (`-race`)
  - **静态分析**: 严格的 GolangCI-Lint 配置
  - **代码覆盖率**: 自动生成 HTML 报告

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

## 🛠 开发指南

### 常用命令 (Makefile)

```bash
# 运行所有测试 (包含 Race Detector)
make test

# 生成代码覆盖率报告 (coverage.html)
make coverage

# 执行静态代码检查
make lint

# 自动格式化代码
make fmt

# 开发模式（热重载）
make dev

# 更新 Swagger 文档
make swagger

# 编译二进制文件
make build
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
- **测试**: [Testify](https://github.com/stretchr/testify)
- **热重载**: [Air](https://github.com/air-verse/air)
- **文档**: [Swag](https://github.com/swaggo/swag)
- **Lint**: [GolangCI-Lint](https://golangci-lint.run/)

## 📄 许可证

MIT
