package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/user/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) domain.IUserRepository {
	return &userRepository{
		db,
	}
}

func (r *userRepository) Find(ctx context.Context, p pagination.Pagination, f domain.User) (orgs domain.Users, count int64, err error) {
	stmt := r.db.WithContext(ctx).Model(&domain.User{}).
		Where(&f).
		Count(&count).
		Scopes(p.Paginate())

	if p.SortBy != "" && p.OrderBy != "" {
		stmt.Order(fmt.Sprintf("%s %s", p.SortBy, p.OrderBy))
	} else {
		stmt.Order("updated_at DESC")
	}

	if err = stmt.Find(&orgs).Error; err != nil {
		return
	}

	return
}

func (r *userRepository) FindOne(ctx context.Context, f domain.User) (org *domain.User, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.User{}).Where(&f).First(&org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return
	}

	return
}

func (r *userRepository) Save(ctx context.Context, d domain.User) (org *domain.User, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.User{}).Create(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.User{ID: d.ID})
}

func (r *userRepository) Update(ctx context.Context, d domain.User) (org *domain.User, err error) {
	if err = r.db.WithContext(ctx).Save(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.User{ID: d.ID})
}

func (r *userRepository) Delete(ctx context.Context, org domain.User) (err error) {
	return r.db.WithContext(ctx).Model(&domain.User{}).Delete(&org).Error
}
