package ip

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name     string
		header   map[string]string
		remote   string
		expected string
	}{
		{
			name:     "X-Forwarded-For",
			header:   map[string]string{"X-Forwarded-For": "10.0.0.1, 10.0.0.2"},
			remote:   "127.0.0.1:12345",
			expected: "10.0.0.1",
		},
		{
			name:     "X-Real-IP",
			header:   map[string]string{"X-Real-IP": "10.0.0.3"},
			remote:   "127.0.0.1:12345",
			expected: "10.0.0.3",
		},
		{
			name:     "Both Headers (X-Forwarded-For preferred)",
			header:   map[string]string{"X-Forwarded-For": "10.0.0.1", "X-Real-IP": "10.0.0.3"},
			remote:   "127.0.0.1:12345",
			expected: "10.0.0.1",
		},
		{
			name:     "No Headers",
			header:   nil,
			remote:   "127.0.0.1:12345",
			expected: "127.0.0.1",
		},
		{
			name:     "No Headers, IPv6",
			header:   nil,
			remote:   "[::1]:12345",
			expected: "::1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			for k, v := range tt.header {
				req.Header.Set(k, v)
			}
			req.RemoteAddr = tt.remote

			ip := GetClientIP(req)
			assert.Equal(t, tt.expected, ip)
		})
	}
}
