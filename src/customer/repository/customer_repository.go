package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/smallbiznis/customer/domain"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"gorm.io/gorm"
)

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(
	db *gorm.DB,
) domain.ICustomerRepository {
	return &customerRepository{
		db,
	}
}

func (r *customerRepository) Find(ctx context.Context, p pagination.Pagination, f domain.Customer) (orders domain.Customers, count int64, err error) {
	stmt := r.db.WithContext(ctx).Model(&domain.Customer{}).
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

func (r *customerRepository) FindOne(ctx context.Context, f domain.Customer) (org *domain.Customer, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.Customer{}).
		Where(&f).First(&org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return
	}

	return
}

func (r *customerRepository) Save(ctx context.Context, d domain.Customer) (org *domain.Customer, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.Customer{}).Create(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.Customer{ID: d.ID})
}

func (r *customerRepository) Update(ctx context.Context, d domain.Customer) (org *domain.Customer, err error) {
	if err = r.db.WithContext(ctx).Save(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.Customer{ID: d.ID})
}

func (r *customerRepository) Delete(ctx context.Context, org domain.Customer) (err error) {
	return r.db.WithContext(ctx).Model(&domain.Customer{}).Delete(&org).Error
}
