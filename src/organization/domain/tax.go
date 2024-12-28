package domain

import (
	"context"
	"time"

	"github.com/smallbiznis/go-genproto/smallbiznis/organization/v1"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type TaxRule struct {
	ID             string         `gorm:"column:tax_rule_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"tax_rule_id"`
	OrganizationID string         `gorm:"column:organization_id;type:uuid;" json:"organization_id"`
	CountryID      string         `gorm:"column:country_id" json:"country_id"`
	Country        Countries      `gorm:"->;foreignKey:CountryID" json:"country"`
	Type           string         `gorm:"column:type;" json:"type"`
	Rate           float32        `gorm:"column:rate;" json:"rate"`
	CreatedAt      time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

func (m *TaxRule) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

func (m *TaxRule) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return
}

func (m *TaxRule) ToProto() *organization.TaxRule {
	return &organization.TaxRule{
		TaxId:          m.ID,
		OrganizationId: m.OrganizationID,
		CountryId:      m.CountryID,
		Rate:           m.Rate,
		CreatedAt:      timestamppb.New(m.CreatedAt),
		UpdatedAt:      timestamppb.New(m.UpdatedAt),
	}
}

type TaxRules []TaxRule

func (m TaxRules) ToProto() (data []*organization.TaxRule) {
	for _, v := range m {
		data = append(data, v.ToProto())
	}
	return
}

type ITaxRulesRepository interface {
	Find(context.Context, pagination.Pagination, TaxRule) (TaxRules, int64, error)
	FindOne(context.Context, TaxRule) (*TaxRule, error)
	Update(context.Context, TaxRule) (*TaxRule, error)
	Save(context.Context, TaxRule) (*TaxRule, error)
}
