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
	"go.uber.org/zap"

	_ "github.com/exiaohu/go-demo/docs" // swagger docs
	"github.com/exiaohu/go-demo/internal/handler"
	"github.com/exiaohu/go-demo/internal/middleware"
	"github.com/exiaohu/go-demo/pkg/logger"
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

	logger.Info("Application starting",
		zap.String("app_name", cfg.AppName),
		zap.String("version", cfg.Version),
	)

	// 创建 HTTP 服务器
	router := http.NewServeMux()
	router.HandleFunc("/", handler.HomeHandler)
	router.HandleFunc("/healthz", handler.HealthCheckHandler)
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

	router.HandleFunc("/add", handler.AddHandler)
	router.HandleFunc("/subtract", handler.SubtractHandler)
	router.HandleFunc("/multiply", handler.MultiplyHandler)
	router.HandleFunc("/divide", handler.DivideHandler)

	// 配置 CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // 生产环境请修改为具体域名
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            cfg.Debug,
	})

	// 应用中间件
	handler := middleware.Chain(router,
		middleware.Recovery,
		middleware.RequestID,
		middleware.LoggerMiddleware,
		middleware.Metrics,
		middleware.RateLimit, // 1. 先限流，防止恶意请求消耗资源
		middleware.Gzip,      // 2. 再压缩，只对允许通过的请求进行处理
	)

	// 包装 CORS
	finalHandler := corsHandler.Handler(handler)

	// 启动服务器
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           finalHandler,
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

	logger.Info("Server exited gracefully")
}
