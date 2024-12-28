package repository

import (
	"context"
	"errors"

	"github.com/smallbiznis/oauth2-server/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IAccountRepository interface {
	Find(context.Context, model.Pagination, model.Account) (model.AggregateAccounts, int64, error)
	FindOne(context.Context, model.Account) (*model.AggregateAccount, error)
	Save(context.Context, model.Account) (*model.AggregateAccount, error)
	Update(context.Context, model.Account) (*model.AggregateAccount, error)
	Delete(context.Context, model.Account) error
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) IAccountRepository {
	return &accountRepository{db}
}

func (r *accountRepository) Find(ctx context.Context, p model.Pagination, f model.Account) (accounts model.AggregateAccounts, count int64, err error) {
	stmt := r.db.Model(&model.Account{}).WithContext(ctx).Where(f).Select("accounts.*, p.profile_uri, p.first_name, p.last_name").Joins("JOIN profiles p on p.account_id = account_id")

	if err = stmt.Count(&count).Error; err != nil {
		return
	}

	if err := stmt.Limit(p.Limit).Offset(p.Offset).Order(clause.OrderByColumn{Column: clause.Column{Name: p.SortBy}, Desc: p.OrderBy.Bool()}).Find(&accounts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, count, nil
		}
		return nil, count, err
	}

	return
}

func (r *accountRepository) FindOne(ctx context.Context, f model.Account) (account *model.AggregateAccount, err error) {
	stmt := r.db.Model(&model.Account{}).WithContext(ctx).Where(f).Select("accounts.*, p.profile_uri, p.first_name, p.last_name").Joins("JOIN profiles p on p.account_id = account_id")

	if err := stmt.First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return
}

func (r *accountRepository) Save(ctx context.Context, u model.Account) (account *model.AggregateAccount, err error) {
	if err = r.db.Model(&model.Account{}).Create(&u).Error; err != nil {
		return
	}

	return r.FindOne(ctx, model.Account{ID: u.ID})
}

func (r *accountRepository) Update(ctx context.Context, u model.Account) (account *model.AggregateAccount, err error) {
	if err = r.db.Model(&model.Account{}).Save(&u).Error; err != nil {
		return
	}

	return r.FindOne(ctx, model.Account{ID: u.ID})
}

func (r *accountRepository) Delete(ctx context.Context, account model.Account) (err error) {
	return r.db.WithContext(ctx).Model(&model.Account{}).Delete(&account).Error
}
