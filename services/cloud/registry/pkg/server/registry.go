package server

import (
	"context"
	"encoding/base64"
	"errors"
	"time"

	"github.com/ukama/ukama/services/cloud/registry/pkg"
	"github.com/ukama/ukama/services/common/msgbus"

	"github.com/jackc/pgconn"
	db2 "github.com/ukama/ukama/services/cloud/registry/pkg/db"
	"github.com/ukama/ukama/services/common/grpc"

	pb "github.com/ukama/ukama/services/cloud/registry/pb/gen"
	"github.com/ukama/ukama/services/cloud/registry/pkg/bootstrap"

	"github.com/goombaio/namegenerator"

	uuid2 "github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/sql"
	"github.com/ukama/ukama/services/common/ukama"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RegistryServer struct {
	pb.UnimplementedRegistryServiceServer
	orgRepo         db2.OrgRepo
	nodeRepo        db2.NodeRepo
	netRepo         db2.NetRepo
	bootstrapClient bootstrap.Client
	deviceGatewayIp string
	nameGenerator   namegenerator.Generator
	queuePub        msgbus.QPub
	baseRoutingKey  msgbus.RoutingKeyBuilder
}

func NewRegistryServer(orgRepo db2.OrgRepo, nodeRepo db2.NodeRepo, netRepo db2.NetRepo, bootstrapClient bootstrap.Client,
	deviceGatewayIp string, publisher msgbus.QPub) *RegistryServer {
	seed := time.Now().UTC().UnixNano()

	return &RegistryServer{
		orgRepo:         orgRepo,
		nodeRepo:        nodeRepo,
		bootstrapClient: bootstrapClient,
		deviceGatewayIp: deviceGatewayIp,
		netRepo:         netRepo,
		nameGenerator:   namegenerator.NewNameGenerator(seed),
		queuePub:        publisher,
		baseRoutingKey:  msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName),
	}
}

const defaultNetworkName = "default"

func (r *RegistryServer) AddOrg(ctx context.Context, request *pb.AddOrgRequest) (*pb.AddOrgResponse, error) {
	logrus.Infof("Adding org %v", request)
	if len(request.Owner) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "owner id cannot be empty")
	}

	owner, err := uuid2.Parse(request.Owner)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of owner id. Error %s", err.Error())
	}

	org := &db2.Org{
		Name:        request.Name,
		Owner:       owner,
		Certificate: generateCertificate(),
	}
	err = r.orgRepo.Add(org, func() error {
		return r.bootstrapClient.AddOrUpdateOrg(org.Name, org.Certificate, r.deviceGatewayIp)
	})
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	_, err = r.netRepo.Add(org.BaseModel.ID, defaultNetworkName)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	orgResp := &pb.Organization{Name: request.Name, Owner: org.Owner.String()}

	r.pubEvent(orgResp, r.baseRoutingKey.SetActionCreate().SetObject("org").MustBuild())

	return &pb.AddOrgResponse{
		Org: orgResp,
	}, nil
}

func generateCertificate() string {
	logrus.Warning("Certificate generation is not yet implemented")
	return base64.StdEncoding.EncodeToString([]byte("Test certificate"))
}

