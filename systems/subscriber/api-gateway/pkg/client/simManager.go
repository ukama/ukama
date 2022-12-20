package client

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/subscriber/api-gateway/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/* temp */
type SMDummyReq struct {
	Dummy string
}

type SMDummyResp struct {
	Dummy string
}

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

func (sm *SimManager) GetSubscriber(req *SMDummyReq) (*SMDummyResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	return sm.client.AddOrg(ctx, req)
}

func (sm *SimManager) AddSubscriber(req *SMDummyReq) (*SMDummyResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	return sm.client.UpdateOrg(ctx, req)
}

func (sm *SimManager) DeleteSubscriber(req *SMDummyReq) (*SMDummyResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	return sm.client.GetOrg(ctx, req)
}

func (sm *SimManager) GetSim(req *SMDummyReq) (*SMDummyResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	return sm.client.AddNodeForOrg(ctx, req)
}

func (sm *SimManager) AllocateSim(req *SMDummyReq) (*SMDummyResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	return sm.client.GetNodeForOrg(ctx, req)
}

func (sm *SimManager) AddPackageToSim(req *SMDummyReq) (*SMDummyResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	return sm.client.DeleteNodeForOrg(ctx, req)
}

func (sm *SimManager) RemovePackageForSim(req *SMDummyReq) (*SMDummyResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	return sm.client.AddSystemForOrg(ctx, req)
}

func (sm *SimManager) DeleteSim(req *SMDummyReq) (*SMDummyResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	return sm.client.UpdateSystemForOrg(ctx, req)
}

func (sm *SimManager) PatchSim(req *SMDummyReq) (*SMDummyResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	return sm.client.GetSystemForOrg(ctx, req)
}

func (sm *SimManager) getAllSubscribers(req *SMDummyReq) (*SMDummyResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	return sm.client.DeleteSystemForOrg(ctx, req)
}

func (sm *SimManager) getAllSims(req *SMDummyReq) (*SMDummyResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	return sm.client.DeleteSystemForOrg(ctx, req)
}
