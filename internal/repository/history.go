package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/exiaohu/go-demo/internal/model"
)

// HistoryRepository 定义历史记录数据访问接口
type HistoryRepository interface {
	Create(ctx context.Context, history *model.CalculationHistory) error
	List(ctx context.Context, limit int) ([]model.CalculationHistory, error)
}

type GormHistoryRepository struct {
	db *gorm.DB
}

// NewHistoryRepository 创建 HistoryRepository 实例
func NewHistoryRepository(db *gorm.DB) *GormHistoryRepository {
	return &GormHistoryRepository{db: db}
}

func (r *GormHistoryRepository) Create(ctx context.Context, history *model.CalculationHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *GormHistoryRepository) List(ctx context.Context, limit int) ([]model.CalculationHistory, error) {
	var history []model.CalculationHistory
	err := r.db.WithContext(ctx).Order("created_at desc").Limit(limit).Find(&history).Error
	return history, err
}
