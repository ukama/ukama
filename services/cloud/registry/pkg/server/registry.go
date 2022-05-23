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
	
	queuePub        msgbus.QPub
	baseRoutingKey  msgbus.RoutingKeyBuilder
}

func NewRegistryServer(orgRepo db2.OrgRepo, nodeRepo db2.NodeRepo, netRepo db2.NetRepo, bootstrapClient bootstrap.Client,
	deviceGatewayIp string, publisher msgbus.QPub) *RegistryServer {
	

	return &RegistryServer{
		orgRepo:         orgRepo,
		nodeRepo:        nodeRepo,
		bootstrapClient: bootstrapClient,
		deviceGatewayIp: deviceGatewayIp,
		netRepo:         netRepo,
		
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

func (r *RegistryServer) UpdateNode(ctx context.Context, req *pb.UpdateNodeRequest) (*pb.UpdateNodeResponse, error) {
	
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

func (r *RegistryServer) pubEvent(payload any, key string) {
	go func() {
		err := r.queuePub.Publish(payload, key)
		if err != nil {
			logrus.Errorf("Failed to publish event. Error %s", err.Error())
		}
	}()
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


func invalidNodeIdError(nodeId string, err error) error {
	return status.Errorf(codes.InvalidArgument, "invalid node id %s. Error %s", nodeId, err.Error())
}
