package repository

import (
	"context"

	"github.com/smallbiznis/oauth2-server/model"
	"gorm.io/gorm"
)

type IProviderRepository interface {
	Find(context.Context) ([]model.Provider, error)
	FindOne(context.Context) (*model.Provider, error)
	Save(context.Context) (*model.Provider, error)
	Update(context.Context) (*model.Provider, error)
}

type providerRepository struct {
	db *gorm.DB
}

func NewProviderRepository(db *gorm.DB) IProviderRepository {
	return &providerRepository{
		db,
	}
}

func (r *providerRepository) Find(ctx context.Context) (providers []model.Provider, err error) {
	return
}

func (r *providerRepository) FindOne(ctx context.Context) (provider *model.Provider, err error) {
	return
}

func (r *providerRepository) Save(ctx context.Context) (provider *model.Provider, err error) {
	return
}

func (r *providerRepository) Update(ctx context.Context) (provider *model.Provider, err error) {
	return
}
