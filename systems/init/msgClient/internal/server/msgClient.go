package server

import (
	"context"

	"github.com/ukama/ukama/systems/init/msgClient/internal/db"
	"github.com/ukama/ukama/systems/init/msgClient/internal/queue"
	pb "github.com/ukama/ukama/systems/init/msgClient/pb/gen"
)

type MsgClientServer struct {
	s  db.ServiceRepo
	r  db.RouteRepo
	mq *queue.MsgBusListener

	pb.UnimplementedMsgClientServiceServer
}

func NewMsgClientServer(serviceRepo db.ServiceRepo, keyRepo db.RouteRepo, mq *queue.MsgBusListener) *MsgClientServer {
	return &MsgClientServer{
		s:  serviceRepo,
		r:  keyRepo,
		mq: mq,
	}
}

func (m *MsgClientServer) RegisterService(context.Context, *pb.RegisterServiceReq) (*pb.RegisterServiceResp, error) {
	/* Add service */

	return nil, nil
}

func (m *MsgClientServer) RegisterRoutes(context.Context, *pb.RegisterRoutesReq) (*pb.RegisterRoutesResp, error) {
	/* Add a route and serviceID */

	/* Restart listener */
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
