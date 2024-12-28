package domain

import (
	"context"
	"time"

	"github.com/smallbiznis/go-genproto/smallbiznis/transaction/v1"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type Order struct {
	ID                string                `gorm:"column:order_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"order_id"`
	OrganizationID    string                `gorm:"column:organization_id;type:uuid" json:"organization_id"`
	CustomerID        *string               `gorm:"column:customer_id;type:uuid;default:NULL" json:"customer_id"`
	BillingAddressID  *string               `gorm:"column:billing_address_id;type:uuid;default:NULL" json:"billing_address_id"`
	BillingAddress    *OrderBillingAddress  `gorm:"foreignKey:OrderID" json:"billing_address"`
	ShippingAddressID *string               `gorm:"column:shipping_address_id;type:uuid;default:NULL" json:"shipping_address_id"`
	ShippingAddress   *OrderShippingAddress `gorm:"foreignKey:OrderID" json:"shipping_address"`
	OrderNo           string                `gorm:"column:order_no" json:"order_no"`
	OrderItems        OrderItems            `gorm:"foreignKey:OrderID" json:"order_items"`
	SubTotal          float32               `gorm:"column:sub_total" json:"sub_total"`
	TaxAmount         float32               `gorm:"column:tax_amount" json:"tax_amount"`
	TotalAmount       float32               `gorm:"column:total_amount" json:"total_amount"`
	Status            string                `gorm:"column:status" json:"status"`
	CreatedAt         time.Time             `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         time.Time             `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt         gorm.DeletedAt        `gorm:"column:deleted_at" json:"-"`
}

func (m *Order) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

func (m *Order) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return
}

func (m *Order) ToProto() *transaction.Order {
	return &transaction.Order{
		OrderId:           m.ID,
		OrganizationId:    m.OrganizationID,
		SalesChannelId:    "",
		CustomerId:        *m.CustomerID,
		BillingAddressId:  "",
		ShippingAddressId: "",
		OrderNo:           m.OrderNo,
		OrderItems:        m.OrderItems.ToProto(),
		TaxAmount:         m.TaxAmount,
		SubTotal:          m.SubTotal,
		TotalAmount:       m.TotalAmount,
		Status:            transaction.OrderStatus_created,
		CreatedAt:         timestamppb.New(m.CreatedAt),
		UpdatedAt:         timestamppb.New(m.UpdatedAt),
	}
}

type Orders []Order

func (m Orders) ToProto() (data []*transaction.Order) {
	for _, v := range m {
		data = append(data, v.ToProto())
	}
	return
}

type ShippingHistory struct {
	ID              string         `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OrderShippingID string         `gorm:"column:order_shipping_id" json:"order_shipping_id"`
	CreatedAt       time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

func (m *ShippingHistory) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

func (m *ShippingHistory) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return
}

type OrderShipping struct {
	ID               string         `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OrderID          string         `gorm:"column:order_id" json:"order_id"`
	ShippingMethodID string         `gorm:"column:shipping_method_id" json:"shipping_method_id"`
	WaybillNumber    string         `gorm:"column:waybill_number" json:"waybill_number"`
	CreatedAt        time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

func (m *OrderShipping) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

func (m *OrderShipping) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return
}

type IOrderRepository interface {
	Find(context.Context, pagination.Pagination, Order) (Orders, int64, error)
	FindOne(context.Context, Order) (*Order, error)
	Save(context.Context, Order) (*Order, error)
	Update(context.Context, Order) (*Order, error)
}
