/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"
	"errors"
	"fmt"
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
	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	cinvent "github.com/ukama/ukama/systems/common/rest/client/inventory"
	pb "github.com/ukama/ukama/systems/registry/node/pb/gen"
	sitepb "github.com/ukama/ukama/systems/registry/site/pb/gen"
)

const (
	Undefined = "-1"
)

type NodeServer struct {
	orgName         string
	org             uuid.UUID
	nodeRepo        db.NodeRepo
	siteRepo        db.SiteRepo
	nodeStatusRepo  db.NodeStatusRepo
	nameGenerator   namegenerator.Generator
	siteService     providers.SiteClientProvider
	pushGateway     string
	msgbus          mb.MsgBusServiceClient
	baseRoutingKey  msgbus.RoutingKeyBuilder
	inventoryClient cinvent.ComponentClient
	pb.UnimplementedNodeServiceServer
}

func NewNodeServer(orgName string, nodeRepo db.NodeRepo, siteRepo db.SiteRepo, nodeStatusRepo db.NodeStatusRepo,
	pushGateway string, msgBus mb.MsgBusServiceClient, siteService providers.SiteClientProvider, org uuid.UUID, inventoryClientProvider cinvent.ComponentClient) *NodeServer {
	seed := time.Now().UTC().UnixNano()

	return &NodeServer{
		orgName:         orgName,
		org:             org,
		nodeRepo:        nodeRepo,
		nodeStatusRepo:  nodeStatusRepo,
		siteRepo:        siteRepo,
		siteService:     siteService,
		nameGenerator:   namegenerator.NewNameGenerator(seed),
		pushGateway:     pushGateway,
		msgbus:          msgBus,
		baseRoutingKey:  msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		inventoryClient: inventoryClientProvider,
	}
}

func (n *NodeServer) AddNode(ctx context.Context, req *pb.AddNodeRequest) (*pb.AddNodeResponse, error) {
	log.Infof("Adding node  %v", req.NodeId)

	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id. Error %s", err.Error())
	}

	if len(req.Name) == 0 {
		req.Name = n.nameGenerator.Generate()
	}

	node := &db.Node{
		Id: nId.StringLowercase(),
		Status: db.NodeStatus{
			NodeId:       nId.StringLowercase(),
			Connectivity: ukama.NodeConnectivityUndefined,
			State:        ukama.NodeStateUnknown,
		},
		Type:      ukama.NodeType(nId.GetNodeType()),
		Name:      req.Name,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	err = n.nodeRepo.Add(node, nil)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	if n.msgbus != nil {
		route := n.baseRoutingKey.SetActionCreate().SetObject("node").MustBuild()

		evt := &epb.EventRegistryNodeCreate{
			NodeId:    nId.StringLowercase(),
			Name:      node.Name,
			Type:      node.Type.String(),
			Latitude:  node.Latitude,
			Longitude: node.Longitude,
		}

		err = n.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
		}
	}

	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfOnlineNodes, pkg.NumberOfOfflineNodes)

	return &pb.AddNodeResponse{Node: dbNodeToPbNode(node)}, nil
}

func (n *NodeServer) GetNode(ctx context.Context, req *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	log.Infof("Get node  %v", req.GetNodeId())

	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err.Error())
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
		log.Errorf("error getting all nodes for site: %s", err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "nodes")
	}

	resp := &pb.GetBySiteResponse{
		SiteId: req.GetSiteId(),
		Nodes:  dbNodesToPbNodes(nodes),
	}

	return resp, nil
}

func (n *NodeServer) GetNodesForNetwork(ctx context.Context, req *pb.GetByNetworkRequest) (*pb.GetByNetworkResponse, error) {
	log.Infof("Getting all nodes on site %v", req.GetNetworkId())

	network, err := uuid.FromString(req.GetNetworkId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of network uuid. Error %s", err.Error())
	}

	nodes, err := n.siteRepo.GetByNetwork(network)
	if err != nil {
		log.Errorf("error getting all nodes for network: %s", err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "nodes")
	}

	resp := &pb.GetByNetworkResponse{
		NetworkId: req.GetNetworkId(),
		Nodes:     dbNodesToPbNodes(nodes),
	}

	return resp, nil
}

