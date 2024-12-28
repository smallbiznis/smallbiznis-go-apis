package main

import (
	"github.com/smallbiznis/inventory/domain"
	"gorm.io/gorm"
)

func Automigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.InventoryItem{},
	)
}
