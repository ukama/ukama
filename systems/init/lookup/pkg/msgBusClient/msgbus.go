package msgBusClient

import (
	"context"
	"time"

	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/sirupsen/logrus"

	pb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
	"google.golang.org/grpc"
)

type Registration struct {
}

type MsgBusClient struct {
	service    string
	system     string
	instanceId string
	msgBusURI  string
	//queuePub       msgbus.QPub
	//baseRoutingKey msgbus.RoutingKeyBuilder
	conn    *grpc.ClientConn
	client  pb.MsgBusClient
	timeout time.Duration
	host    string
	retry   int8
	routes  []string
}

func NewMsgBusClient(timeout time.Duration, system string,
	service string, instanceId string, msgBusURI string,
	msgClientURI string, retry int8, routes []string) *MsgBusClient {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, msgClientURI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewMsgBusClient(conn)

	return &MsgBusClient{
		service:    service,
		system:     system,
		instanceId: instanceId,
		msgBusURI:  msgBusURI,
		conn:       conn,
		client:     client,
		timeout:    timeout,
		retry:      retry,
		host:       msgClientURI,
	}

}

func (m *MsgBusClient) Init() error {

	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	_, err := m.client.Initialize(ctx, &pb.InitRequest{
		ServiceName: m.service,
		SystemName:  m.system,
		QueueURI:    m.msgBusURI,
		InstanceId:  m.instanceId,
		Routes:      m.routes,
	})
	if err != nil {
		return err
	}

	return nil
}

// TODO: Check if we require specific function for each Message type.
func (m *MsgBusClient) PublishRequest(route string, msg protoreflect.ProtoMessage) error {

	// ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	// defer cancel()

	// anyMsg, err := anypb.New(msg)
	// if err != nil {
	// 	return err
	// }

	// _, err = m.client.PusblishMsg(ctx, &pb.PublishMsgRequest{RoutingKey: route, Msg: anyMsg})
	// if err != nil {
	// 	return err
	// }
	logrus.Debugf("Published:\n Message: %v  \n Key: %s \n ", msg, route)
	return nil

}
