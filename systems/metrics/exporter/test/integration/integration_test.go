//go:build integration
// +build integration

package integration

import (
	"context"
	"math/rand"
	"testing"
	"time"

	confr "github.com/num30/config"
	uuid "github.com/satori/go.uuid"
	"github.com/ukama/ukama/systems/common/config"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/anypb"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/metrics/exporter/pb/gen"
	"google.golang.org/grpc"
)

type TestConfig struct {
	ServiceHost string        `default:"localhost:9095"`
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

func Test_FullFlow(t *testing.T) {
	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	log.Infoln("Connecting to service ", tConfig.ServiceHost)
	conn, c, err := CreateEventClient()
	defer conn.Close()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)
		return
	}

	simUsage := pb.SimUsage{
		Id:           "b20c61f1-1c5a-4559-bfff-cd00f746697d",
		SubscriberID: "c214f255-0ed6-4aa1-93e7-e333658c7318",
		NetworkID:    "40987edb-ebb6-4f84-a27c-99db7c136127",
		OrgID:        "880f7c63-eb57-461a-b514-248ce91e9b3e",
		Type:         "test_simple",
		BytesUsed:    uint64(rand.Int63n(4096000)),
		SessionId:    uuid.NewV4().String(),
		StartTime:    time.Now().Unix() - int64(rand.Intn(30000)),
		EndTime:      time.Now().Unix(),
	}

	anyE, err := anypb.New(&simUsage)
	assert.NoError(t, err)

	// Contact the server and print out its response.
	t.Run("SimUsageEvent", func(t *testing.T) {
		_, err := c.EventNotification(ctx, &epb.Event{
			RoutingKey: "event.cloud.simmanager.sim.usage",
			Msg:        anyE,
		})
		assert.NoError(t, err)

	})
}

func CreateEventClient() (*grpc.ClientConn, epb.EventNotificationServiceClient, error) {
	log.Infoln("Connecting to Lookup ", tConfig.ServiceHost)
	context, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(context, tConfig.ServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	c := epb.NewEventNotificationServiceClient(conn)
	return conn, c, nil
}
