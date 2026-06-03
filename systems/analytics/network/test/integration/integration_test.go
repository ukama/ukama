//go:build integration
// +build integration

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package integration

import (
	"context"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tj/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ukama/ukama/systems/common/config"

	rconf "github.com/num30/config"

	pb "github.com/ukama/ukama/systems/analytics/network/pb/gen"
)

var tConfig *TestConfig

func init() {
	// load config
	tConfig = &TestConfig{}

	reader := rconf.NewConfReader("integration")

	err := reader.Read(tConfig)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	log.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	log.Infof("Config: %+v\n", tConfig)
}

type TestConfig struct {
	ServiceHost string        `default:"localhost:9090"`
	Queue       *config.Queue `default:"{}"`
}

func Test_FullFlow(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	log.Infoln("Connecting to analytics network service ", tConfig.ServiceHost)

	conn, err := grpc.NewClient(tConfig.ServiceHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		assert.NoError(t, err, "did not connect: %v", err)

		return
	}
	defer conn.Close()

	c := pb.NewNetworkServiceClient(conn)

	t.Run("GetOverview", func(tt *testing.T) {
		resp, err := c.GetOverview(ctx, &pb.GetOverviewRequest{})
		assert.NoError(tt, err)
		assert.NotNil(tt, resp)
	})

	t.Run("GetTopology", func(tt *testing.T) {
		resp, err := c.GetTopology(ctx, &pb.GetTopologyRequest{})
		assert.NoError(tt, err)
		assert.NotNil(tt, resp)
	})

	t.Run("GetSites", func(tt *testing.T) {
		resp, err := c.GetSites(ctx, &pb.GetSitesRequest{})
		assert.NoError(tt, err)
		assert.NotNil(tt, resp)
	})

	t.Run("GetNodes", func(tt *testing.T) {
		resp, err := c.GetNodes(ctx, &pb.GetNodesRequest{})
		assert.NoError(tt, err)
		assert.NotNil(tt, resp)
	})

	t.Run("GetNodePool", func(tt *testing.T) {
		resp, err := c.GetNodePool(ctx, &pb.GetNodePoolRequest{})
		assert.NoError(tt, err)
		assert.NotNil(tt, resp)
	})

	t.Run("GetAlarms", func(tt *testing.T) {
		resp, err := c.GetAlarms(ctx, &pb.GetAlarmsRequest{})
		assert.NoError(tt, err)
		assert.NotNil(tt, resp)
	})

	t.Run("GetEvents", func(tt *testing.T) {
		resp, err := c.GetEvents(ctx, &pb.GetEventsRequest{})
		assert.NoError(tt, err)
		assert.NotNil(tt, resp)
	})
}
