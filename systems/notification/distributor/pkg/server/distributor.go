/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"fmt"

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/notification/distributor/pkg/db"
	"github.com/ukama/ukama/systems/notification/distributor/pkg/providers"
	"google.golang.org/grpc/status"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"

	"github.com/ukama/ukama/systems/common/notification"
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	sreg "github.com/ukama/ukama/systems/common/rest/client/subscriber"
	"github.com/ukama/ukama/systems/common/roles"
	pb "github.com/ukama/ukama/systems/notification/distributor/pb/gen"
)

type DistributorServer struct {
	pb.UnimplementedDistributorServiceServer
	notify             db.NotifyHandler
	eventNotifyService providers.EventNotifyClientProvider
	orgName            string
	orgId              string
	networkClient      creg.NetworkClient
	memberkClient      creg.MemberClient
	subscriberClient   sreg.SubscriberClient
	nodeClient         creg.NodeClient
}

func NewDistributorServer(nc creg.NetworkClient, nodec creg.NodeClient, mc creg.MemberClient, sc sreg.SubscriberClient, n db.NotifyHandler, orgName string, orgId string, eventNotifyService providers.EventNotifyClientProvider) *DistributorServer {

	d := &DistributorServer{
		notify:             n,
		orgId:              orgId,
		orgName:            orgName,
		eventNotifyService: eventNotifyService,
		networkClient:      nc,
		memberkClient:      mc,
		subscriberClient:   sc,
		nodeClient:         nodec,
	}

	/* start notification handler routine */
	d.notify.Start()

	return d
}

func (n *DistributorServer) validateRequest(req *pb.NotificationStreamRequest) (roles.RoleType, error) {
	roleType := roles.TYPE_INVALID
	if req.GetOrgId() != "" {
		if req.GetOrgId() != n.orgId {
			log.Errorf("Invalid org id %s in request", req.OrgId)
			return roleType, status.Errorf(codes.InvalidArgument, "invalid org id")
		}
	}

	/* validate member of org or member role */
	if req.GetUserId() != "" {
		resp, err := n.memberkClient.GetByUserId(req.GetUserId())
		if err != nil {
			return roleType, status.Errorf(codes.InvalidArgument,
				"invalid user id. Error %s", err.Error())
		}
		roleType = roles.RoleType(upb.RoleType(upb.RoleType_value[resp.Member.Role]))
	}

	if req.GetNetworkId() != "" {
		_, err := n.networkClient.Get(req.GetNetworkId())
		if err != nil {
			return roleType, status.Errorf(codes.InvalidArgument,
				"invalid network id. Error %s", err.Error())
		}
	}
	if req.GetNodeId() != "" {
		nId, err := ukama.ValidateNodeId(req.NodeId)
		if err == nil {
			_, err = n.nodeClient.Get(nId.String())
		}
		if err != nil {
			return roleType, status.Errorf(codes.InvalidArgument,
				"invalid network id. Error %s", err.Error())
		}
	}

	if req.GetSubscriberId() != "" {
		_, err := n.subscriberClient.Get(req.GetSubscriberId())
		if err != nil {
			return roleType, status.Errorf(codes.InvalidArgument,
				"invalid subscriber id. Error %s", err.Error())
		}
	}

	return roleType, nil
}

func (n *DistributorServer) GetNotificationStream(req *pb.NotificationStreamRequest, srv pb.DistributorService_GetNotificationStreamServer) error {
	log.Infof("Notification stream requested for %+v.", req)
	roleType, err := n.validateRequest(req)
	if err != nil {
		return err
	}

	/* Get valid scopes for request */
	commonScopes := []notification.NotificationScope{}
	if roleType != roles.TYPE_INVALID {
		roleScopes := notification.RoleToNotificationScopes[roleType]
		for _, rs := range req.Scopes {
			rsId := notification.NotificationScope(upb.NotificationScope_value[rs])
			if rsId != notification.NotificationScope(upb.NotificationScope_SCOPE_INVALID) {
				for _, vs := range roleScopes {
					if vs == rsId {
						commonScopes = append(commonScopes, vs)
					}
				}
			}
		}
	} else {
		log.Errorf("Invalid roles %+v for user %s", roleType, req.UserId)
		return fmt.Errorf("invalid role for user")
	}

	/* register */
	id, sub := n.notify.Register(req.OrgId, req.NetworkId, req.SubscriberId, req.UserId, req.NodeId, commonScopes)

	defer func() {
		if err := n.notify.Deregister(id); err != nil {
			// Handle the error, for example, log it
			log.Printf("Error deregistering: %v", err)
		}
	}()

	for {
		select {
		case <-srv.Context().Done():
			log.Infof("Client closed connection for request %+v.Error %+v", req, srv.Context().Err())
			goto EXIT

		case data := <-sub.DataChan:
			log.Infof("Sending notification: %+v", data)

			err = srv.Send(data)
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
