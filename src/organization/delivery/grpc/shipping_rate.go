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

func (srv *OrganizationServiceSever) ListShippingRate(ctx context.Context, req *organization.ListShippingRateRequest) (*organization.ListShippingRateResponse, error) {

	f := domain.ShippingRate{
		OrganizationID: req.OrganizationId,
	}

	shipping, count, err := srv.shippingRateRepository.Find(ctx, pagination.Pagination{
		Page: int(req.Page),
		Size: int(req.Size),
	}, f)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &organization.ListShippingRateResponse{
		Data:      shipping.ToProto(),
		TotalData: int32(count),
	}, nil
}

func (srv *OrganizationServiceSever) GetShippingRate(ctx context.Context, req *organization.ShippingRate) (*organization.ShippingRate, error) {

	shipping, err := srv.shippingRateRepository.FindOne(ctx, domain.ShippingRate{
		ID: req.ShippingRateId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if shipping == nil {
		return nil, status.Error(codes.InvalidArgument, "shipping rate not found")
	}

	return shipping.ToProto(), nil
}

func (srv *OrganizationServiceSever) CreateShippingRate(ctx context.Context, req *organization.ShippingRate) (*organization.ShippingRate, error) {

	newShippingRate := domain.ShippingRate{
		ID:             uuid.NewString(),
		OrganizationID: req.OrganizationId,
		Type:           req.Type.String(),
		Name:           req.Name,
		Description:    req.Description,
		Price:          req.Price,
	}

	if _, err := srv.shippingRateRepository.Save(ctx, newShippingRate); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return srv.GetShippingRate(ctx, &organization.ShippingRate{
		ShippingRateId: newShippingRate.ID,
	})
}

func (srv *OrganizationServiceSever) UpdateShippingRate(ctx context.Context, req *organization.ShippingRate) (*organization.ShippingRate, error) {

	return &organization.ShippingRate{}, nil
}
