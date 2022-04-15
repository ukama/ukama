package pkg

import (
	gen "github.com/ukama/ukamaX/cloud/hss/pb/gen/simmgr"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type SimManagerStub struct {
}

func (s *SimManagerStub) ActivateSim(ctx context.Context, in *gen.ActivateSimRequest, opts ...grpc.CallOption) (*gen.ActivateSimResponse, error) {
	return &gen.ActivateSimResponse{}, nil
}

func (s *SimManagerStub) GetSimStatus(ctx context.Context, in *gen.GetSimStatusRequest, opts ...grpc.CallOption) (*gen.GetSimStatusResponse, error) {
	return &gen.GetSimStatusResponse{}, nil
}

func (s *SimManagerStub) GetSimInfo(ctx context.Context, in *gen.GetSimInfoRequest, opts ...grpc.CallOption) (*gen.GetSimInfoResponse, error) {
	return &gen.GetSimInfoResponse{}, nil
}

func (s *SimManagerStub) SetServiceStatus(ctx context.Context, in *gen.SetServiceStatusRequest, opts ...grpc.CallOption) (*gen.SetServiceStatusResponse, error) {
	return &gen.SetServiceStatusResponse{}, nil
}

type CloserStub struct {
}

func (c CloserStub) Close() error {
	return nil
}
