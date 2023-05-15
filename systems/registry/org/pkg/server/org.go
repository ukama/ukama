package server

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/grpc"
	metric "github.com/ukama/ukama/systems/common/metrics"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/registry/org/pb/gen"
	"github.com/ukama/ukama/systems/registry/org/pkg"
	"github.com/ukama/ukama/systems/registry/org/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrgService struct {
	pb.UnimplementedOrgServiceServer
	orgRepo        db.OrgRepo
	userRepo       db.UserRepo
	orgName        string
	baseRoutingKey msgbus.RoutingKeyBuilder
	msgbus         mb.MsgBusServiceClient
	pushgateway    string
}

func NewOrgServer(orgRepo db.OrgRepo, userRepo db.UserRepo, defaultOrgName string, msgBus mb.MsgBusServiceClient, pushgateway string) *OrgService {
	return &OrgService{
		orgRepo:        orgRepo,
		userRepo:       userRepo,
		orgName:        defaultOrgName,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName),
		msgbus:         msgBus,
		pushgateway:    pushgateway,
	}
}

func (o *OrgService) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	log.Infof("Adding org %v", req)

	owner, err := uuid.FromString(req.GetOrg().GetOwner())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of owner uuid. Error %s", err.Error())
	}

	org := &db.Org{
		Name:        req.GetOrg().GetName(),
		Owner:       owner,
		Certificate: req.GetOrg().GetCertificate(),
	}

	err = o.orgRepo.Add(org, func(org *db.Org, tx *gorm.DB) error {
		org.Id = uuid.NewV4()

		txDb := sql.NewDbFromGorm(tx, pkg.IsDebugMode)

		// Adding owner as a member
		user, err := db.NewUserRepo(txDb).Get(owner)
		if err != nil {
			return err
		}

		log.Infof("Adding owner as member")
		member := &db.OrgUser{
			OrgId:  org.Id,
			UserId: user.Id,
			Uuid:   org.Owner,
		}

		err = db.NewOrgRepo(txDb).AddMember(member)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, grpc.SqlErrorToGrpc(err, "owner")
		}

		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	// route := o.baseRoutingKey.SetAction("add").SetObject("org").MustBuild()
	// err = o.msgbus.PublishRequest(route, req)
	// if err != nil {
	// 	log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	// }

	_ = o.pushOrgCountMetric()
	_ = o.pushOrgMemberCountMetric(org.Id)
	_ = o.pushUserCountMetric()

	return &pb.AddResponse{Org: dbOrgToPbOrg(org)}, nil
}

func (o *OrgService) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Infof("Getting org %v", req)

	orgID, err := uuid.FromString(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of org uuid. Error %s", err.Error())
	}

	org, err := o.orgRepo.Get(orgID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	return &pb.GetResponse{Org: dbOrgToPbOrg(org)}, nil
}

func (o *OrgService) GetByName(ctx context.Context, req *pb.GetByNameRequest) (*pb.GetByNameResponse, error) {
	log.Infof("Getting org %v", req.GetName())

	org, err := o.orgRepo.GetByName(req.GetName())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	return &pb.GetByNameResponse{Org: dbOrgToPbOrg(org)}, nil
}

func (o *OrgService) GetByOwner(ctx context.Context, req *pb.GetByOwnerRequest) (*pb.GetByOwnerResponse, error) {
	log.Infof("Getting all orgs owned by %v", req.GetUserUuid())

	owner, err := uuid.FromString(req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of owner uuid. Error %s", err.Error())
	}

	orgs, err := o.orgRepo.GetByOwner(owner)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "orgs")
	}

	resp := &pb.GetByOwnerResponse{
		Owner: req.GetUserUuid(),
		Orgs:  dbOrgsToPbOrgs(orgs),
	}

	return resp, nil
}

