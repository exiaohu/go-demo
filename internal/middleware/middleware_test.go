package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/exiaohu/go-demo/config"
	"github.com/exiaohu/go-demo/pkg/logger"
)

func init() {
	_ = logger.Initialize(true)
	config.C = &config.Config{}
	config.C.RateLimit.Enabled = false
}

func TestRequestID(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := GetRequestID(r.Context())
		assert.NotEmpty(t, reqID, "Request ID should be in context")
		w.Header().Set("X-Test-Request-ID", reqID)
	})

	handler := RequestID(nextHandler)

	// Test with no Request ID header
	req1 := httptest.NewRequest("GET", "/", nil)
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, req1)
	assert.NotEmpty(t, w1.Header().Get(HeaderXRequestID))

	// Test with existing Request ID header
	req2 := httptest.NewRequest("GET", "/", nil)
	existingID := "existing-id-123"
	req2.Header.Set(HeaderXRequestID, existingID)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)
	assert.Equal(t, existingID, w2.Header().Get(HeaderXRequestID))
}

func TestRecovery(t *testing.T) {
	// 捕获日志输出以免污染测试结果（可选）
	// logger.SetLogger(...)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("something went wrong")
	})

	handler := Recovery(nextHandler)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// 确保 panic 不会通过 Recovery 中间件冒泡
	assert.NotPanics(t, func() {
		handler.ServeHTTP(w, req)
	})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGzip(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("hello world"))
	})

	handler := Gzip(nextHandler)

	// Test without Accept-Encoding: gzip
	req1 := httptest.NewRequest("GET", "/", nil)
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, req1)
	assert.Equal(t, "hello world", w1.Body.String())
	assert.Empty(t, w1.Header().Get("Content-Encoding"))

	// Test with Accept-Encoding: gzip
	req2 := httptest.NewRequest("GET", "/", nil)
	req2.Header.Set("Accept-Encoding", "gzip")
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)
	assert.Equal(t, "gzip", w2.Header().Get("Content-Encoding"))
	assert.NotEqual(t, "hello world", w2.Body.String()) // Should be compressed
}

func TestRateLimit(t *testing.T) {
	// Mock config
	// 注意：这里依赖全局配置，单元测试修改全局变量可能会有副作用
	// 更好的方式是 RateLimit 中间件接受配置参数，而不是直接读取 config.C
	// 但为了演示，我们先暂时这样测试，或者跳过依赖配置的部分

	// 由于 RateLimit 内部直接使用了 config.C，且初始化逻辑有 sync.Once，
	// 单元测试比较难覆盖不同配置的情况。
	// 这里仅测试默认未启用或未配置的情况，或者我们需要重构代码以支持测试。
	// 为了不破坏现有结构，暂时跳过复杂的限流测试，只测试基本透传。

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// 假设默认配置下（RateLimit.Enabled 可能为 false 或未初始化），请求应该通过
	handler := RateLimit(nextHandler)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChain(t *testing.T) {
	var steps []string
	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			steps = append(steps, "mw1 start")
			next.ServeHTTP(w, r)
			steps = append(steps, "mw1 end")
		})
	}
	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			steps = append(steps, "mw2 start")
			next.ServeHTTP(w, r)
			steps = append(steps, "mw2 end")
		})
	}

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		steps = append(steps, "handler")
	})

	handler := Chain(finalHandler, mw1, mw2)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	expected := []string{"mw1 start", "mw2 start", "handler", "mw2 end", "mw1 end"}
	assert.Equal(t, expected, steps)
}
