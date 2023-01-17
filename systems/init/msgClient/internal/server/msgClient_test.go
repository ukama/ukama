package server

import (
	"context"
	"testing"

	"github.com/ukama/ukama/systems/init/msgClient/internal"
	"github.com/ukama/ukama/systems/init/msgClient/internal/db"
	mocks "github.com/ukama/ukama/systems/init/msgClient/mocks"
	pb "github.com/ukama/ukama/systems/init/msgClient/pb/gen"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/stretchr/testify/assert"
)

var healthCheck internal.HeathCheckRoutine
var route1 = db.Route{
	Key: "event.cloud.lookup.organization.create",
}

var ServiceUuid = "1ce2fa2f-2997-422c-83bf-92cf2e7334dd"
var service1 = db.Service{
	Name:        "test",
	InstanceId:  "1",
	MsgBusUri:   "amqp://guest:guest@localhost:5672",
	ListQueue:   "",
	PublQueue:   "",
	Exchange:    "amq.topic",
	ServiceUri:  "localhost:9095",
	GrpcTimeout: 5,
}

func TestLookupServer_RegisterService(t *testing.T) {
	serviceRepo := &mocks.ServiceRepo{}
	routeRepo := &mocks.RouteRepo{}

	rt := route1
	svc := service1

	reqPb := pb.RegisterServiceReq{
		SystemName:  internal.SystemName,
		ServiceName: service1.Name,
		Exchange:    service1.Exchange,
		InstanceId:  service1.InstanceId,
		MsgBusURI:   service1.MsgBusUri,
		ListQueue:   service1.ListQueue,
		PublQueue:   service1.PublQueue,
		ServiceURI:  service1.ServiceUri,
		GrpcTimeout: service1.GrpcTimeout,
		Routes:      []string{route1.Key},
	}

	serviceRepo.On("Register", &service1).Return(&svc, nil).Once()
	serviceRepo.On("RemoveRoutes", &service1).Return(nil).Once()
	routeRepo.On("Add", route1.Key).Return(&rt, nil).Once()
	serviceRepo.On("AddRoute", &svc, &rt).Return(nil).Once()

	s := NewMsgClientServer(serviceRepo, routeRepo, nil)
	_, err := s.RegisterService(context.TODO(), &reqPb)

	assert.NoError(t, err)
	serviceRepo.AssertExpectations(t)
	routeRepo.AssertExpectations(t)
}

func TestLookupServer_StartMsgHandler(t *testing.T) {
	serviceRepo := &mocks.ServiceRepo{}
	routeRepo := &mocks.RouteRepo{}
	msgIf := &mocks.MsgBusHandlerInterface{}

	svc := service1
	svc.ServiceUuid = ServiceUuid
	svc.Routes = []db.Route{route1}

	reqStartPb := pb.StartMsgBusHandlerReq{
		ServiceUuid: ServiceUuid,
	}

	serviceRepo.On("Get", ServiceUuid).Return(&svc, nil).Once()
	msgIf.On("UpdateServiceQueueHandler", &svc).Return(nil).Once()

	s := NewMsgClientServer(serviceRepo, routeRepo, msgIf)
	_, err := s.StartMsgBusHandler(context.TODO(), &reqStartPb)

	assert.NoError(t, err)
	serviceRepo.AssertExpectations(t)
}

func TestLookupServer_StoptMsgHandler(t *testing.T) {
	serviceRepo := &mocks.ServiceRepo{}
	routeRepo := &mocks.RouteRepo{}
	msgIf := &mocks.MsgBusHandlerInterface{}

	reqStopPb := pb.StopMsgBusHandlerReq{
		ServiceUuid: ServiceUuid,
	}

	msgIf.On("StopServiceQueueHandler", reqStopPb.ServiceUuid).Return(nil).Once()

	s := NewMsgClientServer(serviceRepo, routeRepo, msgIf)
	_, err := s.StopMsgBusHandler(context.TODO(), &reqStopPb)

	assert.NoError(t, err)
	msgIf.AssertExpectations(t)
}

func TestLookupServer_Publish(t *testing.T) {
	serviceRepo := &mocks.ServiceRepo{}
	routeRepo := &mocks.RouteRepo{}
	msgIf := &mocks.MsgBusHandlerInterface{}

	svc := service1
	svc.ServiceUuid = ServiceUuid
	svc.Routes = []db.Route{route1}

	reqMsg := pb.PublishMsgRequest{
		ServiceUuid: ServiceUuid,
		RoutingKey:  route1.Key,
		Msg:         &anypb.Any{},
	}

	msgIf.On("Publish", reqMsg.ServiceUuid, reqMsg.RoutingKey, reqMsg.Msg).Return(nil).Once()

	s := NewMsgClientServer(serviceRepo, routeRepo, msgIf)
	_, err := s.PublishMsg(context.TODO(), &reqMsg)

	assert.NoError(t, err)
	msgIf.AssertExpectations(t)
}
