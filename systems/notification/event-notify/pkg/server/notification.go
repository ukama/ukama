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

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/event-notify/pkg"
	"github.com/ukama/ukama/systems/notification/event-notify/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	pb "github.com/ukama/ukama/systems/notification/event-notify/pb/gen"
)

const uuidParsingError = "Error parsing UUID"

type EventToNotifyServer struct {
	pb.UnimplementedEventToNotifyServiceServer
	orgName              string
	orgId                string
	notificationRepo     db.NotificationRepo
	userRepo             db.UserRepo
	userNotificationRepo db.UserNotificationRepo
	msgbus               mb.MsgBusServiceClient
	baseRoutingKey       msgbus.RoutingKeyBuilder
}

func NewEventToNotifyServer(orgName string, orgId string, notificationRepo db.NotificationRepo, userRepo db.UserRepo, userNotificationRepo db.UserNotificationRepo, msgBus mb.MsgBusServiceClient) *EventToNotifyServer {
	return &EventToNotifyServer{
		orgName:              orgName,
		orgId:                orgId,
		notificationRepo:     notificationRepo,
		userNotificationRepo: userNotificationRepo,
		userRepo:             userRepo,
		msgbus:               msgBus,
		baseRoutingKey:       msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (n *EventToNotifyServer) NotificationsStream(req *pb.NotificationsStreamRequest, srv pb.EventToNotifyService_NotificationsStreamServer) error {
	log.Info("Notification stream started")
	notification := &pb.Notification{
		Id:           "A8C62136-E56C-4F93-9003-3307329117C2",
		Title:        "Test",
		Description:  "Test",
		Type:         pb.NotificationType_INFO,
		Scope:        pb.NotificationScope_ORG,
		OrgId:        "ORG_ID",
		NetworkId:    "NETWORK_ID",
		SubscriberId: "SUBSCRIBER_ID",
		UserId:       "USER_ID",
		IsRead:       false,
	}

	if err := srv.Send(notification); err != nil {
		log.Info("error generating response")
		return err
	}

	return nil
}

func (n *EventToNotifyServer) UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.UpdateStatusResponse, error) {
	log.Infof("Update notification %v", req)

	nuuid, err := uuid.FromString(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of notification uuid. Error %s", err.Error())
	}
	err = n.notificationRepo.Update(nuuid, req.GetIsRead())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "eventnotify")
	}

	return &pb.UpdateStatusResponse{
		Id: nuuid.String(),
	}, nil
}

func (n *EventToNotifyServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Infof("Getting notification %v", req)

	nuuid, err := uuid.FromString(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of notification uuid. Error %s", err.Error())
	}
	notification, err := n.notificationRepo.Get(nuuid)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "eventnotify")
	}

	return &pb.GetResponse{
		Notification: dbNotificationToPbNotification(notification),
	}, nil
}

func (n *EventToNotifyServer) GetAll(ctx context.Context, req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	log.Infof("Getting notifications %v", req)

	var ouuid, nuuid, suuid, uuuid uuid.UUID
	var err error

	if req.GetOrgId() != "" {
		ouuid, err = uuid.FromString(req.GetOrgId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid format of org uuid. Error %s", err.Error())
		}
	}

	if req.GetNetworkId() != "" {
		nuuid, err = uuid.FromString(req.GetNetworkId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid format of network uuid. Error %s", err.Error())
		}
	}

	if req.GetSubscriberId() != "" {
		suuid, err = uuid.FromString(req.GetSubscriberId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid format of subscriber uuid. Error %s", err.Error())
		}
	}

	if req.GetUserId() != "" {
		uuuid, err = uuid.FromString(req.GetUserId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid format of user uuid. Error %s", err.Error())
		}
	}

	notifications, err := n.notificationRepo.GetAll(ouuid.String(), nuuid.String(), suuid.String(), uuuid.String())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "eventnotify")
	}

	return &pb.GetAllResponse{
		Notifications: dbNotificationsToPbNotifications(notifications),
	}, nil
}

func dbNotificationToPbNotification(notification *db.Notification) *pb.Notification {
	return &pb.Notification{
		Id:           notification.Id.String(),
		Title:        notification.Title,
		Description:  notification.Description,
		Type:         pb.NotificationType(pb.NotificationType_value[string(notification.Type)]),
		Scope:        pb.NotificationScope(pb.NotificationScope_value[string(notification.Scope)]),
		OrgId:        notification.OrgId,
		NetworkId:    notification.NetworkId,
		SubscriberId: notification.SubscriberId,
		UserId:       notification.UserId,
	}
}

func dbNotificationsToPbNotifications(notifications []db.Notification) []*pb.Notification {
	res := []*pb.Notification{}

	for _, i := range notifications {
		res = append(res, dbNotificationToPbNotification(&i))
	}

	return res
}

func (n *EventToNotifyServer) getUsersMatchingNotification(orgId string, networkId string, subscriberId string, userId string, role db.RoleType) ([]*db.Users, error) {
	var users []*db.Users
	var err error

	done := make(chan bool)

	go func() {
		users, err = n.userRepo.GetUsers(orgId, networkId, subscriberId, userId, role)
		done <- true
	}()

	<-done

	return users, err
}

func (n *EventToNotifyServer) eventPbToDBNotification(notification *db.Notification) error {
	err := n.notificationRepo.Add(notification)
	if err != nil {
		log.Errorf("Error adding notification to db %v", err)
	}
	users, err := n.getUsersMatchingNotification(notification.OrgId, notification.NetworkId, notification.SubscriberId, notification.UserId, db.Owner)

	if err != nil {
		log.Errorf("Error getting users from db %v", err)
		return err
	}

	un := []*db.UserNotification{}

	for _, u := range users {

		un = append(un, &db.UserNotification{
			Id:             uuid.NewV4(),
			NotificationId: notification.Id,
			UserId:         u.Id,
			IsRead:         false,
		})
	}

	err = n.userNotificationRepo.Add(un)
	if err != nil {
		log.Errorf("Error adding users notification to db %v", err)
	}
	return nil
}

func (n *EventToNotifyServer) storeUser(user *db.Users) error {
	err := n.userRepo.Add(user)
	if err != nil {
		log.Errorf("Error adding user to db %v", err)
	}
	return nil
}
