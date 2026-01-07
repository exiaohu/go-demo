package service

import (
	"context"
	"errors"
	"testing"

	"github.com/exiaohu/go-demo/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHistoryRepository 模拟 HistoryRepository
type MockHistoryRepository struct {
	mock.Mock
}

func (m *MockHistoryRepository) Create(ctx context.Context, history *model.CalculationHistory) error {
	args := m.Called(ctx, history)
	return args.Error(0)
}

func (m *MockHistoryRepository) List(ctx context.Context, limit int) ([]model.CalculationHistory, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]model.CalculationHistory), args.Error(1)
}

func TestCalculatorService_Add(t *testing.T) {
	mockRepo := new(MockHistoryRepository)
	svc := NewCalculatorService(mockRepo)

	// 预期 Create 会被异步调用，这里我们不强制检查异步调用是否完成
	// 但我们可以通过 WaitGroup 或 channel 在实际代码中控制，或者在测试中简单地忽略异步错误
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	result, err := svc.Add(context.Background(), 1, 2, "127.0.0.1")
	assert.NoError(t, err)
	assert.Equal(t, 3, result)
}

func TestCalculatorService_Divide_Error(t *testing.T) {
	mockRepo := new(MockHistoryRepository)
	svc := NewCalculatorService(mockRepo)

	result, err := svc.Divide(context.Background(), 10, 0, "127.0.0.1")
	assert.Error(t, err)
	assert.Equal(t, 0, result)
	// 发生错误时不应该记录历史
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCalculatorService_GetHistory(t *testing.T) {
	mockRepo := new(MockHistoryRepository)
	svc := NewCalculatorService(mockRepo)

	expectedHistory := []model.CalculationHistory{
		{Operation: "add", A: 1, B: 2, Result: 3},
	}
	mockRepo.On("List", mock.Anything, 10).Return(expectedHistory, nil)

	history, err := svc.GetHistory(context.Background(), 10)
	assert.NoError(t, err)
	assert.Equal(t, expectedHistory, history)
	mockRepo.AssertExpectations(t)
}

func TestCalculatorService_GetHistory_Error(t *testing.T) {
	mockRepo := new(MockHistoryRepository)
	svc := NewCalculatorService(mockRepo)

	mockRepo.On("List", mock.Anything, 10).Return([]model.CalculationHistory{}, errors.New("db error"))

	history, err := svc.GetHistory(context.Background(), 10)
	assert.Error(t, err)
	assert.Empty(t, history)
	mockRepo.AssertExpectations(t)
}
