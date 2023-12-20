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

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ukama/ukama/systems/common/config"

	rconf "github.com/num30/config"
	log "github.com/sirupsen/logrus"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
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
	// prerequisites
	// In order for this to pass without errors, we need to substitute all the uuid.NewV4().String()
	// values with real co-related values comming from subscriber-registry and
	// data-plan/packages with available sims in sim-pool.

	// we need real subscriberID from subscriber-registry
	subscriberID := uuid.NewV4().String()

	// networkID should match with subscriberID
	networkID := uuid.NewV4().String()

	// we need real and active packageID from data-plan/package
	packageID := uuid.NewV4().String()

	// simType should match packageID's package SimType
	simType := "INTER_MNO_ALL"

	simToken := ""
	startDate := time.Now().UTC().AddDate(0, 0, 1) // tomorrow

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	log.Infoln("Connecting to service ", tConfig.ServiceHost)
	conn, c, err := CreateSimManagerClient()
	defer conn.Close()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)
		return
	}

	simResp := &pb.AllocateSimResponse{}

	t.Run("AllocateSim", func(t *testing.T) {
		var err error

		simResp, err = c.AllocateSim(ctx, &pb.AllocateSimRequest{
			SubscriberId: subscriberID,
			NetworkId:    networkID,
			PackageId:    packageID,
			SimToken:     simToken,
			SimType:      simType,
		})
		assert.NoError(t, err)
	})

	t.Run("GetSim", func(t *testing.T) {
		_, err := c.GetSim(ctx, &pb.GetSimRequest{
			SimId: simResp.Sim.Id,
		})
		assert.NoError(t, err)
	})

	t.Run("ToggleSimStatus", func(t *testing.T) {
		_, err := c.ToggleSimStatus(ctx, &pb.ToggleSimStatusRequest{
			SimId:  simResp.Sim.Id,
			Status: "active",
		})
		assert.NoError(t, err)
	})

	t.Run("GetSimsBySubscriber", func(t *testing.T) {
		_, err := c.GetSimsBySubscriber(ctx, &pb.GetSimsBySubscriberRequest{
			SubscriberId: subscriberID,
		})
		assert.NoError(t, err)
	})

	t.Run("GetSimsByNetwork", func(t *testing.T) {
		_, err := c.GetSimsByNetwork(ctx, &pb.GetSimsByNetworkRequest{
			NetworkId: networkID,
		})
		assert.NoError(t, err)
	})

	t.Run("AddPackageForSim", func(t *testing.T) {
		_, err := c.AddPackageForSim(ctx, &pb.AddPackageRequest{
			SimId:     simResp.Sim.Id,
			PackageId: packageID,
		})
		assert.NoError(t, err)
	})

	t.Run("AddPackageForSim", func(t *testing.T) {
		_, err := c.AddPackageForSim(ctx, &pb.AddPackageRequest{
			SimId:     simResp.Sim.Id,
			PackageId: packageID,
			StartDate: timestamppb.New(startDate),
		})
		assert.NoError(t, err)
	})

	pkgResp := &pb.GetPackagesBySimResponse{}
	t.Run("GetPackagesBySim", func(t *testing.T) {
		var err error
		pkgResp, err = c.GetPackagesBySim(ctx, &pb.GetPackagesBySimRequest{
			SimId: simResp.Sim.Id,
		})
		assert.NoError(t, err)
	})

	t.Run("SetActivePackageForSim", func(t *testing.T) {
		_, err := c.SetActivePackageForSim(ctx, &pb.SetActivePackageRequest{
			SimId:     simResp.Sim.Id,
			PackageId: pkgResp.Packages[0].Id,
		})
		assert.NoError(t, err)
	})

	t.Run("RemovePackageForSim", func(t *testing.T) {
		_, err := c.RemovePackageForSim(ctx, &pb.RemovePackageRequest{
			SimId:     simResp.Sim.Id,
			PackageId: pkgResp.Packages[0].Id,
		})
		assert.NoError(t, err)
	})

	t.Run("ToggleSimStatus", func(t *testing.T) {
		_, err := c.ToggleSimStatus(ctx, &pb.ToggleSimStatusRequest{
			SimId:  simResp.Sim.Id,
			Status: "inactive",
		})
		assert.NoError(t, err)
	})

	t.Run("DeleteSim", func(t *testing.T) {
		_, err := c.DeleteSim(ctx, &pb.DeleteSimRequest{
			SimId: simResp.Sim.Id,
		})
		assert.NoError(t, err)
	})
}

func CreateSimManagerClient() (*grpc.ClientConn, pb.SimManagerServiceClient, error) {
	log.Infoln("Connecting to Sim Manager ", tConfig.ServiceHost)
	context, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(context, tConfig.ServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewSimManagerServiceClient(conn)
	return conn, c, nil
}
