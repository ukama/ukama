package server

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/registry/users/pkg/db"

	"github.com/google/uuid"
	"github.com/ukama/ukama/systems/common/grpc"

	pb "github.com/ukama/ukama/systems/registry/users/pb/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

const uuidParsingError = "Error parsing UUID"

type UserService struct {
	pb.UnimplementedUserServiceServer
	userRepo db.UserRepo
}

func NewUserService(userRepo db.UserRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (u *UserService) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	logrus.Infof("Adding user %v", req)

	user := &db.User{
		Email: req.User.Email,
		Name:  req.User.Name,
		Phone: req.User.Phone,
		Uuid:  uuid.New(),
	}

	err := u.userRepo.Add(user)

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}

	return &pb.AddResponse{User: dbUserToPbUser(user)}, nil
}

func (u *UserService) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	uuid, err := uuid.Parse(req.UserUuid)
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
	uuid, err := uuid.Parse(req.UserUuid)
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
	uuid, err := uuid.Parse(req.UserUuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	usr, err := u.userRepo.Get(uuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}

	if usr.Deactivated {
		return nil, status.Errorf(codes.FailedPrecondition, "user already deactivated")
	}

	// set user's status to suspended
	user := &db.User{
		Uuid:        uuid,
		Deactivated: true,
	}

	err = u.userRepo.Update(user, nil)

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}

	return &pb.DeactivateResponse{User: dbUserToPbUser(user)}, nil
}

func (u *UserService) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	userUUID, err := uuid.Parse(req.UserUuid)
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
