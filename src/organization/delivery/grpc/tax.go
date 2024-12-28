package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/smallbiznis/go-genproto/smallbiznis/organization/v1"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/organization/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (svc *OrganizationServiceSever) ListTaxRule(ctx context.Context, req *organization.LisTaxRequest) (*organization.ListTaxResponse, error) {

	f := domain.TaxRule{
		OrganizationID: req.OrganizationId,
	}

	if req.CountryId != "" {
		f.CountryID = req.CountryId
	}

	taxes, count, err := svc.taxRepository.Find(ctx, pagination.Pagination{
		Page: int(req.Page),
		Size: int(req.Size),
	}, f)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &organization.ListTaxResponse{
		Data:      taxes.ToProto(),
		TotalData: int32(count),
	}, nil
}

func (svc *OrganizationServiceSever) GetTaxRule(ctx context.Context, req *organization.TaxRule) (*organization.TaxRule, error) {

	exist, err := svc.taxRepository.FindOne(ctx, domain.TaxRule{
		ID: req.TaxId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if exist == nil {
		return nil, status.Error(codes.InvalidArgument, "tax not found")
	}

	return exist.ToProto(), nil
}

func (svc *OrganizationServiceSever) CreateTaxRule(ctx context.Context, req *organization.TaxRule) (*organization.TaxRule, error) {

	exist, err := svc.taxRepository.FindOne(ctx, domain.TaxRule{
		OrganizationID: req.OrganizationId,
		CountryID:      req.CountryId,
		Type:           req.Type.String(),
		Rate:           req.Rate,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if exist != nil {
		return nil, status.Error(codes.InvalidArgument, "tax already exist")
	}

	newTax := domain.TaxRule{
		ID:             uuid.NewString(),
		OrganizationID: req.OrganizationId,
		CountryID:      req.CountryId,
		Type:           req.Type.String(),
		Rate:           req.Rate,
	}

	if _, err := svc.taxRepository.Save(ctx, newTax); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return svc.GetTaxRule(ctx, &organization.TaxRule{TaxId: newTax.ID})
}

func (svc *OrganizationServiceSever) UpdateTaxRule(ctx context.Context, req *organization.TaxRule) (*organization.TaxRule, error) {

	return svc.GetTaxRule(ctx, &organization.TaxRule{})
}
