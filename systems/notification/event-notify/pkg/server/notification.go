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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/roles"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/event-notify/pkg"
	"github.com/ukama/ukama/systems/notification/event-notify/pkg/db"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	cnotif "github.com/ukama/ukama/systems/common/notification"
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	pb "github.com/ukama/ukama/systems/notification/event-notify/pb/gen"
)

type EventToNotifyServer struct {
	pb.UnimplementedEventToNotifyServiceServer
	orgName              string
	orgId                string
	notificationRepo     db.NotificationRepo
	userRepo             db.UserRepo
	userNotificationRepo db.UserNotificationRepo
	eventMsgRepo         db.EventMsgRepo
	msgbus               mb.MsgBusServiceClient
	memberkClient        creg.MemberClient
	baseRoutingKey       msgbus.RoutingKeyBuilder
}

func NewEventToNotifyServer(orgName string, orgId string, mc creg.MemberClient, notificationRepo db.NotificationRepo, userRepo db.UserRepo, eventMsgRepo db.EventMsgRepo, userNotificationRepo db.UserNotificationRepo, msgBus mb.MsgBusServiceClient) *EventToNotifyServer {
	return &EventToNotifyServer{
		orgName:              orgName,
		orgId:                orgId,
		notificationRepo:     notificationRepo,
		userNotificationRepo: userNotificationRepo,
		userRepo:             userRepo,
		eventMsgRepo:         eventMsgRepo,
		msgbus:               msgBus,
		memberkClient:        mc,
		baseRoutingKey:       msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (n *EventToNotifyServer) UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.UpdateStatusResponse, error) {
	log.Infof("Update notification %v", req)

	nuuid, err := uuid.FromString(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of notification uuid. Error %s", err.Error())
	}
	err = n.userNotificationRepo.Update(nuuid, req.GetIsRead())
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
	} else {
		return nil, status.Errorf(codes.InvalidArgument,
			"no user uuid provided")
	}

	roleType := upb.RoleType_ROLE_INVALID

	/* validate member of org or member role */
	if req.GetUserId() != "" && req.GetSubscriberId() == "" {
		resp, err := n.memberkClient.GetByUserId(req.GetUserId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid user id. Error %s", err.Error())
		}
		roleType = upb.RoleType(upb.RoleType_value[resp.Member.Role])
	} else if req.GetSubscriberId() != "" {
		roleType = upb.RoleType_ROLE_SUBSCRIBER
	}

	user, err := n.userRepo.GetUsers(ouuid.String(), nuuid.String(), suuid.String(), uuuid.String(), uint8(roleType))
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "eventnotify")
	}

	if len(user) == 0 {
		return nil, status.Errorf(codes.FailedPrecondition,
			"Invalid arguments: no user found")
	}

	notifications := []*db.Notifications{}
	for _, u := range user {
		if u.UserId == req.GetUserId() {
			userNotifications, err := n.userNotificationRepo.GetNotificationsByUserID(u.Id.String())
			log.Infof("Notifications for user: %+v", userNotifications)
			if err != nil {
				return nil, grpc.SqlErrorToGrpc(err, "eventnotify")
			}
			notifications = append(notifications, userNotifications...)
			break
		}
	}

	return &pb.GetAllResponse{
		Notifications: dbNotificationsToPbNotifications(notifications),
	}, nil
}

func dbNotificationsToPbNotifications(notifications []*db.Notifications) []*pb.Notifications {
	res := []*pb.Notifications{}
	for _, i := range notifications {
		n := &pb.Notifications{
			Id:          i.Id.String(),
			Title:       i.Title,
			IsRead:      i.IsRead,
			EventKey:    i.EventKey,
			ResourceId:  i.ResourceId,
			Description: i.Description,
			CreatedAt:   timestamppb.New(i.CreatedAt),
			Type:        upb.NotificationType_name[int32(i.Type)],
			Scope:       upb.NotificationScope_name[int32(i.Scope)],
		}
		res = append(res, n)
	}
	return res
}

