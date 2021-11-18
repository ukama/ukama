package server

import (
	"context"
	uuid "github.com/satori/go.uuid"
	pb "github.com/ukama/ukamaX/cloud/hss/pb/gen"
	"github.com/ukama/ukamaX/cloud/hss/pkg/db"
	"github.com/ukama/ukamaX/common/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	userRepo db.UserRepo
	imsiRepo db.ImsiRepo
}

func NewUserService(userRepo db.UserRepo, imsiRepo db.ImsiRepo) *UserService {
	return &UserService{userRepo: userRepo, imsiRepo: imsiRepo}
}

func (u *UserService) Add(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	if len(req.User.Imsi) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "IMSI cannot be empty")
	}

	if len(req.Org) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Org cannot be empty")
	}

	user, err := u.userRepo.Add(&db.User{
		Email:     req.User.Email,
		FirstName: req.User.FirstName,
		LastName:  req.User.LastName,
		Phone:     req.User.Phone,
		UUID:      uuid.NewV4(),
	})
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}

	err = u.imsiRepo.Add(req.Org, &db.Imsi{
		Imsi:     req.User.Imsi,
		UserUuid: user.UUID,
	})
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}
	return &pb.AddUserResponse{
		User: dbUsersToPbUsers(user, req.User.Imsi),
	}, nil
}

func (u *UserService) Delete(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	uuid, err := uuid.FromString(req.UserUuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Error parsing UUID")
	}

	err = u.imsiRepo.Delete(req.UserUuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "imsi")
	}

	err = u.userRepo.Delete(uuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "imsi")
	}
	return &pb.DeleteUserResponse{}, nil
}

func (u *UserService) List(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	users, err := u.userRepo.GetByOrg(req.Org)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "imsi")
	}

	resp := &pb.ListUsersResponse{
		Org:   req.Org,
		Users: dbusersToPbUsers(users),
	}

	return resp, nil
}

func dbusersToPbUsers(users []db.User) []*pb.User {
	res := []*pb.User{}
	for _, u := range users {
		// keep imsi empty for now
		res = append(res, dbUsersToPbUsers(&u, ""))
	}
	return res
}

func dbUsersToPbUsers(user *db.User, imsi string) *pb.User {
	return &pb.User{
		Uuid:      user.UUID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Email:     user.Email,
		Imsi:      imsi,
	}
}
