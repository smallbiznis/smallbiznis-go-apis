package main

import (
	"github.com/smallbiznis/customer/domain"
	"gorm.io/gorm"
)

func Automigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.Customer{},
		&domain.Addreses{},
	)
}

func Migrate(db *gorm.DB) (err error) {
	return db.Transaction(func(tx *gorm.DB) (err error) {
		return
	})
}
