package client

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Lookup struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.LookupServiceClient
	host    string
}

func Newlookup(host string, timeout time.Duration) *Lookup {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewLookupServiceClient(conn)

	return &Lookup{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewLookupFromClient(lookupClient pb.LookupServiceClient) *Lookup {
	return &Lookup{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  lookupClient,
	}
}

func (r *Lookup) Close() {
	r.conn.Close()
}

func (l *Lookup) GetNode(req *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	return l.client.GetNode(ctx, req)
}
