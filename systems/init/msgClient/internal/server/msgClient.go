package server

import (
	"context"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
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

func (m *MsgClientServer) RegisterService(ctx context.Context, req *pb.RegisterServiceReq) (*pb.RegisterServiceResp, error) {
	log.Debugf("Register new listener request for %s", req.ServiceName)
	/* This sholuld be handled as db tx but for now we have two seperate commits one for route
	other for service */
	resp := &pb.RegisterServiceResp{
		State: pb.REGISTRAION_STATUS_NOT_REGISTERED,
	}

	/* Generate UUID */
	id := uuid.NewV4()
	// if err != nil {
	// 	log.Errorf("Service Id genration failed for service %s. Error %s", req.ServiceName, err.Error())
	// 	return &pb.RegisterServiceResp{
	// 		State: pb.REGISTRAION_STATUS_NOT_REGISTERED,
	// 	}, err
	// }

	/* Register service */
	svc := db.Service{
		Name:        req.ServiceName,
		ServiceId:   id.String(),
		ServiceUri:  req.ServiceURI,
		MsgBusUri:   req.MsgBusURI,
		QueueName:   req.QueueName,
		Exchange:    req.Exchange,
		GrpcTimeout: req.GrpcTimeout,
	}

	err := m.s.Register(&svc)
	if err != nil {
		log.Errorf("Failed to register service %s", req.ServiceName)
		return resp, err
	}

	/* Add Routes */
	routes := make([]db.Route, len(req.Routes))
	for i, r := range req.Routes {
		routes[i].Key = r
		rt, err := m.r.Add(r)
		if err != nil {
			/* No need to rollback the already added routes.*/
			log.Errorf("Failed to add route %s for service %s. Error %s", r, req.ServiceName, err.Error())
			return resp, err
		}

		log.Debugf("Adding route %s for %s service", r, svc.Name)
		err = m.s.AddRoute(&svc, rt)
		if err != nil {
			/* No need to rollback the already added routes.*/
			log.Errorf("Failed to add route %s for service %s. Error %s", r, req.ServiceName, err.Error())
			return resp, err
		}
	}

	resp.State = pb.REGISTRAION_STATUS_REGISTERED
	resp.ServiceId = svc.ServiceId
	return resp, nil
}

func (m *MsgClientServer) StartListening(ctx context.Context, req *pb.StartListeningReq) (*pb.StartListeningResp, error) {
	log.Debugf("Start listener request for %s", req.ServiceId)

	svc, err := m.s.Get(req.ServiceId)
	if err != nil {
		log.Errorf("Failed to get listener config for %s", req.ServiceId)
		return nil, err
	}

	/* start listening */
	err = m.mq.UpdateServiceQueueListening(svc)
	if err != nil {
		log.Errorf("Failed to start listener for service %s. Error %s", svc.Name, err.Error())
		return nil, err
	}

	return &pb.StartListeningResp{}, nil
}

func (m *MsgClientServer) StopListening(ctx context.Context, req *pb.StopListeningReq) (*pb.StopListeningResp, error) {

	log.Debugf("Stop listener request for %s", req.ServiceId)
	/* start listening */
	err := m.mq.StopServiceQueueListening(req.ServiceId)
	if err != nil {
		log.Errorf("Failed to stop listener for service %s. Error %s", req.ServiceId, err.Error())
		return nil, err
	}

	return &pb.StopListeningResp{}, nil
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
