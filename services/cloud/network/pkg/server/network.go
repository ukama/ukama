package server

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/pkg/errors"

	bootstrap "github.com/ukama/ukama/services/bootstrap/client"
	"github.com/ukama/ukama/services/cloud/network/pkg"
	"github.com/ukama/ukama/services/common/msgbus"

	db2 "github.com/ukama/ukama/services/cloud/network/pkg/db"
	"github.com/ukama/ukama/services/common/grpc"

	pb "github.com/ukama/ukama/services/cloud/network/pb/gen"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/sql"
	"github.com/ukama/ukama/services/common/ukama"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NetworkServer struct {
	pb.UnimplementedNetworkServiceServer
	orgRepo         db2.OrgRepo
	nodeRepo        db2.NodeRepo
	netRepo         db2.NetRepo
	bootstrapClient bootstrap.Client

	queuePub       msgbus.QPub
	baseRoutingKey msgbus.RoutingKeyBuilder
}

func NewNetworkServer(orgRepo db2.OrgRepo, nodeRepo db2.NodeRepo, netRepo db2.NetRepo, bootstrapClient bootstrap.Client,
	publisher msgbus.QPub) *NetworkServer {

	return &NetworkServer{
		orgRepo:         orgRepo,
		nodeRepo:        nodeRepo,
		netRepo:         netRepo,
		bootstrapClient: bootstrapClient,
		queuePub:        publisher,
		baseRoutingKey:  msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName),
	}
}

const defaultNetworkName = "default"

func (r *NetworkServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	org, err := r.orgRepo.MakeUserOrgExist(req.GetOrgName())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	n, err := r.netRepo.Add(org.ID, req.GetName())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	return &pb.AddResponse{
		Network: &pb.Network{
			Name: n.Name,
		},
		Org: req.GetOrgName(),
	}, nil
}

func (r *NetworkServer) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
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
			// create node type-count map
			nCnt := make(map[string]uint32)
			for t, cnt := range nodecnt {
				nCnt[dbNodeTypeToString(t)] = uint32(cnt)
			}
			orgsResp.Orgs[i].Networks[j] = &pb.ListResponse_Network{
				Name:          nname,
				NumberOfNodes: nCnt,
			}

			j++
		}
		i++
	}

	return orgsResp, nil
}

func (r *NetworkServer) AddNode(ctx context.Context, req *pb.AddNodeRequest) (*pb.AddNodeResponse, error) {
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

func (r *NetworkServer) DeleteNode(ctx context.Context, req *pb.DeleteNodeRequest) (*pb.DeleteNodeResponse, error) {
	logrus.Infof("Deleting the node  %v", req.NodeId)
	nodeId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, invalidNodeIdError(req.GetNodeId(), err)
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

func (n *NetworkServer) UpdateNode(ctx context.Context, req *pb.UpdateNodeRequest) (*pb.UpdateNodeResponse, error) {
	logrus.Infof("Updating the node  %v", req.GetNodeId())

	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, invalidNodeIdError(req.GetNodeId(), err)
	}
	st := pbNodeStateToDb(req.GetNode().GetState())
	err = n.nodeRepo.Update(nodeId, &db2.NodeAttributes{
		State: &st,
		Name:  &req.Node.Name,
	})
	if err != nil {
		duplErr := n.processNodeDuplErrors(err, req.GetNode().GetName(), req.NodeId)
		if duplErr != nil {
			return nil, duplErr
		}

		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp := &pb.UpdateNodeResponse{
		Node: &pb.Node{
			NodeId: req.NodeId,
			Name:   req.Node.Name,
			State:  req.Node.State,
		},
	}
	n.pubEvent(resp, n.baseRoutingKey.SetActionUpdate().SetObject("node").MustBuild())

	return resp, nil
}
func (r *NetworkServer) GetNodes(ctx context.Context, req *pb.GetNodesRequest) (*pb.GetNodesResponse, error) {
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

func (r *NetworkServer) pubEvent(payload any, key string) {
	go func() {
		err := r.queuePub.Publish(payload, key)
		if err != nil {
			logrus.Errorf("Failed to publish event. Error %s", err.Error())
		}
	}()
}

func (n *NetworkServer) processNodeDuplErrors(err error, nodeName string, nodeId string) error {
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

func dbNodeToPbNode(dbn *db2.Node) *pb.Node {
	n := &pb.Node{
		NodeId: dbn.NodeID,
		Type:   dbNodeTypeToPb(dbn.Type),
		Name:   dbn.Name,
		State:  dbNodeStateToPb(dbn.State),
	}

	return n
}

func dbNodeTypeToString(nodeType db2.NodeType) string {
	var pbNodeType string
	switch nodeType {
	case db2.NodeTypeAmplifier:
		pbNodeType = "amplifier"
	case db2.NodeTypeTower:
		pbNodeType = "tower"
	case db2.NodeTypeHome:
		pbNodeType = "home"
	default:
		pbNodeType = "node_type_undefined"
	}
	return pbNodeType
}

func invalidNodeIdError(nodeId string, err error) error {
	return status.Errorf(codes.InvalidArgument, "invalid node id %s. Error %s", nodeId, err.Error())
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
