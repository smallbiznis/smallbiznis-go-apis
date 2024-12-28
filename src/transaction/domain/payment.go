package domain

import (
	"time"

	"gorm.io/gorm"
)

type PaymentState string

var (
	Unpaid PaymentState = "unpaid"
	Paid   PaymentState = "paid"
)

func (m PaymentState) String() string {
	if m == Unpaid ||
		m == Paid {
		return string(m)
	}
	return ""
}

type PaymentMethod string

var (
	Card         PaymentMethod = "card"
	BankTransfer PaymentMethod = "bank_transfer"
)

func (m PaymentMethod) String() string {
	if m == Card ||
		m == BankTransfer {
		return string(m)
	}
	return ""
}

type OrderPayment struct {
	ID                string         `gorm:"column:order_payment_id" json:"order_payment_id"`
	OrderID           string         `gorm:"column:order_id;type:uuid" json:"order_id"`
	PaymentProviderID string         `gorm:"column:payment_provider_id" json:"payment_provider_id"`
	Method            string         `gorm:"column:method" json:"method"`
	Date              *time.Time     `gorm:"column:date" json:"date"`
	DueDate           time.Time      `gorm:"column:due_date" json:"due_date"`
	Status            string         `gorm:"column:status" json:"status"`
	CreatedAt         time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

func (m *OrderPayment) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

func (m *OrderPayment) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return
}
