package main

import (
	"context"
	"fmt"
	"os"
	"time"

	pb "ukamaX/registry/pb/generated"

	uuid2 "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukamaX/common/ukama"
	"google.golang.org/grpc"
)

func main() {
	args := os.Args[1:]
	host := "localhost:9090"
	if len(args) == 2 {
		host = os.Args[1:][0]
	}

	conn, err := grpc.Dial(host, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewRegistryServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	orgName := fmt.Sprintf("self-test-org-%d", time.Now().Unix())
	ownerId := uuid2.NewV1()
	node := ukama.NewVirtualNodeId("HomeNode")

	var r interface{}
	r, err = c.AddOrg(ctx, &pb.AddOrgRequest{Name: orgName, Owner: ownerId.String()})
	handleResponse(err, r)

	r, err = c.GetOrg(ctx, &pb.GetOrgRequest{Name: orgName})
	handleResponse(err, r)

	r, err = c.AddNode(ctx, &pb.AddNodeRequest{
		Node: &pb.Node{
			NodeId: node.String(),
			State:  pb.NodeState_UNDEFINED,
		},
		OrgName: orgName,
	})
	handleResponse(err, r)

	r, err = c.UpdateNode(ctx, &pb.UpdateNodeRequest{NodeId: node.String(), State: pb.NodeState_ONBOARDED})
	handleResponse(err, r)

	nodeResp, err := c.GetNode(ctx, &pb.GetNodeRequest{NodeId: node.String()})
	handleResponse(err, nodeResp)
	if nodeResp.Node.State != pb.NodeState_ONBOARDED {
		logrus.Fatalf("Node state was not updated")
	}
}

func handleResponse(err error, r interface{}) {
	if err != nil {
		logrus.Fatalf("Request failed: %v\n", err)
	}
	fmt.Printf("Response: %v\n", r)
}
