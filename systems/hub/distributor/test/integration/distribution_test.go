//go:build integration
// +build integration

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package integration

import (
	"context"
	"testing"
	"time"

	confr "github.com/num30/config"
	"github.com/ukama/ukama/systems/common/config"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	pb "github.com/ukama/ukama/systems/hub/distributor/pb/gen"
	"google.golang.org/grpc"
)

type TestConfig struct {
	ServiceHost string        `default:"localhost:9090"`
	Queue       *config.Queue `default:"{}"`
	MaxMsgSize  int           `default:"209715200"`
	Store       string        `default:"s3+http://minio:9000/hub-app-local-test/apps"`
}

var tConfig *TestConfig

func init() {
	tConfig = &TestConfig{}
	r := confr.NewConfReader("integration")
	r.Read(tConfig)

	log.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	log.Infof("%+v", tConfig)
}

func Test_FullFlow(t *testing.T) {
	const AppType = "app"
	const AppName = "am-integration-test"
	const AppVer = "0.1.2"
	const FileName = "0.1.2.tar.gz"
	const IndexFile = "0.1.2.caibx"

	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	log.Infoln("Connecting to service ", tConfig.ServiceHost)
	conn, c, err := CreateDistributorClient()
	defer conn.Close()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)
		return
	}

	// Contact the server and print out its response.
	t.Run("CreateChunk", func(t *testing.T) {
		_, err := c.CreateChunk(ctx, &pb.CreateChunkRequest{
			Name:    AppName,
			Type:    AppType,
			Version: AppVer,
			Store:   tConfig.Store,
		})
		assert.NoError(t, err)

	})
}

func CreateDistributorClient() (*grpc.ClientConn, pb.ChunkerServiceClient, error) {
	log.Infoln("Connecting to Lookup ", tConfig.ServiceHost)
	context, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(context, tConfig.ServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(tConfig.MaxMsgSize),
		grpc.MaxCallSendMsgSize(tConfig.MaxMsgSize)))
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewChunkerServiceClient(conn)
	return conn, c, nil
}
