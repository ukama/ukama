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
	"io"
	"os"
	"strings"
	"testing"
	"time"

	confr "github.com/num30/config"
	"github.com/ukama/ukama/systems/common/config"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	pb "github.com/ukama/ukama/systems/hub/artifactmanager/pb/gen"
	"google.golang.org/grpc"
)

type TestConfig struct {
	ServiceHost string        `default:"localhost:9090"`
	Queue       *config.Queue `default:"{}"`
	MaxMsgSize  int           `default:"209715200"`
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
	conn, c, err := CreateArtifactManagerClient()
	defer conn.Close()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)
		return
	}

	f := getFileContent(t)
	defer f.Close()

	cont, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("failed to read testfile: %s", err)
	}

	// Contact the server and print out its response.
	t.Run("StoreArtifact", func(t *testing.T) {
		_, err := c.StoreArtifact(ctx, &pb.StoreArtifactRequest{
			Name:    AppName,
			Type:    pb.ArtifactType(pb.ArtifactType_value[strings.ToUpper(AppType)]),
			Version: AppVer,
			Data:    cont,
		})
		assert.NoError(t, err)

	})

	time.Sleep(10 * time.Second)

	t.Run("GetArtifact", func(t *testing.T) {
		_, err := c.GetArtifact(ctx, &pb.GetArtifactRequest{
			Name:     AppName,
			Type:     pb.ArtifactType(pb.ArtifactType_value[strings.ToUpper(AppType)]),
			FileName: FileName,
		})
		assert.NoError(t, err)

	})

	t.Run("ListArtifacts", func(t *testing.T) {
		_, err := c.ListArtifacts(ctx, &pb.ListArtifactRequest{
			Type: pb.ArtifactType(pb.ArtifactType_value[strings.ToUpper(AppType)]),
		})
		assert.NoError(t, err)

	})

	t.Run("ListVersionsArtifact", func(t *testing.T) {
		_, err := c.GetArtifactVersionList(ctx, &pb.GetArtifactVersionListRequest{
			Name: AppName,
			Type: pb.ArtifactType(pb.ArtifactType_value[strings.ToUpper(AppType)]),
		})
		assert.NoError(t, err)

	})

	t.Run("GetArtifact_IndexFile", func(t *testing.T) {
		_, err := c.GetArtifact(ctx, &pb.GetArtifactRequest{
			Name:     AppName,
			Type:     pb.ArtifactType(pb.ArtifactType_value[strings.ToUpper(AppType)]),
			FileName: IndexFile,
		})
		assert.NoError(t, err)

	})

}

func CreateArtifactManagerClient() (*grpc.ClientConn, pb.ArtifactServiceClient, error) {
	log.Infoln("Connecting to Lookup ", tConfig.ServiceHost)
	context, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(context, tConfig.ServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(tConfig.MaxMsgSize),
		grpc.MaxCallSendMsgSize(tConfig.MaxMsgSize)))
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewArtifactServiceClient(conn)
	return conn, c, nil
}

func getFileContent(t *testing.T) *os.File {
	f, err := os.Open("testdata/capp.tar.gz")
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	return f
}
