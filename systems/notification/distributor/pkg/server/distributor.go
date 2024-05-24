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
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/ukama/ukama/systems/notification/distributor/pkg/providers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	log "github.com/sirupsen/logrus"

	uconf "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/notification/distributor/pb/gen"
	enpb "github.com/ukama/ukama/systems/notification/event-notify/pb/gen"
)

type DistributorServer struct {
	pb.UnimplementedDistributorServiceServer
	DBConfig           *uconf.Database
	eventNotifyService providers.EventNotifyClientProvider
	orgName            string
	orgId              string
}

func NewEventToNotifyServer(orgName string, orgId string, dbConfig *uconf.Database, eventNotifyService providers.EventNotifyClientProvider) *DistributorServer {
	return &DistributorServer{
		DBConfig:           dbConfig,
		orgId:              orgId,
		orgName:            orgName,
		eventNotifyService: eventNotifyService,
	}
}

func (n *DistributorServer) NotificationStream(req *pb.NotificationStreamRequest, srv pb.DistributorService_NotificationStreamServer) error {
	log.Info("Notification stream started")

	svc, err := n.eventNotifyService.GetClient()
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", "postgresql://"+n.DBConfig.Username+":"+n.DBConfig.Password+"@"+n.DBConfig.Host+":"+strconv.Itoa(n.DBConfig.Port)+"/"+n.DBConfig.DbName+"?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbCS := fmt.Sprintf("dbname=%s user=%s password=%s sslmode=disable", n.DBConfig.DbName, n.DBConfig.Username, n.DBConfig.Password)
	listener := pq.NewListener(dbCS, 10*time.Second, time.Minute, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Println(err.Error())
		}
	})

	err = listener.Listen("user_notifications_channel")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		select {
		case notification := <-listener.Notify:
			log.Info(notification)
			params := strings.Split(notification.Extra, ",")
			isRead, _ := strconv.ParseBool(params[2])
			res, err := svc.Get(srv.Context(), &enpb.GetRequest{Id: params[1]})
			if err != nil {
				log.Errorf("Error getting notification: %v", err)
				continue
			}
			un := pb.Notification{
				IsRead:       isRead,
				Id:           res.Notification.Id,
				OrgId:        res.Notification.OrgId,
				Title:        res.Notification.Title,
				UserId:       res.Notification.UserId,
				NetworkId:    res.Notification.NetworkId,
				Description:  res.Notification.Description,
				SubscriberId: res.Notification.SubscriberId,
				ForRole:      pb.RoleType(res.Notification.ForRole),
				Type:         pb.NotificationType(res.Notification.Type),
				Scope:        pb.NotificationScope(res.Notification.Scope),
			}
			log.Infof("Sending notification: %v", un)

			err = srv.Send(&un)
			if err != nil {
				log.Errorf("Error sending notification: %v", err)
				continue
			}
		}
	}
}

func (n *DistributorServer) GetNotifications(ctx context.Context, req *pb.NotificationsRequest) (*pb.NotificationsResponse, error) {
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

	notifications, err := n.s.GetAll(ouuid.String(), nuuid.String(), suuid.String(), uuuid.String())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "eventnotify")
	}

	return &pb.GetAllResponse{
		Notifications: dbNotificationsToPbNotifications(notifications),
	}, nil
}

