package grpc

import (
	"context"

	"github.com/smallbiznis/go-genproto/smallbiznis/member/v1"
	"github.com/smallbiznis/go-genproto/smallbiznis/organization/v1"
	"github.com/smallbiznis/go-genproto/smallbiznis/user/v1"
	"github.com/smallbiznis/go-lib/pkg/pagination"
	"github.com/smallbiznis/member/domain"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MemberServiceServer struct {
	member.UnimplementedMemberServiceServer
	organizationConn organization.ServiceClient
	userConn         user.ServiceClient
	memberRepository domain.IMemberRepository
}

func NewOrganizationServiceServer(
	organizationConn organization.ServiceClient,
	userConn user.ServiceClient,
	memberRepository domain.IMemberRepository,
) *MemberServiceServer {
	return &MemberServiceServer{
		organizationConn: organizationConn,
		userConn:         userConn,
		memberRepository: memberRepository,
	}
}

func (svc *MemberServiceServer) ListMember(ctx context.Context, req *member.ListMemberRequest) (*member.ListMemberResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("ListMember")

	var f domain.Member
	if req.UserId != "" {
		f.UserID = req.UserId
	}

	if req.OrganizationId != "" {
		f.OrganizationID = req.OrganizationId
	}

	members, count, err := svc.memberRepository.Find(ctx, pagination.Pagination{
		Page:    int(req.Page),
		Size:    int(req.Size),
		SortBy:  req.SortBy,
		OrderBy: req.OrderBy.String(),
	}, f)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var newMembers []*member.Member
	for _, v := range members {
		u, err := svc.userConn.Get(ctx, &user.GetUserRequest{UserId: v.UserID})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		org, err := svc.organizationConn.GetOrg(ctx, &organization.GetOrganizationRequest{
			OrganizationId: v.OrganizationID,
		})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		var roles []member.Role
		for _, r := range v.Roles {
			roles = append(roles, member.Role(member.Role_value[r]))
		}

		newMembers = append(newMembers, &member.Member{
			MemberId:     v.ID,
			User:         u,
			Organization: org,
			Roles:        roles,
		})
	}

	return &member.ListMemberResponse{
		TotalData: int32(count),
		Data:      newMembers,
	}, nil
}

func (svc *MemberServiceServer) GetMember(ctx context.Context, req *member.GetMemberRequest) (member *member.Member, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("GetMember")

	exist, err := svc.memberRepository.FindOne(ctx, domain.Member{ID: req.MemberId})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if exist == nil {
		return nil, status.Error(codes.InvalidArgument, "member not found")
	}

	u, err := svc.userConn.Get(ctx, &user.GetUserRequest{UserId: exist.UserID})
	if err != nil {
		return
	}

	org, err := svc.organizationConn.GetOrg(ctx, &organization.GetOrganizationRequest{
		OrganizationId: exist.OrganizationID,
	})
	if err != nil {
		return
	}

	member = exist.ToProto()
	member.User = u
	member.Organization = org

	return
}

func (svc *MemberServiceServer) AddMember(ctx context.Context, req *member.AddMemberRequest) (member *member.Member, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("AddMember")

	u, err := svc.userConn.Get(ctx, &user.GetUserRequest{UserId: req.UserId})
	if err != nil {
		return
	}

	org, err := svc.organizationConn.GetOrg(ctx, &organization.GetOrganizationRequest{
		OrganizationId: req.OrganizatinId,
	})
	if err != nil {
		return
	}

	roles := make([]string, 0)
	for _, r := range req.Roles {
		roles = append(roles, r.String())
	}

	m, err := svc.memberRepository.Save(ctx, domain.Member{
		UserID:         u.UserId,
		OrganizationID: org.Id,
		Roles:          roles,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return m.ToProto(), nil
}

func (svc *MemberServiceServer) UpdateMember(ctx context.Context, req *member.UpdateMemberRequest) (org *member.Member, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("UpdateMember")

	return
}
