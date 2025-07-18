/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/roles"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/invitation/pkg"
	"github.com/ukama/ukama/systems/registry/invitation/pkg/db"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	cnotif "github.com/ukama/ukama/systems/common/rest/client/notification"
	cnucl "github.com/ukama/ukama/systems/common/rest/client/nucleus"
	pb "github.com/ukama/ukama/systems/registry/invitation/pb/gen"
)

type InvitationServer struct {
	pb.UnimplementedInvitationServiceServer
	iRepo                db.InvitationRepo
	orgClient            cnucl.OrgClient
	userClient           cnucl.UserClient
	mailerClient         cnotif.MailerClient
	baseRoutingKey       msgbus.RoutingKeyBuilder
	msgbus               mb.MsgBusServiceClient
	invitationExpiryTime uint
	authLoginbaseURL     string
	orgName              string
	TemplateName         string
}

func NewInvitationServer(iRepo db.InvitationRepo, invitationExpiryTime uint, authLoginbaseURL string, mailerClient cnotif.MailerClient,
	orgClient cnucl.OrgClient, userClient cnucl.UserClient, msgBus mb.MsgBusServiceClient, orgName string, TemplateName string) *InvitationServer {

	return &InvitationServer{
		iRepo:                iRepo,
		mailerClient:         mailerClient,
		invitationExpiryTime: invitationExpiryTime,
		authLoginbaseURL:     authLoginbaseURL,
		orgClient:            orgClient,
		userClient:           userClient,
		msgbus:               msgBus,
		baseRoutingKey:       msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		orgName:              orgName,
		TemplateName:         TemplateName,
	}
}

func (i *InvitationServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	log.Infof("Adding invitation %v", req)
	invitationId := uuid.NewV4()

	if i.orgName == "" || req.GetEmail() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "OrgName, Email, and Name are required")
	}

	expiry := time.Now().Add(time.Hour * time.Duration(i.invitationExpiryTime))

	link, err := generateInvitationLink(i.authLoginbaseURL, uuid.NewV4().String(), expiry)
	if err != nil {
		return nil, err
	}

	orgInfo, err := i.orgClient.Get(i.orgName)
	if err != nil {
		return nil, err
	}

	orgOwnerInfo, err := i.userClient.GetById(orgInfo.Owner)
	if err != nil {
		return nil, err
	}

	err = i.mailerClient.SendEmail(cnotif.SendEmailReq{
		To:           []string{req.GetEmail()},
		TemplateName: i.TemplateName,
		Values: map[string]interface{}{
			"INVITATION": invitationId.String(),
			"LINK":       link,
			"OWNER":      orgOwnerInfo.Name,
			"ORG":        orgInfo.Name,
			"ROLE":       req.GetRole().String(),
		},
	})

	if err != nil {
		return nil, err
	}

	invitedUserInfo, err := i.userClient.GetByEmail(req.GetEmail())
	if err != nil {
		log.Errorf("Failed to get invited user info. Error %s", err.Error())
	}

	userId := "00000000-0000-0000-0000-000000000000"
	if invitedUserInfo != nil && !reflect.DeepEqual(invitedUserInfo, reflect.Zero(reflect.TypeOf(invitedUserInfo)).Interface()) {
		userId = invitedUserInfo.Id
	}

	invite := &db.Invitation{
		Id:        invitationId,
		Name:      req.GetName(),
		Link:      link,
		Email:     strings.ToLower(req.GetEmail()),
		Role:      roles.RoleType(req.Role),
		ExpiresAt: expiry,
		Status:    db.Pending,
		UserId:    userId,
	}

	err = i.iRepo.Add(invite, func(*db.Invitation, *gorm.DB) error {
		invite.Id = uuid.NewV4()
		return nil
	})
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "invitation")
	}

	if i.msgbus != nil {
		route := i.baseRoutingKey.SetActionCreate().SetObject("invitation").MustBuild()
		evt := &epb.EventInvitationCreated{
			Id:        invite.Id.String(),
			Link:      invite.Link,
			Email:     strings.ToLower(req.GetEmail()),
			Name:      invite.Name,
			Role:      upb.RoleType(invite.Role),
			Status:    upb.InvitationStatus(invite.Status),
			UserId:    invite.UserId,
			ExpiresAt: invite.ExpiresAt.String(),
		}
		err = i.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
		}
	}

	return &pb.AddResponse{
		Invitation: dbInvitationToPbInvitation(invite),
	}, nil
}

