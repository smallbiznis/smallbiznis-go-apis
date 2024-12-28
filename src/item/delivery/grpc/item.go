package grpc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/smallbiznis/go-genproto/smallbiznis/inventory/v1"
	"github.com/smallbiznis/go-genproto/smallbiznis/item/v1"
	"github.com/smallbiznis/go-genproto/smallbiznis/organization/v1"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/item/domain"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type ItemService struct {
	item.UnimplementedServiceServer
	db                *gorm.DB
	organizationConn  organization.ServiceClient
	inventoryConn     inventory.ServiceClient
	optionRepository  domain.IOptionRepository
	itemRepository    domain.IItemRepository
	variantRepository domain.IVariantRepository
}

func NewItemService(
	db *gorm.DB,
	organizationConn organization.ServiceClient,
	inventoryConn inventory.ServiceClient,
	optionRepository domain.IOptionRepository,
	itemRepository domain.IItemRepository,
	variantRepository domain.IVariantRepository,
) *ItemService {
	return &ItemService{
		db:                db,
		organizationConn:  organizationConn,
		inventoryConn:     inventoryConn,
		optionRepository:  optionRepository,
		itemRepository:    itemRepository,
		variantRepository: variantRepository,
	}
}

func (svc *ItemService) ListItem(ctx context.Context, req *item.ListItemRequest) (*item.ListItemResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("ListProduct")

	filter := domain.Item{
		OrganizationID: req.OrganizationId,
	}

	if req.Status.String() != "" {
		filter.Status = req.Status.String()
	}

	if req.Type.String() != "" {
		filter.Type = req.Type.String()
	}

	products, count, err := svc.itemRepository.Find(ctx, pagination.Pagination{
		Page:    int(req.Page),
		Size:    int(req.Size),
		SortBy:  req.SortBy,
		OrderBy: req.OrderBy.String(),
	}, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	items := make([]*item.Item, 0)
	for _, it := range products {
		result := &item.Item{
			ItemId:         it.ID,
			OrganizationId: it.OrganizationID,
			Type:           item.Type(item.Type_value[it.Type]),
			Title:          it.Title,
			Slug:           it.Slug,
			BodyHtml:       it.BodyHTML,
			Status:         item.Status(item.Status_value[it.Status]),
			Options:        it.Options.ToProto(),
		}

		variants := make([]*item.Variant, 0)
		for _, v := range it.Variants {
			variant := item.Variant{
				VariantId:      v.ID,
				ItemId:         v.ItemID,
				Sku:            v.SKU,
				Barcode:        v.Barcode,
				Title:          v.Title,
				Taxable:        v.Taxable,
				Price:          v.Price,
				CompareAtPrice: v.CompareAtPrice,
				Cost:           v.Cost,
				Profit:         v.Profit,
				Margin:         v.Margin,
				Weight:         v.Weight,
				WeightUnit:     item.WeightUnit(item.WeightUnit_value[v.WeightUnit]),
				Attributes:     v.Attributes,
			}

			for _, id := range v.InventoryItemIds {
				inv, err := svc.inventoryConn.GetInventory(ctx, &inventory.GetInventoryRequest{
					InventoryItemId: id,
				})
				if err != nil {
					return nil, status.Error(codes.Internal, err.Error())
				}

				variant.Inventories = append(variant.Inventories, inv)
			}

			variants = append(variants, &variant)
		}

		result.Variants = append(result.Variants, variants...)
		items = append(items, result)
	}

	return &item.ListItemResponse{
		TotalData: int32(count),
		Data:      items,
	}, nil
}

func (svc *ItemService) GetItem(ctx context.Context, req *item.GetItemRequest) (*item.Item, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("GetProduct")

	filter := domain.Item{}
	if _, err := uuid.Parse(req.ItemId); err != nil {
		filter.Slug = req.ItemId
	} else {
		filter.ID = req.ItemId
	}

	product, err := svc.itemRepository.FindOne(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if product == nil {
		return nil, status.Error(codes.InvalidArgument, "product not found")
	}

	result := &item.Item{
		ItemId:         product.ID,
		OrganizationId: product.OrganizationID,
		Type:           item.Type(item.Type_value[product.Type]),
		Title:          product.Title,
		Slug:           product.Slug,
		BodyHtml:       product.BodyHTML,
		Status:         item.Status(item.Status_value[product.Status]),
		Options:        product.Options.ToProto(),
	}

	variants := make([]*item.Variant, 0)
	for _, v := range product.Variants {
		variant := item.Variant{
			VariantId:      v.ID,
			ItemId:         v.ItemID,
			Sku:            v.SKU,
			Barcode:        v.Barcode,
			Title:          v.Title,
			Taxable:        v.Taxable,
			Price:          v.Price,
			CompareAtPrice: v.CompareAtPrice,
			Cost:           v.Cost,
			Profit:         v.Profit,
			Margin:         v.Margin,
			Weight:         v.Weight,
			WeightUnit:     item.WeightUnit(item.WeightUnit_value[v.WeightUnit]),
			Attributes:     v.Attributes,
		}

		for _, id := range v.InventoryItemIds {
			inv, err := svc.inventoryConn.GetInventory(ctx, &inventory.GetInventoryRequest{
				InventoryItemId: id,
			})
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}

			variant.Inventories = append(variant.Inventories, inv)
		}

		variants = append(variants, &variant)
	}

	result.Variants = append(result.Variants, variants...)

	return result, nil
}

func (svc *ItemService) AddItem(ctx context.Context, req *item.AddItemRequest) (product *item.Item, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("AddProduct")

	organization, err := svc.organizationConn.GetOrg(ctx, &organization.GetOrganizationRequest{
		OrganizationId: req.OrganizationId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if organization == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid organization")
	}

	newSlug := slug.Make(req.Title)
	exist, err := svc.itemRepository.FindOne(ctx, domain.Item{
		OrganizationID: organization.Id,
		Slug:           newSlug,
	})
	if err != nil {
		zap.L().Error("failed query get item", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	if exist != nil {
		zap.L().Error("item not found")
		return nil, status.Error(codes.InvalidArgument, "item already exist!")
	}

	newProduct := domain.Item{
		ID:             uuid.NewString(),
		OrganizationID: organization.Id,
		Type:           req.Type.String(),
		Slug:           newSlug,
		Title:          req.Title,
		BodyHTML:       req.BodyHtml,
		Status:         req.Status.String(),
	}

	if err := svc.db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {

		if len(req.Options) > 0 {
			productOptions := make(domain.ItemOptions, 0)
			for i, o := range req.Options {

				exist, err := svc.createOrUpdateOption(tx, domain.Option{
					ID:             o.OptionId,
					OrganizationID: newProduct.OrganizationID,
					Name:           o.OptionName,
				})
				if err != nil {
					return err
				}

				productOption := domain.ItemOption{
					ItemID:   newProduct.ID,
					OptionID: exist.ID,
					Position: int32(i + 1),
				}

				for _, v := range o.Values {
					optionValue := domain.OptionValue{
						OptionID: exist.ID,
						Value:    v,
					}

					if err := tx.Where(optionValue).First(&optionValue).Error; err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							if err = tx.Create(&optionValue).Error; err != nil {
								return err
							}
						}
					}

					productOption.Values = append(productOption.Values, v)
				}

				productOptions = append(productOptions, productOption)
			}
			newProduct.Options = productOptions
		}

		if len(req.Variants) > 0 {
			for _, variant := range req.Variants {
				newVariant := domain.Variant{
					ID:             uuid.NewString(),
					OrganizationID: organization.Id,
					ItemID:         newProduct.ID,
					Title:          variant.Title,
					Taxable:        variant.Taxable,
					Price:          variant.Price,
					CompareAtPrice: variant.CompareAtPrice,
					Cost:           variant.Cost,
					Barcode:        variant.Barcode,
					Profit:         variant.Profit,
					Margin:         variant.Margin,
					Weight:         variant.Weight,
					WeightUnit:     variant.WeightUnit.String(),
					Attributes:     variant.Attributes,
				}

				for _, inv := range variant.Inventories {
					result, err := svc.inventoryConn.CreateInventory(ctx, &inventory.Inventory{
						OrganizationId: organization.Id,
						LocationId:     inv.LocationId,
						ItemId:         newVariant.ID,
						Quantity:       inv.Quantity,
					})
					if err != nil {
						return err
					}

					newVariant.InventoryItemIds = append(newVariant.InventoryItemIds, result.InventoryItemId)
				}

				newProduct.Variants = append(newProduct.Variants, newVariant)
			}
		}

		if err = tx.Create(&newProduct).Error; err != nil {
			return err
		}

		return
	}); err != nil {
		zap.L().Error("failed add product", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return svc.GetItem(ctx, &item.GetItemRequest{
		ItemId: newProduct.ID,
	})
}

func (svc *ItemService) buildVariant(product domain.Item) (variants []domain.Variant) {

	// Helper function to recursively generate variants
	var createVariants func(int, []string, string)
	createVariants = func(optionIndex int, currentOptions []string, currentTitle string) {
		if optionIndex == len(product.Options) {
			variant := domain.Variant{
				SKU:        "",
				Title:      currentTitle,
				Attributes: currentOptions,
			}

			for i, v := range currentOptions {
				if int(i+1) == len(currentOptions) {
					variant.SKU += fmt.Sprintf("-%s", v)
				} else {
					variant.SKU += strings.ToUpper(fmt.Sprintf("%s-%s", product.ItemCode(), v))
				}
			}

			variants = append(variants, variant)
			return
		}

		option := product.Options[optionIndex]
		for _, value := range option.Values {
			// Create a new map to avoid mutating the original one
			newOptions := make([]string, 0)
			newOptions = append(newOptions, currentOptions...)
			newOptions = append(newOptions, value)

			var newTitle string
			if int(optionIndex+1) == len(product.Options) {
				newTitle = fmt.Sprintf("%s, %s", currentTitle, value)
			} else {
				newTitle = fmt.Sprintf("%s - %s", currentTitle, value)
			}
			createVariants(optionIndex+1, newOptions, newTitle)
		}
	}

	createVariants(0, []string{}, product.Title)

	return variants
}

func (svc *ItemService) createOrUpdateOption(tx *gorm.DB, req domain.Option) (option *domain.Option, err error) {

	f := domain.Option{
		OrganizationID: req.OrganizationID,
	}

	if req.ID != "" {
		f.ID = req.ID
	}

	if req.Name != "" {
		f.Name = req.Name
	}

	exist, err := svc.optionRepository.FindOne(tx.Statement.Context, f)
	if err != nil {
		return
	}

	if exist == nil {
		if err := tx.Create(&req).Error; err != nil {
			return nil, err
		}

		return &req, nil
	}

	return exist, nil
}

func (svc *ItemService) UpdateItem(ctx context.Context, req *item.Item) (*item.Item, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("UpdateProduct")

	organization, err := svc.organizationConn.GetOrg(ctx, &organization.GetOrganizationRequest{
		OrganizationId: req.OrganizationId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if organization == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid organization")
	}

	exist, err := svc.itemRepository.FindOne(ctx, domain.Item{
		ID:             req.ItemId,
		OrganizationID: req.OrganizationId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if exist == nil {
		return nil, status.Error(codes.InvalidArgument, "product not found")
	}

	exist.Type = req.Type.String()
	exist.Title = req.Title
	exist.BodyHTML = req.BodyHtml
	exist.Status = req.Status.String()

	if err := svc.db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {

		if len(req.Options) > 0 {
			productOptions := make(domain.ItemOptions, 0)
			for i, o := range req.Options {

				exist, err := svc.createOrUpdateOption(tx, domain.Option{
					ID:             o.OptionId,
					OrganizationID: exist.OrganizationID,
					Name:           o.OptionName,
				})
				if err != nil {
					return err
				}

				productOption := domain.ItemOption{
					ItemID:   exist.ID,
					OptionID: exist.ID,
					Position: int32(i + 1),
				}

				for _, v := range o.Values {
					optionValue := domain.OptionValue{
						OptionID: exist.ID,
						Value:    v,
					}

					if err := tx.Where(optionValue).First(&optionValue).Error; err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							if err = tx.Create(&optionValue).Error; err != nil {
								return err
							}
						}

						if err = tx.Save(&optionValue).Error; err != nil {
							return err
						}
					}

					productOption.Values = append(productOption.Values, v)
				}

				productOptions = append(productOptions, productOption)
			}
			exist.Options = productOptions
		}

		if len(req.Variants) > 0 {
			for _, variant := range req.Variants {
				existVariant := domain.Variant{
					ID:             variant.VariantId,
					OrganizationID: organization.Id,
					ItemID:         variant.ItemId,
					Title:          variant.Title,
					Taxable:        variant.Taxable,
					Price:          variant.Price,
					CompareAtPrice: variant.CompareAtPrice,
					Cost:           variant.Cost,
					Barcode:        variant.Barcode,
					Profit:         variant.Profit,
					Margin:         variant.Margin,
					Weight:         variant.Weight,
					WeightUnit:     variant.WeightUnit.String(),
					Attributes:     variant.Attributes,
				}

				for _, inv := range variant.Inventories {
					result, err := svc.inventoryConn.UpdateInventory(ctx, &inventory.Inventory{
						OrganizationId: organization.Id,
						LocationId:     inv.LocationId,
						ItemId:         existVariant.ID,
						Quantity:       inv.Quantity,
					})
					if err != nil {
						return err
					}

					existVariant.InventoryItemIds = append(existVariant.InventoryItemIds, result.InventoryItemId)
				}

				exist.Variants = append(exist.Variants, existVariant)
			}
		}

		if err = tx.Save(exist).Error; err != nil {
			return err
		}

		return
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return svc.GetItem(ctx, &item.GetItemRequest{
		ItemId: exist.ID,
	})
}

func (svc *ItemService) DeleteItem(ctx context.Context, req *item.DeleteItemrequest) (*emptypb.Empty, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("DeleteProduct")

	return &emptypb.Empty{}, nil
}
