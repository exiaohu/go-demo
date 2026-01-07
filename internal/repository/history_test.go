package repository

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/exiaohu/go-demo/internal/model"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// 使用内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// 自动迁移
	err = db.AutoMigrate(&model.CalculationHistory{})
	assert.NoError(t, err)

	return db
}

func TestGormHistoryRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewHistoryRepository(db)

	history := &model.CalculationHistory{
		Operation: "add",
		A:         1,
		B:         2,
		Result:    3,
		ClientIP:  "127.0.0.1",
	}

	err := repo.Create(context.Background(), history)
	assert.NoError(t, err)
	assert.NotZero(t, history.ID)
	assert.NotZero(t, history.CreatedAt)
}

func TestGormHistoryRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewHistoryRepository(db)

	// 插入一些数据
	histories := []model.CalculationHistory{
		{Operation: "add", A: 1, B: 2, Result: 3, ClientIP: "127.0.0.1"},
		{Operation: "sub", A: 5, B: 3, Result: 2, ClientIP: "127.0.0.1"},
		{Operation: "mul", A: 2, B: 3, Result: 6, ClientIP: "127.0.0.1"},
	}

	for _, h := range histories {
		err := repo.Create(context.Background(), &h)
		assert.NoError(t, err)
	}

	// 测试列表查询
	list, err := repo.List(context.Background(), 2)
	assert.NoError(t, err)
	assert.Len(t, list, 2)

	// 验证顺序（最新的在前）
	assert.Equal(t, "mul", list[0].Operation)
	assert.Equal(t, "sub", list[1].Operation)
}
