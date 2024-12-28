package repository

import (
	"context"
	"errors"

	"github.com/smallbiznis/oauth2-server/model"
	"gorm.io/gorm"
)

type IOrganizationRepository interface {
	Find(context.Context) ([]model.Organization, error)
	FindOne(context.Context, model.Organization) (*model.Organization, error)
	Save(context.Context) (*model.Organization, error)
	Update(context.Context) (*model.Organization, error)
}

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) IOrganizationRepository {
	return &organizationRepository{
		db,
	}
}

func (r *organizationRepository) Find(ctx context.Context) (orgs []model.Organization, err error) {
	return
}

func (r *organizationRepository) FindOne(ctx context.Context, f model.Organization) (org *model.Organization, err error) {
	if err = r.db.WithContext(ctx).Model(&model.Organization{}).Where(&f).First(&org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return
}

func (r *organizationRepository) Save(ctx context.Context) (org *model.Organization, err error) {
	return
}

func (r *organizationRepository) Update(ctx context.Context) (org *model.Organization, err error) {
	return
}
