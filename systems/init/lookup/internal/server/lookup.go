package server

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/goombaio/namegenerator"
	"github.com/jackc/pgtype"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/grpc"
	"github.com/ukama/ukama/services/common/msgbus"
	"github.com/ukama/ukama/services/common/ukama"
	"github.com/ukama/ukama/systems/init/lookup/internal"
	"github.com/ukama/ukama/systems/init/lookup/internal/db"
	pb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
	mb "github.com/ukama/ukama/systems/init/lookup/pkg/msgBusClient"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type LookupServer struct {
	systemRepo     db.SystemRepo
	orgRepo        db.OrgRepo
	nodeRepo       db.NodeRepo
	msgbus         mb.MsgBusClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	nameGenerator  namegenerator.Generator
	pb.UnimplementedLookupServiceServer
}

func NewLookupServer(nodeRepo db.NodeRepo, orgRepo db.OrgRepo, systemRepo db.SystemRepo, msgBus *mb.MsgBusClient) *LookupServer {
	seed := time.Now().UTC().UnixNano()
	return &LookupServer{
		nodeRepo:   nodeRepo,
		orgRepo:    orgRepo,
		systemRepo: systemRepo,
		//msgbus:         *msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(internal.ServiceName),
		nameGenerator:  namegenerator.NewNameGenerator(seed),
	}
}

func (l *LookupServer) AddOrg(ctx context.Context, req *pb.AddOrgRequest) (*pb.AddOrgResponse, error) {
	logrus.Infof("Adding Organization %s", req.OrgName)

	var orgIp pgtype.Inet

	err := orgIp.Set(req.Ip)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid ip for Org %s. Error %s", req.OrgName, err.Error())
	}

	org := &db.Org{
		Name:        req.GetOrgName(),
		Certificate: req.GetCertificate(),
		Ip:          orgIp,
	}

	err = l.orgRepo.Upsert(org)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	route := l.baseRoutingKey.SetActionCreate().SetObject("organization").MustBuild()
	err = l.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	return &pb.AddOrgResponse{}, nil
}

func (l *LookupServer) UpdateOrg(ctx context.Context, req *pb.UpdateOrgRequest) (*pb.UpdateOrgResponse, error) {
	logrus.Infof("Updating Organization %s", req.OrgName)

	var orgIp pgtype.Inet

	err := orgIp.Set(req.Ip)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid ip for Org %s. Error %s", req.OrgName, err.Error())
	}

	org := &db.Org{
		Name:        req.GetOrgName(),
		Certificate: req.GetCertificate(),
		Ip:          orgIp,
	}

	err = l.orgRepo.Upsert(org)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	route := l.baseRoutingKey.SetActionUpdate().SetObject("organization").MustBuild()
	err = l.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	return &pb.UpdateOrgResponse{}, nil
}

func (l *LookupServer) GetOrg(ctx context.Context, req *pb.GetOrgRequest) (*pb.GetOrgResponse, error) {
	logrus.Infof("Get Organization %s", req.OrgName)

	org, err := l.orgRepo.GetByName(req.OrgName)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	resp := &pb.GetOrgResponse{
		OrgName:     req.OrgName,
		Certificate: org.Certificate,
		Ip:          org.Ip.IPNet.IP.String(),
	}

	route := l.baseRoutingKey.SetActionDelete().SetObject("organization").MustBuild()
	err = l.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	return resp, nil
}

func (l *LookupServer) AddNodeForOrg(ctx context.Context, req *pb.AddNodeRequest) (*pb.AddNodeResponse, error) {
	logrus.Infof("Updating node %s for org  %s", req.GetNodeId(), req.GetOrgName())

	id, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid node id  %s. Error %s", req.GetNodeId(), err.Error())
	}

	org, err := l.orgRepo.GetByName(req.OrgName)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid node id  %s. Error %s", req.OrgName, err.Error())
	}

	err = l.nodeRepo.AddOrUpdate(&db.Node{NodeID: id.StringLowercase(), OrgID: org.ID})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Unable to add node id %s to %s org. Error %s",
			req.NodeId, req.OrgName, err.Error())
	}

	route := l.baseRoutingKey.SetActionCreate().SetObject("node").MustBuild()
	err = l.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	return &pb.AddNodeResponse{}, nil
}

