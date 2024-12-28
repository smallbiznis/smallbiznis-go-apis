package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/inventory/domain"
	"gorm.io/gorm"
)

type inventoryItemRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) domain.IInventoryItemRepository {
	return &inventoryItemRepository{db}
}

func (r *inventoryItemRepository) Find(ctx context.Context, p pagination.Pagination, f domain.InventoryItem) (inventory domain.InventoryItems, count int64, err error) {
	stmt := r.db.WithContext(ctx).Model(&domain.InventoryItem{}).
		Where(&f).
		Count(&count).
		Scopes(p.Paginate())

	if p.SortBy != "" && p.OrderBy != "" {
		stmt.Order(fmt.Sprintf("%s %s", p.SortBy, p.OrderBy))
	} else {
		stmt.Order("updated_at DESC")
	}

	if err = stmt.Find(&inventory).Error; err != nil {
		return
	}

	return
}

func (r *inventoryItemRepository) FindOne(ctx context.Context, f domain.InventoryItem) (org *domain.InventoryItem, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.InventoryItem{}).Where(&f).First(&org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return
	}

	return
}

func (r *inventoryItemRepository) Save(ctx context.Context, d domain.InventoryItem) (org *domain.InventoryItem, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.InventoryItem{}).Create(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.InventoryItem{ID: d.ID})
}

func (r *inventoryItemRepository) Update(ctx context.Context, d domain.InventoryItem) (org *domain.InventoryItem, err error) {
	if err = r.db.WithContext(ctx).Save(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.InventoryItem{ID: d.ID})
}

func (r *inventoryItemRepository) Delete(ctx context.Context, org domain.InventoryItem) (err error) {
	return r.db.WithContext(ctx).Model(&domain.InventoryItem{}).Delete(&org).Error
}
