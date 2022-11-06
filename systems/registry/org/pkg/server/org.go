package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/ukama/ukama/systems/common/grpc"
	pb "github.com/ukama/ukama/systems/registry/org/pb/gen"
	"github.com/ukama/ukama/systems/registry/org/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrgServer struct {
	pb.UnimplementedOrgServiceServer
	orgRepo db.OrgRepo
}

func NewOrgServer(orgRepo db.OrgRepo) *OrgServer {
	return &OrgServer{
		orgRepo: orgRepo,
	}
}

func (r *OrgServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	logrus.Infof("Adding org %v", req)

	if len(req.GetOrg().GetOwner()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "owner uuid cannot be empty")
	}

	owner, err := uuid.Parse(req.GetOrg().GetOwner())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of owner id. Error %s", err.Error())
	}

	org := &db.Org{
		Name:        req.GetOrg().GetName(),
		Owner:       owner,
		Certificate: req.GetOrg().GetCertificate(),
	}

	// err = r.orgRepo.Add(org, func() error {
	// We need to wrap this call into a transaction to add or update the org in the new init system.
	// see an example with the legacy interface below.
	// return r.bootstrapClient.AddOrUpdateOrg(org.Name, org.Certificate, r.nodeGatewayIP)
	// return nil
	// })
	err = r.orgRepo.Add(org)

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	// pub event async

	return &pb.AddResponse{Org: dbOrgsToPbOrgs(org)}, nil
}

func (r *OrgServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	logrus.Infof("Getting org %v", req)

	org, err := r.orgRepo.Get(int(req.GetId()))
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	return &pb.GetResponse{Org: dbOrgsToPbOrgs(org)}, nil
}

// func (r *OrgServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
// logrus.Infof("Deleting org %v", req)

// err := r.orgRepo.Delete(req.GetName())
// if err != nil {
// return nil, grpc.SqlErrorToGrpc(err, "org")
// }

// // publish event async

// return &pb.DeleteResponse{}, nil
// }

func dbOrgsToPbOrgs(org *db.Org) *pb.Organization {
	return &pb.Organization{
		Id:            uint64(org.ID),
		Name:          org.Name,
		Owner:         org.Owner.String(),
		Certificate:   org.Certificate,
		IsDeactivated: org.Deactivated,
		CreatedAt:     timestamppb.New(org.CreatedAt),
	}
}
