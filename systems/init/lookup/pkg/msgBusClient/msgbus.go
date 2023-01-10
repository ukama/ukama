package msgBusClient

import (
	"context"
	"time"

	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/sirupsen/logrus"

	pb "github.com/ukama/ukama/systems/init/msgClient/pb/gen"
	"google.golang.org/grpc"
)

type Registration struct {
}

type MsgBusClient struct {
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

func NewMsgBusClient(timeout time.Duration, system string,
	service string, instanceId string, msgBusURI string,
	serviceURI string, msgClientURI string, exchange string, lq string, pq string, retry int8, routes []string) *MsgBusClient {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, msgClientURI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewMsgClientServiceClient(conn)

	return &MsgBusClient{
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
		routes:       routes,
		listQueue:    lq,
		publQueue:    pq,
		exchange:     exchange,
	}

}

func (m *MsgBusClient) Register() error {

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
		GrpcTimeout: uint32(m.timeout)})
	if err != nil {
		return err
	}

	if resp.GetState() == pb.REGISTRAION_STATUS_REGISTERED {
		m.uuid = resp.ServiceUuid
	}

	return nil
}

func (m *MsgBusClient) Start() error {

	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	_, err := m.client.StartMsgBusHandler(ctx, &pb.StartMsgBusHandlerReq{
		ServiceUuid: m.uuid})
	if err != nil {
		return err
	}

	return nil
}

func (m *MsgBusClient) Stop() error {

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

func (m *MsgBusClient) PublishRequest(route string, msg protoreflect.ProtoMessage) error {

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
	logrus.Debugf("Published:\n Message: %v  \n Key: %s \n ", msg, route)
	return nil

}
