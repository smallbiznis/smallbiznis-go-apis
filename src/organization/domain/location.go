package domain

import (
	"context"
	"time"

	"github.com/smallbiznis/go-genproto/smallbiznis/organization/v1"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"gorm.io/gorm"
)

type Province struct {
	ID   string `gorm:"column:id;primarykey" json:"provice_id"`
	Name string `gorm:"column:name" json:"name"`
}

type Regencies struct {
	ID        string `gorm:"column:id;primarykey" json:"city_id"`
	ProviceID string `gorm:"column:provice_id" json:"province_id"`
	Name      string `gorm:"column:name" json:"name"`
}

type District struct {
	ID     string `gorm:"column:id;primarykey" json:"district_id"`
	CityID string `gorm:"column:city_id;type:uuid" json:"city_id"`
	Name   string `gorm:"column:name" json:"name"`
}

type Location struct {
	ID             string    `gorm:"column:location_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"location_id"`
	IDs            []string  `gorm:"-" json:"location_ids"`
	OrganizationID string    `gorm:"column:organization_id;type:uuid" json:"organization_id"`
	Name           string    `gorm:"column:location_name" json:"location_name"`
	ContactName    string    `gorm:"column:contact_name" json:"contact_name"`
	ContactPhone   string    `gorm:"column:contact_phone" json:"contact_phone"`
	CountryID      string    `gorm:"column:country_id" json:"country_id"`
	Country        Countries `gorm:"->;foreignKey:CountryID" json:"-"`
	ProvinceID     *string   `gorm:"column:province_id" json:"province_id"`
	Province       Province  `gorm:"foreignKey:ProvinceID" json:"-"`
	RegencyID      *string   `gorm:"column:regency_id" json:"city_id"`
	Regencies      Regencies `gorm:"foreignKey:RegencyID" json:"-"`
	DistrictID     *string   `gorm:"column:district_id" json:"district_id"`
	District       District  `gorm:"foreignKey:DistrictID" json:"-"`
	Address        string    `gorm:"column:address" json:"address"`
	PostalCode     string    `gorm:"column:postal_code" json:"postal_code"`

	IsDefault bool           `gorm:"column:is_default"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

func (m *Location) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

func (m *Location) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return
}

func (m *Location) ToProto() *organization.Location {
	loc := &organization.Location{
		LocationId:     m.ID,
		OrganizationId: m.OrganizationID,
		Name:           m.Name,
		ContactName:    m.ContactName,
		ContactPhone:   m.ContactPhone,
		Country:        m.Country.ToProto(),
		IsDefault:      m.IsDefault,
	}

	if m.ProvinceID != nil {
		loc.ProviceId = *m.ProvinceID
	}

	if m.RegencyID != nil {
		loc.CityId = *m.RegencyID
	}

	if m.DistrictID != nil {
		loc.DistrictId = *m.DistrictID
	}

	return loc
}

type Locations []Location

func (m Locations) ToProto() (data []*organization.Location) {
	for _, v := range m {
		data = append(data, v.ToProto())
	}
	return
}

type ILocationRepository interface {
	Find(context.Context, pagination.Pagination, Location) (Locations, int64, error)
	FindOne(context.Context, Location) (*Location, error)
	Save(context.Context, Location) (*Location, error)
}
