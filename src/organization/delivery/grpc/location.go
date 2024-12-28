package grpc

import (
	"context"

	"github.com/smallbiznis/go-genproto/smallbiznis/organization/v1"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/organization/domain"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *OrganizationServiceSever) ListLocation(ctx context.Context, f *organization.ListLocationRequest) (resp *organization.ListLocationResponse, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("ListLocation")

	locations, count, err := srv.locationRepo.Find(ctx, pagination.Pagination{
		Page:    int(f.Page),
		Size:    int(f.Size),
		SortBy:  f.SortBy,
		OrderBy: f.OrderBy.String(),
	}, domain.Location{
		OrganizationID: f.OrganizationId,
		ID:             f.LocationId,
		IDs:            f.LocationIds,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &organization.ListLocationResponse{
		TotalData: int32(count),
		Data:      locations.ToProto(),
	}, nil
}

func (srv *OrganizationServiceSever) GetLocation(ctx context.Context, f *organization.Location) (resp *organization.Location, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("GetLocation")

	return resp, nil
}

func (srv *OrganizationServiceSever) CreateLocation(ctx context.Context, f *organization.Location) (resp *organization.Location, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("CreateLocation")

	location, err := srv.locationRepo.Save(ctx, domain.Location{
		OrganizationID: f.OrganizationId,
		Name:           f.Name,
		CountryID:      f.Country.CountryCode,
		IsDefault:      f.IsDefault,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return location.ToProto(), nil
}

func (srv *OrganizationServiceSever) UpdateLocation(ctx context.Context, f *organization.Location) (resp *organization.Location, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("UpdateLocation")

	return resp, nil
}
