package client

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/data-plan/package/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PackageClient struct {
	conn          *grpc.ClientConn
	timeout       time.Duration
	packageClient pb.PackagesServiceClient
	host          string
}

func NewPackageClient(packageHost string, timeout time.Duration) *PackageClient {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	packageConn, err := grpc.DialContext(ctx, packageHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewPackagesServiceClient(packageConn)

	return &PackageClient{
		conn:          packageConn,
		packageClient: client,
		timeout:       timeout,
		host:          packageHost,
	}
}

func NewPackageFromClient(client pb.PackagesServiceClient) *PackageClient {
	return &PackageClient{
		host:          "localhost",
		timeout:       1 * time.Second,
		conn:          nil,
		packageClient: client,
	}
}

func (r *PackageClient) Close() {
	r.conn.Close()
}

func (d *PackageClient) AddPackage(req *pb.AddPackageRequest) (*pb.AddPackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	return d.packageClient.Add(ctx, req)
}

func (d *PackageClient) DeletePackage(id string) (*pb.DeletePackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	return d.packageClient.Delete(ctx, &pb.DeletePackageRequest{Uuid: id})
}

func (d *PackageClient) UpdatePackage(req *pb.UpdatePackageRequest) (*pb.UpdatePackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	return d.packageClient.Update(ctx, req)
}

func (d *PackageClient) GetPackage(id string) (*pb.GetPackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	return d.packageClient.Get(ctx, &pb.GetPackageRequest{Uuid: id})
}

func (d *PackageClient) GetPackageDetails(id string) (*pb.GetPackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	return d.packageClient.GetDetails(ctx, &pb.GetPackageRequest{Uuid: id})
}

func (d *PackageClient) GetPackageByOrg(orgId string) (*pb.GetByOrgPackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()
	return d.packageClient.GetByOrg(ctx, &pb.GetByOrgPackageRequest{OrgId: orgId})
}
