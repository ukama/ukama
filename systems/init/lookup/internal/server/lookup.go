package server

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/init/lookup/internal"
	"github.com/ukama/ukama/systems/init/lookup/internal/db"
	pb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LookupServer struct {
	systemRepo     db.SystemRepo
	orgRepo        db.OrgRepo
	nodeRepo       db.NodeRepo
	msgbus         *mb.MsgBusClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedLookupServiceServer
}

func NewLookupServer(nodeRepo db.NodeRepo, orgRepo db.OrgRepo, systemRepo db.SystemRepo, msgBus *mb.MsgBusClient) *LookupServer {
	return &LookupServer{
		nodeRepo:       nodeRepo,
		orgRepo:        orgRepo,
		systemRepo:     systemRepo,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(internal.ServiceName),
	}
}

func (l *LookupServer) getOrg(name string) (*db.Org, error) {
	org, err := l.orgRepo.GetByName(name)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	return org, nil
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

	err = l.orgRepo.Add(org)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	route := l.baseRoutingKey.SetAction("create").SetObject("organization").MustBuild()
	err = l.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	dbOrg, err := l.getOrg(org.Name)
	if err != nil {
		return nil, err
	}

	return &pb.AddOrgResponse{
		OrgName:     dbOrg.Name,
		Certificate: dbOrg.Certificate,
		Ip:          dbOrg.Ip.IPNet.String(),
	}, nil
}

func (l *LookupServer) UpdateOrg(ctx context.Context, req *pb.UpdateOrgRequest) (*pb.UpdateOrgResponse, error) {
	logrus.Infof("Updating Organization %s", req.OrgName)

	org := &db.Org{
		Name:        req.GetOrgName(),
		Certificate: req.GetCertificate(),
	}

	_, err := l.orgRepo.GetByName(req.OrgName)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	if req.Ip != "" {
		err := org.Ip.Set(req.Ip)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid ip for Org %s. Error %s", req.OrgName, err.Error())
		}
	}

	err = l.orgRepo.Update(org)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	route := l.baseRoutingKey.SetActionUpdate().SetObject("organization").MustBuild()
	err = l.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	dbOrg, err := l.getOrg(org.Name)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateOrgResponse{
		OrgName:     dbOrg.Name,
		Certificate: dbOrg.Certificate,
		Ip:          dbOrg.Ip.IPNet.String(),
	}, nil

}

func (l *LookupServer) GetOrg(ctx context.Context, req *pb.GetOrgRequest) (*pb.GetOrgResponse, error) {
	logrus.Infof("Get Organization %s", req.OrgName)

	dbOrg, err := l.getOrg(req.OrgName)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	return &pb.GetOrgResponse{
		OrgName:     dbOrg.Name,
		Certificate: dbOrg.Certificate,
		Ip:          dbOrg.Ip.IPNet.String(),
	}, nil
}

func (l *LookupServer) AddNodeForOrg(ctx context.Context, req *pb.AddNodeRequest) (*pb.AddNodeResponse, error) {
	logrus.Infof("Updating node %s for org  %s", req.GetNodeId(), req.GetOrgName())

	id, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid node id  %s. Error %s", req.GetNodeId(), err.Error())
	}

	org, err := l.orgRepo.GetByName(req.OrgName)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	err = l.nodeRepo.AddOrUpdate(&db.Node{NodeID: id.StringLowercase(), OrgID: org.ID})
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	route := l.baseRoutingKey.SetAction("create").SetObject("node").MustBuild()
	err = l.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	dbNode, err := l.nodeRepo.Get(id)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	return &pb.AddNodeResponse{
		NodeId:  dbNode.NodeID,
		OrgName: dbNode.Org.Name,
	}, nil
}

func (l *LookupServer) GetNode(ctx context.Context, req *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	logrus.Infof("Get node %s.", req.GetNodeId())

	nodeId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, invalidNodeIdError(req.NodeId, err)
	}

	node, err := l.nodeRepo.Get(nodeId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	resp, err := &pb.GetNodeResponse{
		NodeId:      node.NodeID,
		OrgName:     node.Org.Name,
		Certificate: node.Org.Certificate,
		Ip:          node.Org.Ip.IPNet.String(),
	}, nil

	return resp, err
}

func (l *LookupServer) GetNodeForOrg(ctx context.Context, req *pb.GetNodeForOrgRequest) (*pb.GetNodeResponse, error) {
	logrus.Infof("Get node %s for org %s.", req.GetNodeId(), req.GetOrgName())

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
		return nil, grpc.SqlErrorToGrpc(err, "node")
	}

	return &pb.GetNodeResponse{
		NodeId:      node.NodeID,
		OrgName:     node.Org.Name,
		Certificate: node.Org.Certificate,
		Ip:          node.Org.Ip.IPNet.String(),
	}, nil
}

