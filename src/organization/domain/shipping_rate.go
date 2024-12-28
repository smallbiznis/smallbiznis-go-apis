package domain

import (
	"context"
	"time"

	"github.com/smallbiznis/go-genproto/smallbiznis/organization/v1"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"gorm.io/gorm"
)

type ShippingRate struct {
	ID             string         `gorm:"column:shipping_rate_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"shipping_rate_id"`
	OrganizationID string         `gorm:"column:organization_id" json:"organization_id"`
	Type           string         `gorm:"column:type" json:"type"`
	Name           string         `gorm:"column:name" json:"name"`
	Description    string         `gorm:"column:description" json:"description"`
	Price          float32        `gorm:"column:price" json:"price"`
	CreatedAt      time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

func (m *ShippingRate) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

func (m *ShippingRate) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return
}

func (m *ShippingRate) ToProto() *organization.ShippingRate {
	return &organization.ShippingRate{
		ShippingRateId: m.ID,
		OrganizationId: m.OrganizationID,
		Type:           organization.ShippingRate_RateType(organization.ShippingRate_RateType_value[m.Type]),
		Name:           m.Name,
		Description:    m.Description,
		Price:          m.Price,
	}
}

type ShippingRates []ShippingRate

func (m ShippingRates) ToProto() (data []*organization.ShippingRate) {
	for _, v := range m {
		data = append(data, v.ToProto())
	}
	return
}

type IShippingRateRepository interface {
	Find(context.Context, pagination.Pagination, ShippingRate) (ShippingRates, int64, error)
	FindOne(context.Context, ShippingRate) (*ShippingRate, error)
	Update(context.Context, ShippingRate) (*ShippingRate, error)
	Save(context.Context, ShippingRate) (*ShippingRate, error)
}
