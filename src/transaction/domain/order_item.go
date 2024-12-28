package domain

import (
	"time"

	"github.com/smallbiznis/go-genproto/smallbiznis/transaction/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type OrderItem struct {
	ID         string         `gorm:"column:order_item_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"order_item_id"`
	OrderID    string         `gorm:"column:order_id;type:uuid" json:"order_id"`
	Order      Order          `gorm:"foreignKey:OrderID" json:"-"`
	VariantID  string         `gorm:"column:variant_id;type:uuid" json:"variant_id"`
	Quantity   int32          `gorm:"column:quantity" json:"quantity"`
	UnitPrice  float32        `gorm:"column:unit_price" json:"unit_price"`
	TotalPrice float32        `gorm:"column:total_price" json:"total_price"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

func (m *OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

func (m *OrderItem) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return
}

func (m *OrderItem) ToProto() *transaction.OrderItem {
	return &transaction.OrderItem{
		OrderItemId: m.ID,
		OrderId:     m.OrderID,
		ItemId:      m.VariantID,
		Quantity:    m.Quantity,
		UnitPrice:   m.UnitPrice,
		TotalPrice:  m.TotalPrice,
		CreatedAt:   timestamppb.New(m.CreatedAt),
		UpdatedAt:   timestamppb.New(m.UpdatedAt),
	}
}

type OrderItems []OrderItem

func (m OrderItems) ToProto() (data []*transaction.OrderItem) {
	for _, v := range m {
		data = append(data, v.ToProto())
	}
	return
}
