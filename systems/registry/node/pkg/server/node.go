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

func (n *NodeServer) AttachNodes(ctx context.Context, req *pb.AttachNodesRequest) (*pb.AttachNodesResponse, error) {
	nodeID, err := ukama.ValidateNodeId(req.GetParentNodeId())
	if err != nil {
		return nil, invalidNodeIDError(req.GetParentNodeId(), err)
	}

	nds := make([]ukama.NodeID, 0)

	for _, n := range req.GetAttachedNodeIds() {
		nd, err := ukama.ValidateNodeId(n)
		if err != nil {
			return nil, invalidNodeIDError(n, err)
		}

		nds = append(nds, nd)
	}

	net, err := uuid.FromString(req.GetNetwork())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid network id %s. Error %s", req.Network, err.Error())
	}

	err = n.nodeRepo.AttachNodes(nodeID, nds, net)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfActiveNodes, pkg.NumberOfInactiveNodes)

	return &pb.AttachNodesResponse{}, nil
}

func (n *NodeServer) DetachNode(ctx context.Context, req *pb.DetachNodeRequest) (*pb.DetachNodeResponse, error) {
	nodeID, err := ukama.ValidateNodeId(req.DetachedNodeId)
	if err != nil {
		return nil, invalidNodeIDError(req.DetachedNodeId, err)
	}

	err = n.nodeRepo.DetachNode(nodeID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfActiveNodes, pkg.NumberOfInactiveNodes)

	return &pb.DetachNodeResponse{}, nil
}

func (n *NodeServer) UpdateNodeState(ctx context.Context, req *pb.UpdateNodeStateRequest) (*pb.UpdateNodeStateResponse, error) {
	logrus.Infof("Updating node state  %v", req.GetNodeId())

	dbState := db.ParseNodeState(req.State)

	nodeID, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = n.nodeRepo.Update(nodeID, &dbState, nil)
	if err != nil {
		logrus.Error("error updating the node state, ", err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.UpdateNodeStateResponse{
		NodeId: req.GetNodeId(),
		State:  req.State,
	}

	// publish event and return
	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfActiveNodes, pkg.NumberOfInactiveNodes)

	return resp, nil
}

func (n *NodeServer) UpdateNode(ctx context.Context, req *pb.UpdateNodeRequest) (*pb.UpdateNodeResponse, error) {
	logrus.Infof("Updating the node  %v", req.GetNodeId())

	nodeID, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = n.nodeRepo.Update(nodeID, nil, &req.Name)
	if err != nil {
		duplErr := n.processNodeDuplErrors(err, req.NodeId)
		if duplErr != nil {
			return nil, duplErr
		}

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.UpdateNodeResponse{
		Node: &pb.Node{
			NodeId: req.NodeId,
			Name:   req.Name,
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
	logrus.Infof("Get node  %v", req.GetNodeId())

	nodeID, err := ukama.ValidateNodeId(req.GetNodeId())
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

func (n *NodeServer) GetAllNode(ctx context.Context, req *pb.GetAllNodeRequest) (*pb.GetAllNodeResponse, error) {
	logrus.Infof("GetAll Nodes.")

	nodes, err := n.nodeRepo.GetAll()

	if err != nil {
		logrus.Error("error getting all node" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.GetAllNodeResponse{
		Node: dbNodesToPbNodes(nodes),
	}

	return resp, nil
}

func (n *NodeServer) GetFreeNodes(ctx context.Context, req *pb.GetFreeNodeRequest) (*pb.GetFreeNodeResponse, error) {
	logrus.Infof("GetFreeNodes")

	nodes, err := n.nodeRepo.GetFreeNodes()

	if err != nil {
		logrus.Error("error getting the free node" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.GetFreeNodeResponse{
		Node: dbNodesToPbNodes(nodes),
	}

	return resp, nil
}

func (n *NodeServer) AddNode(ctx context.Context, req *pb.AddNodeRequest) (*pb.AddNodeResponse, error) {
	logrus.Infof("Adding node  %v", req.Node)

	nID, err := ukama.ValidateNodeId(req.Node.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id. Error %s", err.Error())
	}

	if req.Node.Type != pb.NodeType_NODE_TYPE_UNDEFINED {
		return nil, status.Errorf(codes.InvalidArgument,
			"type is determined from nodeId and can not be set specifically")
	}

	// Generate random node name if it's missing
	if len(req.Node.Name) == 0 {
		req.Node.Name = n.nameGenerator.Generate()
	}

	node := &db.Node{
		NodeID: req.Node.NodeId,
		State:  db.ParseNodeState(req.Node.State),
		Type:   toDbNodeType(nID.GetNodeType()),
		Name:   req.Node.Name,
	}

	// adding node to DB and bootstrap in transaction
	// Rollback trans if bootstrap fails to add a node
	err = n.nodeRepo.Add(node)

	if err != nil {
		duplErr := n.processNodeDuplErrors(err, node.NodeID)
		if duplErr != nil {
			return nil, duplErr
		}

		logrus.Error("Error adding the node. " + err.Error())

		return nil, status.Errorf(codes.Internal, "error adding the node")
	}

	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfActiveNodes, pkg.NumberOfInactiveNodes)

	return &pb.AddNodeResponse{Node: dbNodeToPbNode(node)}, nil
}

func (n *NodeServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	nID, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, invalidNodeIDError(req.GetNodeId(), err)
	}

	err = n.nodeRepo.Delete(nID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfActiveNodes, pkg.NumberOfInactiveNodes)

	return &pb.DeleteResponse{NodeId: req.GetNodeId()}, nil
}

func invalidNodeIDError(nodeID string, err error) error {
	return status.Errorf(codes.InvalidArgument, "invalid node id %s. Error %s", nodeID, err.Error())
}

func (n *NodeServer) processNodeDuplErrors(err error, nodeID string) error {
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
	n := &pb.Node{
		NodeId: dbn.NodeID,
		State:  dbn.State.String(),
		Type:   dbNodeTypeToPb(dbn.Type),
		Name:   dbn.Name,
	}

	if len(dbn.Attached) > 0 {
		n.Attached = make([]*pb.Node, 0)
	}

	for _, nd := range dbn.Attached {
		n.Attached = append(n.Attached, dbNodeToPbNode(nd))
	}

	return n
}

func dbNodeTypeToPb(nodeType db.NodeType) pb.NodeType {
	var pbNodeType pb.NodeType

	switch nodeType {
	case db.NodeTypeAmplifier:
		pbNodeType = pb.NodeType_AMPLIFIER
	case db.NodeTypeTower:
		pbNodeType = pb.NodeType_TOWER
	case db.NodeTypeHome:
		pbNodeType = pb.NodeType_HOME
	default:
		pbNodeType = pb.NodeType_NODE_TYPE_UNDEFINED
	}

	return pbNodeType
}

func toDbNodeType(nodeType string) db.NodeType {
	switch nodeType {
	case ukama.NODE_ID_TYPE_AMPNODE:
		return db.NodeTypeAmplifier
	case ukama.NODE_ID_TYPE_TOWERNODE:
		return db.NodeTypeTower
	case ukama.NODE_ID_TYPE_HOMENODE:
		return db.NodeTypeHome
	default:
		return db.NodeTypeUnknown
	}
}