/** Deprecated: Use List API instead */
func (n *NodeServer) GetNodes(ctx context.Context, req *pb.GetNodesRequest) (*pb.GetNodesResponse, error) {
	log.Infof("Getting all nodes.")

	nodes, err := n.nodeRepo.GetAll()

	if err != nil {
		log.Error("error getting all nodes" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.GetNodesResponse{
		Nodes: dbNodesToPbNodes(nodes),
	}

	return resp, nil
}

/** Deprecated: Use List API instead */
func (n *NodeServer) GetNodesByState(ctx context.Context, req *pb.GetNodesByStateRequest) (*pb.GetNodesResponse, error) {
	log.Infof("Get nodes by state with connectivity: %v, state: %v", req.GetConnectivity(), req.GetState())

	nodes, err := n.nodeRepo.GetNodesByState(uint8(req.GetConnectivity()), uint8(req.GetState()))
	if err != nil {
		log.Error("error getting all nodes by state" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.GetNodesResponse{
		Nodes: dbNodesToPbNodes(nodes),
	}

	fmt.Printf("Nodes Resp returning %v", resp)
	return resp, nil
}

func (n *NodeServer) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	log.Infof("List nodes by nodeId: %v, siteId: %v, networkId: %v, connectivity: %v, state: %v, type: %v", req.GetNodeId(), req.GetSiteId(), req.GetNetworkId(), req.GetConnectivity().String(), req.GetState().String(), req.GetType())

	var connectivity, state *uint8
	if req.GetConnectivity().String() != Undefined {
		c := uint8(req.GetConnectivity())
		connectivity = &c
	}
	if req.GetState().String() != Undefined {
		s := uint8(req.GetState())
		state = &s
	}
	nodes, err := n.nodeRepo.List(req.GetNodeId(), req.GetSiteId(), req.GetNetworkId(), req.GetType(), connectivity, state)

	if err != nil {
		log.Error("error getting all nodes: " + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.ListResponse{
		Nodes: dbNodesToPbNodes(nodes),
	}

	return resp, nil
}

func (n *NodeServer) UpdateNodeStatus(ctx context.Context, req *pb.UpdateNodeStateRequest) (*pb.UpdateNodeResponse, error) {
	log.Infof("Updating node state  %v", req)

	dbNodeState := ukama.ParseNodeState(req.State)
	dbConnState := ukama.ParseNodeConnectivity(req.Connectivity)

	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err.Error())
	}

	nodeUpdates := &db.NodeStatus{
		NodeId: nodeId.StringLowercase(),
	}

	evt := &epb.EventRegistryNodeStatusUpdate{
		NodeId: nodeUpdates.NodeId,
	}

	if req.Connectivity != "" {
		nodeUpdates.Connectivity = dbConnState
		evt.Status = &epb.EventRegistryNodeStatusUpdate_Connectivity{
			Connectivity: dbConnState.String(),
		}
	}

	if req.State != "" {
		nodeUpdates.State = dbNodeState
		evt.Status = &epb.EventRegistryNodeStatusUpdate_State{
			State: dbNodeState.String(),
		}
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

	if n.msgbus != nil {
		route := n.baseRoutingKey.SetActionUpdate().SetObject("status").MustBuild()

		err = n.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
		}
	}

	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfOnlineNodes, pkg.NumberOfOfflineNodes)

	return &pb.UpdateNodeResponse{Node: dbNodeToPbNode(und)}, nil
}

func (n *NodeServer) UpdateNode(ctx context.Context, req *pb.UpdateNodeRequest) (*pb.UpdateNodeResponse, error) {
	log.Infof("Updating node  %v", req.GetNodeId())

	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err.Error())
	}

	nodeUpdates := &db.Node{
		Id:        nodeId.StringLowercase(),
		Name:      req.Name,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
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
			Id:        req.NodeId,
			Name:      req.Name,
			Latitude:  req.Latitude,
			Longitude: req.Longitude,
		},
	}

	und, err := n.nodeRepo.Get(nodeId)
	if err != nil {
		log.Error("error getting the node, ", err.Error())
		return resp, nil
	}

	if n.msgbus != nil {
		route := n.baseRoutingKey.SetActionUpdate().SetObject("node").MustBuild()

		evt := &epb.EventRegistryNodeUpdate{
			NodeId: nodeUpdates.Id,
			Name:   nodeUpdates.Name,
		}

		err = n.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
		}
	}
	return &pb.UpdateNodeResponse{Node: dbNodeToPbNode(und)}, nil
}

