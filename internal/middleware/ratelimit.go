package middleware

import (
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/exiaohu/go-demo/config"
	"github.com/exiaohu/go-demo/pkg/logger"
	"github.com/exiaohu/go-demo/pkg/util/ip"
)

// IPRateLimiter 存储每个 IP 的限流器
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

// NewIPRateLimiter 创建新的 IP 限流器
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}

	// 启动清理过期 IP 的 goroutine
	go i.cleanup()

	return i
}

// GetLimiter 获取指定 IP 的限流器
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

// cleanup 定期清理（简单实现，生产环境建议使用 LRU 或 Redis）
func (i *IPRateLimiter) cleanup() {
	for {
		time.Sleep(10 * time.Minute)
		i.mu.Lock()
		// 这里简化为清空，实际应判断最后访问时间
		i.ips = make(map[string]*rate.Limiter)
		i.mu.Unlock()
	}
}

var (
	globalLimiter *IPRateLimiter
	once          sync.Once
)

// RateLimit 限流中间件
func RateLimit(next http.Handler) http.Handler {
	if !config.C.RateLimit.Enabled {
		return next
	}

	once.Do(func() {
		globalLimiter = NewIPRateLimiter(
			rate.Limit(config.C.RateLimit.RPS),
			config.C.RateLimit.Burst,
		)
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 获取 IP
		clientIP := ip.GetClientIP(r)

		limiter := globalLimiter.GetLimiter(clientIP)
		if !limiter.Allow() {
			logger.Warn("Rate limit exceeded", zap.String("ip", clientIP))
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
