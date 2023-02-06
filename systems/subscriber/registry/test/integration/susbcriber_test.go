package integration

import (
	"context"
	"testing"
	"time"

	confr "github.com/num30/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/common/config"
	pb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	logrus.Infof("%+v", tConfig)
}

func TestAddSubscriber(t *testing.T) {
	const sysName = "sys"

	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logrus.Infoln("Connecting to service ", tConfig.ServiceHost)
	conn, c, err := CreateSubscriberClient()
	defer conn.Close()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)
		return
	}
	t.Run("Add", func(t *testing.T) {
		_, err := c.Add(ctx, &pb.AddSubscriberRequest{
			FirstName:             "John",
			LastName:              "Doe",
			Email:                 "johndoe@example.com",
			PhoneNumber:           "+1234567890",
			Address:               "123 Main St.",
			IdSerial:              "123456789",
			NetworkID:             "00000000-0000-0000-0000-000000000000",
			ProofOfIdentification: "Drivers License 12345678",
		})
		assert.NoError(t, err)

	})

	t.Run("ListSubscribers", func(t *testing.T) {
		r, err := c.ListSubscribers(ctx, &pb.ListSubscribersRequest{})

		if assert.NoError(t, err) {
			assert.Equal(t, "johndoe@example.com", r.Subscribers[0].Email)
		}
	})
	t.Run("GetByNetwork", func(t *testing.T) {
		r, err := c.GetByNetwork(ctx, &pb.GetByNetworkRequest{
			NetworkID: "00000000-0000-0000-0000-000000000000",
		})

		if assert.NoError(t, err) {
			assert.Equal(t, "john", r.Subscribers[0].FirstName)
		}
	})
	t.Run("Get", func(t *testing.T) {
		_, err := c.Get(ctx, &pb.GetSubscriberRequest{
			SubscriberID: "9dd5b5d8-f9e1-45c3-b5e3-5f5c5b5e9a9f",
		})
		assert.Equal(t, "subscriber not found", err.Error())

	})
	t.Run("Delete", func(t *testing.T) {
		_, err := c.Delete(ctx, &pb.DeleteSubscriberRequest{
			SubscriberID: "9dd5b5d8-f9e1-45c3-b5e3-5f5c5b5e9a9f",
		})
		assert.Equal(t, "subscriber not found", err.Error())

	})

}

func CreateSubscriberClient() (*grpc.ClientConn, pb.RegistryServiceClient, error) {
	logrus.Infoln("Connecting to subsriber-registry ", tConfig.ServiceHost)
	context, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(context, tConfig.ServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewRegistryServiceClient(conn)
	return conn, c, nil
}
