//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	confr "github.com/num30/config"
	"github.com/ukama/ukama/services/common/config"
	"github.com/ukama/ukama/services/common/ukama"
	"google.golang.org/grpc/credentials/insecure"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	pb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
	"google.golang.org/grpc"
)

type TestConfig struct {
	ServiceHost string        `default:"localhost:9090"`
	Queue       *config.Queue `default:"{}"`
}

var tConfig *TestConfig

var testNodeId = ukama.NewVirtualNodeId("HomeNode")

func init() {
	tConfig = &TestConfig{}
	r := confr.NewConfReader("integration")
	r.Read(tConfig)

	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	logrus.Infof("%+v", tConfig)
}

func Test_FullFlow(t *testing.T) {
	orgName := fmt.Sprintf("lookup-integration-self-test-%d", time.Now().Unix())
	const certs = "ukama_certs"
	const ip = "0.0.0.0"
	const sysName = "sys"
	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logrus.Infoln("Connecting to service ", tConfig.ServiceHost)
	conn, c, err := CreateLookupClient()
	defer conn.Close()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)
		return
	}

	// Contact the server and print out its response.
	t.Run("AddOrg", func(t *testing.T) {
		_, err := c.AddOrg(ctx, &pb.AddOrgRequest{
			OrgName:     orgName,
			Certificate: certs,
			Ip:          ip,
		})
		assert.NoError(t, err)

	})

	t.Run("UpdatedOrg", func(T *testing.T) {
		_, err := c.UpdateOrg(ctx, &pb.UpdateOrgRequest{
			OrgName:     orgName,
			Certificate: certs,
			Ip:          "127.0.0.1",
		})
		assert.NoError(t, err)

	})

	t.Run("AddNode", func(t *testing.T) {
		_, err := c.AddNodeForOrg(ctx, &pb.AddNodeRequest{
			NodeId:  testNodeId.String(),
			OrgName: orgName,
		})
		assert.NoError(t, err)
	})

	t.Run("GetNode", func(t *testing.T) {
		r, err := c.GetNodeForOrg(ctx, &pb.GetNodeRequest{
			NodeId:  testNodeId.String(),
			OrgName: orgName,
		})

		if assert.NoError(t, err) {
			assert.Equal(t, testNodeId.StringLowercase(), r.NodeId)
		}
	})

	t.Run("DeleteNode", func(T *testing.T) {
		_, err := c.DeleteNodeForOrg(ctx, &pb.DeleteNodeRequest{
			NodeId:  testNodeId.String(),
			OrgName: orgName,
		})
		assert.NoError(t, err)
	})

	t.Run("AddSystem", func(t *testing.T) {
		r, err := c.UpdateSystemForOrg(ctx, &pb.UpdateSystemRequest{
			SystemName:  sysName,
			OrgName:     orgName,
			Certificate: certs,
			Ip:          "127.0.0.1",
			Port:        100,
		})
		assert.NoError(t, err)

		_, err = uuid.Parse(r.SystemId)
		assert.NoError(t, err)
	})

	t.Run("GetSystem", func(t *testing.T) {
		r, err := c.GetSystemForOrg(ctx, &pb.GetSystemRequest{
			SystemName: sysName,
			OrgName:    orgName,
		})

		if assert.NoError(t, err) {
			assert.Equal(t, strings.ToLower(sysName), r.SystemName)
		}

		_, err = uuid.Parse(r.SystemId)
		assert.NoError(t, err)
	})

	t.Run("DeleteSystem", func(T *testing.T) {
		_, err := c.DeleteSystemForOrg(ctx, &pb.DeleteSystemRequest{
			SystemName: sysName,
			OrgName:    orgName,
		})
		assert.NoError(t, err)
	})

}

func CreateLookupClient() (*grpc.ClientConn, pb.LookupServiceClient, error) {
	logrus.Infoln("Connecting to Lookup ", tConfig.ServiceHost)
	context, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(context, tConfig.ServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewLookupServiceClient(conn)
	return conn, c, nil
}
