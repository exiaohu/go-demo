package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(io.Discard)
	},
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// Gzip 压缩中间件
func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, ok := gzipWriterPool.Get().(*gzip.Writer)
		if !ok {
			// Should not happen if pool is correctly used
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer gzipWriterPool.Put(gz)

		gz.Reset(w)
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")

		// 包装 ResponseWriter
		// 注意：这里需要根据 ResponseWriter 的接口实现更多方法（如 Hijack, Flush 等）
		// 简单起见，我们只实现 Write
		next.ServeHTTP(gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
	})
}
