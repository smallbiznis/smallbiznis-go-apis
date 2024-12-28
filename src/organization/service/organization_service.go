package service

import (
	"context"

	"github.com/gosimple/slug"
	"github.com/smallbiznis/go-lib/pkg/errors"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/organization/domain"
	"go.opentelemetry.io/otel/trace"
)

type OrganizationService struct {
	organizationRepository domain.IOrganizationRepository
}

func NewOrganizationService(
	organizationRepository domain.IOrganizationRepository,
) *OrganizationService {
	return &OrganizationService{
		organizationRepository,
	}
}

func (svc *OrganizationService) Find(ctx context.Context, p pagination.Pagination, f domain.Organization) (orgs domain.Organizations, count int64, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("Find")

	orgs, count, err = svc.organizationRepository.Find(ctx, p, f)
	if err != nil {
		return
	}

	return
}

func (svc *OrganizationService) Get(ctx context.Context, f domain.Organization) (org *domain.Organization, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("Get")

	return svc.organizationRepository.FindOne(ctx, f)
}

func (svc *OrganizationService) Create(ctx context.Context, d domain.Organization) (org *domain.Organization, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("Create")

	exist, err := svc.organizationRepository.FindOne(ctx, domain.Organization{
		ID: slug.Make(d.Title),
	})
	if err != nil {
		return nil, err
	}

	if exist != nil {
		return nil, errors.BadRequest("ORGANIZATION_EXIST", "organization already exist")
	}

	newOrg, err := svc.organizationRepository.Save(ctx, d)
	if err != nil {
		return
	}

	return newOrg, nil
}

func (svc *OrganizationService) Update(ctx context.Context, id string, d domain.Organization) (org *domain.Organization, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("Update")

	exist, err := svc.organizationRepository.FindOne(ctx, domain.Organization{ID: id})
	if err != nil {
		return nil, err
	}

	if exist == nil {
		return nil, errors.BadRequest("ORGANIZATION_NOT_FOUND", "organization not found")
	}

	return svc.organizationRepository.Update(ctx, d)
}
