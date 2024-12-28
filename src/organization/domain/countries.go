package domain

import (
	"context"

	"github.com/smallbiznis/go-genproto/smallbiznis/organization/v1"
	"github.com/smallbiznis/go-lib/pkg/pagination"
)

type Countries struct {
	CountryCode  string `gorm:"column:country_code;primaryKey"`
	CountryName  string `gorm:"column:country_name"`
	Continent    string `gorm:"column:continent"`
	CurrencyCode string `gorm:"column:currency_code"`
}

func (m *Countries) ToProto() *organization.Country {
	return &organization.Country{
		CountryCode:  m.CountryCode,
		CountyName:   m.CountryName,
		Continent:    m.Continent,
		CurrencyCode: m.CurrencyCode,
	}
}

type ICountryRepository interface {
	Find(context.Context, pagination.Pagination, Countries) ([]*Countries, int64, error)
	FindOne(context.Context, Countries) (*Countries, error)
}
