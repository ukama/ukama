package server

import (
	"context"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/messaging/nns/pb/gen"
	"github.com/ukama/ukama/systems/messaging/nns/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type NnsEventServer struct {
	GrpcTimeout time.Duration
	NnsConn     *grpc.ClientConn
	NnsClient   pb.NnsClient
	Registry    client.NodeRegistryClient
	epb.UnimplementedEventNotificationServiceServer
}

func NewNnsEventServer(c client.NodeRegistryClient, nns string, t time.Duration) *NnsEventServer {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()
	conn, err := grpc.DialContext(ctx, nns, grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("Could not connect: %v", err)
	}
	client := pb.NewNnsClient(conn)

	return &NnsEventServer{
		GrpcTimeout: t,
		Registry:    c,
		NnsConn:     conn,
		NnsClient:   client,
	}
}

func (l *NnsEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case "event.cloud.mesh.node.online":
		msg, err := l.unmarshalNodeOnlineEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = l.handleEventNodeOnline(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case "event.cloud.mesh.node.offline":
		msg, err := l.unmarshalNodeOfflineEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = l.handleEventNodeOffline(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (l *NnsEventServer) unmarshalNodeOnlineEvent(msg *anypb.Any) (*epb.NodeOnlineEvent, error) {
	p := &epb.NodeOnlineEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddOrgRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (l *NnsEventServer) handleEventNodeOnline(key string, msg *epb.NodeOnlineEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	log.Infof("Getting org and network for %s", msg.GetNodeId())
	nodeInfo, err := l.Registry.GetNode(msg.GetNodeId())
	if err != nil {
		log.Errorf("Failed to get org and network. Error: %+v", err)
		log.Warningf("Node id %s won't have org and network info", msg.GetNodeId())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(l.GrpcTimeout))
	defer cancel()
	_, err = l.NnsClient.Set(ctx, &pb.SetRequest{
		NodeId:   msg.GetNodeId(),
		NodeIp:   msg.GetMeshIp(),
		MeshIp:   msg.GetMeshIp(),
		NodePort: msg.GetNodePort(),
		MeshPort: msg.GetMeshPort(),
		Org:      nodeInfo.Org,
		Network:  nodeInfo.Network,
	}, grpc_retry.WithMax(3))

	if err != nil {
		log.Errorf("Failed to set node IP. Error: %+v", err)
		return err
	}
	log.Infof("Node %s IP set to %s", msg.GetNodeId(), msg.GetMeshIp())

	return nil
}

func (l *NnsEventServer) unmarshalNodeOfflineEvent(msg *anypb.Any) (*epb.NodeOfflineEvent, error) {
	p := &epb.NodeOfflineEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddOrgRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (l *NnsEventServer) handleEventNodeOffline(key string, msg *epb.NodeOfflineEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)
	return nil
}
