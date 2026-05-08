package adapters

import (
	"context"
	crpc "github.com/ukama/ukama/systems/node/controller/pkg/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type NodeCommandAdapter interface {
	Send(ctx context.Context, nodeID, method, path string, body []byte) error
}
type ControllerAdapter struct {
	client  crpc.CommandServiceClient
	timeout time.Duration
}

func NewControllerAdapter(host string, timeout time.Duration) (*ControllerAdapter, error) {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &ControllerAdapter{client: crpc.NewCommandServiceClient(conn), timeout: timeout}, nil
}
func (a *ControllerAdapter) Send(ctx context.Context, nodeID, method, path string, body []byte) error {
	cctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()
	_, err := a.client.SendNodeCommand(cctx, &crpc.SendNodeCommandRequest{NodeId: nodeID, Method: method, Path: path, Body: body, Source: "site-controller"})
	return err
}
