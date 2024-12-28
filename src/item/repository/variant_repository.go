package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/item/domain"
	"gorm.io/gorm"
)

type variantRepository struct {
	db *gorm.DB
}

func NewVariantRepository(db *gorm.DB) domain.IVariantRepository {
	return &variantRepository{db}
}

func (r *variantRepository) Find(ctx context.Context, p pagination.Pagination, f domain.Variant) (product domain.Variants, count int64, err error) {
	stmt := r.db.WithContext(ctx).Model(&domain.Variant{}).
		Preload("Item").
		Where(&f).
		Count(&count).
		Scopes(p.Paginate())

	if p.SortBy != "" && p.OrderBy != "" {
		stmt.Order(fmt.Sprintf("%s %s", p.SortBy, p.OrderBy))
	} else {
		stmt.Order("updated_at DESC")
	}

	if err = stmt.Find(&product).Error; err != nil {
		return
	}

	return
}

func (r *variantRepository) FindOne(ctx context.Context, f domain.Variant) (org *domain.Variant, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.Variant{}).
		Preload("Item").
		Where(&f).First(&org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return
	}

	return
}

func (r *variantRepository) Save(ctx context.Context, d domain.Variant) (org *domain.Variant, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.Variant{}).Create(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.Variant{ID: d.ID})
}

func (r *variantRepository) BatchSave(ctx context.Context, d []domain.Variant) (err error) {
	return r.db.WithContext(ctx).Model(&domain.Variant{}).Create(&d).Error
}

func (r *variantRepository) Update(ctx context.Context, d domain.Variant) (org *domain.Variant, err error) {
	if err = r.db.WithContext(ctx).Save(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.Variant{ID: d.ID})
}

func (r *variantRepository) Delete(ctx context.Context, org domain.Variant) (err error) {
	return r.db.WithContext(ctx).Model(&domain.Variant{}).Delete(&org).Error
}
