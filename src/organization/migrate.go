package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gosimple/slug"
	"github.com/smallbiznis/go-lib/pkg/env"
	"github.com/smallbiznis/organization/domain"
	"gorm.io/gorm"
)

func Automigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.Countries{},
		&domain.Province{},
		&domain.Regencies{},
		&domain.District{},
		&domain.Organization{},
		&domain.Location{},
		&domain.TaxRule{},
		&domain.ShippingRate{},
	)
}

func Migrate(db *gorm.DB) (err error) {
	return db.Transaction(func(tx *gorm.DB) (err error) {

		if tx.Migrator().HasTable(&domain.Organization{}) {
			defaultOrg := strings.Trim(env.Lookup("DEFAULT_ORGANIZATION", "SmallBiznis"), " ")
			org := domain.Organization{
				ID:             "0648d98a-3efa-47f7-b6ac-d3239a1559fd",
				OrganizationID: slug.Make(defaultOrg),
				Title:          defaultOrg,
				CountryID:      "ID",
				IsDefault:      true,
				Status:         "ACTIVE",
			}

			exist := domain.Organization{}
			if err = tx.Model(&domain.Organization{}).Where("organization_id = ?", org.OrganizationID).First(&exist).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					if err = tx.Create(&org).Error; err != nil {
						fmt.Printf("err: %v", err)
						return
					}
				}
			}
			return
		}
		return
	})
}
