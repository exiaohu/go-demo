package response

import (
	"encoding/json"
	"net/http"

	"github.com/exiaohu/go-demo/internal/middleware"
)

// Response 统一响应结构
type Response struct {
	Code      int         `json:"code"`                 // 业务状态码
	Message   string      `json:"message"`              // 提示信息
	Data      interface{} `json:"data,omitempty"`       // 数据载荷
	RequestID string      `json:"request_id,omitempty"` // 请求 ID
}

// JSON 发送 JSON 响应
func JSON(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := Response{
		Code:      status,
		Message:   http.StatusText(status),
		Data:      data,
		RequestID: middleware.GetRequestID(r.Context()),
	}

	json.NewEncoder(w).Encode(resp)
}

// Error 发送错误响应
func Error(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := Response{
		Code:      status,
		Message:   message,
		RequestID: middleware.GetRequestID(r.Context()),
	}

	json.NewEncoder(w).Encode(resp)
}

// Success 发送成功响应
func Success(w http.ResponseWriter, r *http.Request, data interface{}) {
	JSON(w, r, http.StatusOK, data)
}
