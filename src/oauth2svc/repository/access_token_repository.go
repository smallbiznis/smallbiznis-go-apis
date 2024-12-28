package repository

import (
	"context"
	"errors"

	"github.com/smallbiznis/oauth2-server/model"
	"gorm.io/gorm"
)

type IAccessTokenRepository interface {
	FindOne(ctx context.Context, f model.AccessToken) (*model.AccessToken, error)
	Save(ctx context.Context, a model.AccessToken) (*model.AccessToken, error)
	Update(ctx context.Context, a model.AccessToken) (*model.AccessToken, error)
}

type accessTokenRepository struct {
	db *gorm.DB
}

func NewAccessTokenRepository(
	db *gorm.DB,
) IAccessTokenRepository {
	return &accessTokenRepository{db}
}

func (r *accessTokenRepository) FindOne(ctx context.Context, f model.AccessToken) (accessToken *model.AccessToken, err error) {
	if err := r.db.WithContext(ctx).Model(&model.AccessToken{}).Where(&f).First(&accessToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return
}

func (r *accessTokenRepository) Save(ctx context.Context, a model.AccessToken) (accessToken *model.AccessToken, err error) {
	if err = r.db.WithContext(ctx).Model(&model.AccessToken{}).Create(&a).Error; err != nil {
		return
	}

	return r.FindOne(ctx, model.AccessToken{ID: a.ID})
}

func (r *accessTokenRepository) Update(ctx context.Context, a model.AccessToken) (accessToken *model.AccessToken, err error) {
	return
}
