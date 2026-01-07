package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/exiaohu/go-demo/pkg/logger"
)

type contextKey string

const (
	// RequestIDKey 上下文中 Request ID 的键
	RequestIDKey contextKey = "request_id"
	// HeaderXRequestID HTTP Header 中的 Request ID 键
	HeaderXRequestID = "X-Request-ID"
)

// RequestID 为每个请求生成唯一的 ID
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 优先从请求头获取 Request ID
		reqID := r.Header.Get(HeaderXRequestID)
		if reqID == "" {
			reqID = uuid.New().String()
		}

		// 将 Request ID 设置到响应头
		w.Header().Set(HeaderXRequestID, reqID)

		// 将 Request ID 注入到 Context
		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)

		// 传递带有 Request ID 的 Context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID 从 Context 中获取 Request ID
func GetRequestID(ctx context.Context) string {
	if val, ok := ctx.Value(RequestIDKey).(string); ok {
		return val
	}
	return ""
}

// WithLogger 注入带有 Request ID 的 logger 到 Context（可选增强）
func WithLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// LoggerMiddleware 重写 RequestLogger 以包含 Request ID
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := wrapResponseWriter(w)

		next.ServeHTTP(wrapped, r)

		reqID := GetRequestID(r.Context())

		logger.Info("HTTP Request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", wrapped.status),
			zap.Duration("duration", time.Since(start)),
			zap.String("ip", r.RemoteAddr),
			zap.String("request_id", reqID),
			zap.String("user_agent", r.UserAgent()),
		)
	})
}
