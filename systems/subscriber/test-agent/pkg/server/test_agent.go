package server

import (
	pb "github.com/ukama/ukama/systems/subscriber/test-agent/pb/gen"
)

type TestAgentServer struct {
	pb.UnimplementedTestAgentServiceServer
}

func NewTestAgentServer() *TestAgentServer {
	return &TestAgentServer{}
}
