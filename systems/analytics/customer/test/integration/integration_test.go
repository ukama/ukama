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

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/tj/assert"

	log "github.com/sirupsen/logrus"
	rconf "github.com/num30/config"
	pb "github.com/ukama/ukama/systems/analytics/customer/pb/gen"
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
	ServiceHost string `default:"localhost:9090"`
}

func connect(t *testing.T) (pb.CustomerServiceClient, context.Context, context.CancelFunc) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	log.Infoln("Connecting to customer service ", tConfig.ServiceHost)

	conn, err := grpc.NewClient(tConfig.ServiceHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		assert.NoError(t, err, "did not connect: %v", err)
		cancel()
		t.FailNow()
	}

	return pb.NewCustomerServiceClient(conn), ctx, cancel
}

func TestGetOverview(t *testing.T) {
	c, ctx, cancel := connect(t)
	defer cancel()

	resp, err := c.GetOverview(ctx, &pb.GetOverviewRequest{
		Window: &pb.Window{Period: "today"},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Kpis)
}

func TestGetSimPool(t *testing.T) {
	c, ctx, cancel := connect(t)
	defer cancel()

	resp, err := c.GetSimPool(ctx, &pb.GetSimPoolRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Kpis)
}

func TestList(t *testing.T) {
	c, ctx, cancel := connect(t)
	defer cancel()

	resp, err := c.List(ctx, &pb.ListRequest{Page: 1, PageSize: 10})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Meta)
}
