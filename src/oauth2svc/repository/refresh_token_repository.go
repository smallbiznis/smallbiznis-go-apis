package repository

import (
	"context"
	"errors"

	"github.com/smallbiznis/oauth2-server/model"
	"gorm.io/gorm"
)

type IRefreshTokenRepository interface {
	FindOne(ctx context.Context, f model.RefreshToken) (*model.RefreshToken, error)
	Save(ctx context.Context, a model.RefreshToken) (*model.RefreshToken, error)
	Update(ctx context.Context, a model.RefreshToken) (*model.RefreshToken, error)
}

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(
	db *gorm.DB,
) IRefreshTokenRepository {
	return &refreshTokenRepository{db}
}

func (r *refreshTokenRepository) FindOne(ctx context.Context, f model.RefreshToken) (accessToken *model.RefreshToken, err error) {
	if err := r.db.WithContext(ctx).Model(&model.RefreshToken{}).Where(&f).First(&accessToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return
}

func (r *refreshTokenRepository) Save(ctx context.Context, a model.RefreshToken) (accessToken *model.RefreshToken, err error) {
	if err = r.db.WithContext(ctx).Model(&model.RefreshToken{}).Create(&a).Error; err != nil {
		return
	}

	return r.FindOne(ctx, model.RefreshToken{ID: a.ID})
}

func (r *refreshTokenRepository) Update(ctx context.Context, a model.RefreshToken) (accessToken *model.RefreshToken, err error) {
	return
}
