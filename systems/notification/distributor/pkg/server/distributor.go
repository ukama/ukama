/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	uconf "github.com/ukama/ukama/systems/common/config"
	pb "github.com/ukama/ukama/systems/notification/distributor/pb/gen"
)

type DistributorServer struct {
	pb.UnimplementedDistributorServiceServer
	DBConfig *uconf.Database
	orgName  string
	orgId    string
}

func NewEventToNotifyServer(orgName string, orgId string, dbConfig *uconf.Database) *DistributorServer {
	return &DistributorServer{
		DBConfig: dbConfig,
		orgId:    orgId,
		orgName:  orgName,
	}
}

func (n *DistributorServer) NotificationStream(req *pb.NotificationStreamRequest, srv pb.DistributorService_NotificationStreamServer) error {
	log.Info("Notification stream started")

	db, err := sql.Open("postgres", "postgresql://"+n.DBConfig.Username+":"+n.DBConfig.Password+"@"+n.DBConfig.Host+":"+strconv.Itoa(n.DBConfig.Port)+"/"+n.DBConfig.DbName+"?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Pinged successfully")

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
			log.Println("Received notification:", notification.Extra)
		}
	}
}
