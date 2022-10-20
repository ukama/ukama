package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ukama/ukama/systems/data-plan/pkg/config"
	"github.com/ukama/ukama/systems/data-plan/pkg/db"
	pb "github.com/ukama/ukama/systems/data-plan/pkg/pb"
	services "github.com/ukama/ukama/systems/data-plan/pkg/services"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	h := db.Init(c.DBUrl)

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Rates server running on", c.Port)

	s := services.Server{
		H: h,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterRatesServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