func (l *LookupServer) DeleteNodeForOrg(ctx context.Context, req *pb.DeleteNodeRequest) (*pb.DeleteNodeResponse, error) {
	logrus.Infof("Removing node %s from org  %s", req.GetNodeId(), req.GetOrgName())

	nodeId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, invalidNodeIdError(req.NodeId, err)
	}

	_, err = l.orgRepo.GetByName(req.OrgName)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	err = l.nodeRepo.Delete(nodeId)
	if err != nil {
		if strings.Contains(err.Error(), "node missing") {
			return nil, status.Errorf(codes.NotFound, "Unable to Delete node %s for %s org. Error %s",
				req.NodeId, req.OrgName, err.Error())
		} else {
			return nil, status.Errorf(codes.InvalidArgument, "Unable to Delete node %s for %s org. Error %s",
				req.NodeId, req.OrgName, err.Error())
		}
	}

	route := l.baseRoutingKey.SetActionDelete().SetObject("node").MustBuild()
	err = l.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}
	return &pb.DeleteNodeResponse{}, nil
}

func (l *LookupServer) getSystem(name string) (*db.System, error) {

	system, err := l.systemRepo.GetByName(name)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "system")
	}
	return system, nil

}

func (l *LookupServer) GetSystemForOrg(ctx context.Context, req *pb.GetSystemRequest) (*pb.GetSystemResponse, error) {
	logrus.Infof("Requesting System %s info for org  %s", req.GetSystemName(), req.GetOrgName())
	org := req.GetOrgName()

	_, err := l.orgRepo.GetByName(org)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	system, err := l.getSystem(req.GetSystemName())
	if err != nil {
		return nil, err
	}

	return &pb.GetSystemResponse{
		SystemName:  system.Name,
		SystemId:    system.Uuid,
		Certificate: system.Certificate,
		Ip:          system.Ip.IPNet.IP.String(),
		Port:        system.Port,
		Health:      system.Health,
	}, nil

}

func (l *LookupServer) AddSystemForOrg(ctx context.Context, req *pb.AddSystemRequest) (*pb.AddSystemResponse, error) {
	logrus.Infof("Adding system %s for org  %s", req.GetSystemName(), req.GetOrgName())

	var sysIp pgtype.Inet
	sysId := uuid.New().String()

	org, err := l.orgRepo.GetByName(req.OrgName)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	err = sysIp.Set(req.Ip)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid ip for system %s. Error %s", req.OrgName, err.Error())
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

	err = l.systemRepo.Add(sys)
	if err != nil {
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "system")
		}
	}

	route := l.baseRoutingKey.SetAction("create").SetObject("system").MustBuild()
	err = l.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	resp, err := l.getSystem(sys.Name)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "system")
	}

	return &pb.AddSystemResponse{
		SystemName:  resp.Name,
		SystemId:    resp.Uuid,
		Certificate: resp.Certificate,
		Ip:          resp.Ip.IPNet.IP.String(),
		Port:        resp.Port,
	}, nil
}

func (l *LookupServer) UpdateSystemForOrg(ctx context.Context, req *pb.UpdateSystemRequest) (*pb.UpdateSystemResponse, error) {
	logrus.Infof("Updating system %s for org  %s", req.GetSystemName(), req.GetOrgName())

	_, err := l.orgRepo.GetByName(req.OrgName)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	_, err = l.systemRepo.GetByName(req.SystemName)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "system")
	}

	sys := &db.System{
		Name:        strings.ToLower(req.SystemName),
		Certificate: req.Certificate,
		Port:        req.Port,
	}

	if req.Ip != "" {
		err = sys.Ip.Set(req.Ip)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid ip for Org %s. Error %s", req.OrgName, err.Error())
		}
	}

	logrus.Debugf("System details: %+v", sys)

	err = l.systemRepo.Update(sys)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Unable to add system %s to %s org. Error %s",
			req.SystemName, req.OrgName, err.Error())
	}

	route := l.baseRoutingKey.SetActionUpdate().SetObject("system").MustBuild()
	err = l.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	dbSystem, err := l.getSystem(sys.Name)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "system")
	}

	return &pb.UpdateSystemResponse{
		SystemName:  dbSystem.Name,
		SystemId:    dbSystem.Uuid,
		Certificate: dbSystem.Certificate,
		Ip:          dbSystem.Ip.IPNet.IP.String(),
		Port:        dbSystem.Port,
	}, nil
}

func (l *LookupServer) DeleteSystemForOrg(ctx context.Context, req *pb.DeleteSystemRequest) (*pb.DeleteSystemResponse, error) {
	logrus.Infof("Deleting System %s from org  %s", req.GetSystemName(), req.GetOrgName())

	_, err := l.orgRepo.GetByName(req.OrgName)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "org")
	}

	err = l.systemRepo.Delete(req.SystemName)
	if err != nil {
		if strings.Contains(err.Error(), "system missing") {
			return nil, status.Errorf(codes.NotFound, "Unable to Delete system %s from %s org. Error %s",
				req.SystemName, req.OrgName, err.Error())
		} else {
			return nil, status.Errorf(codes.InvalidArgument, "Unable to Delete system %s from %s org. Error %s",
				req.SystemName, req.OrgName, err.Error())
		}
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
