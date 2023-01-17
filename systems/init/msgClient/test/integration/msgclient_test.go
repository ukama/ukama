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

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/init/msgClient/internal/db"
	pb "github.com/ukama/ukama/systems/init/msgClient/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
)

type TestConfig struct {
	ServiceHost string        `default:"localhost:9095"`
	Queue       *config.Queue `default:"{}"`
}

var tConfig *TestConfig

var route1 = db.Route{
	Key: "event.cloud.msgClient.testintegration.create",
}

var ServiceUuid = "1ce2fa2f-2997-422c-83bf-92cf2e7334dd"
var service1 = db.Service{
	Name:        "test-service",
	InstanceId:  "1",
	MsgBusUri:   "amqp://guest:guest@localhost:5672",
	ListQueue:   "",
	PublQueue:   "",
	Exchange:    "amq.topic",
	ServiceUri:  "localhost:9090",
	GrpcTimeout: 5,
}

func init() {
	tConfig = &TestConfig{}
	r := confr.NewConfReader("integration")
	r.Read(tConfig)

	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	logrus.Infof("%+v", tConfig)
}

func Test_FullFlow(t *testing.T) {

	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logrus.Infoln("Connecting to service ", tConfig.ServiceHost)
	conn, c, err := CreateMsgBusClient()
	defer conn.Close()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)
		return
	}

	var serviceId string
	// Contact the server and print out its response.
	t.Run("Register", func(t *testing.T) {
		resp, err := c.RegisterService(ctx, &pb.RegisterServiceReq{
			SystemName:  "test-msgClient",
			ServiceName: service1.Name,
			Exchange:    service1.Exchange,
			InstanceId:  service1.InstanceId,
			MsgBusURI:   service1.MsgBusUri,
			ListQueue:   service1.ListQueue,
			PublQueue:   service1.PublQueue,
			ServiceURI:  service1.ServiceUri,
			GrpcTimeout: service1.GrpcTimeout,
			Routes:      []string{route1.Key},
		})
		assert.NoError(t, err)

		if resp != nil && resp.State == pb.REGISTRAION_STATUS_REGISTERED {
			serviceId = resp.ServiceUuid
		}

	})

	t.Run("Start", func(T *testing.T) {
		_, err := c.StartMsgBusHandler(ctx, &pb.StartMsgBusHandlerReq{
			ServiceUuid: serviceId,
		})
		assert.NoError(t, err)

	})

	t.Run("Publish", func(T *testing.T) {

		msg := &pb.StartMsgBusHandlerReq{
			ServiceUuid: serviceId,
		}

		anyMsg, err := anypb.New(msg)
		assert.NoError(t, err)

		_, err = c.PublishMsg(ctx, &pb.PublishMsgRequest{
			ServiceUuid: serviceId,
			RoutingKey:  route1.Key,
			Msg:         anyMsg,
		})
		assert.NoError(t, err)

	})

	t.Run("Stop", func(t *testing.T) {
		_, err := c.StopMsgBusHandler(ctx, &pb.StopMsgBusHandlerReq{
			ServiceUuid: serviceId,
		})
		assert.NoError(t, err)
	})

}

func CreateMsgBusClient() (*grpc.ClientConn, pb.MsgClientServiceClient, error) {
	logrus.Infoln("Connecting to MsgBusClientService ", tConfig.ServiceHost)
	context, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(context, tConfig.ServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewMsgClientServiceClient(conn)
	return conn, c, nil
}
