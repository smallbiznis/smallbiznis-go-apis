package domain

import (
	"time"

	"gorm.io/gorm"
)

type OrderShippingAddress struct {
	ShippingAddressID string         `gorm:"column:shipping_address_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"shipping_address_id"`
	OrderID           string         `gorm:"column:order_id" json:"order_id"`
	CreatedAt         time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

func (m *OrderShippingAddress) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

func (m *OrderShippingAddress) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return
}
