package client

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/subscriber/SimPool/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/* temp */
type SPDummyReq struct {
	Dummy string
}

type SPDummyResp struct {
	Dummy string
}

type SimPool struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.SimPoolServiceClient
	host    string
}

func NewSimPool(host string, timeout time.Duration) *SimPool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewSimPoolServiceClient(conn)

	return &SimPool{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewSimPoolFromClient(SimPoolClient pb.SimPoolServiceClient) *SimPool {
	return &SimPool{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  SimPoolClient,
	}
}

func (sp *SimPool) Close() {
	sp.conn.Close()
}

func (sp *SimPool) GetSimPoolStats(req *SPDummyReq) (*SPDummyResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sp.timeout)
	defer cancel()

	return sp.client.AddOrg(ctx, req)
}

func (sp *SimPool) AddSimsToSimPool(req *SPDummyReq) (*SPDummyResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sp.timeout)
	defer cancel()

	return sp.client.UpdateOrg(ctx, req)
}

func (sp *SimPool) UploadSimsToSimPool(req *SPDummyReq) (*SPDummyResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sp.timeout)
	defer cancel()

	return sp.client.GetOrg(ctx, req)
}

func (sp *SimPool) deleteSimFromSimPool(req *SPDummyReq) (*SPDummyResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sp.timeout)
	defer cancel()

	return sp.client.AddNodeForOrg(ctx, req)
}
