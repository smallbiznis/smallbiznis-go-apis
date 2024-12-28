package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Organization
type Organization struct {
	ID             string     `bson:"_id" gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OrganizationID *string    `bson:"organization_id" gorm:"column:organization_id" json:"organization_id"`
	Logo           *string    `bson:"logo" gorm:"column:logo;" json:"logo"`
	Name           string     `bson:"name" gorm:"column:name;" json:"name" validate:"required"`
	DisplayName    string     `bson:"display_name" gorm:"display_name" json:"display_name"`
	IsDefault      bool       `bson:"is_default" gorm:"is_default" json:"is_default"`
	CreatedAt      time.Time  `bson:"created_at" gorm:"column:created_at;default:now();not null" json:"created_at"`
	UpdatedAt      time.Time  `bson:"updated_at" gorm:"column:updated_at;default:now();not null" json:"updated_at"`
	DeletedAt      *time.Time `bson:"deleted_at" gorm:"deleted_at,index" json:"-"`
}

func (m *Organization) BeforeCreate(tx *gorm.DB) (err error) {
	current := time.Now()
	m.CreatedAt = current
	m.UpdatedAt = current
	return
}

func (m *Organization) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()
	return
}

type UserOrganization struct {
	Organization string    `bson:"organization" gorm:"column:organization" jso:"organization"`
	UserID       uuid.UUID `bson:"user_id" gorm:"column:user_id;type:uuid" json:"id"`
}

// FilterOrganization
type FilterOrganization struct {
	Pagination
	Name string `form:"name"`
}
