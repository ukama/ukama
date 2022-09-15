package msgBusClient

import (
	"context"
	"time"

	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/sirupsen/logrus"

	"github.com/ukama/ukama/services/common/msgbus"
	pb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
	"google.golang.org/grpc"
)

type Registration struct {
}

type MsgBusClient struct {
	name           string
	instanceId     string
	msgBusURI      string
	queuePub       msgbus.QPub
	baseRoutingKey msgbus.RoutingKeyBuilder
	conn           *grpc.ClientConn
	client         pb.MsgBusClient
	timeout        time.Duration
	host           string
}

func NewMsgBusClient(timeout time.Duration, name string, instanceId string, msgBusURI string) *MsgBusClient {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, msgBusURI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewMsgBusClient(conn)

	return &MsgBusClient{
		name:       name,
		instanceId: instanceId,
		msgBusURI:  msgBusURI,
		conn:       conn,
		client:     client,
	}

}

func (m *MsgBusClient) Init() error {

	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	_, err := m.client.Initialize(ctx, &pb.InitRequest{ServiceName: m.name, QueueURI: m.msgBusURI, InstanceId: m.instanceId})
	if err != nil {
		return err
	}

	return nil
}

// TODO: Check if we require specific function for each Message type.
func (m *MsgBusClient) PublishRequest(route string, msg protoreflect.ProtoMessage) error {

	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	anyMsg, err := anypb.New(msg)
	if err != nil {
		return err
	}

	_, err = m.client.PusblishMsg(ctx, &pb.PublishMsgRequest{RoutingKey: route, Msg: anyMsg})
	if err != nil {
		return err
	}

	return nil

}
