//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/common/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	rconf "github.com/num30/config"
	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/subscriber/test-agent/pb/gen"
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
	const (
		iccid  = "b8f04217beabf6a19e7eb5b3"
		imsi   = "eabf6a19e7eb5b3"
		status = "inactive"
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	log.Infoln("Connecting to service ", tConfig.ServiceHost)
	conn, c, err := CreateTestAgentClient()
	defer conn.Close()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)
		return
	}

	t.Run("GetSim", func(t *testing.T) {
		_, err := c.GetSim(ctx, &pb.GetSimRequest{
			Iccid: iccid,
		})

		assert.NoError(t, err)
	})

	t.Run("DeactivateSim", func(t *testing.T) {
		_, err := c.DeactivateSim(ctx, &pb.DeactivateSimRequest{
			Iccid: iccid,
		})

		assert.Error(t, err)
	})

	t.Run("ActivateSim", func(t *testing.T) {
		_, err := c.ActivateSim(ctx, &pb.ActivateSimRequest{
			Iccid: iccid,
		})

		assert.NoError(t, err)
	})

	t.Run("TerminateSim", func(t *testing.T) {
		_, err := c.TerminateSim(ctx, &pb.TerminateSimRequest{
			Iccid: iccid,
		})

		assert.Error(t, err)
	})

	t.Run("ActivateSim", func(t *testing.T) {
		_, err := c.ActivateSim(ctx, &pb.ActivateSimRequest{
			Iccid: iccid,
		})

		assert.Error(t, err)
	})

	t.Run("DeactivateSim", func(t *testing.T) {
		_, err := c.DeactivateSim(ctx, &pb.DeactivateSimRequest{
			Iccid: iccid,
		})

		assert.NoError(t, err)
	})

	t.Run("TerminateSim", func(t *testing.T) {
		_, err := c.TerminateSim(ctx, &pb.TerminateSimRequest{
			Iccid: iccid,
		})

		assert.NoError(t, err)
	})
}

func CreateTestAgentClient() (*grpc.ClientConn, pb.TestAgentServiceClient, error) {
	log.Infoln("Connecting to Test Agent ", tConfig.ServiceHost)

	context, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	conn, err := grpc.DialContext(context, tConfig.ServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewTestAgentServiceClient(conn)
	return conn, c, nil
}
