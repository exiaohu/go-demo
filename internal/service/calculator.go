package service

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/exiaohu/go-demo/internal/math"
	"github.com/exiaohu/go-demo/internal/model"
	"github.com/exiaohu/go-demo/internal/repository"
	"github.com/exiaohu/go-demo/pkg/logger"
)

// CalculatorService 定义计算服务接口
type CalculatorService interface {
	Add(ctx context.Context, a, b int, ip string) (int, error)
	Subtract(ctx context.Context, a, b int, ip string) (int, error)
	Multiply(ctx context.Context, a, b int, ip string) (int, error)
	Divide(ctx context.Context, a, b int, ip string) (int, error)
	GetHistory(ctx context.Context, limit int) ([]model.CalculationHistory, error)
	Close() error
}

type StandardCalculatorService struct {
	repo repository.HistoryRepository
	wg   sync.WaitGroup
}

// NewCalculatorService 创建 CalculatorService 实例
func NewCalculatorService(repo repository.HistoryRepository) *StandardCalculatorService {
	return &StandardCalculatorService{repo: repo}
}

func (s *StandardCalculatorService) calculateAndRecord(
	_ context.Context,
	opName string,
	a, b int,
	ip string,
	op func(int, int) (int, error),
) (int, error) {
	result, err := op(a, b)
	if err != nil {
		return 0, err
	}

	// 异步记录历史
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		history := &model.CalculationHistory{
			Operation: opName,
			A:         a,
			B:         b,
			Result:    result,
			ClientIP:  ip,
		}
		// 使用 Background context
		if err := s.repo.Create(context.Background(), history); err != nil {
			logger.Error("Failed to save history", zap.Error(err))
		}
	}()

	return result, nil
}

func (s *StandardCalculatorService) Close() error {
	s.wg.Wait()
	return nil
}

func (s *StandardCalculatorService) Add(ctx context.Context, a, b int, ip string) (int, error) {
	return s.calculateAndRecord(ctx, "add", a, b, ip, math.Add)
}

func (s *StandardCalculatorService) Subtract(ctx context.Context, a, b int, ip string) (int, error) {
	return s.calculateAndRecord(ctx, "subtract", a, b, ip, math.Subtract)
}

func (s *StandardCalculatorService) Multiply(ctx context.Context, a, b int, ip string) (int, error) {
	return s.calculateAndRecord(ctx, "multiply", a, b, ip, math.Multiply)
}

func (s *StandardCalculatorService) Divide(ctx context.Context, a, b int, ip string) (int, error) {
	return s.calculateAndRecord(ctx, "divide", a, b, ip, math.Divide)
}

func (s *StandardCalculatorService) GetHistory(ctx context.Context, limit int) ([]model.CalculationHistory, error) {
	return s.repo.List(ctx, limit)
}
