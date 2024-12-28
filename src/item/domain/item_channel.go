package domain

import (
	"time"

	"gorm.io/gorm"
)

type ItemChannel struct {
	ID        string         `gorm:"column:product_channel_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"product_channel_id"`
	ProductID string         `gorm:"column:product_id;type:uuid" json:"product_id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}
