package model

import (
	"time"

	"gorm.io/gorm"
)

// CalculationHistory 记录计算历史
type CalculationHistory struct {
	ID        uint           `gorm:"primarykey"       json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"            json:"-"`
	Operation string         `gorm:"size:32;not null" json:"operation"`
	A         int            `gorm:"not null"         json:"a"`
	B         int            `gorm:"not null"         json:"b"`
	Result    int            `gorm:"not null"         json:"result"`
	ClientIP  string         `gorm:"size:64"          json:"client_ip"`
}