func (r *RegistryServer) GetOrg(ctx context.Context, request *pb.GetOrgRequest) (*pb.Organization, error) {
	logrus.Infof("Getting org %v", request)
	org, err := r.orgRepo.GetByName(request.Name)
	if err != nil {
		if sql.IsNotFoundError(err) {
			return nil, status.Errorf(codes.NotFound, "Organization not found")
		}

		logrus.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.Organization{Name: org.Name, Owner: org.Owner.String()}, nil
}

func (r *RegistryServer) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	logrus.Infof("Listing orgs")
	orgs, err := r.netRepo.List()
	if err != nil {
		logrus.Error(err)
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	orgsResp := &pb.ListResponse{
		Orgs: make([]*pb.ListResponse_Org, len(orgs)),
	}

	i := 0
	for o, n := range orgs {
		orgsResp.Orgs[i] = &pb.ListResponse_Org{
			Name:     o,
			Networks: make([]*pb.ListResponse_Network, len(n)),
		}

		j := 0
		for nname, nodecnt := range n {
			orgsResp.Orgs[i].Networks[j] = &pb.ListResponse_Network{
				Name:          nname,
				NumberOfNodes: int32(nodecnt),
			}
			j++
		}
		i++
	}

	return orgsResp, nil
}

func (r *RegistryServer) AddNode(ctx context.Context, req *pb.AddNodeRequest) (*pb.AddNodeResponse, error) {
	logrus.Infof("Adding node  %v", req.Node)
	if len(req.OrgName) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "organization name cannot be empty")
	}

	netName := req.Network
	if len(req.Network) == 0 {
		netName = defaultNetworkName
	}

	network, err := r.netRepo.Get(req.OrgName, netName)
	if err != nil {
		logrus.Error(err)
		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	nId, err := ukama.ValidateNodeId(req.Node.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of node id. Error %s", err.Error())
	}

	if req.Node.Type != pb.NodeType_NODE_TYPE_UNDEFINED {
		return nil, status.Errorf(codes.InvalidArgument, "type is determined from nodeId and can not be set specifically")
	}

	// Generate random node name if it's missing
	if len(req.Node.Name) == 0 {
		req.Node.Name = r.nameGenerator.Generate()
	}

	node := &db2.Node{
		NodeID:    req.Node.NodeId,
		State:     pbNodeStateToDb(req.Node.State),
		Type:      toDbNodeType(nId.GetNodeType()),
		Name:      req.Node.Name,
		NetworkID: network.ID,
	}

	// adding node to DB and bootstrap in transaction
	// Rollback trans if bootstrap fails to add a node
	err = r.nodeRepo.Add(node, func() error {
		return r.bootstrapClient.AddNode(network.Org.Name, node.NodeID)
	})

	if err != nil {
		duplErr := r.processNodeDuplErrors(err, node.Name, node.NodeID)
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

func (r *RegistryServer) processNodeDuplErrors(err error, nodeName string, nodeId string) error {
	var pge *pgconn.PgError
	if errors.As(err, &pge) {
		if pge.Code == sql.PGERROR_CODE_UNIQUE_VIOLATION && pge.ConstraintName == "node_name_network_idx" {
			return status.Errorf(codes.AlreadyExists, "node with name %s already exists in network", nodeName)
		} else if pge.Code == sql.PGERROR_CODE_UNIQUE_VIOLATION {
			return status.Errorf(codes.AlreadyExists, "node with node id %s already exist", nodeId)
		}
	}
	return nil
}

func toDbNodeType(nodeType string) db2.NodeType {
	switch nodeType {
	case ukama.NODE_ID_TYPE_AMPNODE:
		return db2.NodeTypeAmplifier
	case ukama.NODE_ID_TYPE_TOWERNODE:
		return db2.NodeTypeTower
	case ukama.NODE_ID_TYPE_HOMENODE:
		return db2.NodeTypeHome
	default:
		return db2.NodeTypeUnknown
	}
}

func (r *RegistryServer) DeleteNode(ctx context.Context, req *pb.DeleteNodeRequest) (*pb.DeleteNodeResponse, error) {
	logrus.Infof("Deleting the node  %v", req.NodeId)
	nodeId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = r.nodeRepo.Delete(nodeId)
	if err != nil {
		logrus.Error("error deleting the node, ", err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.DeleteNodeResponse{
		NodeId: req.NodeId,
	}
	r.pubEvent(resp, r.baseRoutingKey.SetActionDelete().SetObject("node").MustBuild())

	return resp, nil
}

func (r *RegistryServer) GetNode(ctx context.Context, req *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	logrus.Infof("Get node  %v", req.GetNodeId())

	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	node, err := r.nodeRepo.Get(nodeId)
	if err != nil {
		logrus.Error("error getting the node" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.GetNodeResponse{
		Node: dbNodeToPbNode(node),
		Org: &pb.Organization{
			Owner: node.Network.Org.Owner.String(),
			Name:  node.Network.Org.Name,
		},
	}
	if node.Network != nil {
		resp.Network = &pb.Network{
			Name: node.Network.Name,
		}
	} else {
		resp.Network = &pb.Network{
			Name: "default",
		}
	}

	return resp, nil
}

func (r *RegistryServer) UpdateNodeState(ctx context.Context, req *pb.UpdateNodeStateRequest) (*pb.UpdateNodeStateResponse, error) {
	logrus.Infof("Updating node state  %v", req.GetNodeId())

	dbState := pbNodeStateToDb(req.State)

	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = r.nodeRepo.Update(nodeId, &dbState, nil)
	if err != nil {
		logrus.Error("error updating the node state, ", err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.UpdateNodeStateResponse{
		NodeId: req.GetNodeId(),
		State:  req.State,
	}
	r.pubEvent(resp, r.baseRoutingKey.SetActionUpdate().SetObject("node").MustBuild())

	return resp, nil
}

func (r *RegistryServer) UpdateNode(ctx context.Context, req *pb.UpdateNodeRequest) (*pb.UpdateNodeResponse, error) {
	logrus.Infof("Updating the node  %v", req.GetNodeId())

	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = r.nodeRepo.Update(nodeId, nil, &req.Name)
	if err != nil {
		duplErr := r.processNodeDuplErrors(err, req.Name, req.NodeId)
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
	r.pubEvent(resp, r.baseRoutingKey.SetActionUpdate().SetObject("node").MustBuild())

	return resp, nil
}

func (r *RegistryServer) GetNodes(ctx context.Context, req *pb.GetNodesRequest) (*pb.GetNodesResponse, error) {
	logrus.Infof("Get nodes for org %s", req.OrgName)

	var nodes []db2.Node
	var err error
	if len(req.OrgName) != 0 {
		nodes, err = r.nodeRepo.GetByOrg(req.OrgName)
	}
	if err != nil {
		logrus.Error(err)
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	pbNodes := make([]*pb.Node, 0)
	for _, n := range nodes {
		pbNodes = append(pbNodes, dbNodeToPbNode(&n))
	}

	resp := &pb.GetNodesResponse{
		Nodes:   pbNodes,
		OrgName: req.OrgName,
	}

	return resp, nil
}

func (r *RegistryServer) AttachNodes(ctx context.Context, req *pb.AttachNodesRequest) (*pb.AttachNodesResponse, error) {
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
	err = r.nodeRepo.AttachNodes(nodeId, nds)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	r.pubEvent(req, r.baseRoutingKey.SetActionUpdate().SetObject("node").MustBuild())

	return &pb.AttachNodesResponse{}, nil
}

func (r *RegistryServer) DetachNodes(ctx context.Context, req *pb.DetachNodeRequest) (*pb.DetachNodeResponse, error) {
	nodeId, err := ukama.ValidateNodeId(req.DetachedNodeId)
	if err != nil {
		return nil, invalidNodeIdError(req.DetachedNodeId, err)
	}

	err = r.nodeRepo.DetachNode(nodeId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	r.pubEvent(req, r.baseRoutingKey.SetActionUpdate().SetObject("node").MustBuild())

	return &pb.DetachNodeResponse{}, nil
}

func (r *RegistryServer) pubEvent(payload any, key string) {
	go func() {
		err := r.queuePub.Publish(payload, key)
		if err != nil {
			logrus.Errorf("Failed to publish event. Error %s", err.Error())
		}
	}()
}

func pbNodeStateToDb(state pb.NodeState) db2.NodeState {
	var dbState db2.NodeState
	switch state {
	case pb.NodeState_ONBOARDED:
		dbState = db2.Onboarded
	case pb.NodeState_PENDING:
		dbState = db2.Pending
	default:
		dbState = db2.Undefined
	}
	return dbState
}

func dbNodeStateToPb(state db2.NodeState) pb.NodeState {
	var pbState pb.NodeState
	switch state {
	case db2.Onboarded:
		pbState = pb.NodeState_ONBOARDED
	case db2.Pending:
		pbState = pb.NodeState_PENDING
	default:
		pbState = pb.NodeState_UNDEFINED
	}

	return pbState
}

func dbNodeToPbNode(dbn *db2.Node) *pb.Node {
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

func dbNodeTypeToPb(nodeType db2.NodeType) pb.NodeType {
	var pbNodeType pb.NodeType
	switch nodeType {
	case db2.NodeTypeAmplifier:
		pbNodeType = pb.NodeType_AMPLIFIER
	case db2.NodeTypeTower:
		pbNodeType = pb.NodeType_TOWER
	case db2.NodeTypeHome:
		pbNodeType = pb.NodeType_HOME
	default:
		pbNodeType = pb.NodeType_NODE_TYPE_UNDEFINED
	}
	return pbNodeType
}

func invalidNodeIdError(nodeId string, err error) error {
	return status.Errorf(codes.InvalidArgument, "invalid node id %s. Error %s", nodeId, err.Error())
}
