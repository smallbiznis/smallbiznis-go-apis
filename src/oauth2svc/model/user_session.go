package model

import (
	"net/http"
	"time"

	"gorm.io/gorm"
)

type UserSession struct {
	ID        string    `bson:"_id" gorm:"column:id;primaryKey" json:"session_id"`
	UserID    string    `bson:"user_id" gorm:"column:user_id" json:"user_id"`
	User      Account   `bson:"-" gorm:"foreignKey:UserID" json:"-"`
	IP        string    `bson:"ip" gorm:"column:ip" json:"ip"`
	UserAgent string    `bson:"user_agent" gorm:"column:user_agent" json:"user_agent"`
	CreatedAt time.Time `bson:"created_at" gorm:"column:created_at;default:now();not null" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" gorm:"column:updated_at;default:now();not null" json:"updated_at"`
	ExpiredAt time.Time `bson:"expired_at" gorm:"column:expired_at;" json:"expired_at"`
}

func (m *UserSession) TableName() string {
	return "sessions"
}

func (m *UserSession) BeforeCreate(tx *gorm.DB) (err error) {
	current := time.Now()
	m.CreatedAt = current
	m.UpdatedAt = current
	return
}

func (m *UserSession) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()
	return
}

func NewSession(req *http.Request) *UserSession {
	return &UserSession{
		IP:        req.RemoteAddr,
		UserAgent: req.UserAgent(),
		ExpiredAt: time.Now().Add(12 * time.Hour),
	}
}
