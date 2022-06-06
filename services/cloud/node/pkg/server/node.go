package server

import (
	"context"
	"errors"
	"github.com/goombaio/namegenerator"
	"github.com/jackc/pgconn"
	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/services/cloud/node/pb/gen"
	"github.com/ukama/ukama/services/cloud/node/pkg"
	"github.com/ukama/ukama/services/cloud/node/pkg/db"
	"github.com/ukama/ukama/services/common/grpc"
	"github.com/ukama/ukama/services/common/msgbus"
	"github.com/ukama/ukama/services/common/sql"
	"github.com/ukama/ukama/services/common/ukama"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type NodeServer struct {
	nodeRepo       db.NodeRepo
	queuePub       msgbus.QPub
	baseRoutingKey msgbus.RoutingKeyBuilder
	nameGenerator  namegenerator.Generator
	pb.UnimplementedNodeServiceServer
}

func NewNodeServer(nodeRepo db.NodeRepo, queuePub msgbus.QPub) *NodeServer {
	seed := time.Now().UTC().UnixNano()
	return &NodeServer{nodeRepo: nodeRepo,
		queuePub:       queuePub,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName),
		nameGenerator:  namegenerator.NewNameGenerator(seed),
	}
}

func (n *NodeServer) AttachNodes(ctx context.Context, req *pb.AttachNodesRequest) (*pb.AttachNodesResponse, error) {
	nodeId, err := ukama.ValidateNodeId(req.GetParentNodeId())
	if err != nil {
		return nil, invalidNodeIdError(req.GetParentNodeId(), err)
	}

	nds := make([]ukama.NodeID, 0)
	for _, n := range req.GetAttachedNodeIds() {
		nd, err := ukama.ValidateNodeId(n)
		if err != nil {
			return nil, invalidNodeIdError(n, err)
		}
		nds = append(nds, nd)
	}
	err = n.nodeRepo.AttachNodes(nodeId, nds)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	n.pubEvent(req, n.baseRoutingKey.SetActionUpdate().SetObject("node").MustBuild())

	return &pb.AttachNodesResponse{}, nil
}

func (n *NodeServer) DetachNode(ctx context.Context, req *pb.DetachNodeRequest) (*pb.DetachNodeResponse, error) {
	nodeId, err := ukama.ValidateNodeId(req.DetachedNodeId)
	if err != nil {
		return nil, invalidNodeIdError(req.DetachedNodeId, err)
	}

	err = n.nodeRepo.DetachNode(nodeId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	n.pubEvent(req, n.baseRoutingKey.SetActionUpdate().SetObject("node").MustBuild())

	return &pb.DetachNodeResponse{}, nil
}

func (n *NodeServer) UpdateNodeState(ctx context.Context, req *pb.UpdateNodeStateRequest) (*pb.UpdateNodeStateResponse, error) {
	logrus.Infof("Updating node state  %v", req.GetNodeId())

	dbState := pbNodeStateToDb(req.State)

	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = n.nodeRepo.Update(nodeId, &dbState, nil)
	if err != nil {
		logrus.Error("error updating the node state, ", err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.UpdateNodeStateResponse{
		NodeId: req.GetNodeId(),
		State:  req.State,
	}
	n.pubEvent(resp, n.baseRoutingKey.SetActionUpdate().SetObject("node").MustBuild())

	return resp, nil
}

func (n *NodeServer) UpdateNode(ctx context.Context, req *pb.UpdateNodeRequest) (*pb.UpdateNodeResponse, error) {
	logrus.Infof("Updating the node  %v", req.GetNodeId())

	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = n.nodeRepo.Update(nodeId, nil, &req.Name)
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
	n.pubEvent(resp, n.baseRoutingKey.SetActionUpdate().SetObject("node").MustBuild())

	return resp, nil
}

func (n *NodeServer) GetNode(ctx context.Context, req *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	logrus.Infof("Get node  %v", req.GetNodeId())

	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	node, err := n.nodeRepo.Get(nodeId)
	if err != nil {
		logrus.Error("error getting the node" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.GetNodeResponse{
		Node: dbNodeToPbNode(node),
	}

	return resp, nil
}
func (n *NodeServer) AddNode(ctx context.Context, req *pb.AddNodeRequest) (*pb.AddNodeResponse, error) {
	logrus.Infof("Adding node  %v", req.Node)

	nId, err := ukama.ValidateNodeId(req.Node.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of node id. Error %s", err.Error())
	}

	if req.Node.Type != pb.NodeType_NODE_TYPE_UNDEFINED {
		return nil, status.Errorf(codes.InvalidArgument, "type is determined from nodeId and can not be set specifically")
	}

	// Generate random node name if it's missing
	if len(req.Node.Name) == 0 {
		req.Node.Name = n.nameGenerator.Generate()
	}

	node := &db.Node{
		NodeID: req.Node.NodeId,
		State:  pbNodeStateToDb(req.Node.State),
		Type:   toDbNodeType(nId.GetNodeType()),
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

	return &pb.AddNodeResponse{
		Node: dbNodeToPbNode(node),
	}, nil
}

func (n *NodeServer) mustEmbedUnimplementedNodeRegistryServiceServer() {
	//TODO implement me
	panic("implement me")
}

func invalidNodeIdError(nodeId string, err error) error {
	return status.Errorf(codes.InvalidArgument, "invalid node id %s. Error %s", nodeId, err.Error())
}

func (n *NodeServer) pubEvent(payload any, key string) {
	go func() {
		err := n.queuePub.Publish(payload, key)
		if err != nil {
			logrus.Errorf("Failed to publish event. Error %s", err.Error())
		}
	}()
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

func (n *NodeServer) processNodeDuplErrors(err error, nodeId string) error {
	var pge *pgconn.PgError
	if errors.As(err, &pge) && pge.Code == sql.PGERROR_CODE_UNIQUE_VIOLATION {
		return status.Errorf(codes.AlreadyExists, "node with node id %s already exist", nodeId)
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
