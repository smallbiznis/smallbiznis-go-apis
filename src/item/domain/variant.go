package domain

import (
	"context"
	"time"

	"github.com/lib/pq"
	"github.com/smallbiznis/go-genproto/smallbiznis/item/v1"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type Variant struct {
	ID               string         `gorm:"column:variant_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"variant_id"`
	OrganizationID   string         `gorm:"column:organization_id;type:uuid;" json:"organization_id"`
	ItemID           string         `gorm:"column:item_id;type:uuid" json:"-"`
	Item             Item           `gorm:"foreignKey:ItemID" json:"item"`
	SKU              string         `gorm:"column:sku" json:"sku"`
	Title            string         `gorm:"column:title" json:"title"`
	Taxable          bool           `gorm:"column:taxable" json:"taxable"`
	Price            float32        `gorm:"column:price" json:"price"`
	CompareAtPrice   float32        `gorm:"column:compare_at_price" json:"compare_at_price"`
	Cost             float32        `gorm:"column:cost" json:"cost"`
	Barcode          string         `gorm:"column:barcode" json:"barcode"`
	Profit           float32        `gorm:"column:profit" json:"profit"`
	Margin           float32        `gorm:"column:margin" json:"margin"`
	Weight           float32        `gorm:"column:weight" json:"weight"`
	WeightUnit       string         `gorm:"column:weight_unit" json:"weight_unit"`
	Attributes       pq.StringArray `gorm:"column:attributes;type:TEXT;" json:"attributes"`
	PreparationTime  time.Duration  `gorm:"column:preparation_time" json:"preparation_time"`
	InventoryItemIds pq.StringArray `gorm:"column:inventory_item_ids;type:TEXT;" json:"inventory_item_ids"`
	CreatedAt        time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

func (m *Variant) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

func (m *Variant) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return
}

func (m *Variant) ToProto() *item.Variant {
	return &item.Variant{
		VariantId:      m.ID,
		ItemId:         m.ItemID,
		Sku:            m.SKU,
		Title:          m.Title,
		Taxable:        m.Taxable,
		Price:          m.Price,
		CompareAtPrice: m.CompareAtPrice,
		Cost:           m.Cost,
		Barcode:        m.Barcode,
		Profit:         m.Profit,
		Margin:         m.Margin,
		Weight:         m.Weight,
		WeightUnit:     item.WeightUnit(item.WeightUnit_value[m.WeightUnit]),
		Attributes:     m.Attributes,
		CreatedAt:      timestamppb.New(m.CreatedAt),
		UpdatedAt:      timestamppb.New(m.UpdatedAt),
	}
}

type Variants []Variant

func (m Variants) ToProto() (data []*item.Variant) {
	for _, v := range m {
		data = append(data, v.ToProto())
	}
	return
}

type IVariantRepository interface {
	Find(context.Context, pagination.Pagination, Variant) (Variants, int64, error)
	FindOne(context.Context, Variant) (*Variant, error)
	Save(context.Context, Variant) (*Variant, error)
	BatchSave(context.Context, []Variant) error
	Update(context.Context, Variant) (*Variant, error)
	Delete(context.Context, Variant) error
}
