package model

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"github.com/lib/pq"
	s "github.com/smallbiznis/oauth2-server/internal/pkg/strings"
	"gorm.io/gorm"
)

type ApplicationType string

var (
	Web     ApplicationType = "web"
	Android ApplicationType = "android"
	IOS     ApplicationType = "iOS"
	Desktop ApplicationType = "desktop"
)

func (m ApplicationType) String() string {
	if m == Web ||
		m == Android ||
		m == IOS ||
		m == Desktop {
		return string(m)
	}
	return ""
}

type GrantType string

type GrantTypes []GrantType

var (
	GrantClientCredentials GrantType = "client_credentials"
	GrantAuthorizationCode GrantType = "authorization_code"
	GrantImplicit          GrantType = "implicit"
	GrantRefreshToken      GrantType = "refresh_token"
	GrantPassword          GrantType = "password"
)

func (gt GrantType) String() string {
	if gt == GrantClientCredentials ||
		gt == GrantAuthorizationCode ||
		gt == GrantImplicit ||
		gt == GrantPassword ||
		gt == GrantRefreshToken {
		return string(gt)
	}
	return ""
}

func (o *GrantTypes) Scan(src any) error {
	grantTypes := make([]GrantType, 0)
	for _, v := range strings.Split(src.(string), ",") {
		grantTypes = append(grantTypes, GrantType(v))
	}
	*o = grantTypes
	return nil
}

func (o GrantTypes) Value() (driver.Value, error) {
	if len(o) == 0 {
		return nil, nil
	}

	b, err := json.Marshal(o)
	return string(b), err
}

type RedirectUrls []string

func (o *RedirectUrls) Scan(src any) error {
	redirectUrls := make([]string, 0)
	for _, v := range strings.Split(src.(string), ",") {
		redirectUrls = append(redirectUrls, v)
	}
	*o = redirectUrls
	return nil
}

func (o RedirectUrls) Value() (driver.Value, error) {
	if len(o) == 0 {
		return nil, nil
	}
	b, err := json.Marshal(o)
	return string(b), err
}

type Scopes []string

func (o *Scopes) Scan(src any) error {
	ips := make([]string, 0)
	for _, v := range strings.Split(src.(string), ",") {
		ips = append(ips, v)
	}
	*o = ips
	return nil
}

func (o Scopes) Value() (driver.Value, error) {
	if len(o) == 0 {
		return nil, nil
	}
	b, err := json.Marshal(o)
	return string(b), err
}

// Application
type Application struct {
	ID                   string          `bson:"_id" gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey" json:"application_id" form:"app_id"`
	OrganizationID       string          `bson:"organization_id" gorm:"column:organization_id;type:uuid" json:"-"`
	UserID               *string         `bson:"user_id" gorm:"column:user_id;type:uuid" json:"-"`
	Type                 ApplicationType `bson:"type" gorm:"column:type" json:"application_type"`
	AndroidAppID         string          `bson:"android_app_id" gorm:"column:android_app_id" json:"android_app_id"`
	AndroidSHA1          string          `bson:"android_sha1" gorm:"column:android_sha1" json:"android_sha1"`
	AppBundleID          string          `bson:"app_bundle_id" gorm:"column:app_bundle_id" json:"app_bundle_id"`
	AppStoreID           string          `bson:"app_store_id" gorm:"column:app_store_id" json:"app_store_id"`
	LogoURI              string          `bson:"logo_uri" gorm:"column:logo_uri" json:"application_logo_uri"`
	Name                 string          `bson:"name" gorm:"column:name;" json:"application_name" form:"name"`
	DisplayName          string          `bson:"display_name" gorm:"display_name" json:"display_name"`
	ClientID             string          `bson:"client_id" gorm:"client_id" json:"client_id"`
	ClientSecret         string          `bson:"client_secret" gorm:"client_secret" json:"client_secret"`
	GrantTypes           pq.StringArray  `bson:"grant_types" gorm:"column:grant_types;type:TEXT" json:"grant_types"`
	RedirectUrls         pq.StringArray  `bson:"redirect_urls" gorm:"column:redirect_urls;type:TEXT" json:"redirect_urls"`
	Scopes               pq.StringArray  `bson:"scopes" gorm:"column:scopes;type:TEXT" json:"scopes"`
	AccessTokenExpiresIn int64           `bson:"access_token_expires_in" gorm:"column:access_token_expires_in" json:"access_token_expires_in"`
	CreatedAt            time.Time       `bson:"created_at" gorm:"column:created_at;default:now();not null" json:"created_at"`
	UpdatedAt            time.Time       `bson:"updated_at" gorm:"column:updated_at;default:now();not null" json:"updated_at"`
	DeletedAt            *time.Time      `bson:"deleted_at" gorm:"column:deleted_at;index" json:"-"`
}

