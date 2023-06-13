package server

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/goombaio/namegenerator"
	"github.com/jackc/pgconn"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/node/pkg"
	"github.com/ukama/ukama/systems/registry/node/pkg/db"
	"github.com/ukama/ukama/systems/registry/node/pkg/providers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	log "github.com/sirupsen/logrus"
	metric "github.com/ukama/ukama/systems/common/metrics"
	pb "github.com/ukama/ukama/systems/registry/node/pb/gen"
	orgpb "github.com/ukama/ukama/systems/registry/org/pb/gen"
)

type NodeServer struct {
	nodeRepo       db.NodeRepo
	siteRepo       db.SiteRepo
	baseRoutingKey msgbus.RoutingKeyBuilder
	nameGenerator  namegenerator.Generator
	orgService     providers.OrgClientProvider
	pushGateway    string
	pb.UnimplementedNodeServiceServer
}

func NewNodeServer(nodeRepo db.NodeRepo, siteRepo db.SiteRepo,
	pushGateway string, orgService providers.OrgClientProvider) *NodeServer {
	seed := time.Now().UTC().UnixNano()

	return &NodeServer{
		nodeRepo:       nodeRepo,
		siteRepo:       siteRepo,
		orgService:     orgService,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName),
		nameGenerator:  namegenerator.NewNameGenerator(seed),
		pushGateway:    pushGateway,
	}
}

func (n *NodeServer) AddNode(ctx context.Context, req *pb.AddNodeRequest) (*pb.AddNodeResponse, error) {
	logrus.Infof("Adding node  %v", req.NodeId)

	nID, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id. Error %s", err.Error())
	}

	strState := strings.ToLower(req.GetState())
	nodeState := db.ParseNodeState(strState)
	if req.GetState() != "" && nodeState == db.Undefined {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid node type. Error: node type %q not supported", req.GetState())
	}

	orgId, err := uuid.FromString(req.GetOrgId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of org uuid. Error %s", err.Error())
	}

	svc, err := n.orgService.GetClient()
	if err != nil {
		return nil, err
	}

	remoteOrg, err := svc.Get(ctx, &orgpb.GetRequest{Id: orgId.String()})
	if err != nil {
		return nil, err
	}

	// What should we do if the remote org exists but is deactivated?
	// For now we simply abort.
	if remoteOrg.Org.IsDeactivated {
		return nil, status.Errorf(codes.FailedPrecondition,
			"org is deactivated: cannot add network to it")
	}

	if len(req.Name) == 0 {
		req.Name = n.nameGenerator.Generate()
	}

	node := &db.Node{
		Id:    req.NodeId,
		OrgId: orgId,
		State: nodeState,
		Type:  nID.GetNodeType(),
		Name:  req.Name,
	}

	log.Infof("Adding node %s", node.Name)
	err = n.nodeRepo.Add(node, nil)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfActiveNodes, pkg.NumberOfInactiveNodes)

	return &pb.AddNodeResponse{Node: dbNodeToPbNode(node)}, nil
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

