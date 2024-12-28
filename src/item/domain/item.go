package domain

import (
	"context"
	"strings"
	"time"

	"github.com/smallbiznis/go-genproto/smallbiznis/item/v1"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type Type string

var (
	Physical Type = "physical"
	Menu     Type = "menu"
)

func (m Type) String() string {
	if m == Physical ||
		m == Menu {
		return string(m)
	}
	return ""
}

func (m Type) ToProto() item.Type {
	if m == Physical {
		return item.Type_physical
	}
	return item.Type_menu
}

type Status string

var (
	All      Status = ""
	ACTIVE   Status = "active"
	DRAFT    Status = "draft"
	ARCHIVED Status = "archived"
)

func (m Status) String() string {
	if m == ACTIVE ||
		m == DRAFT ||
		m == ARCHIVED {
		return string(m)
	}
	return ""
}

func (m Status) ToProto() item.Status {
	if m == ACTIVE {
		return item.Status_active
	} else if m == DRAFT {
		return item.Status_draft
	} else if m == ARCHIVED {
		return item.Status_archived
	}

	return All.ToProto()
}

type Item struct {
	ID             string         `gorm:"column:item_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"item_id"`
	OrganizationID string         `gorm:"column:organization_id;type:uuid" json:"organization_id"`
	Type           string         `gorm:"column:type" json:"type"`
	Title          string         `gorm:"column:title" json:"title"`
	Slug           string         `gorm:"column:slug" json:"slug"`
	BodyHTML       string         `gorm:"column:body_html" json:"body_html"`
	Status         string         `gorm:"column:status" json:"status"`
	Taxable        bool           `gorm:"column:taxable" json:"taxable"`
	Options        ItemOptions    `gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"options"`
	Variants       Variants       `gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"variants"`
	CreatedAt      time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

func (m *Item) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

func (m *Item) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return
}

func (m *Item) ItemCode() (code string) {
	codes := strings.Split(m.Title, " ")
	for _, value := range codes {
		code += string([]rune(value)[0])
	}
	return
}

func (m *Item) ToProto() *item.Item {
	return &item.Item{
		ItemId:         m.ID,
		Type:           item.Type(item.Type_value[m.Type]),
		OrganizationId: m.OrganizationID,
		Title:          m.Title,
		Slug:           m.Slug,
		BodyHtml:       m.BodyHTML,
		Options:        m.Options.ToProto(),
		Variants:       m.Variants.ToProto(),
		Status:         item.Status(item.Status_value[m.Status]),
		CreatedAt:      timestamppb.New(m.CreatedAt),
		UpdatedAt:      timestamppb.New(m.UpdatedAt),
	}
}

type Items []Item

func (m Items) ToProto() (data []*item.Item) {
	for _, v := range m {
		data = append(data, v.ToProto())
	}
	return
}

type IItemRepository interface {
	Find(context.Context, pagination.Pagination, Item) (Items, int64, error)
	FindOne(context.Context, Item) (*Item, error)
	Save(context.Context, Item) (*Item, error)
	Update(context.Context, Item) (*Item, error)
	Delete(context.Context, Item) error
}
