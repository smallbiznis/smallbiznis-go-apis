package domain

import (
	"context"
	"time"

	"github.com/smallbiznis/go-genproto/smallbiznis/user/v1"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type User struct {
	ID          string         `gorm:"column:user_id;type:uuid;default:uuid_generate_v4();primaryKey" uri:"user_id" json:"user_id"`
	AvatarURI   string         `gorm:"column:avatar_uri" json:"avatar_uri"`
	Username    string         `gorm:"column:username;uniqueIndex" form:"username" json:"username"`
	Email       string         `gorm:"column:email" form:"email" json:"email"`
	PhoneNumber string         `gorm:"column:phone_number" form:"phone_number" json:"phone_number"`
	FirstName   string         `gorm:"column:first_name" json:"first_name"`
	LastName    string         `gorm:"column:last_name" json:"last_name"`
	BirthDate   *string        `gorm:"column:birth_date;type:date" json:"birth_date"`
	VerifiedAt  *time.Time     `gorm:"column:verified_at" json:"verified_at"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

func (m *User) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

func (m *User) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return
}

func (m *User) ToProto() *user.User {
	proto := &user.User{
		UserId:    m.ID,
		AvatarUri: m.AvatarURI,
		Email:     m.Email,
		FirstName: m.FirstName,
		LastName:  m.LastName,
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}
	if m.VerifiedAt != nil {
		proto.VerifiedAt = timestamppb.New(*m.VerifiedAt)
	}

	return proto
}

type Users []User

func (m Users) ToProto() (data []*user.User) {
	for _, v := range m {
		data = append(data, v.ToProto())
	}
	return
}

type IUserRepository interface {
	Find(context.Context, pagination.Pagination, User) (Users, int64, error)
	FindOne(context.Context, User) (*User, error)
	Save(context.Context, User) (*User, error)
	Update(context.Context, User) (*User, error)
	Delete(context.Context, User) error
}
