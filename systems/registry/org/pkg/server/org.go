package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/ukama/ukama/systems/common/grpc"
	pb "github.com/ukama/ukama/systems/registry/org/pb/gen"

	userpb "github.com/ukama/ukama/systems/registry/users/pb/gen"

	"github.com/ukama/ukama/systems/registry/org/pkg/db"
	pkg "github.com/ukama/ukama/systems/registry/org/pkg/providers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrgService struct {
	pb.UnimplementedOrgServiceServer
	orgRepo     db.OrgRepo
	userRepo    db.UserRepo
	userService pkg.UserClientProvider
}

func NewOrgServer(orgRepo db.OrgRepo, userRepo db.UserRepo, userService pkg.UserClientProvider) *OrgService {
	return &OrgService{
		orgRepo:     orgRepo,
		userRepo:    userRepo,
		userService: userService,
	}
}

func (r *OrgService) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	logrus.Infof("Adding org %v", req)

	owner, err := uuid.Parse(req.GetOrg().GetOwner())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of owner id. Error %s", err.Error())
	}

	org := &db.Org{
		Name:        req.GetOrg().GetName(),
		Owner:       owner,
		Certificate: req.GetOrg().GetCertificate(),
	}

	err = r.orgRepo.Add(org)

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	return &pb.AddResponse{Org: dbOrgToPbOrg(org)}, nil
}

func (r *OrgService) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	logrus.Infof("Getting org %v", req)

	org, err := r.orgRepo.Get(uint(req.GetId()))
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	return &pb.GetResponse{Org: dbOrgToPbOrg(org)}, nil
}

func (r *OrgService) GetByName(ctx context.Context, req *pb.GetByNameRequest) (*pb.GetByNameResponse, error) {
	logrus.Infof("Getting org %v", req.GetName())

	org, err := r.orgRepo.GetByName(req.GetName())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	return &pb.GetByNameResponse{Org: dbOrgToPbOrg(org)}, nil
}

func (r *OrgService) GetByOwner(ctx context.Context, req *pb.GetByOwnerRequest) (*pb.GetByOwnerResponse, error) {
	logrus.Infof("Getting all orgs own by %v", req.GetUserUuid())

	owner, err := uuid.Parse(req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of owner uuid. Error %s", err.Error())
	}

	orgs, err := r.orgRepo.GetByOwner(owner)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "orgs")
	}

	resp := &pb.GetByOwnerResponse{
		Owner: req.GetUserUuid(),
		Orgs:  dbOrgsToPbOrgs(orgs),
	}

	return resp, nil
}

func (u *OrgService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	uuid, err := uuid.Parse(req.UserUuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of owner uuid. Error %s", err.Error())
	}

	user, err := u.userRepo.Update(&db.User{
		Uuid:        uuid,
		Deactivated: req.GetAttributes().IsDeactivated,
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}

	return &pb.UpdateUserResponse{User: dbUserToPbUser(user)}, nil
}

func (r *OrgService) AddMember(ctx context.Context, req *pb.MemberRequest) (*pb.MemberResponse, error) {
	// Get the Organization
	org, err := r.orgRepo.Get(uint(req.GetOrgId()))
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	// Get the User
	userUUID, err := uuid.Parse(req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of user uuid. Error %s", err.Error())
	}

	user, err := r.userRepo.Get(userUUID)
	if err != nil {
		logrus.Infof("lookup for remote user %s", userUUID)

		svc, err := r.userService.GetClient()
		if err != nil {
			return nil, err
		}

		remoteUser, err := svc.Get(ctx, &userpb.GetRequest{UserUuid: userUUID.String()})
		if err != nil {
			return nil, err
		}

		// What should we do if the remote user exists, but is already deactivated
		if remoteUser.User.IsDeactivated {
			return nil, status.Errorf(codes.FailedPrecondition, "user is deactivated: cannot be added as member")
		}

		logrus.Infof("Adding remove user %s to local user repo", userUUID)
		user = &db.User{Uuid: userUUID,
			Deactivated: remoteUser.User.IsDeactivated,
		}

		err = r.userRepo.Add(user)
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "user")
		}
	}

	logrus.Infof("Adding member")
	member := &db.OrgUser{
		OrgID:  org.ID,
		UserID: user.ID,
		Uuid:   userUUID,
	}

	err = r.orgRepo.AddMember(member)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	return &pb.MemberResponse{Member: dbMemberToPbMember(member)}, nil
}

