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

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	rconf "github.com/num30/config"
	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/analytics/business/pb/gen"
	uconf "github.com/ukama/ukama/systems/common/config"
)

type TestConfig struct {
	ServiceHost string        `default:"localhost:9090"`
	Queue       *uconf.Queue  `default:"{}"`
	Timeout     time.Duration `default:"3s"`
}

var tConfig *TestConfig

func init() {
	tConfig = &TestConfig{}

	reader := rconf.NewConfReader("integration")
	err := reader.Read(tConfig)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	log.Infof("Config: %+v", tConfig)
}

func connect(t *testing.T) (pb.BusinessServiceClient, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), tConfig.Timeout)

	conn, err := grpc.NewClient(tConfig.ServiceHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cancel()
		assert.NoError(t, err, "did not connect: %v", err)
		return nil, nil, nil
	}

	return pb.NewBusinessServiceClient(conn), ctx, cancel
}

func Test_FullFlow(t *testing.T) {
	client, ctx, cancel := connect(t)
	defer cancel()

	t.Run("GetInventoryReadiness", func(t *testing.T) {
		resp, err := client.GetInventoryReadiness(ctx, &pb.GetInventoryReadinessRequest{})
		assert.NoError(t, err)
		if assert.NotNil(t, resp) {
			assert.Len(t, resp.Kpis, 4)
		}
	})

	t.Run("GetHome", func(t *testing.T) {
		resp, err := client.GetHome(ctx, &pb.GetHomeRequest{
			Window: &pb.Window{Period: "today"},
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("GetSalesOverview", func(t *testing.T) {
		resp, err := client.GetSalesOverview(ctx, &pb.GetSalesOverviewRequest{
			Window: &pb.Window{Period: "month"},
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})
}
