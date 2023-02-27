package server

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/registry/users/pkg/db"

	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/uuid"

	pb "github.com/ukama/ukama/systems/registry/users/pb/gen"

	orgpb "github.com/ukama/ukama/systems/registry/org/pb/gen"

	pkg "github.com/ukama/ukama/systems/registry/users/pkg/providers"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

const uuidParsingError = "Error parsing UUID"
const ukamaOrgID = 1

type UserService struct {
	pb.UnimplementedUserServiceServer
	userRepo   db.UserRepo
	orgService pkg.OrgClientProvider
}

func NewUserService(userRepo db.UserRepo, orgService pkg.OrgClientProvider) *UserService {
	return &UserService{
		userRepo:   userRepo,
		orgService: orgService,
	}
}

func (u *UserService) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	logrus.Infof("Adding user %v", req)

	user := &db.User{
		Email: req.User.Email,
		Name:  req.User.Name,
		Phone: req.User.Phone,
		Uuid:  uuid.NewV4(),
	}

	err := u.userRepo.Add(user, func(user *db.User, tx *gorm.DB) error {
		logrus.Infof("Adding user %s as member of default org", user.Uuid)

		svc, err := u.orgService.GetClient()
		if err != nil {
			return err
		}

		_, err = svc.RegisterUser(ctx, &orgpb.RegisterUserRequest{
			UserUuid: user.Uuid.String(),
			OrgId:    ukamaOrgID,
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}

	return &pb.AddResponse{User: dbUserToPbUser(user)}, nil
}

func (u *UserService) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	uuid, err := uuid.FromString(req.UserUuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	user, err := u.userRepo.Get(uuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}

	return &pb.GetResponse{User: dbUserToPbUser(user)}, nil
}

func (u *UserService) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	uuid, err := uuid.FromString(req.UserUuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	user := &db.User{
		Uuid:  uuid,
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
	userUUID, err := uuid.FromString(req.UserUuid)
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
		Uuid:        userUUID,
		Deactivated: true,
	}

	err = u.userRepo.Update(user, func(user *db.User, tx *gorm.DB) error {
		logrus.Infof("Deactivating remote org user %s", userUUID)

		svc, err := u.orgService.GetClient()
		if err != nil {
			return err
		}

		_, err = svc.UpdateUser(ctx, &orgpb.UpdateUserRequest{UserUuid: user.Uuid.String(),
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

	return &pb.DeactivateResponse{User: dbUserToPbUser(user)}, nil
}

func (u *UserService) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	userUUID, err := uuid.FromString(req.UserUuid)
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

	return &pb.DeleteResponse{}, nil
}

func dbUserToPbUser(user *db.User) *pb.User {
	return &pb.User{
		Uuid:          user.Uuid.String(),
		Name:          user.Name,
		Phone:         user.Phone,
		Email:         user.Email,
		IsDeactivated: user.Deactivated,
		CreatedAt:     timestamppb.New(user.CreatedAt),
	}
}
