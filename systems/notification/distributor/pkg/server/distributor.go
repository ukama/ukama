/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"github.com/ukama/ukama/systems/notification/distributor/pkg/db"
	"github.com/ukama/ukama/systems/notification/distributor/pkg/providers"
	"google.golang.org/grpc/status"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"

	uconf "github.com/ukama/ukama/systems/common/config"
	pb "github.com/ukama/ukama/systems/notification/distributor/pb/gen"
)

type Clients struct {
	Nucleus    providers.NucleusProvider
	Registry   providers.RegistryProvider
	Subscriber providers.SubscriberProvider
}

type DistributorServer struct {
	pb.UnimplementedDistributorServiceServer
	notify             db.NotifyHandler
	eventNotifyService providers.EventNotifyClientProvider
	clients            Clients
	orgName            string
	orgId              string
}

func NewEventToNotifyServer(clients Clients, orgName string, orgId string, dbConfig *uconf.Database, eventNotifyService providers.EventNotifyClientProvider) *DistributorServer {

	return &DistributorServer{
		notify:             db.NewNotifyHandler(dbConfig, eventNotifyService),
		orgId:              orgId,
		orgName:            orgName,
		eventNotifyService: eventNotifyService,
		clients:            clients,
	}
}

func (n *DistributorServer) validateRequest(req *pb.NotificationStreamRequest) error {
	//role := memberpb.RoleType_USERS
	if req.GetOrgId() != "" {
		if req.GetOrgId() != n.orgId {
			log.Errorf("Invalid org id %s in request", req.OrgId)
			return status.Errorf(codes.InvalidArgument, "invalid org id")
		}
	}

	/* validate member of org or member role */
	if req.GetUserId() != "" {
		_, err := n.clients.Registry.GetMember(n.orgName, req.GetUserId())
		if err != nil {
			return status.Errorf(codes.InvalidArgument,
				"invalid user id. Error %s", err.Error())
		}
		//role = resp.Member.Role
	}

	if req.GetNetworkId() != "" {
		_, err := n.clients.Registry.GetNetwork(n.orgName, req.GetNetworkId())
		if err != nil {
			return status.Errorf(codes.InvalidArgument,
				"invalid network id. Error %s", err.Error())
		}
	}

	if req.GetSubscriberId() != "" {
		_, err := n.clients.Subscriber.GetSubscriber(n.orgName, req.GetSubscriberId())
		if err != nil {
			return status.Errorf(codes.InvalidArgument,
				"invalid subscriber id. Error %s", err.Error())
		}
	}

	return nil
}

func (n *DistributorServer) GetNotificationStream(req *pb.NotificationStreamRequest, srv pb.DistributorService_GetNotificationStreamServer) error {
	log.Info("Notification stream started for +v.", req)

	err := n.validateRequest(req)
	if err != nil {
		return err
	}

	/* register */
	id, sub := n.notify.Register(req.OrgId, req.NetworkId, req.SubscriberId, req.UserId, req.Scopes)

	defer n.notify.Deregister(req.OrgId)

	for {
		select {
		case data := <-sub.DataChan:
			log.Infof("Sending notification: %+v", data)

			err := srv.Send(data)
			if err != nil {
				log.Errorf("Error sending notification: %v", err)
				continue
			}

		case <-sub.QuitChan:
			log.Errorf("Quiting Notification stream for sub %s with %+v", id, sub)
			goto EXIT

		}
	}

EXIT:
	return nil
}

// func (n *DistributorServer) GetNotifications(ctx context.Context, req *pb.NotificationsRequest) (*pb.NotificationsResponse, error) {
// 	log.Infof("Getting notifications %v", req)

// 	var ouuid, nuuid, suuid, uuuid uuid.UUID
// 	var err error

// 	if req.GetOrgId() != "" {
// 		ouuid, err = uuid.FromString(req.GetOrgId())
// 		if err != nil {
// 			return nil, status.Errorf(codes.InvalidArgument,
// 				"invalid format of org uuid. Error %s", err.Error())
// 		}
// 	}

// 	if req.GetNetworkId() != "" {
// 		nuuid, err = uuid.FromString(req.GetNetworkId())
// 		if err != nil {
// 			return nil, status.Errorf(codes.InvalidArgument,
// 				"invalid format of network uuid. Error %s", err.Error())
// 		}
// 	}

// 	if req.GetSubscriberId() != "" {
// 		suuid, err = uuid.FromString(req.GetSubscriberId())
// 		if err != nil {
// 			return nil, status.Errorf(codes.InvalidArgument,
// 				"invalid format of subscriber uuid. Error %s", err.Error())
// 		}
// 	}

// 	if req.GetUserId() != "" {
// 		uuuid, err = uuid.FromString(req.GetUserId())
// 		if err != nil {
// 			return nil, status.Errorf(codes.InvalidArgument,
// 				"invalid format of user uuid. Error %s", err.Error())
// 		}
// 	}

// 	notifications, err := n.s.GetAll(ouuid.String(), nuuid.String(), suuid.String(), uuuid.String())
// 	if err != nil {
// 		return nil, grpc.SqlErrorToGrpc(err, "eventnotify")
// 	}

// 	return &pb.GetAllResponse{
// 		Notifications: dbNotificationsToPbNotifications(notifications),
// 	}, nil
// }
