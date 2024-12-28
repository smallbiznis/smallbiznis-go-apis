package model

import "time"

type CategoryProvider string

var (
	OAuth2 CategoryProvider = "oauth2"
	Email  CategoryProvider = "email"
	SMS    CategoryProvider = "sms"
)

func (cp CategoryProvider) String() string {
	if cp == OAuth2 ||
		cp == Email ||
		cp == SMS {
		return string(cp)
	}
	return ""
}

type Provider struct {
	ID        string           `bson:"provider_id" gorm:"column:provider_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"provider_id"`
	LogoURL   string           `bson:"-" gorm:"column:logo_url" json:"logo_url"`
	Title     string           `bson:"title" gorm:"column:title;" json:"title"`
	Category  CategoryProvider `bson:"category" gorm:"column:category;" json:"category"`
	Status    string           `bson:"status" gorm:"column:status;" json:"status"`
	CreatedAt time.Time        `bson:"created_at" gorm:"column:created_at;default:now();not null" json:"created_at"`
	UpdatedAt time.Time        `bson:"updated_at" gorm:"column:updated_at;default:now();not null" json:"updated_at"`
}
