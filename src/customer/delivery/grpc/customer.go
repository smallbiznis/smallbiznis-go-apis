package grpc

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/smallbiznis/customer/domain"
	"github.com/smallbiznis/go-genproto/smallbiznis/customer/v1"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type CustomerService struct {
	customer.UnimplementedCustomerServiceServer
	db                 *gorm.DB
	customerRepository domain.ICustomerRepository
	addressRepository  domain.IAddressRepository
}

func NewCustomerService(
	db *gorm.DB,
	customerRepository domain.ICustomerRepository,
	addressRepository domain.IAddressRepository,
) *CustomerService {
	return &CustomerService{
		db:                 db,
		customerRepository: customerRepository,
		addressRepository:  addressRepository,
	}
}

func (svc *CustomerService) ListCustomer(ctx context.Context, req *customer.ListCustomerRequest) (*customer.ListCustomerResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("ListCustomer")

	f := domain.Customer{
		OrganizationID: req.OrganizationId,
	}

	customers, count, err := svc.customerRepository.Find(ctx, pagination.Pagination{
		Page: int(req.Page),
		Size: int(req.Size),
	}, f)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &customer.ListCustomerResponse{
		TotalData: int32(count),
		Data:      customers.ToProto(),
	}, nil
}

func (svc *CustomerService) CreateCustomer(ctx context.Context, req *customer.Customer) (*customer.Customer, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("CreateCustomer")

	newCustomer := domain.Customer{
		ID:             uuid.NewString(),
		OrganizationID: req.OrganizationId,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
	}

	if req.AccountId != "" {
		newCustomer.AccountID = &req.AccountId
	}

	if _, err := svc.customerRepository.Save(ctx, newCustomer); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return svc.GetCustomer(ctx, &customer.GetCustomerRequest{CustomerId: newCustomer.ID})
}

func (svc *CustomerService) GetCustomer(ctx context.Context, req *customer.GetCustomerRequest) (*customer.Customer, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("GetCustomer")

	exist, err := svc.customerRepository.FindOne(ctx, domain.Customer{ID: req.CustomerId})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if exist == nil {
		return nil, status.Error(codes.InvalidArgument, "customer not found")
	}

	return exist.ToProto(), nil
}

func (svc *CustomerService) UpdateCustomer(ctx context.Context, req *customer.Customer) (*customer.Customer, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("UpdateCustomer")

	return &customer.Customer{}, nil
}

func (svc *CustomerService) DeleteCustomer(ctx context.Context, req *customer.DeleteCustomerRequest) (*emptypb.Empty, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("DeleteCustomer")

	return &emptypb.Empty{}, nil
}

func (svc *CustomerService) ListAddress(ctx context.Context, req *customer.ListAddressRequest) (*customer.ListAddressResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("ListAddress")

	f := domain.Address{
		CustomerID: req.CustomerId,
	}

	address, count, err := svc.addressRepository.Find(ctx, pagination.Pagination{
		Page: int(req.Page),
		Size: int(req.Size),
	}, f)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &customer.ListAddressResponse{
		TotalData: int32(count),
		Data:      address.ToProto(),
	}, nil
}

func (svc *CustomerService) CreateAddress(ctx context.Context, req *customer.Address) (*customer.Address, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("CreateAddress")

	newAddress := domain.Address{
		ID:           uuid.NewString(),
		CustomerID:   req.CustomerId,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Address:      req.Address,
		PostalCode:   req.PostalCode,
		IsDefault:    req.IsDefault,
	}

	fmt.Printf("address: %v\n", req)
	if _, err := svc.addressRepository.Save(ctx, newAddress); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return svc.GetAddress(ctx, &customer.Address{AddressId: newAddress.ID})
}

func (svc *CustomerService) GetAddress(ctx context.Context, req *customer.Address) (*customer.Address, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("GetAddress")

	exist, err := svc.addressRepository.FindOne(ctx, domain.Address{
		ID: req.AddressId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if exist == nil {
		return nil, status.Error(codes.InvalidArgument, "address not found")
	}

	return exist.ToProto(), nil
}

func (svc *CustomerService) UpdateAddress(ctx context.Context, req *customer.Address) (*customer.Address, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("UpdateAddress")

	return &customer.Address{}, nil
}

func (svc *CustomerService) DeleteAddress(ctx context.Context, req *customer.Address) (*emptypb.Empty, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("DeleteAddress")

	return &emptypb.Empty{}, nil
}
