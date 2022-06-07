//go:build integration
// +build integration

package integration

import (
	"context"
	"github.com/ukama/ukama/services/common/config"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	pb "github.com/ukama/ukama/services/cloud/org/pb/gen"
	"google.golang.org/grpc"
)

type TestConfig struct {
	OrgHost string
}

var tConfig *TestConfig

func init() {
	tConfig := &TestConfig{
		OrgHost: "localhost:9090",
	}

	config.LoadConfig("integration", tConfig)
	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	logrus.Infof("%+v", tConfig)
}

func Test_FullFlow(t *testing.T) {
	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logrus.Infoln("Connecting to service ", tConfig.OrgHost)
	conn, err := grpc.DialContext(ctx, tConfig.OrgHost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		assert.NoError(t, err, "did not connect: %v", err)
		return
	}
	defer conn.Close()

	c := pb.NewOrgServiceClient(conn)

	// Contact the server and print out its response.

	r, err := c.GetOrg(ctx, &pb.GetOrgRequest{Name: "some_name"})
	assert.NoError(t, err)
	assert.NotNil(t, r)
	// Asserts go here
}
