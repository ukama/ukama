//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/common/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	rconf "github.com/num30/config"
	"github.com/sirupsen/logrus"
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
		logrus.Fatalf("Failed to read config: %v", err)
	}

	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	logrus.Infof("Config: %+v\n", tConfig)
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

	logrus.Infoln("Connecting to service ", tConfig.ServiceHost)
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
			SubscriberID: subscriberID,
			NetworkID:    networkID,
			PackageID:    packageID,
			SimToken:     simToken,
			SimType:      simType,
		})
		assert.NoError(t, err)
	})

	t.Run("GetSim", func(t *testing.T) {
		_, err := c.GetSim(ctx, &pb.GetSimRequest{
			SimID: simResp.Sim.Id,
		})
		assert.NoError(t, err)
	})

	t.Run("ToggleSimStatus", func(t *testing.T) {
		_, err := c.ToggleSimStatus(ctx, &pb.ToggleSimStatusRequest{
			SimID:  simResp.Sim.Id,
			Status: "active",
		})
		assert.NoError(t, err)
	})

	t.Run("GetSimsBySubscriber", func(t *testing.T) {
		_, err := c.GetSimsBySubscriber(ctx, &pb.GetSimsBySubscriberRequest{
			SubscriberID: subscriberID,
		})
		assert.NoError(t, err)
	})

	t.Run("GetSimsByNetwork", func(t *testing.T) {
		_, err := c.GetSimsByNetwork(ctx, &pb.GetSimsByNetworkRequest{
			NetworkID: networkID,
		})
		assert.NoError(t, err)
	})

	t.Run("AddPackageForSim", func(t *testing.T) {
		_, err := c.AddPackageForSim(ctx, &pb.AddPackageRequest{
			SimID:     simResp.Sim.Id,
			PackageID: packageID,
		})
		assert.NoError(t, err)
	})

	t.Run("AddPackageForSim", func(t *testing.T) {
		_, err := c.AddPackageForSim(ctx, &pb.AddPackageRequest{
			SimID:     simResp.Sim.Id,
			PackageID: packageID,
			StartDate: timestamppb.New(startDate),
		})
		assert.NoError(t, err)
	})

	pkgResp := &pb.GetPackagesBySimResponse{}
	t.Run("GetPackagesBySim", func(t *testing.T) {
		var err error
		pkgResp, err = c.GetPackagesBySim(ctx, &pb.GetPackagesBySimRequest{
			SimID: simResp.Sim.Id,
		})
		assert.NoError(t, err)
	})

	t.Run("SetActivePackageForSim", func(t *testing.T) {
		_, err := c.SetActivePackageForSim(ctx, &pb.SetActivePackageRequest{
			SimID:     simResp.Sim.Id,
			PackageID: pkgResp.Packages[0].Id,
		})
		assert.NoError(t, err)
	})

	t.Run("RemovePackageForSim", func(t *testing.T) {
		_, err := c.RemovePackageForSim(ctx, &pb.RemovePackageRequest{
			SimID:     simResp.Sim.Id,
			PackageID: pkgResp.Packages[0].Id,
		})
		assert.NoError(t, err)
	})

	t.Run("ToggleSimStatus", func(t *testing.T) {
		_, err := c.ToggleSimStatus(ctx, &pb.ToggleSimStatusRequest{
			SimID:  simResp.Sim.Id,
			Status: "inactive",
		})
		assert.NoError(t, err)
	})

	t.Run("DeleteSim", func(t *testing.T) {
		_, err := c.DeleteSim(ctx, &pb.DeleteSimRequest{
			SimID: simResp.Sim.Id,
		})
		assert.NoError(t, err)
	})
}

func CreateSimManagerClient() (*grpc.ClientConn, pb.SimManagerServiceClient, error) {
	logrus.Infoln("Connecting to Sim Manager ", tConfig.ServiceHost)
	context, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(context, tConfig.ServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewSimManagerServiceClient(conn)
	return conn, c, nil
}
