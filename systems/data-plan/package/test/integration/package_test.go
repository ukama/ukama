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

	"github.com/ukama/ukama/systems/common/config"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	rconf "github.com/num30/config"
	log "github.com/sirupsen/logrus"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/data-plan/package/pb/gen"
)

type TestConfig struct {
	ServiceHost string        `default:"localhost:9090"`
	Queue       *config.Queue `default:"{}"`
}

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
func Test_FullFlow(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	log.Infoln("Connecting to service ", tConfig.ServiceHost)
	conn, c, err := CreatePackageClient()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)
		return
	}
	defer conn.Close()

	t.Run("Add", func(t *testing.T) {
		var err error

		_, err = c.Add(ctx, &pb.AddPackageRequest{
			Name:        "Daily-pack",
			OrgID:       "5b5c3f5e-1f3b-4723-8f99-fe0ed6c539d2",
			Active:      true,
			Duration:    1,
			SimType:     "test",
			SmsVolume:   20,
			DataVolume:  12,
			VoiceVolume: 34,
			OrgRatesID:  0,
		})
		assert.NoError(t, err)
	})
	assert.NoError(t, err)

	t.Run("Update", func(t *testing.T) {
		var err error

		_, err = c.Update(ctx, &pb.UpdatePackageRequest{
			PackageID:   uuid.NewV4().String(),
			Name:        "Updated-Daily-pack",
			Duration:    2,
			SmsVolume:   40,
			DataVolume:  24,
			VoiceVolume: 68,
			OrgRatesID:  0,
		})
		assert.NoError(t, err)

	})
	t.Run("Get", func(t *testing.T) {
		packageResp, err := c.Get(ctx, &pb.GetPackageRequest{
			PackageID: uuid.NewV4().String(),
		})
		assert.NoError(t, err)
		assert.Equal(t, packageResp.Package.Name, "Weekly-pack")
		assert.Equal(t, packageResp.Package.Duration, uint64(7))
	})

	t.Run("Delete", func(t *testing.T) {
		_, err := c.Delete(ctx, &pb.DeletePackageRequest{
			PackageID: uuid.NewV4().String(),
		})
		assert.NoError(t, err)
	})
}

func CreatePackageClient() (*grpc.ClientConn, pb.PackagesServiceClient, error) {
	log.Infoln("Connecting to Sim Manager ", tConfig.ServiceHost)
	context, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(context, tConfig.ServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewPackagesServiceClient(conn)
	return conn, c, nil
}
