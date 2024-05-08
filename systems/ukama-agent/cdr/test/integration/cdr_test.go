//go:build integration
// +build integration

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
	pb "github.com/ukama/ukama/systems/ukama-agent/cdr/pb/gen"
	"google.golang.org/grpc"
)

type TestConfig struct {
	ServiceHost string        `default:"localhost:9090"`
	Queue       *config.Queue `default:"{}"`
}

var tConfig *TestConfig

func init() {
	tConfig = &TestConfig{}
	r := confr.NewConfReader("integration")
	r.Read(tConfig)

	log.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	log.Infof("%+v", tConfig)
}

func CreateCDRClient() (*grpc.ClientConn, pb.CDRServiceClient, error) {
	log.Infoln("Connecting to CDR ", tConfig.ServiceHost)
	context, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(context, tConfig.ServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewCDRServiceClient(conn)
	return conn, c, nil
}

func Test_FullFlow(t *testing.T) {
	Imsi := "123456789012345678"
	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	log.Infoln("Connecting to service ", tConfig.ServiceHost)
	conn, c, err := CreateCDRClient()
	defer conn.Close()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)
		return
	}

	t.Run("ReadUsage", func(t *testing.T) {
		_, err := c.GetUsage(ctx, &pb.UsageReq{
			Imsi: Imsi,
		})
		assert.NoError(t, err)
	})

}