func (m *Application) BeforeCreate(tx *gorm.DB) (err error) {
	current := time.Now()
	m.CreatedAt = current
	m.UpdatedAt = current
	return
}

func (m *Application) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()
	return
}

func NewApplication(organization string) *Application {
	return &Application{
		OrganizationID:       organization,
		AccessTokenExpiresIn: 86400,
		GrantTypes: pq.StringArray{
			GrantAuthorizationCode.String(),
			GrantClientCredentials.String(),
		},
		ClientID:     s.RandomHex(16),
		ClientSecret: s.RandomHex(32),
		Scopes: pq.StringArray{
			"openid",
			"email",
			"profile",
		},
	}
}

func (m *Application) SetName(name string) {
	m.Name = name
}

func (m *Application) SetClientID(clientID string) {
	m.ClientID = clientID
}

func (m *Application) SetClientSecret(clientSecret string) {
	m.ClientSecret = clientSecret
}

type RequestListApplication struct {
	Pagination
	Excludes       []string `form:"excludes"`
	OrganizationID string   `form:"organization_id"`
	ClientID       string   `form:"client_id"`
}

type ResponseListApplication struct {
	Data      []Application `json:"data"`
	TotalData int64         `json:"total_data"`
}

type RequestCreateApplication struct {
	OrganizationID string          `json:"-"`
	LogoUrl        string          `json:"logo_url"`
	DisplayName    string          `json:"display_name" validate:"required"`
	Type           ApplicationType `json:"type" validate:"required"`
	GrantTypes     GrantTypes      `json:"grant_types" validate:"required"`
	RedirectUrls   []string        `json:"redirect_urls" validate:"required"`
	AndroidAppID   string          `json:"android_app_id" validate:"required_if=Type android"`
	AndroidSHA1    string          `json:"android_sha1" validate:"required_if=Type android"`
}

type RequestUpdateApplication struct {
	ID             string         `json:"application_id" url:"application_id"`
	OrganizationID string         `json:"-"`
	LogoUrl        string         `json:"logo_url" validate:"required"`
	DisplayName    string         `json:"display_name" validate:"required"`
	RedirectUrls   pq.StringArray `json:"redirect_urls" validate:"required"`
}

type FilterApplication struct {
	Pagination
	Name         string         `form:"name"`
	Names        pq.StringArray `form:"-"`
	Organization string         `form:"organization"`
	ClientID     string         `form:"client_id"`
}

type TokenRequest struct {
	ClientID     string    `form:"client_id"`
	ClientSecret string    `form:"client_secret"`
	Username     string    `form:"username"`
	Password     string    `form:"password"`
	GrantType    GrantType `form:"grant_type" validate:"required"`
	Code         string    `form:"code"`
	RefreshToken string    `form:"refresh_token"`
	RedirectUri  string    `form:"redirect_uri" validate:"required"`
	Scope        string    `form:"scope" validate:"required"`
	State        string    `form:"state"`
	Nonce        string    `form:"nonce"`
}

type ResponseType string

var (
	Token ResponseType = "token"
	Code  ResponseType = "code"
)

func (rt ResponseType) String() string {
	if rt == Token ||
		rt == Code {
		return string(rt)
	}
	return ""
}

type AuthorizationRequest struct {
	UserID              string       `form:"-" url:"-" json:"-"`
	ClientID            string       `form:"client_id" url:"client_id" validate:"required"`
	ResponseType        ResponseType `form:"response_type" url:"response_type" validate:"required,oneof=code token"`
	RedirectUri         string       `form:"redirect_uri" url:"redirect_uri" validate:"required"`
	Scope               string       `form:"scope" url:"scope" validate:"required"`
	CodeChallenge       string       `form:"code_challenge" url:"code_challenge"`
	CodeChallengeMethod string       `form:"code_challenge_method" url:"code_challenge_method"`
	State               string       `form:"state" url:"state"`
	Nonce               string       `form:"nonce" url:"nonce"`
}

type RequestRevoke struct {
	ClientID     string `form:"client_id" json:"client_id" validate:"required"`
	ClientSecret string `form:"client_secret" json:"client_secret"`
	Token        string `form:"token" json:"token" validatE:"required"`
}
