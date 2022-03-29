package server

import (
	"context"
	"encoding/base64"
	db2 "github.com/ukama/ukamaX/cloud/registry/pkg/db"
	"time"

	pb "github.com/ukama/ukamaX/cloud/registry/pb/gen"
	"github.com/ukama/ukamaX/cloud/registry/pkg/bootstrap"

	"github.com/goombaio/namegenerator"

	uuid2 "github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukamaX/common/sql"
	"github.com/ukama/ukamaX/common/ukama"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RegistryServer struct {
	pb.UnimplementedRegistryServiceServer
	orgRepo         db2.OrgRepo
	nodeRepo        db2.NodeRepo
	bootstrapClient bootstrap.Client
	deviceGatewayIp string
	nameGenerator   namegenerator.Generator
}

func NewRegistryServer(orgRepo db2.OrgRepo, nodeRepo db2.NodeRepo, bootstrapClient bootstrap.Client, deviceGatewayIp string) *RegistryServer {
	seed := time.Now().UTC().UnixNano()

	return &RegistryServer{
		orgRepo:         orgRepo,
		nodeRepo:        nodeRepo,
		bootstrapClient: bootstrapClient,
		deviceGatewayIp: deviceGatewayIp,
		nameGenerator:   namegenerator.NewNameGenerator(seed)}
}

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
		if sql.IsDuplicateKeyError(err) {
			return nil, status.Errorf(codes.AlreadyExists, "organization already exist")
		}
		logrus.Errorf("Error adding the org. Error: %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.AddOrgResponse{
		Org: &pb.Organization{Id: org.ID, Name: request.Name, Owner: org.Owner.String()},
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

	return &pb.Organization{Id: org.ID, Name: org.Name, Owner: org.Owner.String()}, nil
}

func (r *RegistryServer) AddNode(ctx context.Context, req *pb.AddNodeRequest) (*pb.AddNodeResponse, error) {
	logrus.Infof("Adding node  %v", req.Node)
	if len(req.OrgName) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "organization name cannot be empty")
	}

	org, err := r.orgRepo.GetByName(req.OrgName)
	if err != nil {
		if sql.IsNotFoundError(err) {
			return nil, status.Errorf(codes.NotFound, "Organization not found")
		}

		logrus.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
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
		NodeID: req.Node.NodeId,
		OrgID:  org.ID,
		State:  pbNodeStateToDb(req.Node.State),
		Type:   toDbNodeType(nId.GetNodeType()),
		Name:   req.Node.Name,
	}

	// adding node to DB and bootstrap in transaction
	// Rollback trans if bootstrap fails to add a node
	err = r.nodeRepo.Add(node, func() error {
		return r.bootstrapClient.AddNode(org.Name, node.NodeID)
	})

	if err != nil {
		if sql.IsDuplicateKeyError(err) {
			return nil, status.Errorf(codes.AlreadyExists, "node already exist")
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
	case ukama.NODE_ID_TYPE_COMPNODE:
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
		return nil, status.Errorf(codes.Internal, "error deleting the node")
	}

	return &pb.DeleteNodeResponse{}, nil
}

func (r *RegistryServer) GetNode(ctx context.Context, req *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	logrus.Infof("Get node  %v", req.GetNodeId())

	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	node, err := r.nodeRepo.Get(nodeId)
	if err != nil {
		if sql.IsNotFoundError(err) {
			return nil, status.Errorf(codes.NotFound, "node not found")
		}

		logrus.Error("error getting the node" + err.Error())
		return nil, status.Errorf(codes.Internal, "error getting the node")
	}

	return &pb.GetNodeResponse{
		Node: dbNodeToPbNode(node),
		Org: &pb.Organization{
			Id:    node.OrgID,
			Owner: node.Org.Owner.String(),
			Name:  node.Org.Name,
		},
	}, nil
}

func (r *RegistryServer) UpdateNode(ctx context.Context, req *pb.UpdateNodeRequest) (*pb.UpdateNodeResponse, error) {
	logrus.Infof("Updating the node  %v", req.GetNodeId())

	dbState := pbNodeStateToDb(req.State)

	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = r.nodeRepo.Update(nodeId, dbState)
	if err != nil {
		if sql.IsNotFoundError(err) {
			return nil, status.Errorf(codes.NotFound, "node not found")
		}
		logrus.Error("error updating the node" + err.Error())
		return nil, status.Errorf(codes.Internal, "error updating the node")
	}

	return &pb.UpdateNodeResponse{}, nil
}

func (r *RegistryServer) GetNodes(ctx context.Context, req *pb.GetNodesRequest) (*pb.GetNodesResponse, error) {
	logrus.Infof("Get nodes for org %s", req.OrgName)

	var nodes []db2.Node
	var err error
	if len(req.OrgName) != 0 {
		nodes, err = r.nodeRepo.GetByOrg(req.OrgName)
	}
	if err != nil {
		logrus.Error(err.Error())
		return nil, status.Errorf(codes.Internal, "error getting nodes")
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
	return &pb.Node{
		NodeId: dbn.NodeID,
		State:  dbNodeStateToPb(dbn.State),
		Type:   dbNodeTypeToPb(dbn.Type),
		Name:   dbn.Name,
	}
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
