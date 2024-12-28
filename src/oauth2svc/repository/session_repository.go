package repository

import (
	"context"
	"errors"
	"time"

	"github.com/smallbiznis/oauth2-server/model"
	"gorm.io/gorm"
)

type ISessionRepository interface {
	FindOne(ctx context.Context, f model.UserSession) (user *model.UserSession, err error)
	Find(ctx context.Context, f model.UserSession) (users []model.UserSession, count int64, err error)
	Create(ctx context.Context, session model.UserSession) (user *model.UserSession, err error)
	Delete(ctx context.Context, id string) (err error)
}

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) ISessionRepository {
	return &sessionRepository{db}
}

func (r *sessionRepository) FindOne(ctx context.Context, f model.UserSession) (user *model.UserSession, err error) {
	if err = r.db.WithContext(ctx).Model(&model.UserSession{}).Where(&f).Where("expired_at > ?", time.Now()).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return
	}
	return
}

func (r *sessionRepository) Find(ctx context.Context, f model.UserSession) (users []model.UserSession, count int64, err error) {
	return
}

func (r *sessionRepository) Create(ctx context.Context, session model.UserSession) (user *model.UserSession, err error) {
	if err = r.db.WithContext(ctx).Model(&model.UserSession{}).Create(&session).Error; err != nil {
		return
	}

	return r.FindOne(ctx, model.UserSession{ID: session.ID})
}

func (r *sessionRepository) Delete(ctx context.Context, id string) (err error) {
	return r.db.WithContext(ctx).Model(&model.UserSession{}).Delete("id = ?", id).Error
}
