package server

import (
	"context"

	"github.com/ukama/ukama/systems/init/msgClient/internal/db"
	"github.com/ukama/ukama/systems/init/msgClient/internal/queue"
	pb "github.com/ukama/ukama/systems/init/msgClient/pb/gen"
)

type MsgClientServer struct {
	serviceRepo    db.ServiceRepo
	routingKeyRepo db.RoutingKeyRepo
	listner        *queue.QueueListener

	pb.UnimplementedMsgClientServiceServer
}

func NewMsgClientServer(serviceRepo db.ServiceRepo, keyRepo db.RoutingKeyRepo, l *queue.QueueListener) *MsgClientServer {
	return &MsgClientServer{
		serviceRepo:    serviceRepo,
		routingKeyRepo: keyRepo,
		listner:        l,
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
