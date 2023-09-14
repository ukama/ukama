package server

import (
	"context"
	"fmt"
	"time"

	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/invitation/pkg/db"
	"github.com/ukama/ukama/systems/registry/invitation/pkg/providers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/registry/invitation/pb/gen"
)

type InvitationServer struct {
	pb.UnimplementedInvitationServiceServer
	iRepo                db.InvitationRepo
	nucleusSystem        providers.NucleusClientProvider
	notification         providers.NotificationClient
	invitationExpiryTime time.Time
	authLoginbaseURL     string
	baseRoutingKey       msgbus.RoutingKeyBuilder
	msgbus               mb.MsgBusServiceClient
	orgName              string
}

func NewInvitationServer(iRepo db.InvitationRepo, invitationExpiryTime time.Time, authLoginbaseURL string, notification providers.NotificationClient, nucleusSystem providers.NucleusClientProvider, msgBus mb.MsgBusServiceClient, orgName string) *InvitationServer {

	return &InvitationServer{
		iRepo:                iRepo,
		notification:         notification,
		invitationExpiryTime: invitationExpiryTime,
		authLoginbaseURL:     authLoginbaseURL,
		nucleusSystem:        nucleusSystem,
		msgbus:               msgBus,
		orgName:              orgName,
	}
}

func (i *InvitationServer) Add(ctx context.Context, req *pb.AddInvitationRequest) (*pb.AddInvitationResponse, error) {
	log.Infof("Adding invitation %v", req)
	invitationId := uuid.NewV4()

	if req.GetOrg() == "" || req.GetEmail() == "" || req.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Org, Email, and Name are required")
	}

	link, err := generateInvitationLink(i.authLoginbaseURL, uuid.NewV4().String(),
		i.invitationExpiryTime)
	if err != nil {
		return nil, err
	}

	res, err := i.nucleusSystem.GetOrgByName(req.GetOrg())
	if err != nil {
		return nil, err
	}

	remoteUserResp, err := i.nucleusSystem.GetUserById(res.Org.Owner)
	if err != nil {
		return nil, err
	}

	err = i.notification.SendEmail(providers.SendEmailReq{
		To:           []string{req.GetEmail()},
		TemplateName: "member-invitation",
		Values: map[string]interface{}{
			"INVITATION": invitationId.String(),
			"LINK":       link,
			"OWNER":      remoteUserResp.User.Name,
			"ORG":        res.Org.Name,
			"ROLE":       req.GetRole().String(),
			"NAME":       req.GetName(),
		},
	})

	if err != nil {
		return nil, err
	}
	route := i.baseRoutingKey.SetAction("add").SetObject("invitation").MustBuild()
	err = i.msgbus.PublishRequest(route, req)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}
	userInfo, err := i.nucleusSystem.GetByEmail(req.GetEmail())
	if err != nil {
		return nil, err
	}

	existingInvitation, _ := i.iRepo.GetInvitationByEmail(req.GetEmail())

	if existingInvitation == nil {
		err = i.iRepo.Add(
			&db.Invitation{
				Id:        invitationId,
				Org:       req.GetOrg(),
				Name:      req.GetName(),
				Link:      link,
				Email:     req.GetEmail(),
				Role:      pbRoleTypeToDb(req.GetRole()),
				ExpiresAt: i.invitationExpiryTime,
				Status:    db.Pending,
				UserId:    userInfo.User.Id,
			},
			nil,
		)

		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "invitation")
		}
	}

	return &pb.AddInvitationResponse{
		Invitation: &pb.Invitation{
			Id:       invitationId.String(),
			Org:      req.GetOrg(),
			Link:     link,
			Email:    req.GetEmail(),
			Role:     req.GetRole(),
			Name:     req.GetName(),
			Status:   pb.StatusType_Pending,
			UserId:   userInfo.User.Id,
			ExpireAt: timestamppb.New(i.invitationExpiryTime),
		},
	}, nil
}
func (i *InvitationServer) Delete(ctx context.Context, req *pb.DeleteInvitationRequest) (*pb.DeleteInvitationResponse, error) {
	log.Infof("Deleting invitation %v", req)

	iuuid, err := uuid.FromString(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of invitation uuid. Error %s", err.Error())
	}

	err = i.iRepo.Delete(iuuid, nil)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "invitation")
	}

	return &pb.DeleteInvitationResponse{
		Id: req.GetId(),
	}, nil
}
func (i *InvitationServer) Update(ctx context.Context, req *pb.UpdateInvitationStatusRequest) (*pb.UpdateInvitationStatusResponse, error) {
	log.Infof("Updating invitation %v", req)

	iuuid, err := uuid.FromString(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of invitation uuid. Error %s", err.Error())
	}

	if req.GetStatus() != pb.StatusType_Unknown {
		return nil, status.Errorf(codes.InvalidArgument, "Status is required")
	}

	err = i.iRepo.UpdateStatus(iuuid, req.GetStatus().String())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "invitation")
	}

	return &pb.UpdateInvitationStatusResponse{
		Id:     req.GetId(),
		Status: req.GetStatus(),
	}, nil
}
func (i *InvitationServer) Get(ctx context.Context, req *pb.GetInvitationRequest) (*pb.GetInvitationResponse, error) {
	log.Infof("Getting invitation %v", req)

	iuuid, err := uuid.FromString(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of invitation uuid. Error %s", err.Error())
	}
	invitation, err := i.iRepo.Get(iuuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "invitation")
	}

	return &pb.GetInvitationResponse{
		Invitation: dbInvitationToPbInvitation(invitation),
	}, nil
}

