package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/member/domain"
	"gorm.io/gorm"
)

type memberRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) domain.IMemberRepository {
	return &memberRepository{
		db,
	}
}

func (r *memberRepository) Find(ctx context.Context, p pagination.Pagination, f domain.Member) (orgs domain.Members, count int64, err error) {
	stmt := r.db.WithContext(ctx).Model(&domain.Member{}).
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

func (r *memberRepository) FindOne(ctx context.Context, f domain.Member) (org *domain.Member, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.Member{}).Where(&f).First(&org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return
	}

	return
}

func (r *memberRepository) Save(ctx context.Context, d domain.Member) (org *domain.Member, err error) {
	if err = r.db.WithContext(ctx).Model(&domain.Member{}).Create(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.Member{ID: d.ID})
}

func (r *memberRepository) Update(ctx context.Context, d domain.Member) (org *domain.Member, err error) {
	if err = r.db.WithContext(ctx).Save(&d).Error; err != nil {
		return
	}

	return r.FindOne(ctx, domain.Member{ID: d.ID})
}

func (r *memberRepository) Delete(ctx context.Context, org domain.Member) (err error) {
	return r.db.WithContext(ctx).Model(&domain.Member{}).Delete(&org).Error
}
