package server

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/goombaio/namegenerator"
	"github.com/jackc/pgconn"
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
	netpb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	pb "github.com/ukama/ukama/systems/registry/node/pb/gen"
	orgpb "github.com/ukama/ukama/systems/registry/org/pb/gen"
)

type NodeServer struct {
	nodeRepo       db.NodeRepo
	siteRepo       db.SiteRepo
	baseRoutingKey msgbus.RoutingKeyBuilder
	nameGenerator  namegenerator.Generator
	orgService     providers.OrgClientProvider
	networkService providers.NetworkClientProvider
	pushGateway    string
	pb.UnimplementedNodeServiceServer
}

func NewNodeServer(nodeRepo db.NodeRepo, siteRepo db.SiteRepo,
	pushGateway string,
	orgService providers.OrgClientProvider,
	networkService providers.NetworkClientProvider) *NodeServer {
	seed := time.Now().UTC().UnixNano()

	return &NodeServer{
		nodeRepo:       nodeRepo,
		siteRepo:       siteRepo,
		orgService:     orgService,
		networkService: networkService,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName),
		nameGenerator:  namegenerator.NewNameGenerator(seed),
		pushGateway:    pushGateway,
	}
}

func (n *NodeServer) AddNode(ctx context.Context, req *pb.AddNodeRequest) (*pb.AddNodeResponse, error) {
	log.Infof("Adding node  %v", req.NodeId)

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
			"org is deactivated: cannot add node to it")
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

	err = n.nodeRepo.Add(node, nil)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfActiveNodes, pkg.NumberOfInactiveNodes)

	return &pb.AddNodeResponse{Node: dbNodeToPbNode(node)}, nil
}

func (n *NodeServer) GetNode(ctx context.Context, req *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	log.Infof("Get node  %v", req.GetNodeId())

	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	node, err := n.nodeRepo.Get(nodeId)

	if err != nil {
		log.Error("error getting the node" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.GetNodeResponse{Node: dbNodeToPbNode(node)}

	return resp, nil
}

func (n *NodeServer) GetBySite(ctx context.Context, req *pb.GetBySiteRequest) (*pb.GetBySiteResponse, error) {
	log.Infof("Getting all node  on  site %v", req.GetSiteId())

	site, err := uuid.FromString(req.GetSiteId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of site uuid. Error %s", err.Error())
	}

	nodes, err := n.siteRepo.GetNodes(site)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "nodes")
	}

	resp := &pb.GetBySiteResponse{
		SiteId: req.GetSiteId(),
		Nodes:  dbNodesToPbNodes(nodes),
	}

	return resp, nil
}

