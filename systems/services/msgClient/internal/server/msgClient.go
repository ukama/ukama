package server

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/services/msgClient/internal/db"
	"github.com/ukama/ukama/systems/services/msgClient/internal/queue"
	pb "github.com/ukama/ukama/systems/services/msgClient/pb/gen"
)

type MsgClientServer struct {
	sys string
	s   db.ServiceRepo
	r   db.RouteRepo
	h   queue.MsgBusHandlerInterface

	pb.UnimplementedMsgClientServiceServer
}

func NewMsgClientServer(serviceRepo db.ServiceRepo, keyRepo db.RouteRepo, h queue.MsgBusHandlerInterface, sys string) *MsgClientServer {
	return &MsgClientServer{
		sys: sys,
		s:   serviceRepo,
		r:   keyRepo,
		h:   h,
	}
}

func (m *MsgClientServer) RegisterService(ctx context.Context, req *pb.RegisterServiceReq) (*pb.RegisterServiceResp, error) {
	log.Debugf("Register new listener request for %s", req.ServiceName)
	/* This sholuld be handled as db tx but for now we have two seperate commits one for route
	other for service */
	resp := &pb.RegisterServiceResp{
		State: pb.REGISTRAION_STATUS_NOT_REGISTERED,
	}

	if !strings.EqualFold(m.sys, req.SystemName) {
		return nil, fmt.Errorf("invalid system name %s in request", req.SystemName)
	}
	/* Register service */
	svc := db.Service{
		Name:        req.ServiceName,
		InstanceId:  req.InstanceId,
		ServiceUri:  req.ServiceURI,
		MsgBusUri:   req.MsgBusURI,
		ListQueue:   req.ListQueue,
		PublQueue:   req.PublQueue,
		Exchange:    req.Exchange,
		GrpcTimeout: req.GrpcTimeout,
	}

	service, err := m.s.Register(&svc)
	if err != nil {
		log.Errorf("Failed to register service %s", req.ServiceName)
		return resp, err
	}

	log.Debugf("Removing old route for %s service", service.Name)
	err = m.s.RemoveRoutes(service)
	if err != nil {
		log.Errorf("Failed to remove old routes for service %s. Error %s", req.ServiceName, err.Error())
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

		log.Debugf("Adding route %s for %s service", r, service.Name)
		err = m.s.AddRoute(service, rt)
		if err != nil {
			/* No need to rollback the already added routes.*/
			log.Errorf("Failed to add route %s for service %s. Error %s", r, req.ServiceName, err.Error())
			return resp, err
		}
	}

	resp.State = pb.REGISTRAION_STATUS_REGISTERED
	resp.ServiceUuid = service.ServiceUuid
	return resp, nil
}

func (m *MsgClientServer) StartMsgBusHandler(ctx context.Context, req *pb.StartMsgBusHandlerReq) (*pb.StartMsgBusHandlerResp, error) {
	log.Debugf("Start handler request for %s", req.ServiceUuid)

	svc, err := m.s.Get(req.ServiceUuid)
	if err != nil {
		log.Errorf("Failed to get listener config for %s", req.ServiceUuid)
		return nil, err
	}

	/* Update Service handler for message queue */
	err = m.h.UpdateServiceQueueHandler(svc)
	if err != nil {
		log.Errorf("Failed to start listener for service %s. Error %s", svc.Name, err.Error())
		return nil, err
	}

	return &pb.StartMsgBusHandlerResp{}, nil
}

func (m *MsgClientServer) StopMsgBusHandler(ctx context.Context, req *pb.StopMsgBusHandlerReq) (*pb.StopMsgBusHandlerResp, error) {

	log.Debugf("Stop handler request for %s", req.ServiceUuid)
	/* start listening */
	err := m.h.StopServiceQueueHandler(req.ServiceUuid)
	if err != nil {
		log.Errorf("Failed to stop listener for service %s. Error %s", req.ServiceUuid, err.Error())
		return nil, err
	}

	return &pb.StopMsgBusHandlerResp{}, nil
}

func (m *MsgClientServer) UnregisterService(ctx context.Context, req *pb.UnregisterServiceReq) (*pb.UnregisterServiceResp, error) {

	log.Debugf("Remove handler request for %s", req.ServiceUuid)

	/* Listener */
	err := m.h.RemoveServiceQueueListening(req.ServiceUuid)
	if err != nil {
		return nil, err
	}

	/* Publisher */
	err = m.h.RemoveServiceQueuePublisher(req.ServiceUuid)
	if err != nil {
		return nil, err
	}

	err = m.s.UnRegister(req.ServiceUuid)
	if err != nil {
		return nil, err
	}
	log.Debugf("listener and publisher removed for service %s", req.ServiceUuid)

	return &pb.UnregisterServiceResp{}, nil
}

func (m *MsgClientServer) PublishMsg(ctx context.Context, req *pb.PublishMsgRequest) (*pb.PublishMsgResponse, error) {
	log.Debugf("Publish request for %s service", req.ServiceUuid)

	err := m.h.Publish(req.ServiceUuid, req.RoutingKey, req.Msg)
	if err != nil {
		return nil, err
	}
	return &pb.PublishMsgResponse{}, nil
}
