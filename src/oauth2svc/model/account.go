package model

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AccountType string

var (
	User           AccountType = "user"
	ServiceAccount AccountType = "service_account"
)

func (p AccountType) String() string {
	if p == User ||
		p == ServiceAccount {
		return string(p)
	}
	return ""
}

type AuthProvider string

var (
	Password    AuthProvider = "password"
	PhoneNumber AuthProvider = "phone_number"
)

func (p AuthProvider) String() string {
	if p == Password ||
		p == PhoneNumber {
		return string(p)
	}
	return ""
}

type AccountRole string

var (
	ROLE_USER            AccountRole = "ROLE_USER"
	ROLE_SERVICE_ACCOUNT AccountRole = "ROLE_SERVICE_ACCOUNT"
)

func (ar AccountRole) String() string {
	if ar == ROLE_USER ||
		ar == ROLE_SERVICE_ACCOUNT {
		return string(ar)
	}
	return ""
}

type Account struct {
	ID                 string         `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey" json:"account_id"`
	OrganizationID     string         `bson:"organization_id" gorm:"column:organization_id;type:uuid" json:"-"`
	Type               AccountType    `bson:"type" gorm:"column:type" json:"type"`
	Provider           AuthProvider   `bson:"provider" gorm:"column:provider;type:varchar(255)" json:"provider"`
	Username           string         `bson:"username" gorm:"column:username" json:"username"`
	Password           string         `gorm:"column:password;type:varchar(255)" json:"-"`
	Profile            Profile        `gorm:"foreignKey:AccountID" json:"-"`
	VerifiedAt         *time.Time     `gorm:"column:verified_at" json:"verified_at"`
	Roles              pq.StringArray `gorm:"column:roles;type:varchar(255)" json:"roles"`
	CreatedAt          time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt          *time.Time     `gorm:"column:deleted_at;index" json:"deleted_at"`
	LastUpdatePassword *time.Time     `gorm:"column:last_update_password" json:"last_update_password"`
}

func NewAccount(orgID string) *Account {
	return &Account{
		OrganizationID: orgID,
		Roles: pq.StringArray{
			ROLE_USER.String(),
		},
	}
}

func NewServiceAccount(orgID string) *Account {
	source := rand.NewSource(99)
	r := rand.New(source)
	return &Account{
		OrganizationID: orgID,
		Username:       fmt.Sprintf("%v@smallbiznis.id", r.Int31()),
		Roles: pq.StringArray{
			ROLE_SERVICE_ACCOUNT.String(),
		},
	}
}

func (m *Account) SetPassword(password string) (err error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	m.Password = string(b)
	return
}

func (m *Account) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(password)) == nil
}

func (m *Account) BeforeCreate(tx *gorm.DB) (err error) {
	current := time.Now()
	m.CreatedAt = current
	m.UpdatedAt = current
	return
}

func (m *Account) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()
	return
}

type AggregateAccount struct {
	Account
	ProfileURI string `gorm:"column:profile_uri" json:"profile_uri"`
	FirstName  string `gorm:"column:first_name" json:"first_name"`
	LastName   string `gorm:"column:last_name" json:"last_name"`
	SessionID  string `gorm:"-" json:"-"`
}

type AggregateAccounts []AggregateAccount

type RequestLookup struct {
	AccountID   string `form:"account_id"`
	Email       string `form:"email"`
	PhoneNumber string `form:"phone_number"`
}

type RequestSignUp struct {
	Provider    AuthProvider `json:"provider" validate:"required,oneof=password google facebook phone_number"`
	Email       string       `json:"email" validate:"required_if=Provider password"`
	Password    string       `json:"password" validate:"required_if=Provider password"`
	PhoneNumber string       `json:"phone_number" validate:"required_if=Provider phone_number"`
	FirstName   string       `json:"first_name" validate:"required"`
	LastName    string       `json:"last_name"`
}

type RequestSignInWithPassword struct {
	Request  *http.Request
	Email    string `form:"email" json:"email" validate:"required,email"`
	Password string `form:"password" json:"password" validate:"required"`
}

type RequestSignInWithPhoneNumber struct {
	Request     http.Request
	PhoneNumber string `form:"phone_number" json:"phone_number" validate:"required"`
	SessionID   string `form:"session_id" json:"session_id" validate:"required"`
}

type RequestSendVerificationCode struct {
	Request     http.Request
	ClientID    string `json:"client_id" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}
