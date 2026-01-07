package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/exiaohu/go-demo/internal/model"
	"github.com/exiaohu/go-demo/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCalculatorService 模拟 CalculatorService
type MockCalculatorService struct {
	mock.Mock
}

func (m *MockCalculatorService) Add(ctx context.Context, a, b int, ip string) (int, error) {
	args := m.Called(ctx, a, b, ip)
	return args.Int(0), args.Error(1)
}

func (m *MockCalculatorService) Subtract(ctx context.Context, a, b int, ip string) (int, error) {
	args := m.Called(ctx, a, b, ip)
	return args.Int(0), args.Error(1)
}

func (m *MockCalculatorService) Multiply(ctx context.Context, a, b int, ip string) (int, error) {
	args := m.Called(ctx, a, b, ip)
	return args.Int(0), args.Error(1)
}

func (m *MockCalculatorService) Divide(ctx context.Context, a, b int, ip string) (int, error) {
	args := m.Called(ctx, a, b, ip)
	return args.Int(0), args.Error(1)
}

func (m *MockCalculatorService) GetHistory(ctx context.Context, limit int) ([]model.CalculationHistory, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]model.CalculationHistory), args.Error(1)
}

func (m *MockCalculatorService) Close() error {
	args := m.Called()
	return args.Error(0)
}

func setupHandler() (*Handler, *MockCalculatorService) {
	mockService := new(MockCalculatorService)
	return NewHandler(mockService), mockService
}

func TestHomeHandler(t *testing.T) {
	h, _ := setupHandler()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.HomeHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Welcome to Playground!", rr.Body.String())
}

func TestHomeHandler_NotFound(t *testing.T) {
	h, _ := setupHandler()
	req, err := http.NewRequest("GET", "/not-found", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.HomeHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestHealthCheckHandler(t *testing.T) {
	h, _ := setupHandler()
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.HealthCheckHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())
}

func TestHistoryHandler(t *testing.T) {
	h, mockService := setupHandler()

	expectedHistory := []model.CalculationHistory{
		{Operation: "add", A: 1, B: 2, Result: 3},
	}
	// 匹配 context
	mockService.On("GetHistory", mock.Anything, 10).Return(expectedHistory, nil)

	req, err := http.NewRequest("GET", "/history", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.HistoryHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])

	mockService.AssertExpectations(t)
}

func TestMathHandlers(t *testing.T) {
	h, mockService := setupHandler()

	// Setup mocks with context matching
	mockService.On("Add", mock.Anything, 1, 2, mock.Anything).Return(3, nil)
	mockService.On("Subtract", mock.Anything, 5, 3, mock.Anything).Return(2, nil)
	mockService.On("Multiply", mock.Anything, 4, 3, mock.Anything).Return(12, nil)
	mockService.On("Divide", mock.Anything, 10, 2, mock.Anything).Return(5, nil)
	mockService.On("Divide", mock.Anything, 10, 0, mock.Anything).Return(0, errors.New(errors.ErrTypeValidation, "Division by zero"))

	tests := []struct {
		name           string
		handler        func(http.ResponseWriter, *http.Request)
		queryParams    string
		method         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Add 1+2",
			handler:        h.AddHandler,
			queryParams:    "?a=1&b=2",
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"code":200,"message":"OK","data":{"result":3}}`,
		},
		{
			name:           "Subtract 5-3",
			handler:        h.SubtractHandler,
			queryParams:    "?a=5&b=3",
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"code":200,"message":"OK","data":{"result":2}}`,
		},
		{
			name:           "Multiply 4*3",
			handler:        h.MultiplyHandler,
			queryParams:    "?a=4&b=3",
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"code":200,"message":"OK","data":{"result":12}}`,
		},
		{
			name:           "Divide 10/2",
			handler:        h.DivideHandler,
			queryParams:    "?a=10&b=2",
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"code":200,"message":"OK","data":{"result":5}}`,
		},
		{
			name:           "Divide by zero",
			handler:        h.DivideHandler,
			queryParams:    "?a=10&b=0",
			method:         "GET",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":400,"message":"[400] Division by zero"}`,
		},
		{
			name:           "Invalid parameter",
			handler:        h.AddHandler,
			queryParams:    "?a=abc&b=2",
			method:         "GET",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":400,"message":"[400] Invalid parameter format: strconv.Atoi: parsing \"abc\": invalid syntax"}`,
		},
		{
			name:           "Missing parameter",
			handler:        h.AddHandler,
			queryParams:    "?a=1",
			method:         "GET",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":400,"message":"[400] Parameter is required"}`,
		},
		{
			name:           "Method Not Allowed",
			handler:        h.AddHandler,
			queryParams:    "?a=1&b=2",
			method:         "POST",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   `{"code":405,"message":"Method not allowed"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/"+tt.queryParams, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(tt.handler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			// 验证响应体
			if tt.expectedStatus == http.StatusOK || tt.expectedStatus == http.StatusBadRequest || tt.expectedStatus == http.StatusMethodNotAllowed {
				var expectedMap, actualMap map[string]interface{}
				err = json.Unmarshal([]byte(tt.expectedBody), &expectedMap)
				assert.NoError(t, err)
				err = json.Unmarshal(rr.Body.Bytes(), &actualMap)
				assert.NoError(t, err)
				assert.Equal(t, expectedMap, actualMap)
			} else {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}
		})
	}
}
