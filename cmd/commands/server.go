package commands

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"

	_ "github.com/exiaohu/go-demo/docs" // swagger docs
	"github.com/exiaohu/go-demo/internal/handler"
	"github.com/exiaohu/go-demo/internal/middleware"
	"github.com/exiaohu/go-demo/internal/model"
	"github.com/exiaohu/go-demo/internal/repository"
	"github.com/exiaohu/go-demo/internal/service"
	"github.com/exiaohu/go-demo/pkg/database"
	"github.com/exiaohu/go-demo/pkg/logger"
	"github.com/exiaohu/go-demo/pkg/tracer"
)

func newServerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Start the HTTP server",
		Run: func(_ *cobra.Command, _ []string) {
			runServer()
		},
	}
}

func runServer() {
	cfg, err := loadConfigHelper()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		_ = logger.Sync()
	}()

	// 初始化 Tracer
	shutdownTracer, err := tracer.InitTracer(cfg.AppName)
	if err != nil {
		logger.Fatal("Failed to initialize tracer", zap.Error(err))
	}
	defer func() {
		if err := shutdownTracer(context.Background()); err != nil {
			logger.Error("Failed to shutdown tracer", zap.Error(err))
		}
	}()

	logger.Info("Application starting",
		zap.String("app_name", cfg.AppName),
		zap.String("version", cfg.Version),
	)

	// 初始化数据库
	if err := database.Initialize(cfg); err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	// 自动迁移
	if err := database.AutoMigrate(&model.CalculationHistory{}); err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}

	// 依赖注入
	historyRepo := repository.NewHistoryRepository(database.DB)
	calcService := service.NewCalculatorService(historyRepo)
	h := handler.NewHandler(calcService)

	// 创建 HTTP 服务器
	router := http.NewServeMux()
	router.HandleFunc("/", h.HomeHandler)
	router.HandleFunc("/healthz", h.HealthCheckHandler)
	router.Handle("/metrics", promhttp.Handler())

	// Pprof (Debug 模式开启)
	if cfg.Debug {
		logger.Info("Pprof enabled at /debug/pprof/")
		router.HandleFunc("/debug/pprof/", pprof.Index)
		router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		router.HandleFunc("/debug/pprof/profile", pprof.Profile)
		router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		router.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	// Swagger 文档
	router.Handle("/swagger/", httpSwagger.WrapHandler)

	// API v1 路由组
	v1 := http.NewServeMux()
	v1.HandleFunc("/add", h.AddHandler)
	v1.HandleFunc("/subtract", h.SubtractHandler)
	v1.HandleFunc("/multiply", h.MultiplyHandler)
	v1.HandleFunc("/divide", h.DivideHandler)
	v1.HandleFunc("/history", h.HistoryHandler)

	// 注册 v1 路由，同时保留根路径以兼容旧版本（可选）
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))
	// 兼容旧路由
	router.HandleFunc("/add", h.AddHandler)
	router.HandleFunc("/subtract", h.SubtractHandler)
	router.HandleFunc("/multiply", h.MultiplyHandler)
	router.HandleFunc("/divide", h.DivideHandler)
	router.HandleFunc("/history", h.HistoryHandler)

	// 配置 CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // 生产环境请修改为具体域名
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            cfg.Debug,
	})

	// 应用中间件
	// 中间件执行顺序（从外到内）：
	// 1. RequestID: 生成请求 ID，方便追踪
	// 2. Logger: 记录请求日志（包括 panic 后的 500）
	// 3. Metrics: 记录监控指标（包括 panic 后的 500）
	// 4. Recovery: 捕获 panic，防止服务崩溃
	// 5. RateLimit: 限流
	// 6. Gzip: 响应压缩
	handler := middleware.Chain(router,
		corsHandler.Handler,
		middleware.RequestID,
		middleware.LoggerMiddleware,
		middleware.Metrics,
		middleware.Recovery,
		middleware.RateLimit,
		middleware.Gzip,
	)

	// 包装 Tracer Middleware
	tracedHandler := otelhttp.NewHandler(handler, "http-server")

	// 启动服务器
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           tracedHandler,
		ReadHeaderTimeout: 3 * time.Second,
	}

	// 启动服务器的 goroutine
	go func() {
		logger.Info("Server starting", zap.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	// 等待异步任务完成
	if err := calcService.Close(); err != nil {
		logger.Error("Failed to close calculator service", zap.Error(err))
	}

	// 关闭数据库连接
	if err := database.Close(); err != nil {
		logger.Error("Failed to close database connection", zap.Error(err))
	} else {
		logger.Info("Database connection closed")
	}

	logger.Info("Server exited gracefully")
}
