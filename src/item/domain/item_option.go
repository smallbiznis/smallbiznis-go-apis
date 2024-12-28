package domain

import (
	"time"

	"github.com/lib/pq"
	"github.com/smallbiznis/go-genproto/smallbiznis/item/v1"
	"gorm.io/gorm"
)

// ProductOption
type ItemOption struct {
	ID        string         `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OptionID  string         `gorm:"column:option_id;type:uuid;" json:"option_id"`
	ItemID    string         `gorm:"column:item_id;type:uuid;" json:"item_id"`
	Option    Option         `gorm:"foreignKey:OptionID" json:"options"`
	Position  int32          `gorm:"column:position" json:"position"`
	Values    pq.StringArray `gorm:"column:values;type:VARCHAR(255);" json:"values"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

func (m *ItemOption) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

func (m *ItemOption) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return
}

func (m *ItemOption) ToProto() *item.Item_Option {
	return &item.Item_Option{
		OptionId:   m.OptionID,
		OptionName: m.Option.Name,
		Position:   m.Position,
		Values:     m.Values,
	}
}

type ItemOptions []ItemOption

func (m ItemOptions) ToProto() (data []*item.Item_Option) {
	for _, v := range m {
		data = append(data, v.ToProto())
	}
	return
}
