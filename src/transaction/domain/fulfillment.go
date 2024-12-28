package domain

import (
	"time"

	"gorm.io/gorm"
)

type OrderFulfillment struct {
	ID        string         `gorm:"column:fulfillment_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OrderID   string         `gorm:"column:order_id" json:"order_id"`
	Status    string         `gorm:"column:status" json:"status"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

func (m *OrderFulfillment) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

func (m *OrderFulfillment) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return
}
