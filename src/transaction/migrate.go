package main

import (
	"github.com/smallbiznis/transaction/domain"
	"gorm.io/gorm"
)

func Automigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.Order{},
		// &domain.OrderBillingAddress{},
		// &domain.OrderShippingAddress{},
		&domain.OrderItem{},
		// &domain.OrderPayment{},
		// &domain.OrderFulfillment{},
		// &domain.OrderShipping{},
		// &domain.ShippingHistory{},
	)
}

func Migrate(db *gorm.DB) (err error) {
	return db.Transaction(func(tx *gorm.DB) (err error) {
		return
	})
}
