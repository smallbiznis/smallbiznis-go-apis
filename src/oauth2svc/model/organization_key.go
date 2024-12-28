package model

import (
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UseType string

var (
	Sign UseType = "sig"
	Enc  UseType = "enc"
)

func (m UseType) String() string {
	if m == Sign ||
		m == Enc {
		return string(m)
	}
	return ""
}

type OrganizationKey struct {
	ID             string                  `gorm:"column:id;type:uuid;default:uuid_generate_v4();" json:"-"`
	OrganizationID string                  `gorm:"column:organization_id;type:uuid" json:"-"`
	Key            string                  `gorm:"column:key;not null" json:"key"`
	KeyID          string                  `gorm:"column:key_id;not null" json:"key_id"`
	Algorithm      jose.SignatureAlgorithm `gorm:"column:algorithm;not null" json:"algorithm"`
	Use            UseType                 `gorm:"column:use;not null" json:"use"`
	Certificates   pq.StringArray          `gorm:"column:certificates;type:VARCHAR(255)" json:"certificates"`
	// CertificatesURL             *URL                    `gorm:"column:certificates_url" json:"certificates_url"`
	CertificateThumbprintSHA256 []byte    `gorm:"column:certificate_thumbprint" json:"certificate_thumbprint"`
	CreatedAt                   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt                   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type OrganizationKeys []OrganizationKey

func (m *OrganizationKey) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

func (m *OrganizationKey) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return
}
