package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Profile struct {
	ID         string    `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey" json:"-"`
	AccountID  string    `gorm:"column:account_id" json:"-"`
	ProfileURI string    `gorm:"column:profile_uri" json:"profile_uri"`
	FirstName  string    `gorm:"column:first_name" json:"first_name"`
	LastName   string    `gorm:"column:last_name" json:"last_name"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (m *Profile) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.NewString()
	return
}
