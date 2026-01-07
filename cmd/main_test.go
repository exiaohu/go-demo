package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/exiaohu/go-demo/internal/handler"
)

func TestHandlers(t *testing.T) {
	tests := []struct {
		name           string
		handler        http.HandlerFunc
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Add 1+2",
			handler:        handler.AddHandler,
			path:           "/add?a=1&b=2",
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"code\":200,\"message\":\"OK\",\"data\":{\"result\":3}}",
		},
		{
			name:           "Subtract 5-3",
			handler:        handler.SubtractHandler,
			path:           "/subtract?a=5&b=3",
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"code\":200,\"message\":\"OK\",\"data\":{\"result\":2}}",
		},
		{
			name:           "Multiply 4*3",
			handler:        handler.MultiplyHandler,
			path:           "/multiply?a=4&b=3",
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"code\":200,\"message\":\"OK\",\"data\":{\"result\":12}}",
		},
		{
			name:           "Divide 10/2",
			handler:        handler.DivideHandler,
			path:           "/divide?a=10&b=2",
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"code\":200,\"message\":\"OK\",\"data\":{\"result\":5}}",
		},
		{
			name:           "Divide by zero",
			handler:        handler.DivideHandler,
			path:           "/divide?a=10&b=0",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"code\":400,\"message\":\"[400] Division by zero\"}",
		},
		{
			name:           "Invalid param",
			handler:        handler.AddHandler,
			path:           "/add?a=abc&b=2",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"code\":400,\"message\":\"[400] Invalid parameter format: strconv.Atoi: parsing \\\"abc\\\": invalid syntax\"}",
		},
		{
			name:           "Missing param",
			handler:        handler.AddHandler,
			path:           "/add?a=1",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"code\":400,\"message\":\"[400] Parameter is required\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := tt.handler

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if body := rr.Body.String(); strings.TrimSpace(body) != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					strings.TrimSpace(body), tt.expectedBody)
			}
		})
	}
}