func (o *OrgService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	uuid, err := uuid.FromString(req.UserUuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of owner uuid. Error %s", err.Error())
	}

	user, err := o.userRepo.Update(&db.User{
		Uuid:        uuid,
		Deactivated: req.GetAttributes().IsDeactivated,
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}

	return &pb.UpdateUserResponse{User: dbUserToPbUser(user)}, nil
}

func (o *OrgService) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.MemberResponse, error) {
	// Get the Organization
	org, err := o.orgRepo.GetByName(o.orgName)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	// Get the User
	userUUID, err := uuid.FromString(req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of user uuid. Error %s", err.Error())
	}

	_, err = o.userRepo.Get(userUUID)
	if err == nil {
		return nil, status.Errorf(codes.FailedPrecondition,
			"user is already registered")
	}

	user := &db.User{Uuid: userUUID}
	member := &db.OrgUser{}

	err = o.userRepo.Add(user, func(user *db.User, tx *gorm.DB) error {
		txDb := sql.NewDbFromGorm(tx, pkg.IsDebugMode)

		member := &db.OrgUser{
			OrgId:  org.Id,
			UserId: user.Id,
			Uuid:   userUUID,
		}

		err = db.NewOrgRepo(txDb).AddMember(member)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	// route := o.baseRoutingKey.SetAction("register").SetObject("user").MustBuild()
	// err = o.msgbus.PublishRequest(route, req)
	// if err != nil {
	// 	log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	// }

	_ = o.pushOrgMemberCountMetric(org.Id)
	_ = o.pushUserCountMetric()

	return &pb.MemberResponse{Member: dbMemberToPbMember(member)}, nil
}

func (o *OrgService) AddMember(ctx context.Context, req *pb.MemberRequest) (*pb.MemberResponse, error) {
	// Get the Organization
	org, err := o.orgRepo.GetByName(req.GetOrgName())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	// Get the User
	userUUID, err := uuid.FromString(req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of user uuid. Error %s", err.Error())
	}

	user, err := o.userRepo.Get(userUUID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}

	log.Infof("Adding member")
	member := &db.OrgUser{
		OrgId:  org.Id,
		UserId: user.Id,
		Uuid:   userUUID,
		Role:   pbRoleTypeToDb(req.GetRole()),
	}

	err = o.orgRepo.AddMember(member)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	// route := o.baseRoutingKey.SetAction("add").SetObject("member").MustBuild()
	// err = o.msgbus.PublishRequest(route, req)
	// if err != nil {
	// 	log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	// }

	_ = o.pushOrgMemberCountMetric(org.Id)

	return &pb.MemberResponse{Member: dbMemberToPbMember(member)}, nil
}

func (o *OrgService) GetMember(ctx context.Context, req *pb.MemberRequest) (*pb.MemberResponse, error) {
	uuid, err := uuid.FromString(req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of user uuid. Error %s", err.Error())
	}

	// Get the Organization
	org, err := o.orgRepo.GetByName(req.GetOrgName())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	member, err := o.orgRepo.GetMember(org.Id, uuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	return &pb.MemberResponse{Member: dbMemberToPbMember(member)}, nil
}

func (o *OrgService) GetMembers(ctx context.Context, req *pb.GetMembersRequest) (*pb.GetMembersResponse, error) {
	org, err := o.orgRepo.GetByName(req.GetOrgName())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	members, err := o.orgRepo.GetMembers(org.Id)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "orgs")
	}

	resp := &pb.GetMembersResponse{
		Org:     org.Name,
		Members: dbMembersToPbMembers(members),
	}

	return resp, nil
}

func (o *OrgService) UpdateMember(ctx context.Context, req *pb.UpdateMemberRequest) (*pb.MemberResponse, error) {
	uuid, err := uuid.FromString(req.GetMember().GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of user uuid. Error %s", err.Error())
	}

	org, err := o.orgRepo.GetByName(req.GetMember().GetOrgName())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	member := &db.OrgUser{
		OrgId:       org.Id,
		Uuid:        uuid,
		Deactivated: req.GetAttributes().IsDeactivated,
		Role:        pbRoleTypeToDb(req.Member.Role),
	}

	err = o.orgRepo.UpdateMember(member.OrgId, member)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	_ = o.pushOrgMemberCountMetric(org.Id)

	return &pb.MemberResponse{Member: dbMemberToPbMember(member)}, nil
}

func (o *OrgService) RemoveMember(ctx context.Context, req *pb.MemberRequest) (*pb.MemberResponse, error) {
	uuid, err := uuid.FromString(req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of user uuid. Error %s", err.Error())
	}

	org, err := o.orgRepo.GetByName(req.GetOrgName())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	member, err := o.orgRepo.GetMember(org.Id, uuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	if org.Owner == member.Uuid {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot remove the current owner of the Organization")
	}

	if !member.Deactivated {
		return nil, status.Errorf(codes.FailedPrecondition,
			"member must be deactivated first")
	}

	err = o.orgRepo.RemoveMember(org.Id, uuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	// route := o.baseRoutingKey.SetAction("remove").SetObject("member").MustBuild()
	// err = o.msgbus.PublishRequest(route, req)
	// if err != nil {
	// 	log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	// }

	_ = o.pushOrgMemberCountMetric(org.Id)

	return &pb.MemberResponse{}, nil
}

func (o *OrgService) GetMemberRole(ctx context.Context, req *pb.MemberRoleRequest) (*pb.GetMemberRoleResponse, error) {

	userId, err := uuid.FromString(req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of user uuid. Error %s", err.Error())
	}
	orgId, err := uuid.FromString(req.GetOrgId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of orgId. Error %s", err.Error())
	}

	memberRole, err := o.orgRepo.GetMemberRole(orgId, userId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	return &pb.GetMemberRoleResponse{Role: memberRole}, nil

}
func dbOrgToPbOrg(org *db.Org) *pb.Organization {
	return &pb.Organization{
		Id:            org.Id.String(),
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
		OrgId:         member.OrgId.String(),
		UserId:        uint64(member.UserId),
		Uuid:          member.Uuid.String(),
		Role:          pb.RoleType(member.Role),
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

func (o *OrgService) pushOrgCountMetric() error {
	actOrg, inactOrg, err := o.orgRepo.GetOrgCount()
	if err != nil {
		log.Errorf("failed to get Org count: %s", err.Error())
		return err
	}

	err = metric.CollectAndPushSimMetrics(o.pushgateway, pkg.OrgMetrics, pkg.NumberOfActiveOrgs, float64(actOrg), nil, pkg.SystemName+"-"+pkg.ServiceName)
	if err != nil {
		log.Errorf("Error while pushing active Org metric to pushgateway %s", err.Error())
		return err
	}

	err = metric.CollectAndPushSimMetrics(o.pushgateway, pkg.OrgMetrics, pkg.NumberOfInactiveOrgs, float64(inactOrg), nil, pkg.SystemName+"-"+pkg.ServiceName)
	if err != nil {
		log.Errorf("Error while pushing inactive Org metric to pushgateway %s", err.Error())
		return err
	}
	return nil
}

func (o *OrgService) pushOrgMemberCountMetric(orgId uuid.UUID) error {

	actMemOrg, inactMemOrg, err := o.orgRepo.GetMemberCount(orgId)
	if err != nil {
		log.Errorf("failed to get member count for org %s.Error: %s", orgId.String(), err.Error())
		return err
	}

	labels := make(map[string]string)
	labels["org"] = orgId.String()

	err = metric.CollectAndPushSimMetrics(o.pushgateway, pkg.OrgMetrics, pkg.NumberOfActiveMembersOfOrgs, float64(actMemOrg), labels, pkg.SystemName+"-"+pkg.ServiceName)
	if err != nil {
		log.Errorf("Error while pushing active members of Org metric to pushgateway %s", err.Error())
		return err
	}

	err = metric.CollectAndPushSimMetrics(o.pushgateway, pkg.OrgMetrics, pkg.NumberOfInactiveMembersOfOrgs, float64(inactMemOrg), labels, pkg.SystemName+"-"+pkg.ServiceName)
	if err != nil {
		log.Errorf("Error while pushing inactive members Org metric to pushgateway %s", err.Error())
		return err
	}
	return nil
}

func (o *OrgService) pushUserCountMetric() error {

	actUser, inactUser, err := o.userRepo.GetUserCount()
	if err != nil {
		log.Errorf("failed to get user count.Error: %s", err.Error())
		return err
	}

	err = metric.CollectAndPushSimMetrics(o.pushgateway, pkg.OrgMetrics, pkg.NumberOfActiveUsers, float64(actUser), nil, pkg.SystemName+"-"+pkg.ServiceName)
	if err != nil {
		log.Errorf("Error while pushing active users of Org metric to pushgateway %s", err.Error())
		return err
	}

	err = metric.CollectAndPushSimMetrics(o.pushgateway, pkg.OrgMetrics, pkg.NumberOfInactiveUsers, float64(inactUser), nil, pkg.SystemName+"-"+pkg.ServiceName)
	if err != nil {
		log.Errorf("Error while pushing inactive users Org metric to pushgateway %s", err.Error())
		return err
	}
	return nil
}

func (o *OrgService) PushMetrics() error {

	_ = o.pushOrgCountMetric()

	_ = o.pushUserCountMetric()

	orgs, err := o.orgRepo.GetAll()
	if err != nil {
		log.Errorf("Error while reading orgs. Error %s", err.Error())
		return err
	}

	for _, org := range orgs {
		_ = o.pushOrgMemberCountMetric(org.Id)
	}

	return nil

}

func pbRoleTypeToDb(role pb.RoleType) db.RoleType {
	var dbRole db.RoleType

	switch role {
	case pb.RoleType_ADMIN:
		dbRole = db.Admin
	case pb.RoleType_MEMBER:
		dbRole = db.Member
	default:
		dbRole = db.Undefined
	}

	return dbRole
}
