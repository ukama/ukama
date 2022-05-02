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
	pb "github.com/ukama/ukama/services/foo/pb/gen"
	"google.golang.org/grpc"
)

type TestConfig struct {
	FooHost string
}

var tConfig *TestConfig

func init() {
	tConfig := &TestConfig{
		FooHost: "localhost:9090",
	}

	config.LoadConfig("integration", tConfig)
	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	logrus.Infof("%+v", tConfig)
}

func Test_FullFlow(t *testing.T) {
	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logrus.Infoln("Connecting to service ", tConfig.FooHost)
	conn, err := grpc.DialContext(ctx, tConfig.FooHost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		assert.NoError(t, err, "did not connect: %v", err)
		return
	}
	defer conn.Close()

	c := pb.NewFooServiceClient(conn)

	// Contact the server and print out its response.

	r, err := c.GetFoo(ctx, &pb.GetFooRequest{Name: "some_name"})
	assert.NoError(t, err)
	assert.NotNil(t, r)
	// Asserts goes here
}
