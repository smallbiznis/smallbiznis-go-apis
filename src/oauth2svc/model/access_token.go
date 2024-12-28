package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type AccessToken struct {
	ID                   string         `bson:"_id" gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OrganizationID       string         `bson:"organization_id" gorm:"column:organization_id;type:uuid" json:"organization_id"`
	ApplicationID        string         `bson:"application_id" gorm:"column:application_id;type:uuid" json:"application_id"`
	UserID               string         `bson:"user_id" gorm:"user_id" json:"user_id"`
	User                 Account        `bson:"-" gorm:"foreignKey:UserID" json:"-"`
	AccessToken          string         `bson:"access_token" gorm:"access_token" json:"access_token"`
	AccessTokenExpiresIn time.Time      `bson:"access_token_expires_in" gorm:"access_token_expires_in" json:"access_token_expires_in"`
	Scopes               pq.StringArray `bson:"scopes" gorm:"column:scopes;type:TEXT" json:"scopes"`
	CreatedAt            time.Time      `bson:"created_at" gorm:"column:created_at;default:now();not null" json:"created_at"`
	UpdatedAt            time.Time      `bson:"updated_at" gorm:"column:updated_at;default:now();not null" json:"updated_at"`
}

func NewAccessToken() *AccessToken {
	return &AccessToken{
		ID: uuid.NewString(),
	}
}

func (m *AccessToken) BeforeCreate(tx *gorm.DB) (err error) {
	current := time.Now()
	m.CreatedAt = current
	m.UpdatedAt = current
	return
}

func (m *AccessToken) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()
	return
}
