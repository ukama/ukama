package server

import (
	"context"
	"errors"
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
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	netpb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	pb "github.com/ukama/ukama/systems/registry/node/pb/gen"
)

type NodeServer struct {
	org            uuid.UUID
	nodeRepo       db.NodeRepo
	siteRepo       db.SiteRepo
	nodeStatusRepo db.NodeStatusRepo
	nameGenerator  namegenerator.Generator
	orgService     providers.OrgClientProvider
	networkService providers.NetworkClientProvider
	pushGateway    string
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedNodeServiceServer
}

func NewNodeServer(nodeRepo db.NodeRepo, siteRepo db.SiteRepo, nodeStatusRepo db.NodeStatusRepo,
	pushGateway string, msgBus mb.MsgBusServiceClient,
	orgService providers.OrgClientProvider,
	networkService providers.NetworkClientProvider,
	org uuid.UUID) *NodeServer {
	seed := time.Now().UTC().UnixNano()
	return &NodeServer{
		org:            org,
		nodeRepo:       nodeRepo,
		nodeStatusRepo: nodeStatusRepo,
		siteRepo:       siteRepo,
		orgService:     orgService,
		networkService: networkService,
		nameGenerator:  namegenerator.NewNameGenerator(seed),
		pushGateway:    pushGateway,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName),
	}
}

func (n *NodeServer) AddNode(ctx context.Context, req *pb.AddNodeRequest) (*pb.AddNodeResponse, error) {
	log.Infof("Adding node  %v", req.NodeId)

	nID, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id. Error %s", err.Error())
	}

	if len(req.Name) == 0 {
		req.Name = n.nameGenerator.Generate()
	}

	node := &db.Node{
		Id:    req.NodeId,
		OrgId: n.org,
		Status: db.NodeStatus{
			NodeId: req.NodeId,
			Conn:   db.Unknown,
			State:  db.Undefined,
		},
		Type: nID.GetNodeType(),
		Name: req.Name,
	}

	err = n.nodeRepo.Add(node, nil)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	route := n.baseRoutingKey.SetAction("create").SetObject("node").MustBuild()

	evt := &epb.NodeCreatedEvent{
		NodeId: node.Id,
		Name:   node.Name,
		Org:    node.OrgId.String(),
		Type:   node.Type,
	}

	err = n.msgbus.PublishRequest(route, evt)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
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

func (n *NodeServer) GetNodesForSite(ctx context.Context, req *pb.GetBySiteRequest) (*pb.GetBySiteResponse, error) {
	log.Infof("Getting all nodes on site %v", req.GetSiteId())

	site, err := uuid.FromString(req.GetSiteId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of site uuid. Error %s", err.Error())
	}

	nodes, err := n.siteRepo.GetNodes(site)
	if err != nil {
		log.Error("error getting all nodes for site" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "nodes")
	}

	resp := &pb.GetBySiteResponse{
		SiteId: req.GetSiteId(),
		Nodes:  dbNodesToPbNodes(nodes),
	}

	return resp, nil
}

func (n *NodeServer) GetNodesForOrg(ctx context.Context, req *pb.GetByOrgRequest) (*pb.GetByOrgResponse, error) {
	if req.Free {
		// return only free nodes for org
		return n.getFreeNodesForOrg(ctx, req)
	}

	// otherwise return all nodes for org
	return n.getNodesForOrg(ctx, req)
}

func (n *NodeServer) GetNodes(ctx context.Context, req *pb.GetNodesRequest) (*pb.GetNodesResponse, error) {
	if req.Free {
		// return only free nodes
		return n.getFreeNodes(ctx, req)
	}

	// otherwise return all nodes
	return n.getAllNodes(ctx, req)
}

