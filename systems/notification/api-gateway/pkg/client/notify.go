package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/notification/notify/pb/gen"
)

type Notify interface {
	Add(nodeId, severity, ntype, serviceName, description, details string, epochTime uint32) (*pb.AddResponse, error)
	Get(id string) (*pb.GetResponse, error)
	List(nodeId, serviceName, nType string, count uint32, sort bool) (*pb.ListResponse, error)
	Delete(id string) (*pb.DeleteResponse, error)
	Purge(nodeId, serviceName, nType string) (*pb.ListResponse, error)
}

type notify struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.NotifyServiceClient
	host    string
}

func NewNotify(host string, timeout time.Duration) (*notify, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)

		return nil, err
	}

	client := pb.NewNotifyServiceClient(conn)

	return &notify{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}, nil
}

func NewNotifyFromClient(notifyClient pb.NotifyServiceClient) *notify {
	return &notify{
		host:    "localhost",
		timeout: 10 * time.Second,
		conn:    nil,
		client:  notifyClient,
	}
}

func (m *notify) Close() {
	m.conn.Close()
}

func (n *notify) Add(nodeId, severity, ntype, serviceName,
	description, details string, epochTime uint32) (*pb.AddResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.Add(ctx,
		&pb.AddRequest{
			NodeId:      nodeId,
			Severity:    severity,
			Type:        ntype,
			ServiceName: serviceName,
			EpochTime:   epochTime,
			Description: description,
			Details:     details,
		})
}

func (n *notify) Get(id string) (*pb.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.Get(ctx, &pb.GetRequest{NotificationId: id})
}

func (n *notify) List(nodeId, serviceName, nType string, count uint32, sort bool) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.List(ctx,
		&pb.ListRequest{
			NodeId:      nodeId,
			Type:        nType,
			ServiceName: serviceName,
			Count:       count,
			Sort:        sort,
		})
}

func (n *notify) Delete(id string) (*pb.DeleteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.Delete(ctx, &pb.GetRequest{NotificationId: id})
}

func (n *notify) Purge(nodeId, serviceName, nType string) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	return n.client.Purge(ctx,
		&pb.PurgeRequest{
			NodeId:      nodeId,
			Type:        nType,
			ServiceName: serviceName,
		})
}
