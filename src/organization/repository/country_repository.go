package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/organization/domain"
	"gorm.io/gorm"
)

type countryRepository struct {
	db *gorm.DB
}

func NewCountryRepository(db *gorm.DB) domain.ICountryRepository {
	return &countryRepository{db}
}

func (r *countryRepository) Find(ctx context.Context, p pagination.Pagination, f domain.Countries) (orgs []*domain.Countries, count int64, err error) {
	stmt := r.db.WithContext(ctx).Model(&domain.Countries{}).
		Where(&f).
		Count(&count).
		Scopes(p.Paginate())

	if p.SortBy != "" && p.OrderBy != "" {
		stmt.Order(fmt.Sprintf("%s %s", p.SortBy, p.OrderBy))
	} else {
		stmt.Order("updated_at DESC")
	}

	if err = stmt.Find(&orgs).Error; err != nil {
		return
	}

	return
}

func (r *countryRepository) FindOne(ctx context.Context, f domain.Countries) (org *domain.Countries, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.Countries{}).Where(&f).First(&org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return
	}

	return
}