func (n *NodeServer) UpdateNodeStatus(ctx context.Context, req *pb.UpdateNodeStateRequest) (*pb.UpdateNodeResponse, error) {
	log.Infof("Updating node state  %v", req.GetNodeId())

	dbNodeState := db.ParseNodeState(req.State)
	dbConnState := db.ParseConnectivityState(req.Connectivity)

	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	nodeUpdates := &db.NodeStatus{
		NodeId: nodeId.StringLowercase(),
	}

	pbStatus := &pb.NodeStatus{}

	if req.State != "" {
		nodeUpdates.State = dbNodeState
		pbStatus.State = dbNodeState.String()
	}

	if req.Connectivity != "" {
		nodeUpdates.Conn = dbConnState
		pbStatus.Connectivity = dbConnState.String()
	}

	err = n.nodeStatusRepo.Update(nodeUpdates)
	if err != nil {
		log.Error("error updating the node state, ", err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	und, err := n.nodeRepo.Get(nodeId)
	if err != nil {
		log.Error("error updating the node state, ", err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

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

	route := n.baseRoutingKey.SetAction("update").SetObject("node").MustBuild()

	evt := &epb.NodeUpdatedEvent{
		NodeId: nodeUpdates.Id,
		Name:   nodeUpdates.Name,
	}

	err = n.msgbus.PublishRequest(route, evt)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
	}

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

	route := n.baseRoutingKey.SetAction("delete").SetObject("node").MustBuild()

	evt := &epb.NodeDeletedEvent{
		NodeId: nodeId.StringLowercase(),
	}

	err = n.msgbus.PublishRequest(route, evt)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
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

func (n *NodeServer) AddNodeToSite(ctx context.Context, req *pb.AddNodeToSiteRequest) (*pb.AddNodeToSiteResponse, error) {
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

	route := n.baseRoutingKey.SetAction("assign").SetObject("node").MustBuild()

	evt := &epb.NodeAssignedEvent{
		NodeId:  nodeId.StringLowercase(),
		Type:    nodeId.GetNodeType(),
		Site:    site.String(),
		Network: net.String(),
	}

	err = n.msgbus.PublishRequest(route, evt)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
	}

	return &pb.AddNodeToSiteResponse{}, nil
}

func (n *NodeServer) ReleaseNodeFromSite(ctx context.Context,
	req *pb.ReleaseNodeFromSiteRequest) (*pb.ReleaseNodeFromSiteResponse, error) {
	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, invalidNodeIDError(req.GetNodeId(), err)
	}

	nd, err := n.siteRepo.RemoveNode(nodeId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	route := n.baseRoutingKey.SetAction("release").SetObject("node").MustBuild()

	evt := &epb.NodeReleasedEvent{
		NodeId:  nodeId.StringLowercase(),
		Type:    nodeId.GetNodeType(),
		Site:    nd.SiteId.String(),
		Network: nd.NetworkId.String(),
	}

	err = n.msgbus.PublishRequest(route, evt)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
	}

	return &pb.ReleaseNodeFromSiteResponse{}, nil
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

func (n *NodeServer) getNodesForOrg(ctx context.Context, req *pb.GetByOrgRequest) (*pb.GetByOrgResponse, error) {
	log.Infof("Getting all nodes for org %v", req.GetOrgId())

	org, err := uuid.FromString(req.GetOrgId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of org uuid. Error %s", err.Error())
	}

	nodes, err := n.nodeRepo.GetForOrg(org)
	if err != nil {
		log.Error("error getting all nodes for org" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "nodes")
	}

	resp := &pb.GetByOrgResponse{
		OrgId: req.GetOrgId(),
		Nodes: dbNodesToPbNodes(nodes),
	}

	return resp, nil
}

func (n *NodeServer) getFreeNodesForOrg(ctx context.Context, req *pb.GetByOrgRequest) (*pb.GetByOrgResponse, error) {
	log.Infof("Getting free nodes for org %v", req.GetOrgId())

	org, err := uuid.FromString(req.GetOrgId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of org uuid. Error %s", err.Error())
	}

	nodes, err := n.siteRepo.GetFreeNodesForOrg(org)
	if err != nil {
		log.Error("error getting free nodes for org" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "nodes")
	}

	resp := &pb.GetByOrgResponse{
		OrgId: req.GetOrgId(),
		Nodes: dbNodesToPbNodes(nodes),
	}

	return resp, nil
}

func (n *NodeServer) getAllNodes(ctx context.Context, req *pb.GetNodesRequest) (*pb.GetNodesResponse, error) {
	log.Infof("Getting all nodes.")

	nodes, err := n.nodeRepo.GetAll()

	if err != nil {
		log.Error("error getting all nodes" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.GetNodesResponse{
		Node: dbNodesToPbNodes(nodes),
	}

	return resp, nil
}

func (n *NodeServer) getFreeNodes(ctx context.Context, req *pb.GetNodesRequest) (*pb.GetNodesResponse, error) {
	log.Infof("Getting all free nodes")

	nodes, err := n.siteRepo.GetFreeNodes()

	if err != nil {
		log.Error("error getting all free nodes" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.GetNodesResponse{
		Node: dbNodesToPbNodes(nodes),
	}

	return resp, nil
}

func (n *NodeServer) pushNodeMeterics(id ukama.NodeID, args ...string) {
	nodesCount, onlineCount, offlineCount, err := n.nodeRepo.GetNodeCount()
	if err != nil {
		log.Errorf("Error while getting node count %s", err.Error())

		return
	}

	log.Infof("Updating metrics for node NodeCount %d Online %d Offline %d", nodesCount, onlineCount, offlineCount)

	for _, arg := range args {
		switch arg {
		case pkg.NumberOfNodes:
			err = metric.CollectAndPushSimMetrics(n.pushGateway, pkg.NodeMetric,
				pkg.NumberOfNodes, float64(nodesCount), nil, pkg.SystemName+"-"+pkg.ServiceName)
		case pkg.NumberOfActiveNodes:
			err = metric.CollectAndPushSimMetrics(n.pushGateway, pkg.NodeMetric,
				pkg.NumberOfActiveNodes, float64(onlineCount), nil, pkg.SystemName+"-"+pkg.ServiceName)
		case pkg.NumberOfInactiveNodes:
			err = metric.CollectAndPushSimMetrics(n.pushGateway, pkg.NodeMetric,
				pkg.NumberOfInactiveNodes, float64(offlineCount), nil, pkg.SystemName+"-"+pkg.ServiceName)
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
		Id: dbn.Id,
		Status: &pb.NodeStatus{
			Connectivity: dbn.Status.Conn.String(),
			State:        dbn.Status.State.String(),
		},
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
