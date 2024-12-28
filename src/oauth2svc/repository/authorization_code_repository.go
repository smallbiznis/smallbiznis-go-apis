package repository

import (
	"context"
	"errors"

	"github.com/smallbiznis/oauth2-server/model"
	"gorm.io/gorm"
)

type IAuthorizationCodeRepository interface {
	FindOne(ctx context.Context, f model.AuthorizationCode) (*model.AuthorizationCode, error)
	Save(ctx context.Context, auth model.AuthorizationCode) (*model.AuthorizationCode, error)
	Update(ctx context.Context, auth model.AuthorizationCode) (*model.AuthorizationCode, error)
}

type authorizationCodeRepository struct {
	db *gorm.DB
}

func NewAuthorizationCodeRepository(db *gorm.DB) IAuthorizationCodeRepository {
	return &authorizationCodeRepository{db}
}

func (r *authorizationCodeRepository) FindOne(ctx context.Context, f model.AuthorizationCode) (code *model.AuthorizationCode, err error) {
	if err := r.db.WithContext(ctx).Model(&model.AuthorizationCode{}).Where(&f).First(&code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return
}

func (r *authorizationCodeRepository) Save(ctx context.Context, auth model.AuthorizationCode) (code *model.AuthorizationCode, err error) {
	if err = r.db.WithContext(ctx).Model(&model.AuthorizationCode{}).Create(&auth).Error; err != nil {
		return
	}

	return r.FindOne(ctx, model.AuthorizationCode{ID: auth.ID})
}

func (r *authorizationCodeRepository) Update(ctx context.Context, auth model.AuthorizationCode) (code *model.AuthorizationCode, err error) {
	if err = r.db.WithContext(ctx).Model(&model.AuthorizationCode{}).Where("id = ?", auth.ID).Save(&auth).Error; err != nil {
		return
	}

	return r.FindOne(ctx, model.AuthorizationCode{ID: auth.ID})
}
