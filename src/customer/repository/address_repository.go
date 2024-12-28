package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/smallbiznis/customer/domain"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"gorm.io/gorm"
)

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(
	db *gorm.DB,
) domain.IAddressRepository {
	return &addressRepository{
		db,
	}
}

func (r *addressRepository) Find(ctx context.Context, p pagination.Pagination, f domain.Address) (orders domain.Addreses, count int64, err error) {
	stmt := r.db.WithContext(ctx).Model(&domain.Address{}).
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

func (r *addressRepository) FindOne(ctx context.Context, f domain.Address) (org *domain.Address, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.Address{}).
		Where(&f).First(&org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return
	}

	return
}

func (r *addressRepository) Save(ctx context.Context, d domain.Address) (org *domain.Address, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.Address{}).Create(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.Address{ID: d.ID})
}

func (r *addressRepository) Update(ctx context.Context, d domain.Address) (org *domain.Address, err error) {
	if err = r.db.WithContext(ctx).Save(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.Address{ID: d.ID})
}

func (r *addressRepository) Delete(ctx context.Context, org domain.Address) (err error) {
	return r.db.WithContext(ctx).Model(&domain.Address{}).Delete(&org).Error
}