func (n *NodeServer) DeleteNode(ctx context.Context, req *pb.DeleteNodeRequest) (*pb.DeleteNodeResponse, error) {
	log.Infof("Deleting node  %v", req.GetNodeId())

	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err.Error())
	}

	err = n.nodeRepo.Delete(nodeId, nil)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	if n.msgbus != nil {
		route := n.baseRoutingKey.SetActionDelete().SetObject("node").MustBuild()

		evt := &epb.EventRegistryNodeDelete{
			NodeId: nodeId.StringLowercase(),
		}

		err = n.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
		}
	}

	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfOnlineNodes, pkg.NumberOfOfflineNodes)

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
		log.Errorf("fail to attach nodes. Errors %s", err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	if n.msgbus != nil {
		route := n.baseRoutingKey.SetAction("attach").SetObject("node").MustBuild()

		evt := &epb.EventRegistryNodeAttach{
			NodeId:    nodeId.StringLowercase(),
			Nodegroup: nds,
		}

		err = n.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
		}
	}

	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfOnlineNodes, pkg.NumberOfOfflineNodes)

	return &pb.AttachNodesResponse{}, nil
}

func (n *NodeServer) DetachNode(ctx context.Context, req *pb.DetachNodeRequest) (*pb.DetachNodeResponse, error) {
	log.Infof("detaching node  %v", req.GetNodeId())

	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, invalidNodeIDError(req.GetNodeId(), err)
	}

	node, err := n.nodeRepo.Get(nodeId)
	if err != nil {
		log.Error("error getting the node, ", err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	attachednodes := make([]string, len(node.Attached))
	for _, an := range node.Attached {
		attachednodes = append(attachednodes, an.Id)
	}

	err = n.nodeRepo.DetachNode(nodeId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	if n.msgbus != nil {
		route := n.baseRoutingKey.SetAction("dettach").SetObject("node").MustBuild()

		evt := &epb.EventRegistryNodeAttach{
			NodeId:    nodeId.StringLowercase(),
			Nodegroup: attachednodes,
		}

		err = n.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
		}
	}

	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfOnlineNodes, pkg.NumberOfOfflineNodes)

	return &pb.DetachNodeResponse{}, nil
}

