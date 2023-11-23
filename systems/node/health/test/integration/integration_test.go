/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package integration

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"

// 	"github.com/ukama/ukama/systems/common/config"
// 	"github.com/ukama/ukama/systems/common/uuid"
// 	pb "github.com/ukama/ukama/systems/node/health/pb/gen"

// 	rconf "github.com/num30/config"
// 	log "github.com/sirupsen/logrus"
// 	grpc "google.golang.org/grpc"
// )

// var tConfig *TestConfig

// func init() {
// 	// load config
// 	tConfig = &TestConfig{}

// 	reader := rconf.NewConfReader("integration")

// 	err := reader.Read(tConfig)
// 	if err != nil {
// 		log.Fatalf("Failed to read config: %v", err)
// 	}

// 	log.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
// 	log.Infof("Config: %+v\n", tConfig)
// }

// type TestConfig struct {
// 	ServiceHost string        `default:"localhost:9090"`
// 	Queue       *config.Queue `default:"{}"`
// 	OrgId       string
// 	OrgName     string
// }

// func Test_FullFlow(t *testing.T) {

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancel()

// 	log.Infoln("Connecting to controller ", tConfig.ServiceHost)

// 	conn, err := grpc.DialContext(ctx, tConfig.ServiceHost, grpc.WithInsecure(), grpc.WithBlock())
// 	if err != nil {
// 		assert.NoError(t, err, "did not connect: %v", err)

// 		return
// 	}

// 	c := pb.NewHealhtServiceClient(conn)

// 	var r interface{}
// 	nodeId := uuid.NewV4()
// 	t.Run("GetRunningApps", func(tt *testing.T) {
// 		r, err = c.GetRunningApps(ctx, &pb.GetRunningAppsRequest{
// 			NodeId: nodeId.String(),
// 		})

// 		handleResponse(tt, err, r)
// 	})

// 	t.Run("StoreRunningAppsInfo", func(_ *testing.T) {
// 		r, err = c.StoreRunningAppsInfo(ctx, &pb.StoreRunningAppsInfoRequest{
// 			NodeId: nodeId.String(),
// 			Timestamp:"12-12-2023",
// 			System:[]*pb.System{},
// 			Capps: []*pb.Capps{},
// 		})
// 		assert.NoError(t, err)
// 	})

// }

// func handleResponse(t *testing.T, err error, r interface{}) {
// 	t.Helper()

// 	log.Printf("Response: %v\n", r)

// 	if err != nil {
// 		assert.FailNow(t, "Request failed: %v\n", err)
// 	}
// }
