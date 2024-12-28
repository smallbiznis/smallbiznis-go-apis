package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/item/domain"
	"gorm.io/gorm"
)

type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) domain.IItemRepository {
	return &itemRepository{db}
}

func (r *itemRepository) Find(ctx context.Context, p pagination.Pagination, f domain.Item) (product domain.Items, count int64, err error) {
	stmt := r.db.WithContext(ctx).Model(&domain.Item{}).Preload("Options", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Option")
	}).Preload("Variants").
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

func (r *itemRepository) FindOne(ctx context.Context, f domain.Item) (org *domain.Item, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.Item{}).Preload("Options", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Option")
	}).Preload("Variants").Where(&f).First(&org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return
	}

	return
}

func (r *itemRepository) Save(ctx context.Context, d domain.Item) (org *domain.Item, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.Item{}).Create(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.Item{ID: d.ID})
}

func (r *itemRepository) Update(ctx context.Context, d domain.Item) (org *domain.Item, err error) {
	if err = r.db.WithContext(ctx).Save(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.Item{ID: d.ID})
}

func (r *itemRepository) Delete(ctx context.Context, org domain.Item) (err error) {
	return r.db.WithContext(ctx).Model(&domain.Item{}).Delete(&org).Error
}
