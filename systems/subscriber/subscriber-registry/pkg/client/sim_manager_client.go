package client

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// SimManagerClientProvider creates a local client to interact with
// a remote instance of  Org service.
type SimManagerClientProvider interface {
	GetSimManagerService() (pb.SimManagerServiceClient, error)
}

type simManagerClientProvider struct {
	simManagerService pb.SimManagerServiceClient
	simManagerHost    string
}

func NewSimManagerClientProvider(simManagerHost string) SimManagerClientProvider {
	return &simManagerClientProvider{simManagerHost: simManagerHost}
}

func (u *simManagerClientProvider) GetSimManagerService() (pb.SimManagerServiceClient, error) {
	if u.simManagerService == nil {
		var conn *grpc.ClientConn

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		log.Infoln("Connecting to SimManager service ", u.simManagerHost)

		conn, err := grpc.DialContext(ctx, u.simManagerHost, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect to SimManager service %s. Error: %v", u.simManagerHost, err)
		}

		u.simManagerService = pb.NewSimManagerServiceClient(conn)
	}

	return u.simManagerService, nil
}
