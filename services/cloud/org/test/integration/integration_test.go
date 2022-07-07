//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	confr "github.com/num30/config"
	"github.com/ukama/ukama/services/cloud/org/pkg/queue"
	"github.com/ukama/ukama/services/common/config"
	"github.com/ukama/ukama/services/common/msgbus"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	pb "github.com/ukama/ukama/services/cloud/org/pb/gen"
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

	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	logrus.Infof("%+v", tConfig)
}

func Test_FullFlow(t *testing.T) {
	orgName := fmt.Sprintf("org-integration-self-test-%d", time.Now().Unix())
	owner := uuid.NewString()

	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logrus.Infoln("Connecting to service ", tConfig.ServiceHost)
	conn, c, err := CreateOrgClient()
	defer conn.Close()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)
		return
	}

	// delete an org in any case
	defer deleteOrg(t, c, orgName)

	// Contact the server and print out its response.
	t.Run("CreateOrg", func(t *testing.T) {
		r, err := c.Add(ctx, &pb.AddRequest{
			Org: &pb.Organization{
				Name:  orgName,
				Owner: owner,
			}})

		if assert.NoError(t, err) {
			assert.Equal(t, orgName, r.GetOrg().GetName())
		}
	})

	t.Run("GetOrg", func(t *testing.T) {
		r, err := c.Get(ctx, &pb.GetRequest{Name: orgName})
		if assert.NoError(t, err) {
			assert.Equal(t, orgName, r.Org.Name)
		}
	})

}

func deleteOrg(t *testing.T, c pb.OrgServiceClient, orgName string) {
	logrus.Info("Deleting org ", orgName)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err := c.Delete(ctx, &pb.DeleteRequest{
		Name: orgName,
	})

	assert.NoError(t, err)
}

func Test_Listener(t *testing.T) {
	// Arrange
	ownerId := uuid.NewString()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, c, err := CreateOrgClient()
	defer conn.Close()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)
		return
	}

	// Act
	err = sendMessageToQueue(ownerId)

	// Assert
	assert.NoError(t, err)
	logrus.Info("Sleeping for 2 seconds")
	time.Sleep(2 * time.Second)

	logrus.Info("Getting org: " + ownerId)
	resp, err := c.Get(ctx, &pb.GetRequest{Name: ownerId})
	if assert.NoError(t, err) {
		assert.Equal(t, ownerId, resp.Org.Owner)
	}

}

func CreateOrgClient() (*grpc.ClientConn, pb.OrgServiceClient, error) {
	logrus.Infoln("Connecting to network ", tConfig.ServiceHost)
	context, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(context, tConfig.ServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewOrgServiceClient(conn)
	return conn, c, nil
}

func sendMessageToQueue(ownerId string) error {
	logrus.Info("Sending message to queue")

	rabbit, err := msgbus.NewQPub(tConfig.Queue.Uri, "network-listener-integration-test", os.Getenv("POD_NAME"))
	if err != nil {
		logrus.Errorf("could not create rabbitmq client %+v", err)
		return err
	}

	err = rabbit.Publish(&queue.UserRegisteredBody{
		Id:    ownerId,
		Email: "org-integration-test@gmail.com",
	}, string(msgbus.UserRegisteredRoutingKey))
	if err != nil {
		logrus.Errorf("could not publish message %+v", err)
	}

	return err
}
