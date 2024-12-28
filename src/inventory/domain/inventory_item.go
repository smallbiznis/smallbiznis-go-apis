package domain

import (
	"context"
	"time"

	"github.com/smallbiznis/go-genproto/smallbiznis/inventory/v1"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type InventoryItem struct {
	ID               string    `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey" uri:"id" json:"id"`
	OrganizationID   string    `gorm:"column:organization_id;type:uuid;" json:"organization_id"`
	LocationID       string    `gorm:"column:location_id;type:uuid;" json:"location_id"`
	ItemID           string    `gorm:"column:item_id;type:uuid" json:"item_id"`
	Quantity         int32     `gorm:"column:quantity" json:"quantity"`
	ReservedQuantity int32     `gorm:"column:reserved_quantity" json:"reserved_quantity"`
	ReorderLevel     int32     `gorm:"column:reorder_level" json:"reorder_level"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (m *InventoryItem) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

func (m *InventoryItem) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return
}

func (m *InventoryItem) ToProto() *inventory.Inventory {
	return &inventory.Inventory{
		InventoryItemId:  m.ID,
		OrganizationId:   m.OrganizationID,
		LocationId:       m.LocationID,
		ItemId:           m.ItemID,
		Quantity:         m.Quantity,
		ReservedQuantity: m.ReservedQuantity,
		CreatedAt:        timestamppb.New(m.CreatedAt),
		UpdatedAt:        timestamppb.New(m.UpdatedAt),
	}
}

type InventoryItems []InventoryItem

func (m InventoryItems) ToProto() (data []*inventory.Inventory) {
	for _, v := range m {
		data = append(data, v.ToProto())
	}
	return
}

type IInventoryItemRepository interface {
	Find(context.Context, pagination.Pagination, InventoryItem) (InventoryItems, int64, error)
	FindOne(context.Context, InventoryItem) (*InventoryItem, error)
	Save(context.Context, InventoryItem) (*InventoryItem, error)
	Update(context.Context, InventoryItem) (*InventoryItem, error)
	Delete(context.Context, InventoryItem) error
}
