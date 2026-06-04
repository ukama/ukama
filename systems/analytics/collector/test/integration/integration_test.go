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
	"github.com/ukama/ukama/systems/common/config"

	rconf "github.com/num30/config"
	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/analytics/collector/pb/gen"
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

func Test_GetRefreshState(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	log.Infoln("Connecting to collector ", tConfig.ServiceHost)

	conn, err := grpc.NewClient(tConfig.ServiceHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		assert.NoError(t, err, "did not connect: %v", err)

		return
	}
	defer conn.Close()

	c := pb.NewCollectorServiceClient(conn)

	resp, err := c.GetRefreshState(ctx, &pb.GetRefreshStateRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
