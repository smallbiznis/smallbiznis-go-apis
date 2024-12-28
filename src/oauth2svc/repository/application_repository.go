package repository

import (
	"context"
	"errors"

	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/oauth2-server/model"
	"gorm.io/gorm"
)

type IApplicationRepository interface {
	Find(context.Context, pagination.Pagination, []string, model.Application) ([]model.Application, int64, error)
	FindOne(context.Context, []string, model.Application) (*model.Application, error)
	Save(context.Context, []string, model.Application) (*model.Application, error)
	Update(context.Context, []string, model.Application) (*model.Application, error)
}

type applicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(
	db *gorm.DB,
) IApplicationRepository {
	return &applicationRepository{
		db,
	}
}

func (r *applicationRepository) Find(ctx context.Context, p pagination.Pagination, cols []string, f model.Application) (apps []model.Application, count int64, err error) {
	stmt := r.db.Model(&model.Application{}).WithContext(ctx)

	if len(cols) > 0 {
		stmt.Omit(cols...)
	}

	if err = stmt.Where(f).Count(&count).Find(&apps).Error; err != nil {
		return
	}

	return
}

func (r *applicationRepository) FindOne(ctx context.Context, cols []string, f model.Application) (app *model.Application, err error) {
	stmt := r.db.Model(&model.Application{}).WithContext(ctx)

	if len(cols) > 0 {
		stmt.Omit(cols...)
	}

	if err := stmt.Where(f).First(&app).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return
}

func (r *applicationRepository) Save(ctx context.Context, cols []string, a model.Application) (app *model.Application, err error) {
	if err = r.db.Model(&model.Application{}).WithContext(ctx).Create(&a).Error; err != nil {
		return
	}

	return r.FindOne(ctx, cols, model.Application{ID: a.ID})
}

func (r *applicationRepository) Update(ctx context.Context, cols []string, a model.Application) (app *model.Application, err error) {
	return
}
