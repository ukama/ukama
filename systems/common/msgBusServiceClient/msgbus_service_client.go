package msgBusServiceClient

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"

	log "github.com/sirupsen/logrus"

	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/services/msgClient/pb/gen"
	"google.golang.org/grpc"
)

type MsgBusServiceClient interface {
	Register() error
	Start() error
	Stop() error
	PublishRequest(route string, msg protoreflect.ProtoMessage) error
}

type msgBusServiceClient struct {
	org          string
	uuid         string
	service      string
	system       string
	instanceId   string
	msgBusURI    string
	msgClientURI string
	exchange     string
	listQueue    string
	publQueue    string
	conn         *grpc.ClientConn
	client       pb.MsgClientServiceClient
	timeout      time.Duration
	host         string
	retry        int8
	routes       []string
}

func NewMsgBusClient(timeout time.Duration, org string, system string,
	service string, instanceId string, msgBusURI string,
	serviceURI string, msgClientURI string, exchange string, lq string, pq string, retry int8, routes []string) *msgBusServiceClient {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, msgClientURI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := pb.NewMsgClientServiceClient(conn)

	return &msgBusServiceClient{
		org:          org,
		service:      service,
		system:       system,
		instanceId:   instanceId,
		msgBusURI:    msgBusURI,
		msgClientURI: msgClientURI,
		conn:         conn,
		client:       client,
		timeout:      timeout,
		retry:        retry,
		host:         serviceURI,
		routes:       msgbus.PrepareRoutes(org, routes),
		listQueue:    lq,
		publQueue:    pq,
		exchange:     exchange,
	}

}

func (m *msgBusServiceClient) Register() error {
	log.Debugf("Registering %s service instance %s with routes %+v to MessageBusClient at %s.", m.service, m.instanceId, m.routes, m.msgClientURI)
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	resp, err := m.client.RegisterService(ctx, &pb.RegisterServiceReq{
		ServiceName: m.service,
		SystemName:  m.system,
		MsgBusURI:   m.msgBusURI,
		ServiceURI:  m.host,
		InstanceId:  m.instanceId,
		Routes:      m.routes,
		ListQueue:   m.listQueue,
		PublQueue:   m.publQueue,
		Exchange:    m.exchange,
		GrpcTimeout: uint32(m.timeout.Seconds())})
	if err != nil {
		return err
	}

	if resp.GetState() == pb.REGISTRAION_STATUS_REGISTERED {
		m.uuid = resp.ServiceUuid
	} else {
		return fmt.Errorf("failed to register %s service instance %s: %s", m.service, m.instanceId, resp.State.String())
	}

	log.Infof("%s service instance %s to MessageBusClient at %s.", m.service, m.instanceId, resp.State.String())
	return nil
}

func (m *msgBusServiceClient) Start() error {
	log.Debugf("Starting MessageClientRoutine for %s service instance %s Routine ID %s.", m.service, m.instanceId, m.uuid)
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	_, err := m.client.StartMsgBusHandler(ctx, &pb.StartMsgBusHandlerReq{
		ServiceUuid: m.uuid})
	if err != nil {
		return err
	}

	msg := &epb.PublishServiceStatusUp{
		OrgName:  m.org,
		System:   m.system,
		Service:  m.service,
		Instance: m.instanceId,
	}

	route := msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(m.system).SetOrgName(m.org).SetService(m.service).SetAction("up").SetObject("instance").MustBuild()
	log.Debugf("Publishing service up message on bus: %v", msg)
	err = m.PublishRequest(route, msg)
	if err != nil {
		log.Warningf("Failed to publish message on bus: %v", err)
	}

	return nil
}

func (m *msgBusServiceClient) Stop() error {
	log.Debugf("Stopping MessageClientRoutine for %s service instance %s Routine ID %s.", m.service, m.instanceId, m.uuid)
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	_, err := m.client.StopMsgBusHandler(ctx, &pb.StopMsgBusHandlerReq{
		ServiceUuid: m.uuid,
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *msgBusServiceClient) PublishRequest(route string, msg protoreflect.ProtoMessage) error {
	log.Debugf("Publishing message %s to MessageClientRoutine for %s service instance %s Routine ID %s", route, m.service, m.instanceId, m.uuid)
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	anyMsg, err := anypb.New(msg)
	if err != nil {
		return err
	}

	_, err = m.client.PublishMsg(ctx, &pb.PublishMsgRequest{
		ServiceUuid: m.uuid,
		RoutingKey:  route,
		Msg:         anyMsg})
	if err != nil {
		return err
	}
	log.Debugf("Published:\n Message: %+v  \n Key: %s \n ", msg, route)
	return nil

}
