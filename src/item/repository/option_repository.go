package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/item/domain"
	"gorm.io/gorm"
)

type optionRepository struct {
	db *gorm.DB
}

func NewOptionRepository(db *gorm.DB) domain.IOptionRepository {
	return &optionRepository{db}
}

func (r *optionRepository) Find(ctx context.Context, p pagination.Pagination, f domain.Option) (option domain.Options, count int64, err error) {
	stmt := r.db.WithContext(ctx).Model(&domain.Option{}).
		Where(&f).
		Count(&count).
		Scopes(p.Paginate())

	if p.SortBy != "" && p.OrderBy != "" {
		stmt.Order(fmt.Sprintf("%s %s", p.SortBy, p.OrderBy))
	} else {
		stmt.Order("updated_at DESC")
	}

	if err = stmt.Find(&option).Error; err != nil {
		return
	}

	return
}

func (r *optionRepository) FindOne(ctx context.Context, f domain.Option) (org *domain.Option, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.Option{}).Where(&f).First(&org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return
	}

	return
}

func (r *optionRepository) Save(ctx context.Context, d domain.Option) (org *domain.Option, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.Option{}).Create(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.Option{ID: d.ID})
}

func (r *optionRepository) Update(ctx context.Context, d domain.Option) (org *domain.Option, err error) {
	if err = r.db.WithContext(ctx).Save(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.Option{ID: d.ID})
}

func (r *optionRepository) Delete(ctx context.Context, org domain.Option) (err error) {
	return r.db.WithContext(ctx).Model(&domain.Option{}).Delete(&org).Error
}
