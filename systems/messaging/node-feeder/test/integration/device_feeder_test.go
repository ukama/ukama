//go:build integration
// +build integration

package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coredns/coredns/plugin/pkg/log"
	amqp "github.com/rabbitmq/amqp091-go"
	uuid2 "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/msgbus"
	pbnet "github.com/ukama/ukama/systems/messaging/net/pb/gen"
	"github.com/ukama/ukama/systems/messging/node-feeder/pkg"
	pb "github.com/ukama/ukama/systems/registry/pb/gen"
	"github.com/wagslane/go-rabbitmq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
	"testing"
	"time"
)

type TestConf struct {
	RegistryHost string
	QueueUri     string
	DevicePort   int
	WaitingTime  int
	DevicesCount int
	NetHost      string
}

var testConf = &TestConf{}

const orgName = "device-feeder-integration-tests-org"

func init() {
	testConf = &TestConf{
		QueueUri:     "amqp://guest:guest@localhost:5672/",
		RegistryHost: "localhost:9090",
		NetHost:      "localhost:9090",
		DevicePort:   8080, // dummy device port
		WaitingTime:  10,   // how long dummy node waits for the request from device feeder
		DevicesCount: 3,    // how many devices to create
	}

	config.LoadConfig("integration", testConf)
	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	logrus.Infof("%+v", testConf)

}

// in this test we create a dummy http server and send requests to it via device feeder
func Test_FullFlow(t *testing.T) {

	log.Info("Preparing data for device-feeder test")
	conn, regClient, nodes := PrepareRegistryData(t)
	defer conn.Close()
	defer cleanupData(regClient, nodes)

	log.Infof("Send message to rabbitmq")
	err := sendMessageToQueue(pkg.DevicesUpdateRequest{
		Target:     orgName + ".*",
		Path:       "testEndpoint",
		HttpMethod: "GET",
	})
	if err != nil {
		assert.FailNow(t, "Failed to send message to rabbitmq", err)
	}

	d := DisposableServer{
		requestsCount: testConf.DevicesCount,
		port:          testConf.DevicePort,
	}

	log.Infof("Wait for device-feeder to process message for %d seconds", testConf.WaitingTime)
	isRequesReceived := d.WaitForRequest(time.Duration(testConf.WaitingTime) * time.Second)

	assert.True(t, isRequesReceived, "Request was not received")
}

func cleanupData(client pb.RegistryServiceClient, nodes []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	log.Infof("Cleanup data")
	for _, n := range nodes {
		log.Infof("Delete node: %s", n)
		_, err := client.DeleteNode(ctx, &pb.DeleteNodeRequest{
			NodeId: n,
		})
		if err != nil {
			log.Errorf("Error deleting node %s: %s", n, err)
		}
	}
}

func PrepareRegistryData(t *testing.T) (*grpc.ClientConn, pb.RegistryServiceClient, []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	logrus.Infoln("Connecting to registry ", testConf.RegistryHost)
	regConn, err := grpc.DialContext(ctx, testConf.RegistryHost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		assert.FailNow(t, "Failed to connect to registry", err)
		return nil, nil, nil
	}

	c := pb.NewRegistryServiceClient(regConn)

	logrus.Infoln("Connecting to net ", testConf.NetHost)
	netConn, err := grpc.DialContext(ctx, testConf.NetHost, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		assert.FailNow(t, "Failed to connect to net service", err)
		return nil, nil, nil
	}
	nt := pbnet.NewNnsClient(netConn)

	ownerId := uuid2.NewV4()

	_, err = c.AddOrg(ctx, &pb.AddOrgRequest{Name: orgName, Owner: ownerId.String()})
	if err != nil {
		log.Warning("error adding org: ", err)
	}
	logrus.Infoln("Added org ", orgName)

	nodes := []string{}
	for i := 0; i < testConf.DevicesCount; i++ {
		node := fmt.Sprintf("UK-TEST10-HNODE-AA-00%02d", i)
		nodes = append(nodes, node)
		log.Infof("Adding node: %s", node)
		_, err = c.AddNode(ctx, &pb.AddNodeRequest{
			Node: &pb.Node{
				NodeId: node,
				State:  pb.NodeState_UNDEFINED,
			},
			OrgName: orgName,
		})

		if err != nil {
			log.Warning("error adding node: ", err)
		}

		// Set IP address for node
		var ip net.IP
		ip, err = getCurrentPodIp()
		log.Infof("Setting %s node ip to %s", node, ip)
		if err != nil {
			assert.FailNow(t, "error getting current pod ip: ", err)
			return nil, nil, nil
		}

		_, err = nt.Set(ctx, &pbnet.SetRequest{
			NodeId: node,
			Ip:     ip.String(),
		})

		if err != nil {
			assert.FailNow(t, "error setting node ip: ", err)
			return nil, nil, nil
		}
	}

	return regConn, c, nodes
}

func sendMessageToQueue(msg pkg.DevicesUpdateRequest) error {
	rabbit, err := rabbitmq.NewPublisher(testConf.QueueUri, amqp.Config{})
	if err != nil {
		return err
	}

	message, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = rabbit.Publish(message, []string{string(msgbus.DeviceFeederRequestRoutingKey)},
		rabbitmq.WithPublishOptionsExchange(msgbus.DefaultExchange),
		rabbitmq.WithPublishOptionsExpiration("10000")) //  10 sec, not to flood the queue

	if err != nil {
		log.Fatal("Error publishing message: ", err)
	}

	return err
}

func getCurrentPodIp() (net.IP, error) {
	hostName := os.Getenv("HOSTNAME")
	addr, err := net.LookupIP(hostName)
	if err != nil {
		return nil, fmt.Errorf("Error getting IP for host %s: %v", hostName, err)
	}

	if len(addr) < 1 {
		return nil, fmt.Errorf("No IP found for host %s", hostName)
	}

	return addr[0], nil
}