func (r *OrgService) GetMember(ctx context.Context, req *pb.MemberRequest) (*pb.MemberResponse, error) {
	uuid, err := uuid.Parse(req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of user uuid. Error %s", err.Error())
	}

	member, err := r.orgRepo.GetMember(uint(req.GetOrgId()), uuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	return &pb.MemberResponse{Member: dbMemberToPbMember(member)}, nil
}

func (r *OrgService) GetMembers(ctx context.Context, req *pb.GetMembersRequest) (*pb.GetMembersResponse, error) {
	_, err := r.orgRepo.Get(uint(req.GetOrgId()))
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	members, err := r.orgRepo.GetMembers(uint(req.GetOrgId()))
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "orgs")
	}

	resp := &pb.GetMembersResponse{
		OrgId:   req.GetOrgId(),
		Members: dbMembersToPbMembers(members),
	}

	return resp, nil
}

func (r *OrgService) UpdateMember(ctx context.Context, req *pb.UpdateMemberRequest) (*pb.MemberResponse, error) {
	uuid, err := uuid.Parse(req.GetMember().GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of user uuid. Error %s", err.Error())
	}

	member := &db.OrgUser{
		OrgID:       uint(req.GetMember().GetOrgId()),
		Uuid:        uuid,
		Deactivated: req.GetAttributes().IsDeactivated,
	}

	err = r.orgRepo.UpdateMember(member.OrgID, member)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	return &pb.MemberResponse{Member: dbMemberToPbMember(member)}, nil
}

func (r *OrgService) RemoveMember(ctx context.Context, req *pb.MemberRequest) (*pb.MemberResponse, error) {
	uuid, err := uuid.Parse(req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of user uuid. Error %s", err.Error())
	}

	member, err := r.orgRepo.GetMember(uint(req.GetOrgId()), uuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	if !member.Deactivated {
		return nil, status.Errorf(codes.FailedPrecondition, "member must be deactivated first")
	}

	err = r.orgRepo.RemoveMember(uint(req.GetOrgId()), uuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	return &pb.MemberResponse{}, nil
}

func dbOrgToPbOrg(org *db.Org) *pb.Organization {
	return &pb.Organization{
		Id:            uint64(org.ID),
		Name:          org.Name,
		Owner:         org.Owner.String(),
		Certificate:   org.Certificate,
		IsDeactivated: org.Deactivated,
		CreatedAt:     timestamppb.New(org.CreatedAt),
	}
}

func dbOrgsToPbOrgs(orgs []db.Org) []*pb.Organization {
	res := []*pb.Organization{}

	for _, o := range orgs {
		res = append(res, dbOrgToPbOrg(&o))
	}

	return res
}

func dbUserToPbUser(user *db.User) *pb.User {
	return &pb.User{
		Uuid:          user.Uuid.String(),
		IsDeactivated: user.Deactivated,
	}
}

func dbMemberToPbMember(member *db.OrgUser) *pb.OrgUser {
	return &pb.OrgUser{
		OrgId:         uint64(member.OrgID),
		Uuid:          member.Uuid.String(),
		IsDeactivated: member.Deactivated,
		CreatedAt:     timestamppb.New(member.CreatedAt),
	}
}

func dbMembersToPbMembers(members []db.OrgUser) []*pb.OrgUser {
	res := []*pb.OrgUser{}

	for _, m := range members {
		res = append(res, dbMemberToPbMember(&m))
	}

	return res
}
