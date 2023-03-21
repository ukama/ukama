package client

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SimManager struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.SimManagerServiceClient
	host    string
}

func NewSimManager(host string, timeout time.Duration) *SimManager {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewSimManagerServiceClient(conn)

	return &SimManager{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewSimManagerFromClient(SimManagerClient pb.SimManagerServiceClient) *SimManager {
	return &SimManager{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  SimManagerClient,
	}
}

func (sm *SimManager) Close() {
	sm.conn.Close()
}

func (sm *SimManager) AllocateSim(req *pb.AllocateSimRequest) (*pb.AllocateSimResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()
	return sm.client.AllocateSim(ctx, req)
}

func (sm *SimManager) GetSim(simId string) (*pb.GetSimResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()
	return sm.client.GetSim(ctx, &pb.GetSimRequest{SimId: simId})
}

func (sm *SimManager) GetSimsBySub(subscriberId string) (*pb.GetSimsBySubscriberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()
	return sm.client.GetSimsBySubscriber(ctx, &pb.GetSimsBySubscriberRequest{SubscriberId: subscriberId})
}

func (sm *SimManager) ToggleSimStatus(simId string, status string) (*pb.ToggleSimStatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()
	return sm.client.ToggleSimStatus(ctx, &pb.ToggleSimStatusRequest{SimId: simId, Status: status})
}

func (sm *SimManager) AddPackageToSim(req *pb.AddPackageRequest) (*pb.AddPackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()
	return sm.client.AddPackageForSim(ctx, req)
}

func (sm *SimManager) RemovePackageForSim(req *pb.RemovePackageRequest) (*pb.RemovePackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()
	return sm.client.RemovePackageForSim(ctx, req)
}

func (sm *SimManager) DeleteSim(simId string) (*pb.DeleteSimResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()
	return sm.client.DeleteSim(ctx, &pb.DeleteSimRequest{SimId: simId})
}

func (sm *SimManager) GetSimsByNetwork(networkId string) (*pb.GetSimsByNetworkResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()
	return sm.client.GetSimsByNetwork(ctx, &pb.GetSimsByNetworkRequest{NetworkId: networkId})
}

func (sm *SimManager) GetPackagesForSim(simId string) (*pb.GetPackagesBySimResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()
	return sm.client.GetPackagesBySim(ctx, &pb.GetPackagesBySimRequest{SimId: simId})
}

func (sm *SimManager) SetActivePackageForSim(req *pb.SetActivePackageRequest) (*pb.SetActivePackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()
	return sm.client.SetActivePackageForSim(ctx, req)
}
