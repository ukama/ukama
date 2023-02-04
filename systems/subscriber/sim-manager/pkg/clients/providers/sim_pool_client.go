package providers

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// SimPoolClientProvider creates a local client to interact with
// a remote instance of  Sim Pool service.
type SimPoolClientProvider interface {
	GetClient() (pb.SimServiceClient, error)
}

type simPoolClientProvider struct {
	simPoolService pb.SimServiceClient
	simPoolHost    string
}

func NewSimPoolClientProvider(simPoolHost string) SimPoolClientProvider {
	return &simPoolClientProvider{simPoolHost: simPoolHost}
}

func (p *simPoolClientProvider) GetClient() (pb.SimServiceClient, error) {
	if p.simPoolService == nil {
		var conn *grpc.ClientConn

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		log.Infoln("Connecting to Sim Pool service ", p.simPoolHost)

		conn, err := grpc.DialContext(ctx, p.simPoolHost, grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Errorf("Failed to connect to Sim Pool service %s. Error: %v", p.simPoolHost, err)

			return nil, fmt.Errorf("failed to connect to remote Sim Pool service: %w", err)
		}

		p.simPoolService = pb.NewSimServiceClient(conn)
	}

	return p.simPoolService, nil
}
