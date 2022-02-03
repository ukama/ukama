//go:build integration
// +build integration

package integration

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	uuid2 "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/ukama/ukamaX/cloud/device-feeder/pkg"
	pbnet "github.com/ukama/ukamaX/cloud/net/pb/gen"
	pb "github.com/ukama/ukamaX/cloud/registry/pb/gen"
	"github.com/ukama/ukamaX/common/msgbus"
	"github.com/wagslane/go-rabbitmq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
	"time"
)

type TestConfig struct {
	RegistryHost string
	QueueUri     string
	DevicePort   int
	WaitingTime  int
	DevicesCount int
	NetHost      string
}

type IntegrationTestSuite struct {
	suite.Suite
	config  *TestConfig
	orgName string
}

func NewIntegrationTestSuite(config *TestConfig) *IntegrationTestSuite {
	return &IntegrationTestSuite{config: config, orgName: "device-feeder-integration-tests-org"}
}

func (i *IntegrationTestSuite) Test_FullFlow() {

	log.Info("Preparing data for device-feeder test")
	conn, regClient, nodes := i.PrepareRegistryData()
	defer conn.Close()
	defer cleanupData(regClient, nodes)

	log.Infof("Send message to rabbitmq")
	err := i.sendMessageToQueue(pkg.DevicesUpdateRequest{
		Target:     i.orgName + ".*",
		Path:       "testEndpoint",
		HttpMethod: "GET",
	})
	if err != nil {
		i.FailNow("Failed to send message to rabbitmq", err)
	}

	d := DisposableServer{
		requestsCount: i.config.DevicesCount,
		port:          i.config.DevicePort,
	}

	log.Infof("Wait for device-feeder to process message for %d seconds", i.config.WaitingTime)
	isRequesReceived := d.WaitForRequest(time.Duration(i.config.WaitingTime) * time.Second)

	i.True(isRequesReceived, "Request was not received")
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

func (is *IntegrationTestSuite) PrepareRegistryData() (*grpc.ClientConn, pb.RegistryServiceClient, []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	log.Infoln("Connecting to registry ", is.config.RegistryHost)
	regConn, err := grpc.DialContext(ctx, is.config.RegistryHost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		is.FailNow("Failed to connect to registry", err)
		return nil, nil, nil
	}

	c := pb.NewRegistryServiceClient(regConn)

	log.Infoln("Connecting to net ", is.config.NetHost)
	netConn, err := grpc.DialContext(ctx, is.config.NetHost, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		is.FailNow("Failed to connect to net service", err)
		return nil, nil, nil
	}
	nt := pbnet.NewNnsClient(netConn)

	ownerId := uuid2.NewV4()

	_, err = c.AddOrg(ctx, &pb.AddOrgRequest{Name: is.orgName, Owner: ownerId.String()})
	if err != nil {
		log.Warning("error adding org: ", err)
	}

	nodes := []string{}
	for i := 0; i < is.config.DevicesCount; i++ {
		node := fmt.Sprintf("UK-TEST10-HNODE-AA-00%02d", i)
		nodes = append(nodes, node)
		log.Infof("Adding node: %s", node)
		_, err = c.AddNode(ctx, &pb.AddNodeRequest{
			Node: &pb.Node{
				NodeId: node,
				State:  pb.NodeState_UNDEFINED,
			},
			OrgName: is.orgName,
		})

		if err != nil {
			log.Warning("error adding node: ", err)
		}

		// Set IP address for node
		var ip net.IP
		ip, err = is.getCurrentPodIp()
		log.Infof("Setting %s node ip to %s", node, ip)
		if err != nil {
			is.FailNow("error getting current pod ip: ", err)
			return nil, nil, nil
		}

		_, err = nt.Set(ctx, &pbnet.SetRequest{
			NodeId: node,
			Ip:     ip.String(),
		})

		if err != nil {
			is.FailNow("error setting node ip: ", err)
			return nil, nil, nil
		}
	}

	return regConn, c, nodes
}

func (i *IntegrationTestSuite) handleResponse(err error, r interface{}) {
	fmt.Printf("Response: %v\n", r)
	i.Assert().NoErrorf(err, "Request failed: %v\n", err)
}

func (i *IntegrationTestSuite) sendMessageToQueue(msg pkg.DevicesUpdateRequest) error {
	rabbit, err := rabbitmq.NewPublisher(i.config.QueueUri, amqp.Config{})
	i.Assert().NoError(err)
	if err != nil {
		return err
	}

	message, err := json.Marshal(msg)
	i.Assert().NoError(err)

	err = rabbit.Publish(message, []string{string(msgbus.DeviceFeederRequestRoutingKey)},
		rabbitmq.WithPublishOptionsExchange(msgbus.DefaultExchange),
		rabbitmq.WithPublishOptionsExpiration("10000")) //  10 sec, not to flood the queue

	if err != nil {
		log.Fatal("Error publishing message: ", err)
	}

	return err
}

func (i *IntegrationTestSuite) getCurrentPodIp() (net.IP, error) {
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
