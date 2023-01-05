package server

import (
	"context"

	"github.com/ukama/ukama/systems/init/msgClient/internal/db"
	"github.com/ukama/ukama/systems/init/msgClient/internal/queue"
	pb "github.com/ukama/ukama/systems/init/msgClient/pb/gen"
)

type MsgClientServer struct {
	s db.ServiceRepo
	r db.RoutingKeyRepo
	l *queue.QueueListener

	pb.UnimplementedMsgClientServiceServer
}

func NewMsgClientServer(serviceRepo db.ServiceRepo, keyRepo db.RoutingKeyRepo, l *queue.QueueListener) *MsgClientServer {
	return &MsgClientServer{
		s: serviceRepo,
		r: keyRepo,
		l: l,
	}
}

func (m *MsgClientServer) RegisterService(context.Context, *pb.RegisterServiceReq) (*pb.RegisterServiceResp, error) {
	/* Add service */

	return nil, nil
}

func (m *MsgClientServer) RegisterRoutes(context.Context, *pb.RegisterRoutesReq) (*pb.RegisterRoutesResp, error) {
	/* Add a route and serviceID */

	/* Restart listener */
	err := m.l.RetstartListening()
	return nil, err
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
