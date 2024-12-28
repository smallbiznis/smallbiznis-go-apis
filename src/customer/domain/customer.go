package domain

import (
	"context"
	"time"

	"github.com/smallbiznis/go-genproto/smallbiznis/customer/v1"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"gorm.io/gorm"
)

type Customer struct {
	ID             string         `gorm:"column:customer_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"customer_id"`
	OrganizationID string         `gorm:"column:organization_id;type:uuid;" json:"organization_id"`
	AccountID      *string        `gorm:"column:account_id" json:"account_id"`
	FirstName      string         `gorm:"column:first_name" json:"first_name"`
	LastName       string         `gorm:"column:last_name" json:"last_name"`
	Email          string         `gorm:"column:email" json:"email"`
	CountryID      string         `gorm:"column:country_id;default:ID;" json:"country_id"`
	CreatedAt      time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

func (m *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

func (m *Customer) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return
}

func (m *Customer) ToProto() *customer.Customer {
	return &customer.Customer{
		CustomerId:     m.ID,
		OrganizationId: m.OrganizationID,
		FirstName:      m.FirstName,
		LastName:       m.LastName,
		Email:          m.Email,
		CountryId:      m.CountryID,
	}
}

type Customers []Customer

func (m Customers) ToProto() (data []*customer.Customer) {
	for _, v := range m {
		data = append(data, v.ToProto())
	}
	return
}

type ICustomerRepository interface {
	Find(context.Context, pagination.Pagination, Customer) (Customers, int64, error)
	FindOne(context.Context, Customer) (*Customer, error)
	Save(context.Context, Customer) (*Customer, error)
}
