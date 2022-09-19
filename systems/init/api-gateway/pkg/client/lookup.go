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

func (r *Lookup) Close() {
	r.conn.Close()
}

func (l *Lookup) AddOrg(req *pb.AddOrgRequest) (*pb.AddOrgResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	return l.client.AddOrg(ctx, req)
}

func (l *Lookup) GetOrg(req *pb.GetOrgRequest) (*pb.GetOrgResponse, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	// defer cancel()

	return nil, nil
}

func (l *Lookup) AddNodeForOrg(req *pb.AddNodeRequest) (*pb.AddNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	return l.client.AddNodeForOrg(ctx, req)
}

func (l *Lookup) GetNodeForOrg(req *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	return l.client.GetNodeForOrg(ctx, req)
}

func (l *Lookup) DeleteNodeForOrg(req *pb.DeleteNodeRequest) (*pb.DeleteNodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	return l.client.DeleteNodeForOrg(ctx, req)
}

func (l *Lookup) UpdateSystemForOrg(req *pb.UpdateSystemRequest) (*pb.UpdateSystemResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	return l.client.UpdateSystemForOrg(ctx, req)
}

func (l *Lookup) GetSystemForOrg(req *pb.GetSystemRequest) (*pb.GetSystemResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	return l.client.GetSystemForOrg(ctx, req)
}

func (l *Lookup) DeleteSystemForOrg(req *pb.DeleteSystemRequest) (*pb.DeleteSystemResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	return l.client.DeleteSystemForOrg(ctx, req)
}
