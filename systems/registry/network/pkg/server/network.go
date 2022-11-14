package server

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/pkg/errors"

	"github.com/ukama/ukama/systems/common/grpc"
	db2 "github.com/ukama/ukama/systems/registry/network/pkg/db"

	pb "github.com/ukama/ukama/systems/registry/network/pb/gen"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NetworkServer struct {
	pb.UnimplementedNetworkServiceServer
	orgRepo db2.OrgRepo
	netRepo db2.NetRepo
}

func NewNetworkServer(orgRepo db2.OrgRepo, netRepo db2.NetRepo) *NetworkServer {
	return &NetworkServer{
		orgRepo: orgRepo,
		netRepo: netRepo,
	}
}

func (n *NetworkServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	org, err := n.orgRepo.MakeUserOrgExist(req.GetOrgName())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	nt, err := n.netRepo.Add(org.ID, req.GetName())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	return &pb.AddResponse{
		Network: &pb.Network{
			Name: nt.Name,
		},
		Org: req.GetOrgName(),
	}, nil
}

func (n *NetworkServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	logrus.Infof("Deleting network %s", req.Name)

	err := n.netRepo.Delete(req.OrgName, req.Name)
	if err != nil {
		logrus.Error(err)

		return nil, grpc.SqlErrorToGrpc(err, "network")
	}

	resp := &pb.DeleteResponse{}

	// publish event before returning resp

	return resp, nil
}

func (n *NetworkServer) processNodeDuplErrors(err error, nodeName string, nodeID string) error {
	var pge *pgconn.PgError

	if errors.As(err, &pge) {
		if pge.Code == sql.PGERROR_CODE_UNIQUE_VIOLATION && pge.ConstraintName == "node_name_network_idx" {
			return status.Errorf(codes.AlreadyExists, "node with name %s already exists in network", nodeName)
		} else if pge.Code == sql.PGERROR_CODE_UNIQUE_VIOLATION {
			return status.Errorf(codes.AlreadyExists, "node with node id %s already exist", nodeID)
		}
	}

	return nil
}
