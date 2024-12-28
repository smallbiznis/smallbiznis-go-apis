package domain

import (
	"context"
	"time"

	"github.com/smallbiznis/go-lib/pkg/pagination"
	"gorm.io/gorm"
)

type Option struct {
	ID             string         `gorm:"column:option_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"option_id"`
	OrganizationID string         `gorm:"column:organization_id" json:"organization_id"`
	Name           string         `gorm:"column:option_name" json:"option_name"`
	Values         OptionValues   `gorm:"foreignKey:OptionID" json:"values"`
	CreatedAt      time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

func (m *Option) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

func (m *Option) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return
}

type Options []Option

type OptionValue struct {
	ID        string    `gorm:"column:option_value_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"option_value_id"`
	OptionID  string    `gorm:"column:option_id;type:uuid" json:"option_id"`
	Value     string    `gorm:"column:value" json:"value"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (m *OptionValue) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

func (m *OptionValue) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return
}

type OptionValues []OptionValue

type IOptionRepository interface {
	Find(context.Context, pagination.Pagination, Option) (Options, int64, error)
	FindOne(context.Context, Option) (*Option, error)
	Save(context.Context, Option) (*Option, error)
	Update(context.Context, Option) (*Option, error)
	Delete(context.Context, Option) error
}
