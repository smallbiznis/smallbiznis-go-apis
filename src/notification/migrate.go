package main

import (
	"gorm.io/gorm"
)

func Automigrate(db *gorm.DB) error {
	return db.AutoMigrate()
}

func Migrate(db *gorm.DB) (err error) {
	return db.Transaction(func(tx *gorm.DB) (err error) {
		return
	})
}
