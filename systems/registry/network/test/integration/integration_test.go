//go:build integration
// +build integration

package integration

import (
	"github.com/ukama/ukama/systems/common/config"

	rconf "github.com/num30/config"
	"github.com/sirupsen/logrus"
)

var tConfig *TestConfig

func init() {
	// load config
	tConfig = &TestConfig{}

	reader := rconf.NewConfReader("integration")

	err := reader.Read(tConfig)
	if err != nil {
		logrus.Fatalf("Failed to read config: %v", err)
	}

	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	logrus.Infof("Config: %+v\n", tConfig)
}

type TestConfig struct {
	ServiceHost string        `default:"localhost:9090"`
	Queue       *config.Queue `default:"{}"`
}

// func Test_FullFlow(t *testing.T) {
// const networkName = "test-network"

// orgName := fmt.Sprintf("network-integration-self-test-org-%d", time.Now().Unix())

// ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
// defer cancel()

// logrus.Infoln("Connecting to network ", tConfig.ServiceHost)

// conn, err := grpc.DialContext(ctx, tConfig.ServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
// if err != nil {
// assert.NoError(t, err, "did not connect: %v", err)

// return
// }

// c := pb.NewNetworkServiceClient(conn)
// defer deleteNetworks(t, c, orgName, networkName)

// // Contact the server and print out its response.
// node := ukama.NewVirtualNodeId("HomeNode")
// nodeToDelete := ukama.NewVirtualHomeNodeId()

// var r interface{}

// t.Run("AddNetwork", func(tt *testing.T) {
// r, err = c.Add(ctx, &pb.AddRequest{
// Name:    networkName,
// OrgName: orgName,
// })

// handleResponse(tt, err, r)
// })

// t.Run("AddAndUpdateNode", func(tt *testing.T) {
// nodeName := "HomeNodeX"

// addResp, err := c.AddNode(ctx, &pb.AddNodeRequest{
// Node: &pb.Node{
// NodeId: node.String(),
// State:  pb.NodeState_UNDEFINED,
// Name:   nodeName,
// },
// OrgName: orgName,
// Network: networkName,
// })

// if handleResponse(tt, err, addResp) {
// assert.NotNil(tt, addResp.Node)
// assert.Equal(tt, nodeName, addResp.Node.Name)
// }

// r, err = c.UpdateNode(ctx, &pb.UpdateNodeRequest{NodeId: node.String(), Node: &pb.Node{
// State: pb.NodeState_ONBOARDED,
// }})

// handleResponse(tt, err, r)

// // add second node
// add2, err := c.AddNode(ctx, &pb.AddNodeRequest{
// Node: &pb.Node{
// NodeId: nodeToDelete.String(),
// State:  pb.NodeState_UNDEFINED,
// Name:   "nodeToDelete",
// },
// OrgName: orgName,
// Network: networkName,
// })

// handleResponse(tt, err, add2)

// nodeResp, err := c.GetNodes(ctx, &pb.GetNodesRequest{OrgName: orgName})
// if handleResponse(tt, err, nodeResp) {
// nd := slices.FindPointer(nodeResp.Nodes, func(n *pb.Node) bool {
// return n.NodeId == node.String()
// })
// assert.Equal(tt, pb.NodeState_ONBOARDED, nd.State)
// assert.Equal(tt, pb.NodeType_HOME, nd.Type)
// }
// })

// t.Run("DeleteNode", func(tt *testing.T) {
// r, err = c.DeleteNode(ctx, &pb.DeleteNodeRequest{NodeId: nodeToDelete.String()})
// handleResponse(tt, err, r)
// })

// t.Run("GetNodes", func(tt *testing.T) {
// nodesResp, err := c.GetNodes(ctx, &pb.GetNodesRequest{OrgName: orgName})

// handleResponse(t, err, nodesResp)

// // we added 2 node and deleted 1
// if assert.Equal(tt, 1, len(nodesResp.Nodes)) {
// cont := false
// for _, n := range nodesResp.Nodes {
// if node.String() == n.NodeId {
// cont = true

