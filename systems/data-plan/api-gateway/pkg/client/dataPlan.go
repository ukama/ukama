package client

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	pbBaseRate "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"
	pb "github.com/ukama/ukama/systems/data-plan/package/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DataPlan struct {
	conn           *grpc.ClientConn
	baseRateConn   *grpc.ClientConn
	timeout        time.Duration
	packageClient  pb.PackagesServiceClient
	host           string
	baseRateClient pbBaseRate.BaseRatesServiceClient
}

func NewDataPlan(packageHost, baseRateHost string, timeout time.Duration) *DataPlan {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	packageConn, err := grpc.DialContext(ctx, packageHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewPackagesServiceClient(packageConn)

	baseRateConn, err := grpc.DialContext(ctx, baseRateHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	baseRateClient := pbBaseRate.NewBaseRatesServiceClient(baseRateConn)

	return &DataPlan{
		conn:           packageConn,
		packageClient:  client,
		timeout:        timeout,
		host:           packageHost,
		baseRateClient: baseRateClient,
	}
}

func NewPackageFromClient(client pb.PackagesServiceClient, baseRateClient pbBaseRate.BaseRatesServiceClient) *DataPlan {
	return &DataPlan{
		host:           "localhost",
		timeout:        1 * time.Second,
		conn:           nil,
		packageClient:  client,
		baseRateClient: baseRateClient,
	}
}

func (r *DataPlan) Close() {
	r.conn.Close()
	r.baseRateConn.Close()
}

func (d *DataPlan) AddPackage(req *pb.AddPackageRequest) (*pb.AddPackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	return d.packageClient.Add(ctx, req)
}

func (d *DataPlan) DeletePackage(id string) (*pb.DeletePackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	return d.packageClient.Delete(ctx, &pb.DeletePackageRequest{Uuid: id})
}

func (d *DataPlan) UpdatePackage(req *pb.UpdatePackageRequest) (*pb.UpdatePackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	return d.packageClient.Update(ctx, req)
}
func (d *DataPlan) GetPackage(id string) (*pb.GetPackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	return d.packageClient.Get(ctx, &pb.GetPackageRequest{Uuid: id})
}

func (d *DataPlan) UploadBaseRates(req *pbBaseRate.UploadBaseRatesRequest) (*pbBaseRate.UploadBaseRatesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	return d.baseRateClient.UploadBaseRates(ctx, req)
}
func (d *DataPlan) GetBaseRates(req *pbBaseRate.GetBaseRatesRequest) (*pbBaseRate.GetBaseRatesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()
	return d.baseRateClient.GetBaseRates(ctx, req)
}
func (d *DataPlan) GetBaseRate(id string) (*pbBaseRate.GetBaseRateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()
	return d.baseRateClient.GetBaseRate(ctx, &pbBaseRate.GetBaseRateRequest{Uuid: id})
}
func (d *DataPlan) GetPackageByOrg(orgId string) (*pb.GetByOrgPackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()
	return d.packageClient.GetByOrg(ctx, &pb.GetByOrgPackageRequest{OrgId: orgId})
}
