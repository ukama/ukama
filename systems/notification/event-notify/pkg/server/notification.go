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
	orgName          string
	notificationRepo db.NotificationRepo
	msgbus           mb.MsgBusServiceClient
	baseRoutingKey   msgbus.RoutingKeyBuilder
}

func NewEventToNotifyServer(orgName string, notificationRepo db.NotificationRepo, msgBus mb.MsgBusServiceClient) *EventToNotifyServer {
	return &EventToNotifyServer{
		orgName:          orgName,
		notificationRepo: notificationRepo,
		msgbus:           msgBus,
		baseRoutingKey:   msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
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
		return nil, grpc.SqlErrorToGrpc(err, "event-notify")
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
		return nil, grpc.SqlErrorToGrpc(err, "event-notify")
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
		return nil, grpc.SqlErrorToGrpc(err, "event-notify")
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
		IsRead:       notification.IsRead,
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
