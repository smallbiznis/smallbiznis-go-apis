package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/organization/domain"
	"gorm.io/gorm"
)

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) domain.IOrganizationRepository {
	return &organizationRepository{
		db,
	}
}

func (r *organizationRepository) Find(ctx context.Context, p pagination.Pagination, f domain.Organization) (orgs domain.Organizations, count int64, err error) {
	stmt := r.db.WithContext(ctx).Model(&domain.Organization{}).Preload("Country").
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

func (r *organizationRepository) FindOne(ctx context.Context, f domain.Organization) (org *domain.Organization, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.Organization{}).Preload("Country").Where(&f).First(&org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return
	}

	return
}

func (r *organizationRepository) Save(ctx context.Context, d domain.Organization) (org *domain.Organization, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.Organization{}).Create(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.Organization{ID: d.ID})
}

func (r *organizationRepository) Update(ctx context.Context, d domain.Organization) (org *domain.Organization, err error) {
	if err = r.db.WithContext(ctx).Save(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.Organization{ID: d.ID})
}

func (r *organizationRepository) Delete(ctx context.Context, org domain.Organization) (err error) {
	return r.db.WithContext(ctx).Model(&domain.Organization{}).Delete(&org).Error
}
