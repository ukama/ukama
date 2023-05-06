package server

import (
	"context"
	"errors"
	"time"

	"github.com/goombaio/namegenerator"
	"github.com/jackc/pgconn"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	metric "github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/registry/node/pb/gen"
	"github.com/ukama/ukama/systems/registry/node/pkg"
	"github.com/ukama/ukama/systems/registry/node/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NodeServer struct {
	nodeRepo       db.NodeRepo
	baseRoutingKey msgbus.RoutingKeyBuilder
	nameGenerator  namegenerator.Generator
	pushGateway    string
	pb.UnimplementedNodeServiceServer
}

func NewNodeServer(nodeRepo db.NodeRepo, pushGateway string) *NodeServer {
	seed := time.Now().UTC().UnixNano()

	return &NodeServer{nodeRepo: nodeRepo,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName),
		nameGenerator:  namegenerator.NewNameGenerator(seed),
		pushGateway:    pushGateway,
	}
}

func (n *NodeServer) AddNodeToNetwork(ctx context.Context, req *pb.AddNodeToNetworkRequest) (*pb.AddNodeToNetworkResponse, error) {
	nodeID, err := ukama.ValidateNodeId(req.GetNode())
	if err != nil {
		return nil, invalidNodeIDError(req.GetNode(), err)
	}

	net, err := uuid.FromString(req.GetNetwork())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid network id %s. Error %s", req.Network, err.Error())
	}

	err = n.nodeRepo.AddNodeToNetwork(nodeID, net)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	return &pb.AddNodeToNetworkResponse{}, nil
}

func (n *NodeServer) RemoveNodeFromNetwork(ctx context.Context, req *pb.ReleaseNodeFromNetworkRequest) (*pb.ReleaseNodeFromNetworkResponse, error) {
	nodeID, err := ukama.ValidateNodeId(req.GetNode())
	if err != nil {
		return nil, invalidNodeIDError(req.GetNode(), err)
	}

	err = n.nodeRepo.RemoveNodeFromNetwork(nodeID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	return &pb.ReleaseNodeFromNetworkResponse{}, nil
}

func (n *NodeServer) AttachNodes(ctx context.Context, req *pb.AttachNodesRequest) (*pb.AttachNodesResponse, error) {
	nodeID, err := ukama.ValidateNodeId(req.GetParentNode())
	if err != nil {
		return nil, invalidNodeIDError(req.GetParentNode(), err)
	}

	nds := make([]ukama.NodeID, 0)

	for _, n := range req.GetAttachedNodes() {
		nd, err := ukama.ValidateNodeId(n)
		if err != nil {
			return nil, invalidNodeIDError(n, err)
		}

		nds = append(nds, nd)
	}

	err = n.nodeRepo.AttachNodes(nodeID, nds)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfActiveNodes, pkg.NumberOfInactiveNodes)

	return &pb.AttachNodesResponse{}, nil
}

func (n *NodeServer) DetachNode(ctx context.Context, req *pb.DetachNodeRequest) (*pb.DetachNodeResponse, error) {
	nodeID, err := ukama.ValidateNodeId(req.Node)
	if err != nil {
		return nil, invalidNodeIDError(req.Node, err)
	}

	err = n.nodeRepo.DetachNode(nodeID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfActiveNodes, pkg.NumberOfInactiveNodes)

	return &pb.DetachNodeResponse{}, nil
}

