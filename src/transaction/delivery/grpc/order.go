package grpc

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"github.com/smallbiznis/go-genproto/smallbiznis/customer/v1"
	"github.com/smallbiznis/go-genproto/smallbiznis/inventory/v1"
	"github.com/smallbiznis/go-genproto/smallbiznis/organization/v1"
	"github.com/smallbiznis/go-genproto/smallbiznis/transaction/v1"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/transaction/domain"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type TransactionService struct {
	transaction.UnimplementedTransactionServiceServer
	db               *gorm.DB
	organizationConn organization.ServiceClient
	customerConn     customer.CustomerServiceClient
	inventoryConn    inventory.ServiceClient
	orderRepository  domain.IOrderRepository
}

func NewTransactionService(
	db *gorm.DB,
	organizationConn organization.ServiceClient,
	customerConn customer.CustomerServiceClient,
	inventoryConn inventory.ServiceClient,
	orderRepository domain.IOrderRepository,
) *TransactionService {
	return &TransactionService{
		db:               db,
		organizationConn: organizationConn,
		customerConn:     customerConn,
		inventoryConn:    inventoryConn,
		orderRepository:  orderRepository,
	}
}

func (svc *TransactionService) ListOrder(ctx context.Context, req *transaction.ListOrderRequest) (*transaction.ListOrderResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("ListOrder")

	f := domain.Order{
		OrganizationID: req.OrganizationId,
	}

	if req.Status != "" {
		f.Status = req.Status
	}

	orders, count, err := svc.orderRepository.Find(ctx, pagination.Pagination{
		Page: int(req.Page),
		Size: int(req.Size),
	}, f)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var newOrders []*transaction.Order
	for _, v := range orders {
		// cust, err := svc.customerConn.GetCustomer(ctx, &customer.GetCustomerRequest{
		// 	CustomerId: v.CustomerID,
		// })
		// if err != nil {
		// 	return nil, err
		// }

		// addr, err := svc.customerConn.GetAddress(ctx, &customer.Address{
		// 	AddressId:  v.ShippingAddressID,
		// 	CustomerId: v.CustomerID,
		// })
		// if err != nil {
		// 	return nil, err
		// }

		newOrders = append(newOrders, &transaction.Order{
			OrderId:        v.ID,
			OrderNo:        v.OrderNo,
			OrganizationId: v.OrganizationID,
			SubTotal:       v.SubTotal,
			TaxAmount:      v.TaxAmount,
			TotalAmount:    v.TotalAmount,
			OrderItems:     v.OrderItems.ToProto(),
			CreatedAt:      v.ToProto().CreatedAt,
			UpdatedAt:      v.ToProto().UpdatedAt,
		})
	}

	return &transaction.ListOrderResponse{
		TotalData: int32(count),
		Data:      newOrders,
	}, nil
}

func (svc *TransactionService) GetOrder(ctx context.Context, req *transaction.GetOrderRequest) (*transaction.Order, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("GetOrder")

	filter := domain.Order{}
	if _, err := uuid.Parse(req.OrderId); err != nil {
		filter.OrderNo = req.OrderId
	} else {
		filter.ID = req.OrderId
	}

	order, err := svc.orderRepository.FindOne(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if order == nil {
		return nil, status.Error(codes.InvalidArgument, "order not found")
	}

	return order.ToProto(), nil
}

func (svc *TransactionService) CreateOrder(ctx context.Context, req *transaction.CreateOrderRequest) (*transaction.Order, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("CreateOrder")

	org, err := svc.organizationConn.GetOrg(ctx, &organization.GetOrganizationRequest{
		OrganizationId: req.OrganizationId,
	})
	if err != nil {
		return nil, err
	}

	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	newOrder := domain.Order{
		ID:             uuid.NewString(),
		OrganizationID: org.Id,
		OrderNo:        node.Generate().String(),
		Status:         transaction.OrderStatus_created.String(),
	}

	if req.OrderNo != "" {
		exist, err := svc.orderRepository.FindOne(ctx, domain.Order{
			OrganizationID: org.Id,
			OrderNo:        req.OrderNo,
		})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		if exist == nil {
			return nil, status.Error(codes.InvalidArgument, "order_no alredy exist")
		}

		newOrder.OrderNo = req.OrderNo
	}

	if req.CustomerId != "" {
		newOrder.CustomerID = &req.CustomerId
	}

	if req.BillingAddressId != "" {
		newOrder.BillingAddressID = &req.BillingAddressId
	}

	if req.ShippingAddressId != "" {
		newOrder.ShippingAddressID = &req.ShippingAddressId
	}

	for _, orderItems := range req.OrderItems {

		// variant, err := svc.itemConn.GetVariant(ctx, &item.GetVariantRequest{
		// 	VariantId: orderItems.ItemId,
		// })
		// if err != nil {
		// 	return nil, err
		// }

		// orderItemPrice := variant.Price * float32(orderItems.Quantity)
		// if variant.Taxable {
		// 	taxAmount := orderItemPrice * 0.11
		// 	newOrder.TaxAmount += taxAmount
		// 	newOrder.GrandTotal += taxAmount
		// }

		// newOrder.SubTotal += orderItemPrice
		// newOrder.GrandTotal += orderItemPrice
		newOrder.OrderItems = append(newOrder.OrderItems, domain.OrderItem{
			Quantity: orderItems.Quantity,
		})
	}

	if req.PaymentProvider != nil {
		// current := time.Now()
		// newOrder.Payment = domain.OrderPayment{
		// 	OrderID:           newOrder.ID,
		// 	PaymentProviderID: req.PaymentProvider.PaymentProviderId,
		// 	Status:            transaction.PaymentStatus_Pending,
		// }
	}

	if _, err := svc.orderRepository.Save(ctx, newOrder); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return svc.GetOrder(ctx, &transaction.GetOrderRequest{
		OrderId: newOrder.ID,
	})
}

func (svc *TransactionService) UpdateOrder(ctx context.Context, req *transaction.Order) (*transaction.Order, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("UpdateOrder")

	return &transaction.Order{}, nil
}
