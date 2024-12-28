package grpc

import (
	"context"

	"github.com/smallbiznis/go-genproto/smallbiznis/item/v1"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/item/domain"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (svc *ItemService) ListVariant(ctx context.Context, req *item.ListVariantRequest) (*item.ListVariantResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	f := domain.Variant{
		OrganizationID: req.OrganizationId,
	}

	variants, count, err := svc.variantRepository.Find(ctx, pagination.Pagination{
		Page:    int(req.Page),
		Size:    int(req.Size),
		OrderBy: req.OrderBy.String(),
		SortBy:  req.SortBy,
	}, f)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &item.ListVariantResponse{
		TotalData: int32(count),
		Data:      variants.ToProto(),
	}, nil
}

func (svc *ItemService) GetVariant(ctx context.Context, req *item.GetVariantRequest) (*item.Variant, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("GetVariant")

	variant, err := svc.variantRepository.FindOne(ctx, domain.Variant{
		ID: req.VariantId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if variant == nil {
		return nil, status.Error(codes.InvalidArgument, "variant not found")
	}

	return variant.ToProto(), nil
}

func (svc *ItemService) UpdateVariant(ctx context.Context, req *item.Variant) (*item.Variant, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("GetVariant")

	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}
