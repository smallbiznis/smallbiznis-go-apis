package domain

import (
	"context"
	"time"

	"github.com/lib/pq"
	"github.com/smallbiznis/go-genproto/smallbiznis/member/v1"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type Member struct {
	ID             string         `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey" uri:"id" json:"id"`
	OrganizationID string         `gorm:"column:organization_id;type:uuid;" json:"organization_id"`
	UserID         string         `gorm:"column:user_id;type:uuid;" json:"user_id"`
	Roles          pq.StringArray `gorm:"column:roles;type:VARCHAR(255);" json:"roles"`
	CreatedAt      time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

func (m *Member) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

func (m *Member) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return
}

func (m *Member) ToProto() *member.Member {
	proto := &member.Member{
		MemberId:  m.ID,
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}

	for _, v := range m.Roles {
		proto.Roles = append(proto.Roles, member.Role(member.Role_value[v]))
	}

	return proto
}

type Members []Member

type IMemberRepository interface {
	Find(context.Context, pagination.Pagination, Member) (Members, int64, error)
	FindOne(context.Context, Member) (*Member, error)
	Save(context.Context, Member) (*Member, error)
	Update(context.Context, Member) (*Member, error)
	Delete(context.Context, Member) error
}