func (n *NodeServer) GetFreeNodes(ctx context.Context, req *pb.GetFreeNodesRequest) (*pb.GetFreeNodesResponse, error) {
	log.Infof("Get free nodes")

	nodes, err := n.siteRepo.GetFreeNodes()

	if err != nil {
		log.Error("error getting the free node" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.GetFreeNodesResponse{
		Node: dbNodesToPbNodes(nodes),
	}

	return resp, nil
}

func (n *NodeServer) GetAllNodes(ctx context.Context, req *pb.GetAllNodesRequest) (*pb.GetAllNodesResponse, error) {
	log.Infof("Get all nodes.")

	nodes, err := n.nodeRepo.GetAll()

	if err != nil {
		log.Error("error getting all node" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.GetAllNodesResponse{
		Node: dbNodesToPbNodes(nodes),
	}

	return resp, nil
}

func (n *NodeServer) UpdateNodeState(ctx context.Context, req *pb.UpdateNodeStateRequest) (*pb.UpdateNodeResponse, error) {
	log.Infof("Updating node state  %v", req.GetNodeId())

	dbState := db.ParseNodeState(req.State)

	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	nodeUpdates := &db.Node{
		Id:    nodeId.StringLowercase(),
		State: dbState,
	}

	err = n.nodeRepo.Update(nodeUpdates, nil)
	if err != nil {
		log.Error("error updating the node state, ", err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.UpdateNodeResponse{
		Node: &pb.Node{
			Id:    req.GetNodeId(),
			State: req.State,
		},
	}
	und, err := n.nodeRepo.Get(nodeId)
	if err != nil {
		log.Error("error getting the node, ", err.Error())

		return resp, nil
	}

	// publish event and return
	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfActiveNodes, pkg.NumberOfInactiveNodes)

	return &pb.UpdateNodeResponse{Node: dbNodeToPbNode(und)}, nil
}

func (n *NodeServer) UpdateNode(ctx context.Context, req *pb.UpdateNodeRequest) (*pb.UpdateNodeResponse, error) {
	log.Infof("Updating node  %v", req.GetNodeId())

	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	nodeUpdates := &db.Node{
		Id:   nodeId.StringLowercase(),
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

	und, err := n.nodeRepo.Get(nodeId)
	if err != nil {
		log.Error("error getting the node, ", err.Error())

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

func (n *NodeServer) AttachNodes(ctx context.Context, req *pb.AttachNodesRequest) (*pb.AttachNodesResponse, error) {
	log.Infof("Attaching nodes %v to parent node %s", req.GetAttachedNodes(), req.GetNodeId())

	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, invalidNodeIDError(req.GetNodeId(), err)
	}

	nds := req.GetAttachedNodes()

	err = n.nodeRepo.AttachNodes(nodeId, nds)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfActiveNodes, pkg.NumberOfInactiveNodes)

	return &pb.AttachNodesResponse{}, nil
}

func (n *NodeServer) DetachNode(ctx context.Context, req *pb.DetachNodeRequest) (*pb.DetachNodeResponse, error) {
	log.Infof("detaching node  %v", req.GetNodeId())

	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, invalidNodeIDError(req.GetNodeId(), err)
	}

	err = n.nodeRepo.DetachNode(nodeId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfActiveNodes, pkg.NumberOfInactiveNodes)

	return &pb.DetachNodeResponse{}, nil
}

func (n *NodeServer) AddNodeToNetwork(ctx context.Context, req *pb.AddNodeToNetworkRequest) (*pb.AddNodeToNetworkResponse, error) {
	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, invalidNodeIDError(req.GetNodeId(), err)
	}

	net, err := uuid.FromString(req.GetNetworkId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid network id %s. Error %s", req.GetNetworkId(), err.Error())
	}

	// TODO: update RPC handlers for missing site_id (default site for network)
	site, err := uuid.FromString(req.GetSiteId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid site id %s. Error %s", req.GetSiteId(), err.Error())
	}

	svc, err := n.networkService.GetClient()
	if err != nil {
		return nil, err
	}

	remoteSite, err := svc.GetSite(ctx, &netpb.GetSiteRequest{SiteId: site.String()})
	if err != nil {
		return nil, err
	}

	if remoteSite.Site.NetworkId != net.String() {
		return nil, status.Errorf(codes.FailedPrecondition,
			"provided networkId and site's networkId mismatch")
	}

	node := &db.Site{
		NodeId:    nodeId.StringLowercase(),
		SiteId:    site,
		NetworkId: net,
	}

	err = n.siteRepo.AddNode(node, nil)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	return &pb.AddNodeToNetworkResponse{}, nil
}

func (n *NodeServer) RemoveNodeFromNetwork(ctx context.Context, req *pb.ReleaseNodeFromNetworkRequest) (*pb.ReleaseNodeFromNetworkResponse, error) {
	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, invalidNodeIDError(req.GetNodeId(), err)
	}

	err = n.siteRepo.RemoveNode(nodeId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	return &pb.ReleaseNodeFromNetworkResponse{}, nil
}

func invalidNodeIDError(nodeId string, err error) error {
	return status.Errorf(codes.InvalidArgument, "invalid node id %s. Error %s", nodeId, err.Error())
}

func processNodeDuplErrors(err error, nodeId string) error {
	var pge *pgconn.PgError

	if errors.As(err, &pge) && pge.Code == sql.PGERROR_CODE_UNIQUE_VIOLATION {
		return status.Errorf(codes.AlreadyExists, "node with node id %s already exist", nodeId)
	}

	return grpc.SqlErrorToGrpc(err, "node")
}

func (n *NodeServer) pushNodeMeterics(id ukama.NodeID, args ...string) {
	nodesCount, actCount, inactCount, err := n.nodeRepo.GetNodeCount()
	if err != nil {
		log.Errorf("Error while getting node count %s", err.Error())

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
		log.Errorf("Error while pushing node metric to pushgateway %s", err.Error())

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

// func AddNodeToOrg(repo db.NodeRepo, node *db.Node) error {

// // Generate random node name if it's missing

// // adding node to DB and bootstrap in transaction
// // Rollback trans if bootstrap fails to add a node
// err := repo.Add(node, nil)

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
