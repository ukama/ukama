package client

import (
	"context"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/node/software/pb/gen"
	"google.golang.org/grpc"
)

type SoftwareManager struct {
	conn    *grpc.ClientConn
	client  pb.SoftwareServiceClient
	timeout time.Duration
	host    string
}

func NewSoftwareManager(softwareManagerHost string, timeout time.Duration) *SoftwareManager {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, softwareManagerHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewSoftwareServiceClient(conn)

	return &SoftwareManager{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    softwareManagerHost,
	}
}

func NewSoftwareManagerFromClient(mClient pb.SoftwareServiceClient) *SoftwareManager {
	return &SoftwareManager{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  mClient,
	}
}

func (r *SoftwareManager) Close() {
	r.conn.Close()
}

func (r *SoftwareManager) UpdateSoftware(space string, name string, tag string, nodeId string) (*pb.UpdateSoftwareResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.UpdateSoftware(ctx, &pb.UpdateSoftwareRequest{
		NodeId: nodeId,
		Space: space,
		Name:  name,
		Tag:   tag,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
