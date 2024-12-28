package domain

import (
	"context"
	"time"

	"github.com/smallbiznis/go-genproto/smallbiznis/organization/v1"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type OrganizationStatus string

var (
	ACTIVE   OrganizationStatus = "ACTIVE"
	INACTIVE OrganizationStatus = "INACTIVE"
)

func (m OrganizationStatus) String() organization.Organization_Status {
	if m == ACTIVE {
		return organization.Organization_ACTIVE
	}

	return organization.Organization_INACTIVE
}

type Organization struct {
	ID                    string             `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OrganizationID        string             `gorm:"column:organization_id" uri:"organization_id" json:"organization_id"`
	StripeCustomerID      *string            `gorm:"column:stripe_customer_id" json:"stripe_customer_id"`
	StripeSubscriptionID  *string            `gorm:"column:stripe_subscription_id" json:"stripe_subscription_id"`
	StripePaymentMethodID *string            `gorm:"column:stripe_payment_method_id" json:"stripe_payment_method_id"`
	LogoUrl               string             `gorm:"column:logo_url" json:"logo_url"`
	Title                 string             `gorm:"column:title" form:"title" json:"title" validate:"required,max=50"`
	CountryID             string             `gorm:"column:country_id" json:"country_id"`
	Country               Countries          `gorm:"->;foreignKey:CountryID" json:"country"`
	IsDefault             bool               `gorm:"column:is_default" json:"-"`
	Status                OrganizationStatus `gorm:"column:status" json:"status"`
	CreatedAt             time.Time          `gorm:"column:created_at" json:"created_at"`
	UpdatedAt             time.Time          `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt             gorm.DeletedAt     `gorm:"column:deleted_at" json:"-"`
}

func (m *Organization) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

func (m *Organization) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return
}

func (m *Organization) ToProto() *organization.Organization {
	org := &organization.Organization{
		Id:             m.ID,
		OrganizationId: m.OrganizationID,
		Title:          m.Title,
		LogoUrl:        m.LogoUrl,
		CountryId:      m.CountryID,
		Country:        m.Country.ToProto(),
		IsDefault:      m.IsDefault,
		Status:         m.Status.String(),
		CreatedAt:      timestamppb.New(m.CreatedAt),
		UpdatedAt:      timestamppb.New(m.UpdatedAt),
	}

	if m.StripeCustomerID != nil {
		org.CustomerId = *m.StripeCustomerID
	}

	if m.StripeSubscriptionID != nil {
		org.SubscriptionId = *m.StripeSubscriptionID
	}

	if m.StripePaymentMethodID != nil {
		org.PaymentMethodId = *m.StripePaymentMethodID
	}

	return org
}

type Organizations []Organization

func (m Organizations) ToProto() (data []*organization.Organization) {
	for _, org := range m {
		data = append(data, org.ToProto())
	}
	return
}

type IOrganizationRepository interface {
	Find(context.Context, pagination.Pagination, Organization) (Organizations, int64, error)
	FindOne(context.Context, Organization) (*Organization, error)
	Save(context.Context, Organization) (*Organization, error)
	Update(context.Context, Organization) (*Organization, error)
	Delete(context.Context, Organization) error
}
