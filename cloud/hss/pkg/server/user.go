package server

import (
	"context"
	uuid2 "github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	pbclient "github.com/ukama/ukamaX/cloud/hss/pb/client/gen"
	pb "github.com/ukama/ukamaX/cloud/hss/pb/gen"
	"github.com/ukama/ukamaX/cloud/hss/pkg"
	"github.com/ukama/ukamaX/cloud/hss/pkg/db"
	"github.com/ukama/ukamaX/cloud/hss/pkg/sims"
	"github.com/ukama/ukamaX/common/grpc"
	"github.com/ukama/ukamaX/common/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	userRepo       db.UserRepo
	imsiRepo       db.ImsiRepo
	simManager     pbclient.SimManagerServiceClient
	simManagerName string
	simRepo        db.SimcardRepo
	simProvider    sims.SimProvider
}

func NewUserService(userRepo db.UserRepo, imsiRepo db.ImsiRepo, simRepo db.SimcardRepo,
	simProvider sims.SimProvider, simManager pbclient.SimManagerServiceClient, simManagerName string) *UserService {
	return &UserService{userRepo: userRepo, imsiRepo: imsiRepo,
		simRepo:        simRepo,
		simManager:     simManager,
		simManagerName: simManagerName,
		simProvider:    simProvider,
	}
}

func (u *UserService) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	iccid := ""
	var err error
	isPhysical := false

	if len(req.SimToken) != 0 {
		if req.SimToken == "I_DO_NOT_NEED_A_SIM" { // for debug purpose
			iccid = sims.GetDubugIccid()
		} else { // physical sim
			iccid, err = u.simProvider.GetICCIDWithCode(req.SimToken)
			if err != nil {
				return nil, grpc.SqlErrorToGrpc(err, "iccid")
			}
			isPhysical = true
		}
	} else { // eSims
		iccid, err = u.simProvider.GetICCIDFromPool()
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "iccid")
		}
	}

	user, err := u.addUserWithIccid(ctx, req.User, iccid, isPhysical, req.Org)
	if err != nil {
		return nil, err
	}

	return &pb.AddResponse{User: dbUsersToPbUsers(user),
		Iccid: iccid}, nil
}

func (u *UserService) AddInternal(ctx context.Context, req *pb.AddInternalRequest) (*pb.AddInternalResponse, error) {
	user, err := u.addUserWithIccid(ctx, req.User, req.Iccid, req.IsPhysicalSim, req.Org)
	if err != nil {
		return nil, err
	}

	return &pb.AddInternalResponse{
		User:  dbUsersToPbUsers(user),
		Iccid: req.Iccid,
	}, nil
}

func (u *UserService) addUserWithIccid(ctx context.Context, reqUser *pb.User, iccid string, isPhysicalSim bool, org string) (*db.User, error) {
	user, err := u.userRepo.Add(&db.User{
		Email: reqUser.Email,
		Name:  reqUser.Name,
		Phone: reqUser.Phone,
		Uuid:  uuid2.New(),
	}, org, func(usr *db.User, tx *gorm.DB) error {

		txDb := sql.NewDbFromGorm(tx, pkg.IsDebugMode)

		isOverLimit, err := db.NewUserRepo(txDb).IsOverTheLimit(org)
		if err != nil {
			logrus.Errorf("Error while checking if user is over limit: %v", err)
			return status.Errorf(codes.Internal, "Internal error")
		}
		if isOverLimit {
			return status.Errorf(codes.PermissionDenied, "limit of sim cards reached")
		}

		// call get sim info to make sure the ICCID exist
		var sim *pbclient.GetSimInfoResponse
		if !sims.IsDebugIdentifier(iccid) {
			sim, err = u.simManager.GetSimInfo(ctx, &pbclient.GetSimInfoRequest{
				Iccid: iccid,
			})
			if err != nil {
				return errors.Wrap(err, "failed to get sim info")
			}
		} else { // for debug purpose
			sim = &pbclient.GetSimInfoResponse{
				Iccid: iccid,
				Imsi:  sims.GetDubugImsi(iccid),
			}
		}

		logrus.Debugf("Adding new sim")
		err = db.NewSimcardRepo(txDb).Add(&db.Simcard{
			Iccid:      iccid,
			Source:     u.simManagerName,
			UserID:     usr.ID,
			IsPhysical: isPhysicalSim,
		})
		if err != nil {
			return errors.Wrap(err, "failed to add simcard")
		}

		logrus.Debugf("Adding new imsi %s", sim.Imsi)
		// we create a new imsi repository to make a call in transaction
		return db.NewImsiRepo(txDb).Add(org, &db.Imsi{
			Imsi:     sim.Imsi,
			UserUuid: usr.Uuid,
		})
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}
	return user, nil
}

func (u *UserService) Delete(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	uuid, err := uuid2.Parse(req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Error parsing UUID")
	}

	err = u.imsiRepo.DeleteByUserId(uuid, func(tx *gorm.DB) error {
		err = db.NewUserRepo(sql.NewDbFromGorm(tx, pkg.IsDebugMode)).Delete(uuid)
		if err != nil {
			return grpc.SqlErrorToGrpc(err, "user")
		}

		err = db.NewSimcardRepo(sql.NewDbFromGorm(tx, pkg.IsDebugMode)).DeleteByUser(uuid)
		if err != nil {
			return grpc.SqlErrorToGrpc(err, "sim")
		}
		return nil
	})

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

func (u *UserService) GenerateSimToken(ctx context.Context, in *pb.GenerateSimTokenRequest) (*pb.GenerateSimTokenResponse, error) {
	iccid := ""
	var err error
	if in.FromPool {
		iccid, err = u.simProvider.GetICCIDFromPool()
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "iccid")
		}

	} else {
		if len(in.Iccid) == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "iccid is required when fromPool is false")
		}
		iccid = in.Iccid
	}

	token, err := u.simProvider.GetSimToken(iccid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.GenerateSimTokenResponse{
		SimToken: token,
	}, nil

}

