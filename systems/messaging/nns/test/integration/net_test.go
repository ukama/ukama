//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/messaging/nns/pb/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"math/rand"
	"net"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
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
			Uri: "amqp://guest:guest@192.168.0.14:5672/",
		},
		NodeBaseDomain: "node.mesh",
		DnsHost:        "localhost:5053",
	}

	config.LoadConfig("integration", testConf)
	log.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	log.Infof("%+v", testConf)
}

func Test_FullFlow(t *testing.T) {
	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	log.Infoln("Connecting to service ", testConf.NetHost)
	conn, err := grpc.DialContext(ctx, testConf.NetHost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		assert.NoError(t, err, "did not connect: %v", err)
		return
	}
	defer conn.Close()
	c := pb.NewNnsClient(conn)
	nodeId := ukama.NewVirtualHomeNodeId().String()
	const meshIp = "1.1.1.1"
	const meshPort int32 = 1000
	const nodeIp = "2.2.2.2"
	const nodePort int32 = 2000

	t.Run("SetIp", func(tt *testing.T) {
		r, err := c.Set(ctx, &pb.SetNodeIPRequest{NodeId: nodeId, MeshIp: meshIp, MeshPort: meshPort, NodeIp: nodeIp, NodePort: nodePort})
		handleResponse(tt, err, r)
	})

	t.Run("ResolveIp", func(tt *testing.T) {
		r, err := c.Get(ctx, &pb.GetNodeIPRequest{NodeId: nodeId})
		handleResponse(tt, err, r)
		assert.Equal(tt, meshIp, r.Ip)
	})

	t.Run("ResolveMissingIp", func(tt *testing.T) {
		_, err := c.Get(ctx, &pb.GetNodeIPRequest{NodeId: ukama.NewVirtualHomeNodeId().String()})
		s, ok := status.FromError(err)
		assert.True(tt, ok)
		assert.Equal(tt, codes.NotFound, s.Code())
	})

	t.Run("GetIpList", func(tt *testing.T) {
		_, err := c.Set(ctx, &pb.SetNodeIPRequest{NodeId: ukama.NewVirtualHomeNodeId().String(), MeshIp: meshIp, MeshPort: meshPort, NodeIp: nodeIp, NodePort: nodePort})
		assert.NoError(t, err)
		_, err = c.Set(ctx, &pb.SetNodeIPRequest{NodeId: ukama.NewVirtualHomeNodeId().String(), MeshIp: "1.1.1.2", MeshPort: meshPort, NodeIp: nodeIp, NodePort: nodePort})
		assert.NoError(t, err)
		r, err := c.List(ctx, &pb.ListNodeIPRequest{})
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
		_, err := c.Delete(ctx, &pb.DeleteNodeIPRequest{NodeId: nodeId})
		if assert.NoError(t, err) {
			_, err := c.Get(ctx, &pb.GetNodeIPRequest{NodeId: nodeId})
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

	nw := uuid.NewV4().String()
	site := uuid.NewV4().String()

	ip := fmt.Sprintf("%d.%d.%d.%d",
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(256))
	var port int32 = 1000
	nIp := fmt.Sprintf("%d.%d.%d.%d",
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(256))

	var nPort int32 = 2000
	// Act
	err := sendOnlineEventToQueue(t, nodeId, ip, port, nIp, nPort)
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

	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(ips))
		assert.Equal(t, ip, ips[0].String())
	}

	err = sendAssignedEventToQueue(t, nodeId, site, nw)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	getOrgNetMap(t, nodeId, site, nw)
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

func sendOnlineEventToQueue(t *testing.T, nodeId string, ip string, port int32, nIp string, nPort int32) error {
	rabbit, err := msgbus.NewPublisherClient(testConf.Queue.Uri)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	msg := &epb.NodeOnlineEvent{NodeId: nodeId, MeshIp: ip, MeshPort: port, NodeIp: nIp, NodePort: nPort}

	anyMsg, err := anypb.New(msg)
	if err != nil {
		return err
	}

	payload, err := proto.Marshal(anyMsg)
	if err != nil {
		return err
	}

	err = rabbit.Publish(payload, "", "amq.topic", "event.cloud.mesh.node.online", "topic")
	assert.NoError(t, err)

	return err
}

func sendAssignedEventToQueue(t *testing.T, nodeId, site, nw string) error {
	rabbit, err := msgbus.NewPublisherClient(testConf.Queue.Uri)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	msg := &epb.NodeAssignedEvent{NodeId: nodeId, Network: nw, Site: site}

	anyMsg, err := anypb.New(msg)
	if err != nil {
		return err
	}

	payload, err := proto.Marshal(anyMsg)
	if err != nil {
		return err
	}

	err = rabbit.Publish(payload, "", "amq.topic", "event.cloud.registry.node.assigned", "topic")
	assert.NoError(t, err)

	return err
}

func getOrgNetMap(t *testing.T, id, site, nw string) {
	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	log.Infoln("Connecting to service ", testConf.NetHost)
	conn, err := grpc.DialContext(ctx, testConf.NetHost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		assert.NoError(t, err, "did not connect: %v", err)
		return
	}
	defer conn.Close()
	c := pb.NewNnsClient(conn)

	t.Run("GetNodeOrgMapList", func(t *testing.T) {
		r, err := c.GetNodeOrgMapList(ctx, &pb.NodeOrgMapListRequest{})
		log.Infof("NodeOrgMap is %+v", r)
		assert.NoError(t, err)
		// just make sure it's unique list
		assert.Greater(t, len(r.Map), 0)
		found := false
		for _, n := range r.Map {
			if id == n.NodeId {
				assert.Equal(t, nw, n.Network)
				assert.Equal(t, site, n.Site)
				found = true
			}
		}

		if !found {
			assert.Fail(t, "Id not found.")
		}
	})

}
func handleResponse(t *testing.T, err error, r interface{}) {
	fmt.Printf("Response: %v\n", r)
	assert.NoError(t, err, "Request failed: %v\n", err)
	if err != nil {
		t.FailNow()
	}
}
