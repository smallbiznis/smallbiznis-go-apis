package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/organization/domain"
	"gorm.io/gorm"
)

type locationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) domain.ILocationRepository {
	return &locationRepository{db}
}

func (r *locationRepository) Find(ctx context.Context, p pagination.Pagination, f domain.Location) (orgs domain.Locations, count int64, err error) {
	stmt := r.db.WithContext(ctx).Model(&domain.Location{}).Preload("Country")

	if f.OrganizationID != "" {
		stmt.Where("organization_id = ?", f.OrganizationID)
	}

	if f.ID != "" {
		stmt.Where("id = ?", f.ID)
	}

	if len(f.IDs) > 0 {
		stmt.Where("id in (?)", f.IDs)
	}

	if p.SortBy != "" && p.OrderBy != "" {
		stmt.Order(fmt.Sprintf("%s %s", p.SortBy, p.OrderBy))
	} else {
		stmt.Order("updated_at DESC")
	}

	if err = stmt.Count(&count).Scopes(p.Paginate()).Find(&orgs).Error; err != nil {
		return
	}

	return
}

func (r *locationRepository) FindOne(ctx context.Context, f domain.Location) (org *domain.Location, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.Location{}).Preload("Country").Where(&f).First(&org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return
	}

	return
}

func (r *locationRepository) Save(ctx context.Context, d domain.Location) (org *domain.Location, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.Location{}).Preload("Country").Create(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.Location{ID: d.ID})
}
