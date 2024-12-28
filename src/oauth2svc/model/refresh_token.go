package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID                    string         `bson:"_id" gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OrganizationID        string         `bson:"organization_id" gorm:"column:organization_id;type:uuid" json:"organization_id"`
	ApplicationID         string         `bson:"application_id" gorm:"column:application_id;type:uuid" json:"application_id"`
	UserID                string         `bson:"user_id" gorm:"user_id" json:"user_id"`
	User                  Account        `bson:"-" gorm:"foreignKey:UserID" json:"-"`
	RefreshToken          *string        `bson:"refresh_token" gorm:"refresh_token" json:"refresh_token"`
	RefreshTokenExpiresIn *time.Time     `bson:"refresh_token_expires_in" gorm:"refresh_token_expires_in" json:"refresh_token_expires_in"`
	Scopes                pq.StringArray `bson:"scopes" gorm:"column:scopes;type:TEXT" json:"scopes"`
	Revoke                bool           `bson:"revoke" gorm:"revoke" json:"revoke"`
	CreatedAt             time.Time      `bson:"created_at" gorm:"column:created_at;default:now();not null" json:"created_at"`
	UpdatedAt             time.Time      `bson:"updated_at" gorm:"column:updated_at;default:now();not null" json:"updated_at"`
}

func NewRefreshToken() *RefreshToken {
	return &RefreshToken{
		ID: uuid.NewString(),
	}
}

func (m *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	current := time.Now()
	m.CreatedAt = current
	m.UpdatedAt = current
	return
}

func (m *RefreshToken) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()
	return
}
