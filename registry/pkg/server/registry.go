package server

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	uuid2 "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukamaX/common/sql"
	"github.com/ukama/ukamaX/common/ukama"
	"ukamaX/registry/internal/db"
	pb "ukamaX/registry/pb/generated"
)

type RegistryServer struct {
	pb.UnimplementedRegistryServiceServer
	orgRepo  db.OrgRepo
	nodeRepo db.NodeRepo
}

func NewRegistryServer(orgRepo db.OrgRepo, nodeRepo db.NodeRepo) *RegistryServer {
	return &RegistryServer{orgRepo: orgRepo, nodeRepo: nodeRepo}
}

const NOT_FOUND_ERROR = "not found"

func (r *RegistryServer) AddOrg(ctx context.Context, request *pb.AddOrgRequest) (*pb.AddOrgResponse, error) {
	logrus.Infof("Adding org %v", request)
	if len(request.Owner) == 0 {
		return nil, fmt.Errorf("owner id cannot be empty")
	}

	owner, err := uuid2.FromString(request.Owner)
	if err != nil {
		return nil, errors.Wrap(err, "invalid format of owner id")
	}

	org := &db.Org{
		Name:  request.Name,
		Owner: owner,
	}
	err = r.orgRepo.Add(org)
	if err != nil {
		if sql.IsDuplicateKeyError(err) {
			return nil, fmt.Errorf("organization already exist")
		}
		return nil, err
	}

	return &pb.AddOrgResponse{
		Org: &pb.Organization{Id: org.ID, Name: request.Name, Owner: org.Owner.String()},
	}, nil
}

func (r *RegistryServer) GetOrg(ctx context.Context, request *pb.GetOrgRequest) (*pb.Organization, error) {
	logrus.Infof("Getting org %v", request)
	org, err := r.orgRepo.GetByName(request.Name)
	if err != nil {
		if sql.IsNotFoundError(err) {
			return nil, errors.New(NOT_FOUND_ERROR)
		}

		return nil, err
	}

	return &pb.Organization{Id: org.ID, Name: org.Name, Owner: org.Owner.String()}, nil
}

func (r *RegistryServer) AddNode(ctx context.Context, req *pb.AddNodeRequest) (*pb.AddNodeResponse, error) {
	logrus.Infof("Adding node  %v", req.Node)
	if len(req.OrgName) == 0 {
		return nil, fmt.Errorf("org name cannot be empty")
	}

	org, err := r.orgRepo.GetByName(req.OrgName)
	if err != nil {
		return nil, errors.Wrap(err, "error getting org")
	}

	err = r.nodeRepo.Add(&db.Node{
		NodeID: req.Node.NodeId,
		OrgID:  org.ID,
		State:  pbNodeStateToDb(req.Node.State),
	})

	if err != nil {
		if sql.IsDuplicateKeyError(err) {
			return nil, fmt.Errorf("node already exist")
		}

		return nil, errors.Wrap(err, "error adding the node")
	}

	return &pb.AddNodeResponse{}, nil
}

func (r *RegistryServer) DeleteNode(ctx context.Context, req *pb.DeleteNodeRequest) (*pb.DeleteNodeResponse, error) {
	logrus.Infof("Deleting the node  %v", req.NodeId)
	nodeId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, err
	}

	err = r.nodeRepo.Delete(nodeId)
	if err != nil {
		return nil, errors.Wrap(err, "error deleting the node")
	}

	return &pb.DeleteNodeResponse{}, nil
}

func (r *RegistryServer) GetNode(ctx context.Context, req *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	logrus.Infof("Get node  %v", req.GetNodeId())

	nodeId, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, err
	}

	node, err := r.nodeRepo.Get(nodeId)
	if err != nil {
		if sql.IsNotFoundError(err) {
			return nil, errors.New(NOT_FOUND_ERROR)
		}
		return nil, errors.Wrap(err, "error getting the node")
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
		return nil, err
	}

	err = r.nodeRepo.Update(nodeId, dbState)
	if err != nil {
		return nil, errors.Wrap(err, "error updating the node")
	}

	return &pb.UpdateNodeResponse{}, nil
}

func (r *RegistryServer) GetNodes(ctx context.Context, req *pb.GetNodesRequest) (*pb.GetNodesResponse, error) {
	logrus.Infof("Get nodes for org %s", req.OrgName)

	owner, err := uuid2.FromString(req.Owner)
	if err != nil {
		return nil, errors.Wrap(err, "invalid format of owner id")
	}
	pbNodes := []*pb.Node{}
	nodes, err := r.nodeRepo.GetByOrg(req.OrgName, owner)
	if err != nil {
		return nil, err
	}
	for _, n := range nodes {
		pbNodes = append(pbNodes, dbNodeToPbNode(&n))
	}

	return &pb.GetNodesResponse{
		Nodes:   pbNodes,
		OrgName: req.GetOrgName(),
	}, nil
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

func dbNodeToPbNode(n *db.Node) *pb.Node {
	return &pb.Node{
		NodeId: n.NodeID,
		State:  dbNodeStateToPb(n.State),
	}
}
