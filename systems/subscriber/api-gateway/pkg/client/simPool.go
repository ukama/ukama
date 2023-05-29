package client

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SimPool struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.SimServiceClient
	host    string
}

func NewSimPool(host string, timeout time.Duration) *SimPool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewSimServiceClient(conn)

	return &SimPool{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewSimPoolFromClient(SimPoolClient pb.SimServiceClient) *SimPool {
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

func (sp *SimPool) Get(iccid string) (*pb.GetByIccidResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sp.timeout)
	defer cancel()

	return sp.client.GetByIccid(ctx, &pb.GetByIccidRequest{Iccid: iccid})
}

func (sp *SimPool) GetSims(simType string) (*pb.GetSimsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sp.timeout)
	defer cancel()

	return sp.client.GetSims(ctx, &pb.GetSimsRequest{SimType: simType})
}

func (sp *SimPool) GetStats(simType string) (*pb.GetStatsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sp.timeout)
	defer cancel()

	return sp.client.GetStats(ctx, &pb.GetStatsRequest{SimType: simType})
}

func (sp *SimPool) AddSimsToSimPool(req *pb.AddRequest) (*pb.AddResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sp.timeout)
	defer cancel()

	return sp.client.Add(ctx, req)
}

func (sp *SimPool) UploadSimsToSimPool(req *pb.UploadRequest) (*pb.UploadResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sp.timeout)
	defer cancel()

	return sp.client.Upload(ctx, req)
}

func (sp *SimPool) DeleteSimFromSimPool(id []uint64) (*pb.DeleteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sp.timeout)
	defer cancel()

	return sp.client.Delete(ctx, &pb.DeleteRequest{Id: id})
}
