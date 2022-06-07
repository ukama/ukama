package server

import (
	"context"
	"fmt"

	"github.com/ukama/ukama/services/cloud/users/pkg"
	"github.com/ukama/ukama/services/cloud/users/pkg/db"
	"github.com/ukama/ukama/services/cloud/users/pkg/sims"
	"github.com/ukama/ukama/services/common/msgbus"

	uuid2 "github.com/google/uuid"
	"github.com/sirupsen/logrus"
	hsspb "github.com/ukama/ukama/services/cloud/hss/pb/gen"
	pb "github.com/ukama/ukama/services/cloud/users/pb/gen"
	pbclient "github.com/ukama/ukama/services/cloud/users/pb/gen/simmgr"
	"github.com/ukama/ukama/services/common/errors"
	"github.com/ukama/ukama/services/common/grpc"
	"github.com/ukama/ukama/services/common/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
)
import qrcode "github.com/skip2/go-qrcode"

const uuidParsingError = "Error parsing UUID"

type UserService struct {
	pb.UnimplementedUserServiceServer
	userRepo       db.UserRepo
	imsiService    pkg.ImsiClientProvider
	simManager     pbclient.SimManagerServiceClient
	simManagerName string
	simRepo        db.SimcardRepo
	simProvider    sims.SimProvider
	queuePub       msgbus.QPub
}

