package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/registry/users/pkg/db"

	"github.com/ukama/ukama/systems/common/grpc"
	metric "github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/registry/users/pb/gen"
	"github.com/ukama/ukama/systems/registry/users/pkg"

	orgpb "github.com/ukama/ukama/systems/registry/org/pb/gen"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	pkgP "github.com/ukama/ukama/systems/registry/users/pkg/providers"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

const uuidParsingError = "Error parsing UUID"

type UserService struct {
	pb.UnimplementedUserServiceServer
	userRepo        db.UserRepo
	orgService      pkgP.OrgClientProvider
	baseRoutingKey  msgbus.RoutingKeyBuilder
	msgbus          mb.MsgBusServiceClient
	pushGatewayHost string
}

func NewUserService(userRepo db.UserRepo, orgService pkgP.OrgClientProvider, msgBus mb.MsgBusServiceClient, pushGatewayHost string) *UserService {
	return &UserService{
		userRepo:        userRepo,
		orgService:      orgService,
		baseRoutingKey:  msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName),
		msgbus:          msgBus,
		pushGatewayHost: pushGatewayHost,
	}
}

func (u *UserService) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	log.Infof("Adding user %v", req)

	authId, err := uuid.FromString(req.User.AuthId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	user := &db.User{
		Email:  req.User.Email,
		Name:   req.User.Name,
		Phone:  req.User.Phone,
		AuthId: authId,
	}

	err = u.userRepo.Add(user, func(user *db.User, tx *gorm.DB) error {
		log.Infof("Adding user %s as member of default org", user.Id)

		user.Id = uuid.NewV4()

		svc, err := u.orgService.GetClient()
		if err != nil {
			return err
		}

		_, err = svc.RegisterUser(ctx, &orgpb.RegisterUserRequest{
			UserUuid: user.Id.String(),
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}

	route := u.baseRoutingKey.SetAction("add").SetObject("user").MustBuild()
	err = u.msgbus.PublishRequest(route, req)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	userCount, inActiveUser, err := u.userRepo.GetUserCount()
	if err != nil {
		log.Errorf("failed to get User count: %s", err.Error())
	}

	err = metric.CollectAndPushSimMetrics(u.pushGatewayHost, pkg.UserMetric, pkg.NumberOfActiveUsers, float64(userCount-inActiveUser), nil, pkg.SystemName+"-"+pkg.ServiceName)
	if err != nil {
		log.Errorf("Error while pushing subscriberCount metric to pushgaway %s", err.Error())
	}

	err = metric.CollectAndPushSimMetrics(u.pushGatewayHost, pkg.UserMetric, pkg.NumberOfInactiveUsers, float64(inActiveUser), nil, pkg.SystemName+"-"+pkg.ServiceName)
	if err != nil {
		log.Errorf("Error while pushing subscriberCount metric to pushgaway %s", err.Error())
	}

	return &pb.AddResponse{User: dbUserToPbUser(user)}, nil
}

func (u *UserService) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	uuid, err := uuid.FromString(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	user, err := u.userRepo.Get(uuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}

	return &pb.GetResponse{User: dbUserToPbUser(user)}, nil
}

func (u *UserService) GetByAuthId(ctx context.Context, req *pb.GetByAuthIdRequest) (*pb.GetResponse, error) {
	authId, err := uuid.FromString(req.AuthId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	user, err := u.userRepo.GetByAuthId(authId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}

	return &pb.GetResponse{User: dbUserToPbUser(user)}, nil
}

func (u *UserService) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	uuid, err := uuid.FromString(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	user := &db.User{
		Id:    uuid,
		Name:  req.User.Name,
		Email: req.User.Email,
		Phone: req.User.Phone,
	}

	err = u.userRepo.Update(user, nil)

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}

	return &pb.UpdateResponse{User: dbUserToPbUser(user)}, nil
}

func (u *UserService) Deactivate(ctx context.Context, req *pb.DeactivateRequest) (*pb.DeactivateResponse, error) {
	userUUID, err := uuid.FromString(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	usr, err := u.userRepo.Get(userUUID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}

	if usr.Deactivated {
		return nil, status.Errorf(codes.FailedPrecondition, "user already deactivated")
	}

	// set user's status to suspended
	user := &db.User{
		Id:          userUUID,
		Deactivated: true,
	}

	err = u.userRepo.Update(user, func(user *db.User, tx *gorm.DB) error {
		log.Infof("Deactivating remote org user %s", userUUID)

		svc, err := u.orgService.GetClient()
		if err != nil {
			return err
		}

		_, err = svc.UpdateUser(ctx, &orgpb.UpdateUserRequest{UserUuid: user.Id.String(),
			Attributes: &orgpb.UserAttributes{IsDeactivated: user.Deactivated},
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}
	route := u.baseRoutingKey.SetAction("deactivate").SetObject("user").MustBuild()
	err = u.msgbus.PublishRequest(route, req)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	u.pushUserCountMetrics()

	return &pb.DeactivateResponse{User: dbUserToPbUser(user)}, nil
}

func (u *UserService) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	userUUID, err := uuid.FromString(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	usr, err := u.userRepo.Get(userUUID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}

	if !usr.Deactivated {
		return nil, status.Errorf(codes.FailedPrecondition, "user must be deactivated first")
	}

	// delete user
	err = u.userRepo.Delete(userUUID, func(userUUID uuid.UUID, tx *gorm.DB) error {
		// Perform any linked transation
		return nil
	})
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}

	route := u.baseRoutingKey.SetAction("delete").SetObject("user").MustBuild()

	err = u.msgbus.PublishRequest(route, req)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	u.pushUserCountMetrics()

	return &pb.DeleteResponse{}, nil
}

func (u *UserService) Whoami(ctx context.Context, req *pb.GetRequest) (*pb.WhoamiResponse, error) {
	uuid, err := uuid.FromString(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	user, err := u.userRepo.Get(uuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}

	svc, err := u.orgService.GetClient()
	if err != nil {
		return nil, err
	}

	userOrgs, err := svc.GetByUser(ctx, &orgpb.GetByOwnerRequest{
		UserUuid: user.Id.String()})
	if err != nil {
		return nil, err
	}

	res := &pb.WhoamiResponse{
		User:     dbUserToPbUser(user),
		OwnerOf:  orgOwnOgrsToUserOwnOrgs(userOrgs.OwnerOf),
		MemberOf: orgMmbOgrsToUserMnbOrgs(userOrgs.MemberOf),
	}

	return res, nil
}

func (u *UserService) pushUserCountMetrics() {
	userCount, inActiveUser, err := u.userRepo.GetUserCount()
	if err != nil {
		log.Errorf("failed to get User count: %s", err.Error())
	}

	err = metric.CollectAndPushSimMetrics(u.pushGatewayHost, pkg.UserMetric, pkg.NumberOfActiveUsers, float64(userCount-inActiveUser), nil, pkg.SystemName+"-"+pkg.ServiceName)
	if err != nil {
		log.Errorf("Error while pushing subscriberCount metric to pushgaway %s", err.Error())
	}

	err = metric.CollectAndPushSimMetrics(u.pushGatewayHost, pkg.UserMetric, pkg.NumberOfInactiveUsers, float64(inActiveUser), nil, pkg.SystemName+"-"+pkg.ServiceName)
	if err != nil {
		log.Errorf("Error while pushing subscriberCount metric to pushgaway %s", err.Error())
	}
}

func (u *UserService) PushMetrics() {
	u.pushUserCountMetrics()
}

func dbUserToPbUser(user *db.User) *pb.User {
	return &pb.User{
		Id:            user.Id.String(),
		Name:          user.Name,
		Phone:         user.Phone,
		Email:         user.Email,
		IsDeactivated: user.Deactivated,
		AuthId:        user.AuthId.String(),
		CreatedAt:     timestamppb.New(user.CreatedAt),
	}
}

func orgOwnOgrsToUserOwnOrgs(orgs []*orgpb.Organization) []*pb.Organization {
	res := []*pb.Organization{}

	for _, o := range orgs {
		org := &pb.Organization{
			Id:            o.Id,
			Name:          o.Name,
			Owner:         o.Owner,
			Certificate:   o.Certificate,
			IsDeactivated: o.IsDeactivated,
			CreatedAt:     o.CreatedAt,
		}

		res = append(res, org)
	}

	return res
}

func orgMmbOgrsToUserMnbOrgs(orgs []*orgpb.OrgUser) []*pb.OrgUser {
	res := []*pb.OrgUser{}

	for _, o := range orgs {
		org := &pb.OrgUser{
			OrgId:         o.OrgId,
			Uuid:          o.Uuid,
			Role:          pb.RoleType(o.Role),
			IsDeactivated: o.IsDeactivated,
			CreatedAt:     o.CreatedAt,
		}

		res = append(res, org)
	}

	return res
}