func (u *InvitationServer) GetInvitationByEmail(ctx context.Context, req *pb.GetInvitationByEmailRequest) (*pb.GetInvitationByEmailResponse, error) {
	log.Infof("Getting invitation %v", req)

	if req.GetEmail() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Email is required")
	}

	invitation, err := u.iRepo.GetInvitationByEmail(req.GetEmail())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "invitation")
	}

	return &pb.GetInvitationByEmailResponse{
		Invitation: dbInvitationToPbInvitation(invitation),
	}, nil
}
func (i *InvitationServer) GetByOrg(ctx context.Context, req *pb.GetInvitationByOrgRequest) (*pb.GetInvitationByOrgResponse, error) {
	log.Infof("Getting invitation %v", req)

	if req.GetOrg() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Org is required")
	}

	invitations, err := i.iRepo.GetByOrg(req.GetOrg())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "invitation")
	}

	return &pb.GetInvitationByOrgResponse{
		Invitations: dbInvitationsToPbInvitations(invitations),
	}, nil
}

func pbRoleTypeToDb(role pb.RoleType) db.RoleType {
	switch role {
	case pb.RoleType_Admin:
		return db.Admin
	case pb.RoleType_Vendor:
		return db.Vendor
	case pb.RoleType_Users:
		return db.Users
	case pb.RoleType_Owner:
		return db.Owner
	case pb.RoleType_Employee:
		return db.Employee
	default:
		return db.Users
	}
}

func dbInvitationToPbInvitation(invitation *db.Invitation) *pb.Invitation {
	return &pb.Invitation{
		Id:       invitation.Id.String(),
		Org:      invitation.Org,
		Link:     invitation.Link,
		Email:    invitation.Email,
		Role:     pb.RoleType(invitation.Role),
		Name:     invitation.Name,
		Status:   pb.StatusType(invitation.Status),
		UserId:   invitation.UserId,
		ExpireAt: timestamppb.New(invitation.ExpiresAt),
	}
}

func dbInvitationsToPbInvitations(invitations []*db.Invitation) []*pb.Invitation {
	res := []*pb.Invitation{}

	for _, i := range invitations {
		res = append(res, dbInvitationToPbInvitation(i))
	}

	return res
}

func generateInvitationLink(authLoginbaseURL string, linkID string, expirationTime time.Time) (string, error) {
	link := fmt.Sprintf("%s?linkId=%s", authLoginbaseURL, linkID)

	expiringLink := fmt.Sprintf("%s&expires=%d", link, expirationTime.Unix())

	return expiringLink, nil
}