func (u *UserService) Get(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	uuid, err := uuid2.Parse(req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Error parsing UUID")
	}
	user, err := u.userRepo.Get(uuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}

	simCard := dbSimcardsToPbSimcards(user.Simcards)

	if simCard != nil {
		u.pullSimCardStatuses(ctx, simCard)
	}
	return &pb.GetUserResponse{
		User: dbUsersToPbUsers(user),
		Sim:  simCard,
	}, nil
}

func (u *UserService) pullSimCardStatuses(ctx context.Context, simCard *pb.Sim) {
	logrus.Infof("Get sim card status for %s", simCard.Iccid)
	r, err := u.simManager.GetSimStatus(ctx, &pbclient.GetSimStatusRequest{
		Iccid: simCard.Iccid,
	})

	if err != nil {
		logrus.Errorf("Error getting sim status. Error: %s", err.Error())
		return
	}
	simCard.Carrier = &pb.SimStatus{}

	switch r.Status {
	case pbclient.GetSimStatusResponse_INACTIVE:
		simCard.Carrier.Status = pb.SimStatus_INACTIVE
	case pbclient.GetSimStatusResponse_ACTIVE:
		simCard.Carrier.Status = pb.SimStatus_ACTIVE
	default:
		logrus.Errorf("Unknown sim status %s", r.Status.String())
		simCard.Carrier.Status = pb.SimStatus_UNKNOWN
	}
	simCard.Carrier.Services = &pb.Services{
		Sms:   r.Services.Sms,
		Data:  r.Services.Data,
		Voice: r.Services.Voice,
	}

	simCard.Ukama = &pb.SimStatus{
		// Hardcode unique we implement it
		Status: pb.SimStatus_INACTIVE,
		Services: &pb.Services{
			Sms:   false,
			Data:  false,
			Voice: false,
		},
	}
}

func dbSimcardsToPbSimcards(simcards []db.Simcard) (res *pb.Sim) {
	if len(simcards) > 0 {
		res = &pb.Sim{
			Iccid:      simcards[0].Iccid,
			IsPhysical: simcards[0].IsPhysical,
		}
	}

	if len(simcards) > 1 {
		logrus.Errorf("More then one simcard found for a user")
		for i := 1; i < len(simcards); i++ {
			logrus.Errorf("ICCID %s for user %d not returned to user", simcards[i].Iccid, simcards[i].UserID)
		}
	}

	return res
}

func dbusersToPbUsers(users []db.User) []*pb.User {
	res := []*pb.User{}
	for _, u := range users {
		res = append(res, dbUsersToPbUsers(&u))
	}
	return res
}

func dbUsersToPbUsers(user *db.User) *pb.User {
	return &pb.User{
		Uuid:  user.Uuid.String(),
		Name:  user.Name,
		Phone: user.Phone,
		Email: user.Email,
	}
}
