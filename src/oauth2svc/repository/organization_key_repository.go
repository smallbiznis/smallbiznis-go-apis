package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/smallbiznis/oauth2-server/model"
	"gorm.io/gorm"
)

type IOrganizationKeyRepository interface {
	Find(ctx context.Context, p *model.Pagination, f model.OrganizationKey) (keys []model.OrganizationKey, count int64, err error)
	FindOne(ctx context.Context, f model.OrganizationKey) (user *model.OrganizationKey, err error)
	Create(ctx context.Context, u model.OrganizationKey) (user *model.OrganizationKey, err error)
}

type organizationKeyRepository struct {
	db *gorm.DB
}

func NewOrganizationKeyRepository(db *gorm.DB) IOrganizationKeyRepository {
	return &organizationKeyRepository{
		db,
	}
}

func (r *organizationKeyRepository) Find(ctx context.Context, p *model.Pagination, f model.OrganizationKey) (keys []model.OrganizationKey, count int64, err error) {
	stmt := r.db.WithContext(ctx).Model(&model.OrganizationKey{}).Where(&f).Count(&count)

	if p != nil {
		stmt.Limit(10)
		if p.Limit > 0 {
			stmt.Limit(int(p.Limit))
		}

		stmt.Offset(0)
		if p.Offset > 0 {
			stmt.Offset(int(p.Offset))
		}

		if p.SortBy == "" {
			stmt.Order("created_at DESC")
		} else {
			stmt.Order(fmt.Sprintf("%s %s", p.SortBy, p.OrderBy))
		}
	}

	if err = stmt.Find(&keys).Error; err != nil {
		return
	}

	return
}

func (r *organizationKeyRepository) FindOne(ctx context.Context, f model.OrganizationKey) (user *model.OrganizationKey, err error) {
	if err := r.db.WithContext(ctx).Model(&model.OrganizationKey{}).Where(&f).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *organizationKeyRepository) Create(ctx context.Context, u model.OrganizationKey) (user *model.OrganizationKey, err error) {
	if err := r.db.WithContext(ctx).Model(&model.OrganizationKey{}).Create(&u).Error; err != nil {
		return nil, err
	}
	return r.FindOne(ctx, model.OrganizationKey{ID: u.ID})
}

func (r *organizationKeyRepository) Update(ctx context.Context, u model.OrganizationKey) (user *model.OrganizationKey, err error) {
	if err := r.db.WithContext(ctx).Model(&model.OrganizationKey{}).Save(&u).Error; err != nil {
		return nil, err
	}
	return r.FindOne(ctx, model.OrganizationKey{ID: u.ID})
}