func (n *NodeServer) UpdateNodeState(ctx context.Context, req *pb.UpdateNodeStateRequest) (*pb.UpdateNodeStateResponse, error) {
	logrus.Infof("Updating node state  %v", req.GetNode())

	dbState := db.ParseNodeState(req.State)

	nodeID, err := ukama.ValidateNodeId(req.GetNode())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = n.nodeRepo.Update(nodeID, &dbState, nil)
	if err != nil {
		logrus.Error("error updating the node state, ", err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.UpdateNodeStateResponse{
		Node:  req.GetNode(),
		State: req.State,
	}

	// publish event and return
	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfActiveNodes, pkg.NumberOfInactiveNodes)

	return resp, nil
}

func (n *NodeServer) UpdateNode(ctx context.Context, req *pb.UpdateNodeRequest) (*pb.UpdateNodeResponse, error) {
	logrus.Infof("Updating the node  %v", req.GetNode())

	nodeID, err := ukama.ValidateNodeId(req.GetNode())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = n.nodeRepo.Update(nodeID, nil, &req.Name)
	if err != nil {
		duplErr := processNodeDuplErrors(err, req.Node)
		if duplErr != nil {
			return nil, duplErr
		}

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.UpdateNodeResponse{
		Node: &pb.Node{
			Node: req.Node,
			Name: req.Name,
		},
	}

	und, err := n.nodeRepo.Get(nodeID)
	if err != nil {
		logrus.Error("error getting the node, ", err.Error())

		return resp, nil
	}

	// publish event and return

	return &pb.UpdateNodeResponse{Node: dbNodeToPbNode(und)}, nil
}

func (n *NodeServer) GetNode(ctx context.Context, req *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	logrus.Infof("Get node  %v", req.GetNode())

	nodeID, err := ukama.ValidateNodeId(req.GetNode())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	node, err := n.nodeRepo.Get(nodeID)

	if err != nil {
		logrus.Error("error getting the node" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.GetNodeResponse{Node: dbNodeToPbNode(node)}

	return resp, nil
}

func (n *NodeServer) GetAllNodes(ctx context.Context, req *pb.GetAllNodesRequest) (*pb.GetAllNodesResponse, error) {
	logrus.Infof("GetAll Nodes.")

	nodes, err := n.nodeRepo.GetAll()

	if err != nil {
		logrus.Error("error getting all node" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.GetAllNodesResponse{
		Node: dbNodesToPbNodes(nodes),
	}

	return resp, nil
}

func (n *NodeServer) GetFreeNodes(ctx context.Context, req *pb.GetFreeNodesRequest) (*pb.GetFreeNodesResponse, error) {
	logrus.Infof("GetFreeNodes")

	nodes, err := n.nodeRepo.GetFreeNodes()

	if err != nil {
		logrus.Error("error getting the free node" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.GetFreeNodesResponse{
		Node: dbNodesToPbNodes(nodes),
	}

	return resp, nil
}

func (n *NodeServer) AddNode(ctx context.Context, req *pb.AddNodeRequest) (*pb.AddNodeResponse, error) {
	logrus.Infof("Adding node  %v", req.Node)

	nID, err := ukama.ValidateNodeId(req.Node.Node)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id. Error %s", err.Error())
	}

	if len(req.Node.Name) == 0 {
		req.Node.Name = n.nameGenerator.Generate()
	}

	node := &db.Node{
		NodeID: req.Node.Node,
		State:  db.ParseNodeState(req.Node.State),
		Type:   nID.GetNodeType(),
		Name:   req.Node.Name,
	}

	if node.State == db.Undefined {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node state %s", req.Node.State)
	}
	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfActiveNodes, pkg.NumberOfInactiveNodes)

	err = AddNodeToOrg(n.nodeRepo, node)
	if err != nil {
		return nil, err
	}

	return &pb.AddNodeResponse{Node: dbNodeToPbNode(node)}, nil

}

func AddNodeToOrg(repo db.NodeRepo, node *db.Node) error {

	// Generate random node name if it's missing

	// adding node to DB and bootstrap in transaction
	// Rollback trans if bootstrap fails to add a node
	err := repo.Add(node)

	if err != nil {
		duplErr := processNodeDuplErrors(err, node.NodeID)
		if duplErr != nil {
			return duplErr
		}

		logrus.Error("Error adding the node. " + err.Error())

		return status.Errorf(codes.Internal, "error adding the node")
	}

	return nil
}

func (n *NodeServer) Delete(ctx context.Context, req *pb.DeleteNodeRequest) (*pb.DeleteNodeResponse, error) {
	nID, err := ukama.ValidateNodeId(req.GetNode())
	if err != nil {
		return nil, invalidNodeIDError(req.GetNode(), err)
	}

	err = n.nodeRepo.Delete(nID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfActiveNodes, pkg.NumberOfInactiveNodes)

	return &pb.DeleteNodeResponse{Node: req.GetNode()}, nil
}

func invalidNodeIDError(nodeID string, err error) error {
	return status.Errorf(codes.InvalidArgument, "invalid node id %s. Error %s", nodeID, err.Error())
}

func processNodeDuplErrors(err error, nodeID string) error {
	var pge *pgconn.PgError

	if errors.As(err, &pge) && pge.Code == sql.PGERROR_CODE_UNIQUE_VIOLATION {
		return status.Errorf(codes.AlreadyExists, "node with node id %s already exist", nodeID)
	}

	return grpc.SqlErrorToGrpc(err, "node")
}

func (n *NodeServer) pushNodeMeterics(id ukama.NodeID, args ...string) {
	nodesCount, actCount, inactCount, err := n.nodeRepo.GetNodeCount()
	if err != nil {
		logrus.Errorf("Error while getting node count %s", err.Error())
		return
	}

	for _, arg := range args {
		switch arg {
		case pkg.NumberOfNodes:
			err = metric.CollectAndPushSimMetrics(n.pushGateway, pkg.NodeMetric, pkg.NumberOfNodes, float64(nodesCount), nil, pkg.SystemName+"-"+pkg.ServiceName)
		case pkg.NumberOfActiveNodes:
			err = metric.CollectAndPushSimMetrics(n.pushGateway, pkg.NodeMetric, pkg.NumberOfActiveNodes, float64(actCount), nil, pkg.SystemName+"-"+pkg.ServiceName)
		case pkg.NumberOfInactiveNodes:
			err = metric.CollectAndPushSimMetrics(n.pushGateway, pkg.NodeMetric, pkg.NumberOfInactiveNodes, float64(inactCount), nil, pkg.SystemName+"-"+pkg.ServiceName)
		}
	}

	if err != nil {
		logrus.Errorf("Error while pushing node metric to pushgateway %s", err.Error())
	}
}

func (n *NodeServer) PushMetrics() {
	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfActiveNodes, pkg.NumberOfInactiveNodes)
}

func dbNodesToPbNodes(nodes *[]db.Node) []*pb.Node {
	pbNodes := []*pb.Node{}
	for _, n := range *nodes {
		pbNodes = append(pbNodes, dbNodeToPbNode(&n))
	}
	return pbNodes
}

func dbNodeToPbNode(dbn *db.Node) *pb.Node {
	var net string
	if dbn.Network.Valid {
		net = dbn.Network.UUID.String()
	}

	n := &pb.Node{
		Node:      dbn.NodeID,
		State:     dbn.State.String(),
		Type:      dbn.Type,
		Name:      dbn.Name,
		Network:   net,
		Allocated: dbn.Allocation,
	}

	if len(dbn.Attached) > 0 {
		n.Attached = make([]*pb.Node, 0)
	}

	for _, nd := range dbn.Attached {
		n.Attached = append(n.Attached, dbNodeToPbNode(nd))
	}

	return n
}
