package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/exiaohu/go-demo/pkg/logger"
)

// responseWriter 包装 http.ResponseWriter 以捕获状态码
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

// RequestLogger 记录每个请求的日志
// Deprecated: Use LoggerMiddleware in request_id.go instead for Request ID support
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := wrapResponseWriter(w)

		next.ServeHTTP(wrapped, r)

		logger.Info("HTTP Request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", wrapped.status),
			zap.Duration("duration", time.Since(start)),
			zap.String("ip", r.RemoteAddr),
			zap.String("user_agent", r.UserAgent()),
		)
	})
}

// Recovery 恢复 panic 并记录错误
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", r.URL.Path),
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Chain 链式应用中间件
func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
