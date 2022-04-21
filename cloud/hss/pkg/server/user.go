package server

import (
	"context"
	uuid2 "github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukamaX/cloud/hss/pb/gen"
	pbclient "github.com/ukama/ukamaX/cloud/hss/pb/gen/simmgr"
	"github.com/ukama/ukamaX/cloud/hss/pkg"
	"github.com/ukama/ukamaX/cloud/hss/pkg/db"
	"github.com/ukama/ukamaX/cloud/hss/pkg/sims"
	"github.com/ukama/ukamaX/common/grpc"
	"github.com/ukama/ukamaX/common/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
)

const uuidParsingError = "Error parsing UUID"

type UserService struct {
	pb.UnimplementedUserServiceServer
	userRepo       db.UserRepo
	imsiService    pkg.ImsiClientProvider
	simManager     pbclient.SimManagerServiceClient
	simManagerName string
	simRepo        db.SimcardRepo
	simProvider    sims.SimProvider
}

func NewUserService(userRepo db.UserRepo, imsiProvider pkg.ImsiClientProvider, simRepo db.SimcardRepo,
	simProvider sims.SimProvider, simManager pbclient.SimManagerServiceClient, simManagerName string) *UserService {
	return &UserService{userRepo: userRepo,
		imsiService:    imsiProvider,
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
		sim, err = u.simManager.GetSimInfo(ctx, &pbclient.GetSimInfoRequest{
			Iccid: iccid,
		})
		if err != nil {
			return errors.Wrap(err, "failed to get sim info")
		}

		logrus.Infof("Adding new sim")
		err = db.NewSimcardRepo(txDb).Add(&db.Simcard{
			Iccid:      iccid,
			Source:     u.simManagerName,
			UserID:     usr.ID,
			IsPhysical: isPhysicalSim,
		})
		if err != nil {
			return errors.Wrap(err, "failed to add simcard")
		}

		logrus.Infof("Adding new imsi %s", sim.Imsi)
		s, err := u.imsiService.GetClient()
		if err != nil {
			return errors.Wrap(err, "failed to connect to hss")
		}
		_, err = s.Add(ctx, &pb.AddImsiRequest{
			Imsi: &pb.ImsiRecord{
				Imsi:   sim.Imsi,
				UserId: usr.Uuid.String(),
				Apn:    &pb.Apn{Name: "default"},
			},
			Org: org,
		})
		if err != nil {
			return errors.Wrap(err, "failed to add imsi")
		}
		return nil
	})
	// end of transaction

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}
	return user, nil
}

func (u *UserService) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	users, err := u.userRepo.GetByOrg(req.Org)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "imsi")
	}

	resp := &pb.ListResponse{
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

func (u *UserService) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	uuid, err := uuid2.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}
	user, err := u.userRepo.Get(uuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}

	simCard := dbSimcardsToPbSimcards(user.Simcards)

	if simCard != nil {
		u.pullSimCardStatuses(ctx, simCard)
	}
	return &pb.GetResponse{
		User: dbUsersToPbUsers(user),
		Sim:  simCard,
	}, nil
}

func (u *UserService) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	uuid, err := uuid2.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	node, err := u.userRepo.Update(&db.User{
		Uuid:  uuid,
		Name:  req.User.Name,
		Email: req.User.Email,
		Phone: req.User.Phone,
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}

	return &pb.UpdateResponse{
		User: dbUsersToPbUsers(node),
	}, nil
}

func (u *UserService) SetSimStatus(ctx context.Context, req *pb.SetSimStatusRequest) (*pb.SetSimStatusResponse, error) {
	if req.Carrier.Sms != nil && req.Carrier.Sms.Value {
		return nil, status.Errorf(codes.InvalidArgument, "enabling SMS service is not supported")
	}

	if req.Carrier.Voice != nil && req.Carrier.Voice.Value {
		return nil, status.Errorf(codes.InvalidArgument, "enabling VOICE service is not supported")
	}

	if req.Carrier != nil {
		_, err := u.simManager.SetServiceStatus(ctx, &pbclient.SetServiceStatusRequest{
			Iccid: req.Iccid,
			Services: &pbclient.Services{
				Sms:   req.Carrier.Sms,
				Data:  req.Carrier.Data,
				Voice: req.Carrier.Voice,
			},
		})
		if err != nil {
			return nil, err
		}
	}

	if req.Ukama != nil {
		return nil, status.Errorf(codes.Unimplemented, "Configuring status in ukama network is not yet implemented")
	}

	return &pb.SetSimStatusResponse{}, nil
}

func (u *UserService) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	uuid, err := uuid2.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	user, err := u.userRepo.Get(uuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}

	// terminate sim cards asynchronously
	// would be better to trigger an AMQP event
	go u.terminateSimCard(ctx, user.Simcards)

	// delete user
	err = u.userRepo.Delete(uuid, func(uuid uuid2.UUID, tx *gorm.DB) error {
		err = db.NewSimcardRepo(sql.NewDbFromGorm(tx, pkg.IsDebugMode)).DeleteByUser(uuid)
		if err != nil {
			return grpc.SqlErrorToGrpc(err, "sim")
		}
		return nil
	})
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}
	return &pb.DeleteResponse{}, nil
}

func (u *UserService) terminateSimCard(ctx context.Context, simCards []db.Simcard) {
	for _, sim := range simCards {
		logrus.Infof("Terminating sim. Iccid %s", sim.Iccid)

		_, err := u.simManager.TerminateSim(ctx, &pbclient.TerminateSimRequest{
			Iccid: sim.Iccid,
		})

		if err != nil {
			if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
				logrus.Warning("Simcard not found in sim manager")
				continue
			}

			logrus.Errorf("Error terminating simcard %s: %s", sim.Iccid, err)
		}
	}

}

func (u *UserService) DeactivateUser(ctx context.Context, req *pb.DeactivateUserRequest) (*pb.DeactivateUserResponse, error) {
	uuid, err := uuid2.Parse(req.UserId)
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
	_, err = u.userRepo.Update(&db.User{
		Uuid:        uuid,
		Deactivated: true,
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "user")
	}

	// Deactivate sim cards
	u.terminateSimCard(ctx, usr.Simcards)

	// Delete imsi record from HSS
	s, err := u.imsiService.GetClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to hss")
	}

	_, err = s.Delete(ctx, &pb.DeleteImsiRequest{
		IdOneof: &pb.DeleteImsiRequest_UserId{
			UserId: req.UserId,
		},
	})
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "imsi")
	}

	return &pb.DeactivateUserResponse{}, nil
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
		Sms:   getBoolVal(r.Services.Sms),
		Data:  getBoolVal(r.Services.Data),
		Voice: getBoolVal(r.Services.Voice),
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

func getBoolVal(val *wrapperspb.BoolValue) bool {
	if val == nil {
		return false
	}
	return val.Value
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
		Uuid:          user.Uuid.String(),
		Name:          user.Name,
		Phone:         user.Phone,
		Email:         user.Email,
		IsDeactivated: user.Deactivated,
	}
}
