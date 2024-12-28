package model

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type AuthorizationCode struct {
	ID                  string         `bson:"_id" gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OrganizationID      string         `bson:"organization_id" gorm:"column:organization_id;type:uuid" json:"organization_id"`
	ApplicationID       string         `bson:"application_id" gorm:"column:application_id;type:uuid" json:"application_id"`
	UserID              string         `bson:"user_id" gorm:"user_id" json:"user_id"`
	User                Account        `bson:"-" gorm:"foreignKey:user_id" json:"-"`
	Code                string         `bson:"code" gorm:"code" json:"code"`
	RedirectUri         string         `bson:"redirect_uri" gorm:"redirect_uri" json:"redirect_uri"`
	CodeChallenge       string         `bson:"code_challenge" gorrm:"code_challenge" json:"code_challenge"`
	CodeChallengeMethod string         `bson:"code_challenge_method" gorrm:"code_challenge_method" json:"code_challenge_method"`
	Scopes              pq.StringArray `bson:"scopes" gorm:"column:scopes;type:VARCHAR(255)" json:"scopes"`
	Revoke              bool           `bson:"revoke" gorm:"revoke" json:"revoke"`
	CreatedAt           time.Time      `bson:"created_at" gorm:"column:created_at;default:now();not null" json:"created_at"`
	UpdatedAt           time.Time      `bson:"updated_at" gorm:"column:updated_at;default:now();not null" json:"updated_at"`
	DeletedAt           *time.Time     `bson:"deleted_at" gorm:"column:deleted_at;index" json:"-"`
}

func (m *AuthorizationCode) BeforeCreate(tx *gorm.DB) (err error) {
	current := time.Now()
	m.CreatedAt = current
	m.UpdatedAt = current
	return
}

func (m *AuthorizationCode) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()
	return
}

func (m *AuthorizationCode) TableName() string {
	return "authorization_codes"
}

func NewAuthorizationCode() *AuthorizationCode {
	return &AuthorizationCode{
		ID: uuid.NewString(),
	}
}

type TokenGeneration struct {
	Organization string
	Application  string
	User         AggregateAccount
	Scope        pq.StringArray
	Request      *http.Request
}
