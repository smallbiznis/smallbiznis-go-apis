package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/transaction/domain"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(
	db *gorm.DB,
) domain.IOrderRepository {
	return &orderRepository{
		db,
	}
}

func (r *orderRepository) Find(ctx context.Context, p pagination.Pagination, f domain.Order) (orders domain.Orders, count int64, err error) {
	stmt := r.db.WithContext(ctx).Model(&domain.Order{}).
		Preload("OrderItems").
		Where(&f).
		Count(&count).
		Scopes(p.Paginate())

	if p.SortBy != "" && p.OrderBy != "" {
		stmt.Order(fmt.Sprintf("%s %s", p.SortBy, p.OrderBy))
	} else {
		stmt.Order("updated_at DESC")
	}

	if err = stmt.Find(&orders).Error; err != nil {
		return
	}

	return
}

func (r *orderRepository) FindOne(ctx context.Context, f domain.Order) (org *domain.Order, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.Order{}).
		Preload("OrderItems").
		Where(&f).First(&org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return
	}

	return
}

func (r *orderRepository) Save(ctx context.Context, d domain.Order) (org *domain.Order, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.Order{}).Create(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.Order{ID: d.ID})
}

func (r *orderRepository) Update(ctx context.Context, d domain.Order) (org *domain.Order, err error) {
	if err = r.db.WithContext(ctx).Save(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.Order{ID: d.ID})
}

func (r *orderRepository) Delete(ctx context.Context, org domain.Order) (err error) {
	return r.db.WithContext(ctx).Model(&domain.Order{}).Delete(&org).Error
}
