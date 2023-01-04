package server

import (
	"context"

	"github.com/ukama/ukama/systems/init/msgClient/internal/db"
	pb "github.com/ukama/ukama/systems/init/msgClient/pb/gen"
)

type MsgClientServer struct {
	serviceRepo    db.ServiceRepo
	routingKeyRepo db.RoutingKeyRepo

	pb.UnimplementedMsgClientServiceServer
}

func NewMsgClientServer(serviceRepo db.ServiceRepo, keyRepo db.RoutingKeyRepo) *MsgClientServer {
	return &MsgClientServer{
		serviceRepo:    serviceRepo,
		routingKeyRepo: keyRepo,
	}
}

func (m *MsgClientServer) RegisterService(context.Context, *pb.RegisterServiceReq) (*pb.RegisterServiceResp, error) {
	return nil, nil
}

func (m *MsgClientServer) RegisterRoutes(context.Context, *pb.RegisterRoutesReq) (*pb.RegisterRoutesResp, error) {
	return nil, nil
}

func (m *MsgClientServer) UnregisterService(context.Context, *pb.UnregisterServiceReq) (*pb.UnregisterServiceResp, error) {
	return nil, nil
}

func (m *MsgClientServer) UnregisterRoutes(context.Context, *pb.UnregisterRoutesReq) (*pb.UnregisterRoutesResp, error) {
	return nil, nil
}

func (m *MsgClientServer) PusblishMsg(context.Context, *pb.PublishMsgRequest) (*pb.PublishMsgResponse, error) {
	return nil, nil
}
