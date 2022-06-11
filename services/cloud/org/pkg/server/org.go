package server

import (
	"context"
	"encoding/base64"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	bootstrap "github.com/ukama/ukama/services/bootstrap/client"
	pb "github.com/ukama/ukama/services/cloud/org/pb/gen"
	"github.com/ukama/ukama/services/cloud/org/pkg/db"
	"github.com/ukama/ukama/services/common/grpc"
	"github.com/ukama/ukama/services/common/msgbus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrgServer struct {
	pb.UnimplementedOrgServiceServer
	orgRepo db.OrgRepo

	bootstrapClient bootstrap.Client
	deviceGatewayIp string

	queuePub msgbus.QPub
}

func NewOrgServer(orgRepo db.OrgRepo, bootstrapClient bootstrap.Client,
	deviceGatewayIp string, publisher msgbus.QPub) *OrgServer {

	return &OrgServer{
		orgRepo:         orgRepo,
		bootstrapClient: bootstrapClient,
		deviceGatewayIp: deviceGatewayIp,
		queuePub:        publisher,
	}
}

func (r *OrgServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	logrus.Infof("Adding org %v", req)
	if len(req.GetOrg().GetOwner()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "owner id cannot be empty")
	}

	owner, err := uuid.Parse(req.GetOrg().GetOwner())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of owner id. Error %s", err.Error())
	}

	org := &db.Org{
		Name:        req.GetOrg().GetName(),
		Owner:       owner,
		Certificate: generateCertificate(),
	}
	err = r.orgRepo.Add(org, func() error {
		return r.bootstrapClient.AddOrUpdateOrg(org.Name, org.Certificate, r.deviceGatewayIp)
	})
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	orgResp := &pb.Organization{Name: req.GetOrg().GetName(), Owner: req.GetOrg().GetOwner()}

	r.pubEventAsync(&msgbus.OrgCreatedBody{
		Name:  orgResp.Name,
		Owner: orgResp.Owner,
	}, msgbus.OrgCreatedRoutingKey)

	return &pb.AddResponse{
		Org: orgResp,
	}, nil
}

func (r *OrgServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	logrus.Infof("Getting org %v", req)
	org, err := r.orgRepo.GetByName(req.GetName())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	return &pb.GetResponse{
		Org: &pb.Organization{
			Owner: org.Owner.String(),
			Name:  org.Name,
		},
	}, nil
}

func (r *OrgServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	logrus.Infof("Deleting org %v", req)
	err := r.orgRepo.Delete(req.GetName())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	r.pubEventAsync(&pb.Organization{Name: req.GetName()}, msgbus.OrgDeletedRoutingKey)
	return &pb.DeleteResponse{}, nil
}

func generateCertificate() string {
	logrus.Warning("Certificate generation is not yet implemented")
	return base64.StdEncoding.EncodeToString([]byte("Test certificate"))
}

func (r *OrgServer) pubEventAsync(payload any, key msgbus.RoutingKey) {
	go func() {
		err := r.queuePub.Publish(payload, string(key))
		if err != nil {
			logrus.Errorf("Failed to publish event. Error %s", err.Error())
		}
	}()
}