func (n *NodeServer) AddNodeToSite(ctx context.Context, req *pb.AddNodeToSiteRequest) (*pb.AddNodeToSiteResponse, error) {
	log.Infof("Add node req : %s", req)
	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, invalidNodeIDError(req.GetNodeId(), err)
	}

	netID, err := uuid.FromString(req.GetNetworkId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of network uuid. Error %s", err.Error())
	}

	siteID, err := uuid.FromString(req.GetSiteId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of site uuid. Error %s", err.Error())

	}

	svc, err := n.siteService.GetClient()
	if err != nil {
		return nil, err
	}

	remoteSite, err := svc.Get(ctx, &sitepb.GetRequest{SiteId: siteID.String()})
	if err != nil {

		return nil, err
	}

	if remoteSite.Site.NetworkId != netID.String() {
		return nil, status.Errorf(codes.FailedPrecondition,
			"provided networkId and site's networkId mismatch")
	}

	node := &db.Site{
		NodeId:    nodeId.StringLowercase(),
		SiteId:    siteID,
		NetworkId: netID,
	}

	err = n.siteRepo.AddNode(node, nil)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	if n.msgbus != nil {
		route := n.baseRoutingKey.SetAction("assign").SetObject("node").MustBuild()

		evt := &epb.EventRegistryNodeAssign{
			NodeId:  nodeId.StringLowercase(),
			Type:    nodeId.GetNodeType(),
			Site:    siteID.String(),
			Network: netID.String(),
		}

		err = n.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
		}
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

	if n.msgbus != nil {
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
	}

	return &pb.ReleaseNodeFromSiteResponse{}, nil
}
func (n *NodeServer) addNodeToSite(nodeId, siteId, networkId string) error {
	log.Infof("Add node to site %s", nodeId)
	r, err := n.inventoryClient.Get(nodeId)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "node not found in inventory against component id: %s, Error %s", nodeId, err.Error())
	}

	netID, err := uuid.FromString(networkId)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid format of network uuid. Error %s", err.Error())
	}

	siteID, err := uuid.FromString(siteId)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid format of site uuid. Error %s", err.Error())
	}

	site := &db.Site{
		NodeId:    r.PartNumber,
		SiteId:    siteID,
		NetworkId: netID,
	}

	err = n.siteRepo.AddNode(site, nil)
	if err != nil {
		return grpc.SqlErrorToGrpc(err, "node")
	}

	if n.msgbus != nil {
		route := n.baseRoutingKey.SetAction("assign").SetObject("node").MustBuild()

		evt := &epb.EventRegistryNodeAssign{
			NodeId:  r.PartNumber,
			Type:    r.Type,
			Site:    siteId,
			Network: networkId,
		}

		err = n.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
		}
	}

	return nil
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
		case pkg.NumberOfOnlineNodes:
			err = metric.CollectAndPushSimMetrics(n.pushGateway, pkg.NodeMetric,
				pkg.NumberOfOnlineNodes, float64(onlineCount), nil, pkg.SystemName+"-"+pkg.ServiceName)
		case pkg.NumberOfOfflineNodes:
			err = metric.CollectAndPushSimMetrics(n.pushGateway, pkg.NodeMetric,
				pkg.NumberOfOfflineNodes, float64(offlineCount), nil, pkg.SystemName+"-"+pkg.ServiceName)
		}
	}

	if err != nil {
		log.Errorf("Error while pushing node metric to pushgateway %s", err.Error())
	}
}

func (n *NodeServer) PushMetrics() {
	n.pushNodeMeterics(pkg.NumberOfNodes, pkg.NumberOfOnlineNodes, pkg.NumberOfOfflineNodes)
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
			Connectivity: cpb.NodeConnectivity(dbn.Status.Connectivity),
			State:        cpb.NodeState(dbn.Status.State),
		},
		Type:      dbn.Type.String(),
		Name:      dbn.Name,
		Latitude:  dbn.Latitude,
		Longitude: dbn.Longitude,
	}

	if dbn.Site.NodeId != "" {
		n.Site = &pb.Site{}
		n.Site.NodeId = dbn.Site.NodeId
		n.Site.SiteId = dbn.Site.SiteId.String()
		n.Site.NetworkId = dbn.Site.NetworkId.String()
		n.Site.AddedAt = timestamppb.New(dbn.Site.CreatedAt)
	}

	if len(dbn.Attached) > 0 {
		n.Attached = make([]*pb.Node, 0)
	}

	for _, nd := range dbn.Attached {
		n.Attached = append(n.Attached, dbNodeToPbNode(nd))
	}

	return n
}
