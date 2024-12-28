package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/organization/domain"
	"gorm.io/gorm"
)

type shippingRateRepository struct {
	db *gorm.DB
}

func NewShippingRateRepository(
	db *gorm.DB,
) domain.IShippingRateRepository {
	return &shippingRateRepository{
		db,
	}
}

func (r *shippingRateRepository) Find(ctx context.Context, p pagination.Pagination, f domain.ShippingRate) (rules domain.ShippingRates, count int64, err error) {
	stmt := r.db.WithContext(ctx).Model(&domain.ShippingRate{}).
		Where(&f).
		Count(&count).
		Scopes(p.Paginate())

	if p.SortBy != "" && p.OrderBy != "" {
		stmt.Order(fmt.Sprintf("%s %s", p.SortBy, p.OrderBy))
	} else {
		stmt.Order("updated_at DESC")
	}

	if err = stmt.Find(&rules).Error; err != nil {
		return
	}

	return
}

func (r *shippingRateRepository) FindOne(ctx context.Context, f domain.ShippingRate) (org *domain.ShippingRate, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.ShippingRate{}).
		Where(&f).First(&org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return
	}

	return
}

func (r *shippingRateRepository) Save(ctx context.Context, d domain.ShippingRate) (org *domain.ShippingRate, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.ShippingRate{}).Create(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.ShippingRate{ID: d.ID})
}

func (r *shippingRateRepository) Update(ctx context.Context, d domain.ShippingRate) (org *domain.ShippingRate, err error) {
	if err = r.db.WithContext(ctx).Save(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.ShippingRate{ID: d.ID})
}
