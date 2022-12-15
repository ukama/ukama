package server

import (
	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
)

type SimManagerServer struct {
	pb.UnimplementedSimManagerServiceServer
}

func NewSimManagerServer() *SimManagerServer {
	return &SimManagerServer{}
}
