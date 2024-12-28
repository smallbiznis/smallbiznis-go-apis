package service

import (
	"context"
	"fmt"

	"github.com/gosimple/slug"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/oauth2-server/internal/pkg/errors"
	"github.com/smallbiznis/oauth2-server/model"
	"github.com/smallbiznis/oauth2-server/repository"
)

type IApplicationService interface {
	HandleList(context.Context, *model.RequestListApplication) ([]model.Application, int64, error)
	HandleGet(context.Context, *model.Application) (*model.Application, error)
	HandleCreate(context.Context, *model.RequestCreateApplication) (*model.Application, error)
	HandleUpdate(context.Context, *model.Application) (*model.Application, error)
	HandleDelete(context.Context, *model.Application) error
}

type applicationService struct {
	applicationRepository repository.IApplicationRepository
	accountRepository     repository.IAccountRepository
}

func NewApplicationService(
	applicationRepository repository.IApplicationRepository,
	accountRepository repository.IAccountRepository,
) IApplicationService {
	return &applicationService{
		applicationRepository,
		accountRepository,
	}
}

func (s *applicationService) HandleList(ctx context.Context, req *model.RequestListApplication) (result []model.Application, count int64, err error) {
	return s.applicationRepository.Find(ctx, pagination.Pagination{
		PerPage: req.Pagination.Limit,
		Page:    req.Pagination.Offset,
	}, req.Excludes, model.Application{
		OrganizationID: req.OrganizationID,
		ClientID:       req.ClientID,
	})
}

func (s *applicationService) HandleGet(ctx context.Context, req *model.Application) (*model.Application, error) {

	exist, err := s.applicationRepository.FindOne(ctx, []string{}, *req)
	if err != nil {
		return nil, err
	}

	if exist == nil {
		return nil, errors.ErrApplicationNotFound
	}

	return exist, nil
}

func (s *applicationService) HandleCreate(ctx context.Context, req *model.RequestCreateApplication) (*model.Application, error) {

	newApp := model.NewApplication(req.OrganizationID)
	newServiceAccount := model.NewServiceAccount(req.OrganizationID)

	newApp.LogoURI = req.LogoUrl
	newApp.Name = slug.Make(req.DisplayName)
	newApp.DisplayName = req.DisplayName
	newApp.Type = req.Type
	newApp.AndroidAppID = req.AndroidAppID
	newApp.AndroidSHA1 = req.AndroidSHA1

	if _, err := s.accountRepository.Save(ctx, *newServiceAccount); err != nil {
		return nil, err
	}

	if _, err := s.applicationRepository.Save(ctx, []string{}, *newApp); err != nil {
		if err := s.accountRepository.Delete(ctx, *newServiceAccount); err != nil {
			fmt.Printf("failed delete service account: %v\n", err)
		}
		return nil, err
	}

	return s.HandleGet(ctx, newApp)
}

func (s *applicationService) HandleUpdate(ctx context.Context, req *model.Application) (*model.Application, error) {
	return &model.Application{}, nil
}

func (s *applicationService) HandleDelete(ctx context.Context, req *model.Application) error {
	return nil
}
