package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/google/uuid"
	"github.com/smallbiznis/oauth2-server/internal/pkg/token"
	"github.com/smallbiznis/oauth2-server/model"
	"gorm.io/gorm"
)

var (
	organizationID   = uuid.NewString()
	serviceAccountID = uuid.NewString()
)

func SeedData(db *gorm.DB) {
	if db.Migrator().HasTable(&model.Organization{}) {
		defaultOrganization(db)
	}

	if db.Migrator().HasTable(&model.Provider{}) {
		defaultProvider(db)
	}

	if db.Migrator().HasTable(&model.Account{}) {
		defaultUser(db)
	}

	if db.Migrator().HasTable(&model.Application{}) {
		defaultApplication(db)
	}

}

func defaultOrganization(db *gorm.DB) (err error) {
	organization := model.Organization{
		ID:          "48649ef7-5d70-496c-a887-aa77377d0739",
		Name:        "accounts",
		DisplayName: "Default Organization",
		IsDefault:   true,
	}

	if err = db.Model(&model.Organization{}).Where(&model.Organization{
		Name:      organization.Name,
		IsDefault: organization.IsDefault,
	}).First(&organization).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = db.Model(&model.Organization{}).Create(&organization).Error; err != nil {
				return
			}

			return defaultJwks(db, organization)
		}
		return
	}

	return
}

func defaultJwks(db *gorm.DB, organization model.Organization) (err error) {
	var keys []model.OrganizationKey
	if err = db.Model(&model.OrganizationKey{}).Where(&model.OrganizationKey{
		OrganizationID: organization.ID,
	}).Find(&keys).Error; err != nil {
		return
	}

	if len(keys) == 0 {

		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return err
		}

		// Encode to PEM
		privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
		privateKeyPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		})

		// key := strings.RandomHex(32)
		keyId, _ := token.New(32)
		newKeys := []model.OrganizationKey{
			{
				OrganizationID:              organization.ID,
				Algorithm:                   jose.RS256,
				Key:                         string(privateKeyPEM),
				KeyID:                       keyId,
				Use:                         model.Sign,
				Certificates:                []string{},
				CertificateThumbprintSHA256: []byte(""),
			},
		}

		return db.Model(&model.OrganizationKey{}).Create(&newKeys).Error
	}

	return
}

func defaultProvider(db *gorm.DB) (err error) {
	providers := []model.Provider{
		{
			Title:    "Google",
			Category: model.OAuth2,
			Status:   "active",
		},
		{
			Title:    "Facebook",
			Category: model.OAuth2,
			Status:   "active",
		},
		{
			Title:    "Github",
			Category: model.OAuth2,
			Status:   "active",
		},
		{
			Title:    "Twillio",
			Category: model.SMS,
			Status:   "active",
		},
		{
			Title:    "Twillio",
			Category: model.Email,
			Status:   "active",
		},
	}

	for _, provider := range providers {
		if err = db.Model(&model.Provider{}).Where(&model.Provider{
			Title:    provider.Title,
			Category: provider.Category,
		}).First(&provider).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err = db.Model(&model.Provider{}).Create(&provider).Error; err != nil {
					return
				}
			}
			continue
		}
	}

	return
}

func defaultApplication(db *gorm.DB) (err error) {

	now := time.Now()
	application := model.Application{
		Type:           model.Web,
		Name:           "app-built-in",
		OrganizationID: organizationID,
		DisplayName:    "Application Default",
		UserID:         &serviceAccountID,
		ClientID:       "c1e61dc053b1ec9ded053bd90a6bdbc9",
		ClientSecret:   "21635f4b55e3729895aaa53881717742d3d57a0a026ec959d62956d577238fa6",
		GrantTypes: []string{
			model.GrantAuthorizationCode.String(),
			model.GrantClientCredentials.String(),
		},
		RedirectUrls: []string{
			"http://localhost:3000/api/callback",
		},
		AccessTokenExpiresIn: 86400,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	if err = db.Model(&model.Application{}).Where(&model.Application{
		Name:           application.Name,
		OrganizationID: application.OrganizationID,
	}).First(&application).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return db.Model(&model.Application{}).Create(&application).Error
		}
		return
	}

	return
}

func defaultUser(db *gorm.DB) (err error) {
	// pwd := "$2a$10$SD2p8PvDBjLp1mHlH64xHuq/2xMuuKlqQmu.beFXmEqDo5d6chzd6"
	defaultUsername := "f10dc37244ad@smallbiznis.test"
	users := model.Account{
		ID:             serviceAccountID,
		OrganizationID: organizationID,
		Type:           model.ServiceAccount,
		Provider:       "",
		Username:       defaultUsername,
		Password:       "",
		Roles:          []string{},
	}

	if err = db.Model(&model.Account{}).First(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return db.Model(&model.Account{}).Create(&users).Error
		}
		return
	}

	return
}
