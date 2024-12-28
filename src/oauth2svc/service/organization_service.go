package service

import (
	"context"
	"errors"

	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/oauth2-server/model"
	"github.com/smallbiznis/oauth2-server/repository"
)

type IOrganizationService interface {
	HandleList(context.Context, pagination.Pagination, *model.Organization) ([]model.Organization, int64, error)
	HandleGet(context.Context, *model.Organization) (*model.Organization, error)
	HandleCreate(context.Context, *model.Organization) (*model.Organization, error)
	HandleUpdate(context.Context, *model.Organization) (*model.Organization, error)
	HandleDelete(context.Context, *model.Organization) error
}

type organizationService struct {
	organizationRepository    repository.IOrganizationRepository
	organizationKeyRepository repository.IOrganizationKeyRepository
}

func NewOrganizationService(
	organizationRepository repository.IOrganizationRepository,
	organizationKeyRepository repository.IOrganizationKeyRepository,
) IOrganizationService {
	return &organizationService{
		organizationRepository,
		organizationKeyRepository,
	}
}

func (s *organizationService) HandleList(ctx context.Context, page pagination.Pagination, req *model.Organization) ([]model.Organization, int64, error) {
	return []model.Organization{}, 0, errors.New("unimplement")
}

func (s *organizationService) HandleGet(ctx context.Context, req *model.Organization) (*model.Organization, error) {
	return nil, errors.New("unimplement")
}

func (s *organizationService) HandleCreate(ctx context.Context, req *model.Organization) (*model.Organization, error) {
	return nil, errors.New("unimplement")
}

func (s *organizationService) HandleUpdate(ctx context.Context, req *model.Organization) (*model.Organization, error) {
	return nil, errors.New("unimplement")
}

func (s *organizationService) HandleDelete(ctx context.Context, req *model.Organization) error {
	return errors.New("unimplement")
}
