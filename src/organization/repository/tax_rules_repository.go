package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/organization/domain"
	"gorm.io/gorm"
)

type taxRulesRepository struct {
	db *gorm.DB
}

func NewRulesRepository(
	db *gorm.DB,
) domain.ITaxRulesRepository {
	return &taxRulesRepository{
		db,
	}
}

func (r *taxRulesRepository) Find(ctx context.Context, p pagination.Pagination, f domain.TaxRule) (rules domain.TaxRules, count int64, err error) {
	stmt := r.db.WithContext(ctx).Model(&domain.TaxRule{}).Preload("Country").
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

func (r *taxRulesRepository) FindOne(ctx context.Context, f domain.TaxRule) (org *domain.TaxRule, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.TaxRule{}).Preload("Country").
		Where(&f).First(&org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return
	}

	return
}

func (r *taxRulesRepository) Save(ctx context.Context, d domain.TaxRule) (org *domain.TaxRule, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.TaxRule{}).Create(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.TaxRule{ID: d.ID})
}

func (r *taxRulesRepository) Update(ctx context.Context, d domain.TaxRule) (org *domain.TaxRule, err error) {
	if err = r.db.WithContext(ctx).Save(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.TaxRule{ID: d.ID})
}
