package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/smallbiznis/go-genproto/smallbiznis/inventory/v1"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/inventory/domain"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type InventoryService struct {
	inventory.UnimplementedServiceServer
	db                  *gorm.DB
	inventoryRepository domain.IInventoryItemRepository
}

func NewInventoryService(
	db *gorm.DB,
	inventoryRepository domain.IInventoryItemRepository,
) *InventoryService {
	return &InventoryService{
		db:                  db,
		inventoryRepository: inventoryRepository,
	}
}

func (svc *InventoryService) ListInventory(ctx context.Context, req *inventory.ListInventoryRequest) (*inventory.ListInventoryResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("ListInventoryItem")

	filter := domain.InventoryItem{
		OrganizationID: req.OrganizationId,
	}

	inventories, count, err := svc.inventoryRepository.Find(ctx, pagination.Pagination{
		Page:    int(req.Page),
		Size:    int(req.Size),
		OrderBy: req.OrderBy.String(),
		SortBy:  req.SortBy,
	}, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &inventory.ListInventoryResponse{
		TotalData: int32(count),
		Data:      inventories.ToProto(),
	}, nil
}

func (svc *InventoryService) GetInventory(ctx context.Context, req *inventory.GetInventoryRequest) (*inventory.Inventory, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("GetInventoryItem")

	exist, err := svc.inventoryRepository.FindOne(ctx, domain.InventoryItem{
		ID: req.InventoryItemId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if exist == nil {
		return nil, status.Error(codes.InvalidArgument, "inventory not found")
	}

	return exist.ToProto(), nil
}

func (svc *InventoryService) CreateInventory(ctx context.Context, req *inventory.Inventory) (*inventory.Inventory, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("CreateInventoryItem")

	exist, err := svc.inventoryRepository.FindOne(ctx, domain.InventoryItem{
		OrganizationID: req.OrganizationId,
		LocationID:     req.LocationId,
		ItemID:         req.LocationId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if exist != nil {
		return nil, status.Error(codes.InvalidArgument, "inventory already exist")
	}

	newInventory := domain.InventoryItem{
		ID:             uuid.NewString(),
		OrganizationID: req.OrganizationId,
		LocationID:     req.LocationId,
		ItemID:         req.ItemId,
		Quantity:       req.Quantity,
	}

	if _, err := svc.inventoryRepository.Save(ctx, newInventory); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return svc.GetInventory(ctx, &inventory.GetInventoryRequest{
		InventoryItemId: newInventory.ID,
	})
}

func (svc *InventoryService) UpdateInventory(ctx context.Context, req *inventory.Inventory) (*inventory.Inventory, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("UpdateInventoryItem")

	return &inventory.Inventory{}, nil
}

func (svc *InventoryService) ReservedStock(ctx context.Context, req *inventory.ReservedStockRequest) (*emptypb.Empty, error) {

	exist, err := svc.GetInventory(ctx, &inventory.GetInventoryRequest{InventoryItemId: req.InventoryItemId})
	if err != nil {
		return nil, err
	}

	if _, err := svc.inventoryRepository.Update(ctx, domain.InventoryItem{
		ID:               exist.InventoryItemId,
		OrganizationID:   exist.OrganizationId,
		ItemID:           exist.ItemId,
		Quantity:         exist.Quantity,
		ReservedQuantity: exist.ReservedQuantity + req.Body.ReservedStock,
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (svc *InventoryService) ReleaseStock(ctx context.Context, req *inventory.ReleaseStockRequest) (*emptypb.Empty, error) {

	exist, err := svc.GetInventory(ctx, &inventory.GetInventoryRequest{InventoryItemId: req.InventoryItemId})
	if err != nil {
		return nil, err
	}

	if _, err := svc.inventoryRepository.Update(ctx, domain.InventoryItem{
		ID:               exist.InventoryItemId,
		OrganizationID:   exist.OrganizationId,
		ItemID:           exist.ItemId,
		Quantity:         exist.Quantity,
		ReservedQuantity: exist.ReservedQuantity + req.Body.ReservedStock,
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