func (l *LookupServer) GetNodeForOrg(ctx context.Context, req *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	logrus.Infof("Updating node %s for org  %s", req.GetNodeId(), req.GetOrgName())

	_, err := l.orgRepo.GetByName(req.OrgName)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	nodeId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, invalidNodeIdError(req.NodeId, err)
	}

	node, err := l.nodeRepo.Get(nodeId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	return &pb.GetNodeResponse{
		NodeId:      node.NodeID,
		OrgName:     node.Org.Name,
		Certificate: node.Org.Certificate,
	}, nil
}

func (l *LookupServer) DeleteNodeForOrg(ctx context.Context, req *pb.DeleteNodeRequest) (*pb.DeleteNodeResponse, error) {
	logrus.Infof("Removing node %s from org  %s", req.GetNodeId(), req.GetOrgName())

	nodeId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, invalidNodeIdError(req.NodeId, err)
	}

	err = l.nodeRepo.Delete(nodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid org name %s. Error %s", req.OrgName, err.Error())
	}

	route := l.baseRoutingKey.SetActionDelete().SetObject("node").MustBuild()
	err = l.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}
	return &pb.DeleteNodeResponse{}, nil
}

func (l *LookupServer) GetSystemForOrg(ctx context.Context, req *pb.GetSystemRequest) (*pb.GetSystemResponse, error) {
	logrus.Infof("Requesting System %s info for org  %s", req.GetSystemName(), req.GetOrgName())

	_, err := l.orgRepo.GetByName(req.OrgName)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid node id  %s. Error %s", req.OrgName, err.Error())
	}

	system, err := l.systemRepo.GetByName(req.GetSystemName())
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	return &pb.GetSystemResponse{
		SystemName:  system.Name,
		SystemId:    system.Uuid,
		Certificate: system.Certificate,
		Ip:          system.Ip.IPNet.IP.String(),
		Port:        system.Port,
	}, nil

}

func (l *LookupServer) UpdateSystemForOrg(ctx context.Context, req *pb.UpdateSystemRequest) (*pb.UpdateSystemResponse, error) {
	logrus.Infof("Updating System %s for org  %s", req.GetSystemName(), req.GetOrgName())

	var sysIp pgtype.Inet
	sysId := uuid.New().String()

	org, err := l.orgRepo.GetByName(req.OrgName)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid node id  %s. Error %s", req.OrgName, err.Error())
	}

	system, err := l.systemRepo.GetByName(req.GetSystemName())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Debugf("No system with name %s found. Error %s", req.GetSystemName(), err.Error())
		} else {
			return nil, status.Errorf(codes.Internal, "Something worng with db. Error %s", err.Error())
		}
	}

	if system != nil {
		_, err = uuid.Parse(system.Uuid)
		if err == nil {
			sysId = system.Uuid
		}
	}

	req.SystemId = sysId

	err = sysIp.Set(req.Ip)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid ip for Org %s. Error %s", req.OrgName, err.Error())
	}

	sys := &db.System{
		Name:        strings.ToLower(req.SystemName),
		Certificate: req.Certificate,
		Uuid:        sysId,
		Ip:          sysIp,
		Port:        req.Port,
		OrgID:       org.ID,
	}

	logrus.Debugf("System details: %+v", sys)

	err = l.systemRepo.AddOrUpdate(sys)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Unable to add system %s to %s org. Error %s",
			req.SystemName, req.OrgName, err.Error())
	}

	route := l.baseRoutingKey.SetActionUpdate().SetObject("system").MustBuild()
	err = l.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	return &pb.UpdateSystemResponse{
		SystemId: req.SystemId,
	}, nil

}

func (l *LookupServer) DeleteSystemForOrg(ctx context.Context, req *pb.DeleteSystemRequest) (*pb.DeleteSystemResponse, error) {
	logrus.Infof("Deleting System %s from org  %s", req.GetSystemName(), req.GetOrgName())

	org, err := l.orgRepo.GetByName(req.OrgName)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "System not found. Error %s", req.OrgName, err.Error())
	}

	err = l.systemRepo.Delete(req.SystemName, org.ID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Unable to Delete system %s to %s org. Error %s",
			req.SystemName, req.OrgName, err.Error())
	}

	route := l.baseRoutingKey.SetActionDelete().SetObject("system").MustBuild()
	err = l.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	return &pb.DeleteSystemResponse{}, nil
}

func invalidNodeIdError(nodeId string, err error) error {
	return status.Errorf(codes.InvalidArgument, "invalid node id %s. Error %s", nodeId, err.Error())
}
