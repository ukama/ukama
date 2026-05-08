package server

import (
	"context"
	"fmt"
	"github.com/ukama/ukama/systems/common/ukama"
	crpc "github.com/ukama/ukama/systems/node/controller/pkg/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func (c *ControllerServer) SendNodeCommand(ctx context.Context, req *crpc.SendNodeCommandRequest) (*crpc.SendNodeCommandResponse, error) {
	if req.NodeId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "node_id is required")
	}
	if req.Method == "" {
		return nil, status.Errorf(codes.InvalidArgument, "method is required")
	}
	if req.Path == "" || !strings.HasPrefix(req.Path, "/") {
		return nil, status.Errorf(codes.InvalidArgument, "path must start with /")
	}
	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid node id: %s", err.Error())
	}
	if _, err := c.nRepo.Get(nId.String()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "node has not been registered yet: %s", err.Error())
	}
	target := fmt.Sprintf("%s...%s", c.orgName, nId.String())
	if err := c.publishMessage(target, strings.ToUpper(req.Method), req.Path, nId.String(), req.Body); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to publish node command: %s", err.Error())
	}
	return &crpc.SendNodeCommandResponse{}, nil
}
