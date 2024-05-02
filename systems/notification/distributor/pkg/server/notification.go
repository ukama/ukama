/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	log "github.com/sirupsen/logrus"

	pb "github.com/ukama/ukama/systems/notification/distributor/pb/gen"
)

type DistributorServer struct {
	pb.UnimplementedDistributorServiceServer
	orgName string
	orgId   string
}

func NewEventToNotifyServer(orgName string, orgId string) *DistributorServer {
	return &DistributorServer{
		orgName: orgName,
		orgId:   orgId,
	}
}

func (n *DistributorServer) NotificationStream(req *pb.NotificationStreamRequest, srv pb.DistributorService_NotificationStreamServer) error {
	log.Info("Notification stream started")
	// pubsub := n.rdb.Subscribe(context.Background(), "user-notification")
	// for {
	// 	msg, err := pubsub.ReceiveMessage(context.Background())
	// 	if err != nil {
	// 		log.Info("error subscribing to channel")

	// 		return err
	// 	}
	// 	log.Infof("Streaming message %s", msg)
	// }

	// notification := &pb.Notification{
	// 	Id:           "A8C62136-E56C-4F93-9003-3307329117C2",
	// 	Title:        "Test",
	// 	Description:  "Test",
	// 	Type:         pb.NotificationType_INFO,
	// 	Scope:        pb.NotificationScope_ORG,
	// 	OrgId:        "ORG_ID",
	// 	NetworkId:    "NETWORK_ID",
	// 	SubscriberId: "SUBSCRIBER_ID",
	// 	UserId:       "USER_ID",
	// 	IsRead:       false,
	// }

	// if err := srv.Send(notification); err != nil {
	// 	log.Info("error generating response")
	// 	return err
	// }

	return nil
}
