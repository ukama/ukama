//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	pb "github.com/ukama/ukamaX/cloud/net/pb/gen"
	"github.com/ukama/ukamaX/common/config"
	"github.com/ukama/ukamaX/common/msgbus"
	commonpb "github.com/ukama/ukamaX/common/pb/gen/ukamaos/mesh"
	"github.com/ukama/ukamaX/common/ukama"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

/// Integration test suite for net service

var testConf *TestConf

type TestConf struct {
	NetHost        string
	Queue          config.Queue
	NodeBaseDomain string
	DnsHost        string
}

func init() {
	testConf = &TestConf{
		NetHost: "localhost:9090",
		Queue: config.Queue{
			Uri: "amqp://guest:guest@localhost:5672/",
		},
		NodeBaseDomain: "node.mesh",
	}

	config.LoadConfig("integration", testConf)
	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	logrus.Infof("%+v", testConf)
}

func Test_FullFlow(t *testing.T) {
	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logrus.Infoln("Connecting to service ", testConf.NetHost)
	conn, err := grpc.DialContext(ctx, testConf.NetHost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		assert.NoError(t, err, "did not connect: %v", err)
		return
	}
	defer conn.Close()
	c := pb.NewNnsClient(conn)
	nodeId := ukama.NewVirtualHomeNodeId().String()
	const ip = "1.1.1.1"

	t.Run("SetIp", func(tt *testing.T) {
		r, err := c.Set(ctx, &pb.SetRequest{NodeId: nodeId, Ip: ip})
		handleResponse(tt, err, r)
	})

	t.Run("ResolevIp", func(tt *testing.T) {
		r, err := c.Get(ctx, &pb.GetRequest{NodeId: nodeId})
		handleResponse(tt, err, r)
		assert.Equal(tt, ip, r.Ip)
	})

	t.Run("ResolevMissingIp", func(tt *testing.T) {
		_, err := c.Get(ctx, &pb.GetRequest{NodeId: ukama.NewVirtualHomeNodeId().String()})
		s, ok := status.FromError(err)
		assert.True(tt, ok)
		assert.Equal(tt, codes.NotFound, s.Code())
	})

	t.Run("GetIpList", func(tt *testing.T) {
		_, err := c.Set(ctx, &pb.SetRequest{NodeId: ukama.NewVirtualHomeNodeId().String(), Ip: ip})
		assert.NoError(t, err)
		_, err = c.Set(ctx, &pb.SetRequest{NodeId: ukama.NewVirtualHomeNodeId().String(), Ip: "1.1.1.2"})
		assert.NoError(t, err)
		r, err := c.List(ctx, &pb.ListRequest{})
		assert.NoError(t, err)
		// just make sure it's unique list
		assert.Greater(tt, len(r.Ips), 1)
		un := make(map[string]bool)
		for _, i := range r.Ips {
			if _, ok := un[i]; ok {
				assert.Fail(tt, "Duplicate ip")
			}
			un[i] = true
		}
	})

	t.Run("Delete", func(tt *testing.T) {
		_, err := c.Delete(ctx, &pb.DeleteRequest{NodeId: nodeId})
		if assert.NoError(t, err) {
			_, err := c.Get(ctx, &pb.GetRequest{NodeId: nodeId})
			e, ok := status.FromError(err)
			assert.True(tt, ok)
			assert.Equal(tt, codes.NotFound.String(), e.Code().String())
		}
	})
}

// This tests uses DNS to verify that message was received make sure you set DNSHOST config when running locally
func TestListener(t *testing.T) {
	// Arrange
	nodeId := "UK-000000-HNODE-A0-0001"
	ip := fmt.Sprintf("%d.%d.%d.%d",
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(256))

	// Act
	err := sendMessageToQueue(t, nodeId, ip)
	assert.NoError(t, err)

	// Assert
	time.Sleep(2 * time.Second)

	ips := []net.IP{}
	nodeHost := nodeId + "." + testConf.NodeBaseDomain
	if testConf.DnsHost == "" {
		ips, err = net.LookupIP(nodeHost)
	} else {
		ips, err = resolveWithCustomeDns(nodeHost)
	}

	assert.NoError(t, err)

	assert.Equal(t, 1, len(ips))
	assert.Equal(t, ip, ips[0].String())
}

func resolveWithCustomeDns(host string) ([]net.IP, error) {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, testConf.DnsHost)
		},
	}
	ips, err := r.LookupHost(context.Background(), host)
	if err != nil {
		return nil, err
	}
	return []net.IP{
		net.ParseIP(ips[0]),
	}, nil
}

func sendMessageToQueue(t *testing.T, nodeId string, ip string) error {
	rabbit, err := msgbus.NewPublisherClient(testConf.Queue.Uri)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	message, err := proto.Marshal(&commonpb.Link{NodeId: &nodeId, Ip: &ip})
	assert.NoError(t, err)
	err = rabbit.Publish(message, "", msgbus.DeviceQ.Exchange, msgbus.DeviceConnectedRoutingKey, "topic")
	assert.NoError(t, err)
	return err
}

func handleResponse(t *testing.T, err error, r interface{}) {
	fmt.Printf("Response: %v\n", r)
	assert.NoError(t, err, "Request failed: %v\n", err)
	if err != nil {
		t.FailNow()
	}
}