// break
// }
// }
// assert.True(tt, cont, "Can't find added node %s", node.String())
// }
// })
// }

// func Test_Listener(t *testing.T) {
// // Arrange
// org := fmt.Sprintf("network-listener-integration-test-%d", time.Now().Unix())
// nodeID := "UK-INTEGR-HNODE-A1-NETT"
// netName := "default"
// nodeName := "network-listener-integration"

// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// defer cancel()

// rabbit, err := msgbus.NewQPub(tConfig.Queue.Uri, "network-listener-integration-test", os.Getenv("POD_NAME"))
// if err != nil {
// assert.NoErrorf(t, err, "could not create rabbitmq client %+v", err)

// return
// }

// conn, c, err := CreateNetworkClient()
// defer conn.Close()
// defer deleteNetworks(t, c, org, netName)

// if err != nil {
// assert.NoErrorf(t, err, "did not connect: %+v\n", err)

// return
// }

// t.Run("NetworkAddedEvent", func(tt *testing.T) {
// rabbit.Publish(&msgbus.OrgCreatedBody{
// Name:  org,
// Owner: uuid.NewString(),
// }, string(msgbus.OrgCreatedRoutingKey))

// time.Sleep(2 * time.Second)

// nodeResp, err := c.AddNode(ctx, &pb.AddNodeRequest{Node: &pb.Node{
// NodeId: nodeID, State: pb.NodeState_UNDEFINED,
// }, OrgName: org,
// Network: netName})

// // we want to check that the network is there. If adding node fails that means then network is not read
// // we will reuse this node in next test
// if !handleResponse(t, err, nodeResp) {
// assert.FailNow(tt, "Node should not be added")
// }
// })

// t.Run("NodeUpdateEvent", func(tt *testing.T) {
// err := rabbit.Publish(&msgbus.NodeUpdateBody{
// NodeId: nodeID,
// State:  pb.NodeState_name[int32(pb.NodeState_ONBOARDED)],
// Name:   nodeName,
// }, string(msgbus.NodeUpdatedRoutingKey))

// if !assert.NoError(tt, err, "Publish failed") {
// tt.FailNow()
// }

// assert.NoError(tt, err)
// time.Sleep(2 * time.Second)

// nodeResp, err := c.GetNode(ctx, &pb.GetNodeRequest{
// NodeId: nodeID,
// })

// if handleResponse(tt, err, nodeResp) {
// assert.Equal(tt, pb.NodeState_ONBOARDED, nodeResp.GetNode().GetState())
// assert.Equal(tt, nodeName, nodeResp.GetNode().GetName())
// assert.Equal(tt, netName, nodeResp.GetNetwork().GetName())
// }
// })
// }

// func CreateNetworkClient() (*grpc.ClientConn, pb.NetworkServiceClient, error) {
// logrus.Infoln("Connecting to network ", tConfig.ServiceHost)

// context, cancel := context.WithTimeout(context.Background(), time.Second*3)
// defer cancel()

// conn, err := grpc.DialContext(context, tConfig.ServiceHost, grpc.WithInsecure())
// if err != nil {
// return nil, nil, err
// }

// c := pb.NewNetworkServiceClient(conn)

// return conn, c, nil
// }

// func handleResponse(t *testing.T, err error, r interface{}) bool {
// t.Helper()

// fmt.Printf("Response: %v\n", r)

// if err != nil {
// assert.FailNow(t, "Request failed: %v\n", err)

// return false
// }

// return true
// }

// func deleteNetworks(t *testing.T, c pb.NetworkServiceClient, org string, network string) {
// t.Helper()

// logrus.Infoln("Deleting network ", network)

// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// defer cancel()

// _, err := c.Delete(ctx, &pb.DeleteRequest{OrgName: org, Name: network})
// if err != nil {
// assert.FailNowf(t, "Delete network %s from org %s failed: %v\n", network, org, err)
// }
// }
