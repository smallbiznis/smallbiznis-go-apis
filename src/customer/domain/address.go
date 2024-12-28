package domain

import (
	"context"
	"time"

	"github.com/lib/pq"
	"github.com/smallbiznis/go-genproto/smallbiznis/customer/v1"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type GeoPoint struct {
	Latitude  float32 `gorm:"column:latitude"`
	Longitude float32 `gorm:"column:longitude"`
}

type Address struct {
	ID           string         `gorm:"column:address_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"address_id"`
	CustomerID   string         `gorm:"column:customer_id;type:uuid;" json:"customer_id"`
	Customer     Customer       `gorm:"foreignKey:CustomerID" json:"-"`
	ContactName  string         `gorm:"column:contact_name" json:"contact_name"`
	ContactPhone string         `gorm:"column:contact_phone" json:"contact_phone"`
	ProvinceID   *string        `gorm:"column:province_id;type:uuid;" json:"province_id"`
	CityID       *string        `gorm:"column:city_id;type:uuid" json:"city_id"`
	DistrictID   *string        `gorm:"column:district_id;type:uuid;" json:"district_id"`
	Address      string         `gorm:"column:address" json:"address"`
	PostalCode   string         `gorm:"column:postal_code" json:"postal_code"`
	GeoPoint     pq.StringArray `gorm:"column:geo_point;type:VARCHAR(255);" json:"geo_point"`
	IsDefault    bool           `gorm:"column:is_default" json:"is_default"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

func (m *Address) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

func (m *Address) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return
}

func (m *Address) ToProto() *customer.Address {
	addr := &customer.Address{
		AddressId:    m.ID,
		CustomerId:   m.CustomerID,
		ContactName:  m.ContactName,
		ContactPhone: m.ContactPhone,
		Address:      m.Address,
		PostalCode:   m.PostalCode,
		IsDefault:    m.IsDefault,
		CreatedAt:    timestamppb.New(m.CreatedAt),
		UpdatedAt:    timestamppb.New(m.UpdatedAt),
	}

	if m.ProvinceID != nil {
		addr.ProviceId = *m.ProvinceID
	}

	if m.CityID != nil {
		addr.CityId = *m.CityID
	}

	if m.DistrictID != nil {
		addr.DistrictId = *m.DistrictID
	}

	return addr
}

type Addreses []Address

func (m Addreses) ToProto() (data []*customer.Address) {
	for _, v := range m {
		data = append(data, v.ToProto())
	}
	return
}

type IAddressRepository interface {
	Find(context.Context, pagination.Pagination, Address) (Addreses, int64, error)
	FindOne(context.Context, Address) (*Address, error)
	Save(context.Context, Address) (*Address, error)
}
