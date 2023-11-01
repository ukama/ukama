/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	confr "github.com/num30/config"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/ukama"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	pb "github.com/ukama/ukama/systems/node/configurator/pb/gen"
	"google.golang.org/grpc"
)

type TestConfig struct {
	ServiceHost string        `default:"localhost:9090"`
	Queue       *config.Queue `default:"{}"`
}

var tConfig *TestConfig

var testNodeId = ukama.NewVirtualNodeId("HomeNode")

func init() {
	tConfig = &TestConfig{}
	r := confr.NewConfReader("integration")
	r.Read(tConfig)

	log.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	log.Infof("%+v", tConfig)
}

func Test_FullFlow(t *testing.T) {

	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	log.Infoln("Connecting to service ", tConfig.ServiceHost)
	conn, c, err := CreateConfiguratorClient()
	defer conn.Close()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)
		return
	}

	t.Run("ConfigEvent", func(T *testing.T) {
		_, err := c.ConfigEvent(ctx, &pb.ConfigStoreEvent{})
		assert.NoError(t, err)
	})

	// Contact the server and print out its response.
	t.Run("GetConfigVersion", func(t *testing.T) {
		_, err := c.GetConfigVersion(ctx, &pb.ConfigVersionRequest{
			NodeId: testNodeId.String(),
		})
		assert.NoError(t, err)

	})

}

func CreateConfiguratorClient() (*grpc.ClientConn, pb.ConfiguratorServiceClient, error) {
	log.Infoln("Connecting to Configurator ", tConfig.ServiceHost)
	context, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	conn, err := grpc.DialContext(context, tConfig.ServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewConfiguratorServiceClient(conn)
	return conn, c, nil
}
