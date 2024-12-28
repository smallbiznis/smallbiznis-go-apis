package service

import (
	"context"
	"strings"

	"github.com/smallbiznis/go-genproto/smallbiznis/user/v1"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/user/domain"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceServer struct {
	user.UnimplementedServiceServer
	userRepository domain.IUserRepository
}

func NewUserServiceServer(
	userRepository domain.IUserRepository,
) *UserServiceServer {
	return &UserServiceServer{
		userRepository: userRepository,
	}
}

func (svc *UserServiceServer) Find(ctx context.Context, p pagination.Pagination, f *user.ListUserRequest) (*user.ListUserResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("Find")

	users, count, err := svc.userRepository.Find(ctx, p, domain.User{
		Email: f.Email,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &user.ListUserResponse{
		TotalData: int32(count),
		Data:      users.ToProto(),
	}, nil
}

func (svc *UserServiceServer) Get(ctx context.Context, f domain.User) (org *user.User, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("Get")

	exist, err := svc.userRepository.FindOne(ctx, domain.User{})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if exist == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return exist.ToProto(), nil
}

func (svc *UserServiceServer) Create(ctx context.Context, d *user.User) (org *domain.User, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("Create")

	exist, err := svc.userRepository.FindOne(ctx, domain.User{
		Email: strings.Trim(d.Email, " "),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if exist != nil {
		return nil, status.Error(codes.AlreadyExists, "user already exist")
	}

	newOrg, err := svc.userRepository.Save(ctx, domain.User{
		ID:        d.UserId,
		AvatarURI: d.AvatarUri,
		Email:     strings.Trim(d.Email, " "),
		FirstName: d.FirstName,
		LastName:  d.LastName,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return newOrg, nil
}

func (svc *UserServiceServer) Update(ctx context.Context, id string, d *user.User) (org *user.User, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("Update")

	exist, err := svc.userRepository.FindOne(ctx, domain.User{ID: id})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if exist == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	exist.AvatarURI = d.AvatarUri
	exist.FirstName = d.FirstName
	exist.LastName = d.LastName
	user, err := svc.userRepository.Update(ctx, *exist)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return user.ToProto(), nil
}
