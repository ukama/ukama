package client

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/data-plan/package/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DataPlan struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.PackagesServiceClient
	host    string
}

func NewDataPlan(host string, timeout time.Duration) *DataPlan {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewPackagesServiceClient(conn)

	return &DataPlan{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewPackageFromClient(packageClient pb.PackagesServiceClient) *DataPlan {
	return &DataPlan{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  packageClient,
	}
}

func (r *DataPlan) Close() {
	r.conn.Close()
}

func (d *DataPlan) AddPackage(req *pb.AddPackageRequest) (*pb.AddPackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	return d.client.Add(ctx, req)
}

func (d *DataPlan) DeletePackage(req *pb.DeletePackageRequest) (*pb.DeletePackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	return d.client.Delete(ctx, req)
}

func (d *DataPlan) GetPackage(req *pb.GetPackagesRequest) (*pb.GetPackagesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	return d.client.Get(ctx, req)
}

func (d *DataPlan) UpdatePackage(req *pb.UpdatePackageRequest) (*pb.UpdatePackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	return d.client.Update(ctx, req)
}

