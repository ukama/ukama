package adapters

import (
	"context"

	pb "github.com/ukama/ukama/systems/subscriber/test-agent/pb/gen"
	"google.golang.org/grpc"
)

type TestAgentAdapter struct {
	conn   *grpc.ClientConn
	host   string
	client pb.TestAgentServiceClient
}

func (t *TestAgentAdapter) ActivateSim(ctx context.Context, iccid string) error {
	_, err := t.client.ActivateSim(ctx, &pb.ActivateSimRequest{Iccid: iccid})

	return err
}

func (t *TestAgentAdapter) DeactivateSim(ctx context.Context, iccid string) error {
	_, err := t.client.DeactivateSim(ctx, &pb.DeactivateSimRequest{Iccid: iccid})

	return err
}

func (t *TestAgentAdapter) TerminateSim(ctx context.Context, iccid string) error {
	_, err := t.client.TerminateSim(ctx, &pb.TerminateSimRequest{Iccid: iccid})

	return err
}

func (t *TestAgentAdapter) Close() {
	t.conn.Close()
}
