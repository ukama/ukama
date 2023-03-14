package server

import (
	"context"
	"errors"
	"time"

	"github.com/goombaio/namegenerator"
	"github.com/jackc/pgconn"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	pmetric "github.com/ukama/ukama/systems/common/pushgatewayMetrics"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/ukama"
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
	pb.UnimplementedNodeServiceServer
	org string
	pushMetricHost string
}

func NewNodeServer(nodeRepo db.NodeRepo,org string,pushMetricHost string ) *NodeServer {
	seed := time.Now().UTC().UnixNano()

	return &NodeServer{nodeRepo: nodeRepo,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName),
		nameGenerator:  namegenerator.NewNameGenerator(seed),
		org:org,
		pushMetricHost:pushMetricHost,
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

	err = n.nodeRepo.AttachNodes(nodeID, nds)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	// publish event and return

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

	// publish event and return

	return &pb.DetachNodeResponse{}, nil
}

func (n *NodeServer) UpdateNodeState(ctx context.Context, req *pb.UpdateNodeStateRequest) (*pb.UpdateNodeStateResponse, error) {
	logrus.Infof("Updating node state  %v", req.GetNodeId())

	dbState := pbNodeStateToDb(req.State)

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
	_ ,activeNodeCount,inactiveNodeCount ,err := n.nodeRepo.GetNodeCount()
	if err != nil {
		logrus.Errorf("Error while getting node count %s", err.Error())
	}
	if (dbState ==db.Onboarded){
		err = pmetric.CollectAndPushSimMetrics(n.pushMetricHost, pkg.NodeMetric, pkg.NumberOfActiveNodes, float64(activeNodeCount), map[string]string{"network": "", "org": n.org},pkg.SystemName)
		if err != nil {
			logrus.Errorf("Error while pushing subscriberCount metric to pushgaway %s", err.Error())
		}
	}else if (dbState ==db.Pending){
		err = pmetric.CollectAndPushSimMetrics(n.pushMetricHost, pkg.NodeMetric, pkg.NumberOfInactiveNodes, float64(inactiveNodeCount), map[string]string{"network": "", "org": n.org},pkg.SystemName)
		if err != nil {
			logrus.Errorf("Error while pushing subscriberCount metric to pushgaway %s", err.Error())
		}
	
	}
	
	// publish event and return

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
		State:  pbNodeStateToDb(req.Node.State),
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
	nodesCount ,_,_ ,err := n.nodeRepo.GetNodeCount()
	if err != nil {
		logrus.Errorf("Error while getting node count %s", err.Error())
	}
	err = pmetric.CollectAndPushSimMetrics(n.pushMetricHost, pkg.NodeMetric, pkg.NumberOfNodes, float64(nodesCount), map[string]string{"network": "", "org": n.org},pkg.SystemName)
	if err != nil {
		logrus.Errorf("Error while pushing subscriberCount metric to pushgaway %s", err.Error())
	}

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
	nodesCount ,_,_ ,err := n.nodeRepo.GetNodeCount()
	if err != nil {
		logrus.Errorf("Error while getting node count %s", err.Error())
	}
	err = pmetric.CollectAndPushSimMetrics(n.pushMetricHost, pkg.NodeMetric, pkg.NumberOfNodes, float64(nodesCount), map[string]string{"network": "", "org": n.org},pkg.SystemName)
	if err != nil {
		logrus.Errorf("Error while pushing subscriberCount metric to pushgaway %s", err.Error())
	}

	return &pb.DeleteResponse{NodeId: req.GetNodeId()}, nil
}

func invalidNodeIDError(nodeID string, err error) error {
	return status.Errorf(codes.InvalidArgument, "invalid node id %s. Error %s", nodeID, err.Error())
}

func pbNodeStateToDb(state pb.NodeState) db.NodeState {
	var dbState db.NodeState

	switch state {
	case pb.NodeState_ONBOARDED:
		dbState = db.Onboarded
	case pb.NodeState_PENDING:
		dbState = db.Pending
	default:
		dbState = db.Undefined
	}

	return dbState
}

func (n *NodeServer) processNodeDuplErrors(err error, nodeID string) error {
	var pge *pgconn.PgError

	if errors.As(err, &pge) && pge.Code == sql.PGERROR_CODE_UNIQUE_VIOLATION {
		return status.Errorf(codes.AlreadyExists, "node with node id %s already exist", nodeID)
	}

	return grpc.SqlErrorToGrpc(err, "node")
}

func dbNodeToPbNode(dbn *db.Node) *pb.Node {
	n := &pb.Node{
		NodeId: dbn.NodeID,
		State:  dbNodeStateToPb(dbn.State),
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

func dbNodeStateToPb(state db.NodeState) pb.NodeState {
	var pbState pb.NodeState

	switch state {
	case db.Onboarded:
		pbState = pb.NodeState_ONBOARDED
	case db.Pending:
		pbState = pb.NodeState_PENDING
	default:
		pbState = pb.NodeState_UNDEFINED
	}

	return pbState
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
