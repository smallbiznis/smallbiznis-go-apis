package main

import (
	"github.com/smallbiznis/member/domain"
	"gorm.io/gorm"
)

func Automigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.Member{},
	)
}