func (n *NodeServer) GetFreeNodes(ctx context.Context, req *pb.GetFreeNodesRequest) (*pb.GetFreeNodesResponse, error) {
	logrus.Infof("GetFreeNodes")

	nodes, err := n.siteRepo.GetFreeNodes()

	if err != nil {
		logrus.Error("error getting the free node" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.GetFreeNodesResponse{
		Node: dbNodesToPbNodes(nodes),
	}

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

func (n *NodeServer) UpdateNodeState(ctx context.Context, req *pb.UpdateNodeStateRequest) (*pb.UpdateNodeResponse, error) {
	logrus.Infof("Updating node state  %v", req.GetNodeId())

	dbState := db.ParseNodeState(req.State)

	nodeID, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	nodeUpdates := &db.Node{
		Id:    nodeID.StringLowercase(),
		State: dbState,
	}

	err = n.nodeRepo.Update(nodeUpdates, nil)
	if err != nil {
		logrus.Error("error updating the node state, ", err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.UpdateNodeResponse{
		Node: &pb.Node{
			Id:    req.GetNodeId(),
			State: req.State,
		},
	}
	und, err := n.nodeRepo.Get(nodeID)
	if err != nil {
		logrus.Error("error getting the node, ", err.Error())

		return resp, nil
	}

	// publish event and return
	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfActiveNodes, pkg.NumberOfInactiveNodes)

	return &pb.UpdateNodeResponse{Node: dbNodeToPbNode(und)}, nil
}

func (n *NodeServer) UpdateNode(ctx context.Context, req *pb.UpdateNodeRequest) (*pb.UpdateNodeResponse, error) {
	logrus.Infof("Updating the node  %v", req.GetNodeId())

	nodeID, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	nodeUpdates := &db.Node{
		Id:   nodeID.StringLowercase(),
		Name: req.Name,
	}

	err = n.nodeRepo.Update(nodeUpdates, nil)
	if err != nil {
		duplErr := processNodeDuplErrors(err, req.NodeId)
		if duplErr != nil {
			return nil, duplErr
		}

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.UpdateNodeResponse{
		Node: &pb.Node{
			Id:   req.NodeId,
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

func (n *NodeServer) DeleteNode(ctx context.Context, req *pb.DeleteNodeRequest) (*pb.DeleteNodeResponse, error) {
	log.Infof("Deleting node  %v", req.GetNodeId())

	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = n.nodeRepo.Delete(nodeId, nil)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfActiveNodes, pkg.NumberOfInactiveNodes)

	return &pb.DeleteNodeResponse{}, nil
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
			err = metric.CollectAndPushSimMetrics(n.pushGateway, pkg.NodeMetric,
				pkg.NumberOfNodes, float64(nodesCount), nil, pkg.SystemName+"-"+pkg.ServiceName)
		case pkg.NumberOfActiveNodes:
			err = metric.CollectAndPushSimMetrics(n.pushGateway, pkg.NodeMetric,
				pkg.NumberOfActiveNodes, float64(actCount), nil, pkg.SystemName+"-"+pkg.ServiceName)
		case pkg.NumberOfInactiveNodes:
			err = metric.CollectAndPushSimMetrics(n.pushGateway, pkg.NodeMetric,
				pkg.NumberOfInactiveNodes, float64(inactCount), nil, pkg.SystemName+"-"+pkg.ServiceName)
		}
	}

	if err != nil {
		logrus.Errorf("Error while pushing node metric to pushgateway %s", err.Error())
	}
}

func (n *NodeServer) PushMetrics() {
	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfActiveNodes, pkg.NumberOfInactiveNodes)
}

func dbNodesToPbNodes(nodes []db.Node) []*pb.Node {
	pbNodes := []*pb.Node{}
	for _, n := range nodes {
		pbNodes = append(pbNodes, dbNodeToPbNode(&n))
	}
	return pbNodes
}

func dbNodeToPbNode(dbn *db.Node) *pb.Node {
	n := &pb.Node{
		Id:        dbn.Id,
		State:     dbn.State.String(),
		Type:      dbn.Type,
		Name:      dbn.Name,
		OrgId:     dbn.OrgId.String(),
		CreatedAt: timestamppb.New(dbn.CreatedAt),
	}

	if len(dbn.Attached) > 0 {
		n.Attached = make([]*pb.Node, 0)
	}

	for _, nd := range dbn.Attached {
		n.Attached = append(n.Attached, dbNodeToPbNode(nd))
	}

	return n
}

// func (n *NodeServer) AddNodeToNetwork(ctx context.Context, req *pb.AddNodeToNetworkRequest) (*pb.AddNodeToNetworkResponse, error) {
// nodeID, err := ukama.ValidateNodeId(req.GetNode())
// if err != nil {
// return nil, invalidNodeIDError(req.GetNode(), err)
// }

// net, err := uuid.FromString(req.GetNetwork())
// if err != nil {
// return nil, status.Errorf(codes.InvalidArgument, "invalid network id %s. Error %s", req.Network, err.Error())
// }

// err = n.nodeRepo.AddNodeToNetwork(nodeID, net)
// if err != nil {
// return nil, grpc.SqlErrorToGrpc(err, "node")
// }

// return &pb.AddNodeToNetworkResponse{}, nil
// }

// func (n *NodeServer) RemoveNodeFromNetwork(ctx context.Context, req *pb.ReleaseNodeFromNetworkRequest) (*pb.ReleaseNodeFromNetworkResponse, error) {
// nodeID, err := ukama.ValidateNodeId(req.GetNode())
// if err != nil {
// return nil, invalidNodeIDError(req.GetNode(), err)
// }

// err = n.nodeRepo.RemoveNodeFromNetwork(nodeID)
// if err != nil {
// return nil, grpc.SqlErrorToGrpc(err, "node")
// }

// return &pb.ReleaseNodeFromNetworkResponse{}, nil
// }

// func (n *NodeServer) AttachNodes(ctx context.Context, req *pb.AttachNodesRequest) (*pb.AttachNodesResponse, error) {
// nodeID, err := ukama.ValidateNodeId(req.GetParentNode())
// if err != nil {
// return nil, invalidNodeIDError(req.GetParentNode(), err)
// }

// nds := make([]ukama.NodeID, 0)

// for _, n := range req.GetAttachedNodes() {
// nd, err := ukama.ValidateNodeId(n)
// if err != nil {
// return nil, invalidNodeIDError(n, err)
// }

// nds = append(nds, nd)
// }

// err = n.nodeRepo.AttachNodes(nodeID, nds)
// if err != nil {
// return nil, grpc.SqlErrorToGrpc(err, "node")
// }

// n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfActiveNodes, pkg.NumberOfInactiveNodes)

// return &pb.AttachNodesResponse{}, nil
// }

// func AddNodeToOrg(repo db.NodeRepo, node *db.Node) error {

// // Generate random node name if it's missing

// // adding node to DB and bootstrap in transaction
// // Rollback trans if bootstrap fails to add a node
// err := repo.Add(node)

// if err != nil {
// duplErr := processNodeDuplErrors(err, node.Id)
// if duplErr != nil {
// return duplErr
// }

// logrus.Error("Error adding the node. " + err.Error())

// return status.Errorf(codes.Internal, "error adding the node")
// }

// return nil
// }