func (i *InvitationServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	log.Infof("Deleting invitation %v", req)

	iuuid, err := uuid.FromString(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of invitation uuid. Error %s", err.Error())
	}

	invite, err := i.iRepo.Get(iuuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "invitation")
	}

	err = i.iRepo.Delete(iuuid, nil)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "invitation")
	}

	if i.msgbus != nil {
		route := i.baseRoutingKey.SetActionDelete().SetObject("invitation").MustBuild()
		evt := &epb.EventInvitationDeleted{
			Id:     invite.Id.String(),
			Email:  invite.Email,
			Name:   invite.Name,
			Role:   upb.RoleType(invite.Role),
			UserId: invite.UserId,
		}
		err = i.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
		}
	}

	return &pb.DeleteResponse{
		Id: req.GetId(),
	}, nil
}

func (i *InvitationServer) UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.UpdateStatusResponse, error) {
	log.Infof("Updating invitation %v", req)

	iuuid, err := uuid.FromString(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of invitation uuid. Error %s", err.Error())
	}

	userInfo, err := i.userClient.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	uuuid, err := uuid.FromString(userInfo.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of invitation uuid. Error %s", err.Error())
	}

	err = i.iRepo.UpdateUserId(iuuid, uuuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "invitation")
	}

	err = i.iRepo.UpdateStatus(iuuid, uint8(req.GetStatus().Number()))
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "invitation")
	}

	invite, err := i.iRepo.Get(iuuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "invitation")
	}

	if i.msgbus != nil {
		route := i.baseRoutingKey.SetActionUpdate().SetObject("invitation").MustBuild()
		log.Infof("Route %s", route)
		evt := &epb.EventInvitationUpdated{
			Id:        invite.Id.String(),
			Link:      invite.Link,
			Email:     invite.Email,
			Name:      invite.Name,
			Role:      upb.RoleType(invite.Role),
			Status:    upb.InvitationStatus(invite.Status),
			UserId:    userInfo.Id,
			ExpiresAt: invite.ExpiresAt.String(),
		}
		err = i.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
		}
	}

	return &pb.UpdateStatusResponse{
		Id:     req.GetId(),
		Status: *req.GetStatus().Enum(),
	}, nil
}

func (i *InvitationServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
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

	return &pb.GetResponse{
		Invitation: dbInvitationToPbInvitation(invitation),
	}, nil
}

func (u *InvitationServer) GetByEmail(ctx context.Context, req *pb.GetByEmailRequest) (*pb.GetByEmailResponse, error) {
	log.Infof("Getting invitation %v", req)

	if req.GetEmail() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Email is required")
	}

	invitation, err := u.iRepo.GetByEmail(req.GetEmail())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "invitation")
	}

	return &pb.GetByEmailResponse{
		Invitation: dbInvitationToPbInvitation(invitation),
	}, nil
}

func (i *InvitationServer) GetAll(ctx context.Context, req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	log.Infof("Getting invitations")

	invitations, err := i.iRepo.GetAll()
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "invitations")
	}

	return &pb.GetAllResponse{
		Invitations: dbInvitationsToPbInvitations(invitations),
	}, nil
}

func dbInvitationToPbInvitation(invitation *db.Invitation) *pb.Invitation {
	return &pb.Invitation{
		Id:       invitation.Id.String(),
		Link:     invitation.Link,
		Email:    invitation.Email,
		Role:     upb.RoleType(invitation.Role),
		Name:     invitation.Name,
		Status:   upb.InvitationStatus(invitation.Status),
		UserId:   invitation.UserId,
		ExpireAt: invitation.ExpiresAt.String(),
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
