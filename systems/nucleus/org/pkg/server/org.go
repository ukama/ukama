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
	pb "github.com/ukama/ukama/systems/nucleus/org/pb/gen"
	"github.com/ukama/ukama/systems/nucleus/org/pkg"
	"github.com/ukama/ukama/systems/nucleus/org/pkg/db"
	"github.com/ukama/ukama/systems/nucleus/org/pkg/providers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrgService struct {
	pb.UnimplementedOrgServiceServer
	orgRepo             db.OrgRepo
	userRepo            db.UserRepo
	orchestratorService providers.OrchestratorProvider
	userService         providers.UserClientProvider
	registrySystem      providers.RegistryProvider
	orgName             string
	baseRoutingKey      msgbus.RoutingKeyBuilder
	msgbus              mb.MsgBusServiceClient
	pushgateway         string
	debug               bool
}

func NewOrgServer(orgName string, orgRepo db.OrgRepo, userRepo db.UserRepo, orch providers.OrchestratorProvider, user providers.UserClientProvider, registry providers.RegistryProvider, msgBus mb.MsgBusServiceClient, pushgateway string, debug bool) *OrgService {
	return &OrgService{
		orgRepo:             orgRepo,
		userRepo:            userRepo,
		orchestratorService: orch,
		userService:         user,
		registrySystem:      registry,
		orgName:             orgName,
		baseRoutingKey:      msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		msgbus:              msgBus,
		pushgateway:         pushgateway,
		debug:               debug,
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
		_, err := o.orchestratorService.DeployOrg(providers.DeployOrgRequest{
			OrgId:   org.Id.String(),
			OrgName: org.Name,
			OwnerId: org.Owner.String(),
		})

		if err != nil {
			log.Errorf("Failed to send deploy org request %v sent to orchestrator. Error: %s", org, err.Error())
		} else {
			log.Infof("Deploy org request %v sent to orchestrator", org)
		}
		/* Not required here:  Adding owner as member is done in member service on init */
		return err
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, grpc.SqlErrorToGrpc(err, "owner")
		}

		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	route := o.baseRoutingKey.SetAction("add").SetObject("org").MustBuild()
	err = o.msgbus.PublishRequest(route, req)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	_ = o.pushOrgCountMetric()

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

func (o *OrgService) GetByUser(ctx context.Context, req *pb.GetByOwnerRequest) (*pb.GetByUserResponse, error) {
	log.Infof("Getting all orgs both of membership or owned by %v", req.GetUserUuid())

	userId, err := uuid.FromString(req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of user uuid. Error %s", err.Error())
	}

	user, err := o.userRepo.Get(userId)
	if err != nil {
		if sql.IsNotFoundError(err) {
			return nil, status.Errorf(codes.InvalidArgument,
				"user doesn't exist.")
		} else {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid format of user uuid. Error %s", err.Error())
		}
	}

	ownedOrgs, err := o.orgRepo.GetByOwner(userId)
	if err != nil {
		if !sql.IsNotFoundError(err) {
			return nil, grpc.SqlErrorToGrpc(err, "owned orgs")
		}
	}

	log.Infof("looking for orgs with member %s", userId.String())
	membOrgs, err := o.orgRepo.GetByMember(user.Id)
	if err != nil {
		if !sql.IsNotFoundError(err) {
			return nil, grpc.SqlErrorToGrpc(err, "member orgs")
		}
	}

	log.Infof("found %d owned orgs and %d member orgs", len(ownedOrgs), len(membOrgs))

	resp := &pb.GetByUserResponse{
		User:     req.GetUserUuid(),
		OwnerOf:  dbOrgsToPbOrgs(ownedOrgs),
		MemberOf: dbOrgsToPbOrgs(membOrgs),
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

func (o *OrgService) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
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
	if err != nil {
		if !sql.IsNotFoundError(err) {
			return nil, status.Errorf(codes.FailedPrecondition,
				"failed to check user")
		}
	} else {
		/* TODO: Error is nil why this needs to be done */
		if !sql.IsNotFoundError(err) {
			return nil, status.Errorf(codes.FailedPrecondition,
				"user already exist")
		}
	}

	/* Registering user */
	user := &db.User{Uuid: userUUID}

	err = o.userRepo.Add(user, func(user *db.User, tx *gorm.DB) error {
		/* Add user to members db of org */
		return o.registrySystem.AddMember(org.Name, user.Uuid.String())

	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "member")
	}

	route := o.baseRoutingKey.SetAction("register").SetObject("user").MustBuild()
	err = o.msgbus.PublishRequest(route, req)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	_ = o.pushUserCountMetric()

	return &pb.RegisterUserResponse{}, nil
}

func (o *OrgService) UpdateOrgForUser(ctx context.Context, in *pb.UpdateOrgForUserRequest) (*pb.UpdateOrgForUserResponse, error) {
	// Get the Organization
	orgUUID, err := uuid.FromString(in.GetOrgId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of user uuid. Error %s", err.Error())
	}

	org, err := o.orgRepo.Get(orgUUID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	// Get the User
	userUUID, err := uuid.FromString(in.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of user uuid. Error %s", err.Error())
	}

	user, err := o.userRepo.Get(userUUID)
	if err != nil {
		if !sql.IsNotFoundError(err) {
			return nil, status.Errorf(codes.FailedPrecondition,
				"failed to check user")
		}
	}

	log.Infof("Adding org %s to user %s.", org.Name, userUUID.String())
	// err = o.userRepo.AddOrgToUser(user, org)
	// if err != nil {
	// 	return nil, grpc.SqlErrorToGrpc(err, "user to org")
	// }
	err = o.orgRepo.AddUser(org, user)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user to org")
	}

	route := o.baseRoutingKey.SetActionUpdate().SetObject("user").MustBuild()
	err = o.msgbus.PublishRequest(route, in)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", in, route, err.Error())
	}

	_ = o.pushUserCountMetric()

	return &pb.UpdateOrgForUserResponse{}, nil
}

func (o *OrgService) RemoveOrgForUser(ctx context.Context, in *pb.RemoveOrgForUserRequest) (*pb.RemoveOrgForUserResponse, error) {
	// Get the Organization
	orgUUID, err := uuid.FromString(in.GetOrgId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of user uuid. Error %s", err.Error())
	}

	org, err := o.orgRepo.Get(orgUUID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	// Get the User
	userUUID, err := uuid.FromString(in.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of user uuid. Error %s", err.Error())
	}

	user, err := o.userRepo.Get(userUUID)
	if err != nil {
		if !sql.IsNotFoundError(err) {
			return nil, status.Errorf(codes.FailedPrecondition,
				"failed to check user")
		}
	}

	err = o.userRepo.RemoveOrgFromUser(user, org)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org from user")
	}

	route := o.baseRoutingKey.SetActionUpdate().SetObject("user").MustBuild()
	err = o.msgbus.PublishRequest(route, in)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", in, route, err.Error())
	}

	_ = o.pushUserCountMetric()

	return &pb.RemoveOrgForUserResponse{}, nil
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

	return nil

}
