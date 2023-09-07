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
	timeout        time.Duration
}

func NewSimPoolClientProvider(simPoolHost string, timeout time.Duration) SimPoolClientProvider {
	return &simPoolClientProvider{simPoolHost: simPoolHost, timeout: timeout}
}

func (sp *simPoolClientProvider) GetClient() (pb.SimServiceClient, error) {
	if sp.simPoolService == nil {
		var conn *grpc.ClientConn

		ctx, cancel := context.WithTimeout(context.Background(), sp.timeout)
		defer cancel()

		log.Infoln("Connecting to Sim Pool service ", sp.simPoolHost)

		conn, err := grpc.DialContext(ctx, sp.simPoolHost, grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Errorf("Failed to connect to Sim Pool service %s. Error: %v", sp.simPoolHost, err)

			return nil, fmt.Errorf("failed to connect to remote Sim Pool service: %w", err)
		}

		sp.simPoolService = pb.NewSimServiceClient(conn)
	}

	return sp.simPoolService, nil
}
