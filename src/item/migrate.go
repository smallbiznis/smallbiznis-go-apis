package main

import (
	"github.com/smallbiznis/item/domain"
	"gorm.io/gorm"
)

func Automigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.Option{},
		&domain.OptionValue{},
		&domain.Item{},
		&domain.ItemOption{},
		&domain.Variant{},
	)
}