func NewUserService(userRepo db.UserRepo, imsiProvider pkg.ImsiClientProvider, simRepo db.SimcardRepo,
	simProvider sims.SimProvider, simManager pbclient.SimManagerServiceClient, simManagerName string,
	queuePub msgbus.QPub) *UserService {
	return &UserService{userRepo: userRepo,
		imsiService:    imsiProvider,
		simRepo:        simRepo,
		simManager:     simManager,
		simManagerName: simManagerName,
		simProvider:    simProvider,
		queuePub:       queuePub,
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
			Services: []*db.Service{
				&db.Service{
					Sms:   false,
					Data:  true,
					Voice: false,
					Type:  db.ServiceTypeUkama,
				},
				&db.Service{
					Sms:   false,
					Data:  true,
					Voice: false,
					Type:  db.ServiceTypeCarrier,
				},
			},
		})
		if err != nil {
			return errors.Wrap(err, "failed to add simcard")
		}

		logrus.Infof("Adding new imsi %s", sim.Imsi)
		s, err := u.imsiService.GetClient()
		if err != nil {
			return errors.Wrap(err, "failed to connect to hss")
		}
		_, err = s.Add(ctx, &hsspb.AddImsiRequest{
			Imsi: &hsspb.ImsiRecord{
				Imsi:   sim.Imsi,
				UserId: usr.Uuid.String(),
				Apn:    &hsspb.Apn{Name: "default"},
			},
			Org: org,
		})
		if err != nil {
			return errors.Wrap(err, "failed to add imsi")
		}
		return nil
	})
	// end of transaction

	if !isPhysicalSim {
		err = u.sendEmailToUser(ctx, user.Email, user.Name, iccid)
		if err != nil {
			logrus.Errorf("Error while sending email to user: %v", err)
		}
	}

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

	simCard := dbSimcardsToPbSimcards(user.Simcard)

	if simCard != nil && !user.Deactivated {
		u.pullSimCardStatuses(ctx, simCard)
		u.pullUsage(ctx, simCard)
	}

	return &pb.GetResponse{
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

	switch r.Status {
	case pbclient.GetSimStatusResponse_INACTIVE:
		simCard.Carrier.Status = pb.SimStatus_INACTIVE
	case pbclient.GetSimStatusResponse_ACTIVE:
		simCard.Carrier.Status = pb.SimStatus_ACTIVE
	case pbclient.GetSimStatusResponse_TERMINATED:
		simCard.Carrier.Status = pb.SimStatus_TERMINATED

	default:
		logrus.Errorf("Unknown sim status %s", r.Status.String())
		simCard.Carrier.Status = pb.SimStatus_UNKNOWN
	}
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
	if req.Carrier != nil && req.Carrier.Sms != nil && req.Carrier.Sms.Value {
		return nil, status.Errorf(codes.InvalidArgument, "enabling SMS service is not supported")
	}

	if req.Carrier != nil && req.Carrier.Voice != nil && req.Carrier.Voice.Value {
		return nil, status.Errorf(codes.InvalidArgument, "enabling VOICE service is not supported")
	}

	sim, err := u.simRepo.Get(req.GetIccid())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	ukamaS := sim.GetServices(db.ServiceTypeUkama)
	if ukamaS == nil {
		ukamaS = &db.Service{
			Type:  db.ServiceTypeUkama,
			Iccid: req.Iccid,
		}
	}
	u.updateService(req.Ukama, ukamaS)

	carrierS := sim.GetServices(db.ServiceTypeCarrier)
	if carrierS == nil {
		carrierS = &db.Service{
			Type:  db.ServiceTypeCarrier,
			Iccid: req.Iccid,
		}
	}
	u.updateService(req.Carrier, carrierS)

	err = u.simRepo.UpdateServices(ukamaS, carrierS, func() error {
		if req.Carrier != nil {
			r := &pbclient.SetServiceStatusRequest{
				Iccid: req.Iccid,
				Services: &pbclient.Services{
					Sms:   wrapperspb.Bool(ukamaS.Sms && carrierS.Sms),
					Data:  wrapperspb.Bool(ukamaS.Data && carrierS.Data),
					Voice: wrapperspb.Bool(ukamaS.Voice && carrierS.Voice),
				},
			}
			logrus.Infof("Setting carrier sim status to: %v", r)
			_, err := u.simManager.SetServiceStatus(ctx, r)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	return &pb.SetSimStatusResponse{
		Sim: dbSimcardsToPbSimcards(*sim),
	}, nil
}

func (u *UserService) updateService(req *pb.SetSimStatusRequest_SetServices, ukamaS *db.Service) {
	if req == nil {
		return
	}
	if req.GetSms() != nil {
		ukamaS.Sms = req.GetSms().GetValue()
	}

	if req.GetData() != nil {
		ukamaS.Data = req.GetData().GetValue()
	}

	if req.GetVoice() != nil {
		ukamaS.Voice = req.GetVoice().GetValue()
	}
}

func (u *UserService) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	uuid, err := uuid2.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	_, err = u.DeactivateUser(ctx, &pb.DeactivateUserRequest{
		UserId: req.UserId,
	})

	if err != nil {
		return nil, err
	}

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

func (u *UserService) terminateSimCard(ctx context.Context, sim db.Simcard) {

	logrus.Infof("Terminating sim. Iccid %s", sim.Iccid)

	_, err := u.simManager.TerminateSim(ctx, &pbclient.TerminateSimRequest{
		Iccid: sim.Iccid,
	})

	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
			logrus.Warning("Simcard not found in sim manager")
		} else {
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

	// Deactivate sim cards in sim manager
	u.terminateSimCard(ctx, usr.Simcard)

	err = u.deleteImsiFromHss(ctx, req.UserId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "imsi")
	}

	return &pb.DeactivateUserResponse{}, nil
}

func (u *UserService) deleteImsiFromHss(ctx context.Context, userId string) error {
	// Delete imsi record from HSS
	s, err := u.imsiService.GetClient()
	if err != nil {
		return errors.Wrap(err, "failed to connect to hss")
	}

	_, err = s.Delete(ctx, &hsspb.DeleteImsiRequest{
		IdOneof: &hsspb.DeleteImsiRequest_UserId{
			UserId: userId,
		},
	})
	return err
}

func (u *UserService) GetQrCode(ctx context.Context, req *pb.GetQrCodeRequest) (*pb.GetQrCodeResponse, error) {
	resp, err := u.simManager.GetQrCode(ctx, &pbclient.GetQrCodeRequest{
		Iccid: req.Iccid,
	})
	return &pb.GetQrCodeResponse{
		QrCode: resp.QrCode,
	}, err
}

// add usage to simcard. In case of error, simcard is not updated silently
func (u *UserService) pullUsage(ctx context.Context, simCard *pb.Sim) {
	logrus.Infof("Get sim card usage for %s", simCard.Iccid)
	r, err := u.simManager.GetUsage(ctx, &pbclient.GetUsageRequest{
		Iccid: simCard.Iccid,
	})

	if err != nil {
		logrus.Errorf("Error getting sim status. Error: %s", err.Error())
		return
	}

	simCard.Carrier.Usage = &pb.Usage{
		DataAllowanceBytes: r.DataTotalInBytes,
		DataUsedBytes:      r.DataUsageInBytes,
	}
}
func generateQrcode(qrcodeId string,qrcodeName string) {

	qrCodeImageData, qrGenerateError := qrcode.Encode(qrcodeId,qrcode.medium,256)
	if qrGenerateError != nil {
		return errors.Wrap(qrGenerateError, "failed to generate qrcode")
	 }
	 encodedData := base64.StdEncoding.EncodeToString(qrCodeImageData)
	return encodedData 
}


func (u *UserService) sendEmailToUser(ctx context.Context, email string, name string, iccid string) error {
	logrus.Infof("Sending email to %s", email)
	logrus.Infof("Getting qr code")
	resp, err := u.simManager.GetQrCode(ctx, &pbclient.GetQrCodeRequest{
		Iccid: iccid,
	})
	if err != nil {
		return errors.Wrap("failed to get qr code",%v,err)
	}

	logrus.Infof("Publishing queue message")
	err = u.queuePub.PublishToQueue("mailer", &msgbus.MailMessage{
		To:           email,
		TemplateName: "users-qr-code",
		Values: map[string]any{
			"Name": name,
			"Qr":   generateQrcode(resp.QrCode,name),
			"QrCodeLink":resp.QrCode
		},
	})

	if err != nil {
		return errors.Wrap(err, "failed to send email")
	}
	return nil
}

func dbSimcardsToPbSimcards(simcard db.Simcard) (res *pb.Sim) {
	res = &pb.Sim{
		Iccid:      simcard.Iccid,
		IsPhysical: simcard.IsPhysical,
		Carrier: &pb.SimStatus{
			Services: pbServices(simcard.GetServices(db.ServiceTypeCarrier)),
		},
		Ukama: &pb.SimStatus{
			Services: pbServices(simcard.GetServices(db.ServiceTypeUkama)),
		},
	}

	return res
}

func pbServices(services *db.Service) *pb.Services {
	if services == nil {
		return nil
	}

	return &pb.Services{
		Sms:   services.Sms,
		Data:  services.Data,
		Voice: services.Voice,
	}
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