func dbNotificationToPbNotification(notification *db.Notification) *pb.Notification {
	return &pb.Notification{
		Id:           notification.Id.String(),
		Title:        notification.Title,
		Description:  notification.Description,
		Type:         upb.NotificationType(notification.Type),
		Scope:        upb.NotificationScope(notification.Scope),
		OrgId:        notification.OrgId,
		NetworkId:    notification.NetworkId,
		EventKey:     notification.EventMsg.Key,
		SubscriberId: notification.SubscriberId,
		UserId:       notification.UserId,
		CreatedAt:    timestamppb.New(notification.CreatedAt),
		ResourceId:   notification.ResourceId,
		EventMsg:     notification.EventMsg.Data.Bytes,
	}
}

func removeDuplicatesIfAny(users []*db.Users) []*db.Users {
	m := map[db.Users]struct{}{}
	usersList := []*db.Users{}

	for _, u := range users {
		if _, ok := m[*u]; !ok {
			usersList = append(usersList, u)
			m[*u] = struct{}{}
		}
	}

	return usersList
}

func (n *EventToNotifyServer) filterUsersForNotification(orgId string, subscriberId string, userId string, scope cnotif.NotificationScope) ([]*db.Users, error) {
	var userList []*db.Users
	var err error
	roleTypes := cnotif.NotificationScopeToRoles[scope]
	done := make(chan bool)

	go func() {

		/* user specific notification */
		if userId != "" && userId != db.EmptyUUID {
			log.Debugf("Getting user with id: %s", userId)
			user, err := n.userRepo.GetUser(userId)
			if err != nil {
				log.Errorf("Failed to get user with userID %s.Error: %+v", userId, err)
			} else {
				userList = append(userList, user)
			}

		}

		/* subscriber specific notification */
		if subscriberId != "" && subscriberId != db.EmptyUUID {
			log.Debugf("Getting subscriber with id: %s", userId)
			user, err := n.userRepo.GetSubscriber(subscriberId)
			if err != nil {
				log.Errorf("Failed to get user with subscriberID %s.Error: %+v", subscriberId, err)
			} else {
				userList = append(userList, user)
			}
		}

		/* Get user based on notification scope
		this would work for OWNER, ADMIN and VENDOR */
		log.Debugf("Getting user with roles: %+v", roleTypes)
		users, err := n.userRepo.GetUserWithRoles(orgId, roleTypes)
		if err != nil {
			log.Errorf("Failed to get user with roles %+v.Error: %+v", roleTypes, err)
		} else {
			userList = append(userList, users...)
		}

		done <- true
	}()

	<-done

	return removeDuplicatesIfAny(userList), err
}

func (n *EventToNotifyServer) storeNotification(dn *db.Notification) error {
	err := n.notificationRepo.Add(dn)
	if err != nil {
		log.Errorf("Error adding notification to db %v", err)
	}

	users, err := n.filterUsersForNotification(dn.OrgId, dn.SubscriberId, dn.UserId, dn.Scope)

	if err != nil {
		log.Errorf("Error getting users from db %v", err)
		return err
	}

	un := []*db.UserNotification{}

	for _, u := range users {

		/* Only add vaild notifcation scope for the User */
		if IsValidNotificationScopeForRole(u.Role, dn.Scope) {
			un = append(un, &db.UserNotification{
				Id:             uuid.NewV4(),
				NotificationId: dn.Id,
				UserId:         u.Id,
				IsRead:         false,
			})
		}
	}

	err = n.userNotificationRepo.Add(un)
	return err
}

func (n *EventToNotifyServer) storeUser(user *db.Users) error {
	err := n.userRepo.Add(user)
	if err != nil {
		log.Errorf("Error adding user to db %v", err)
	}
	return nil
}

func (n *EventToNotifyServer) storeEvent(event *db.EventMsg) (uint, error) {
	id, err := n.eventMsgRepo.Add(event)
	if err != nil {
		log.Errorf("Error adding event to db %v", err)
		return 0, err
	}
	return id, nil
}

func IsValidNotificationScopeForRole(r roles.RoleType, s cnotif.NotificationScope) bool {
	valid := false
	for _, v := range cnotif.RoleToNotificationScopes[r] {
		if v == s {
			valid = true
		}
	}
	return valid
}
