package main

import (
	"github.com/smallbiznis/user/domain"
	"gorm.io/gorm"
)

func Automigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
	)
}

func Migrate(db *gorm.DB) (err error) {
	return db.Transaction(func(tx *gorm.DB) (err error) {

		if tx.Migrator().HasTable(&domain.User{}) {
			// org := domain.Organization{
			// 	OrganizationID: "smallbiznis",
			// 	Title:          env.Lookup("DEFAULT_ORGANIZATION", "SmallBiznis"),
			// 	IsDefault:      true,
			// 	Status:         "ACTIVE",
			// }

			// if err = tx.Where("organization_id = ?", org.ID).First(&org).Error; err != nil {
			// 	if errors.Is(err, gorm.ErrRecordNotFound) {
			// 		if err = tx.Create(&org).First(&org).Error; err != nil {
			// 			fmt.Printf("err: %v", err)
			// 			return
			// 		}
			// 	}
			// }

			// apps := []domain.App{
			// 	{
			// 		ID:             "online-store",
			// 		OrganizationID: org.ID,
			// 		Title:          "Online Store",
			// 		IsDefault:      true,
			// 	},
			// 	{
			// 		ID:             "point-of-sales",
			// 		OrganizationID: org.ID,
			// 		Title:          "Point Of Sales",
			// 		IsDefault:      true,
			// 	},
			// }

			// var count int64
			// if err = tx.Model(&domain.App{}).Where("app_id in (?)", []string{
			// 	"online-store",
			// 	"point-of-sales",
			// }).Count(&count).Error; err != nil {
			// 	fmt.Printf("err: %v", err)
			// 	return
			// }

			// if count == 0 {
			// 	if err = tx.Create(&apps).Error; err != nil {
			// 		fmt.Printf("err: %v", err)
			// 		return
			// 	}
			// }

			return
		}

		return
	})
}
