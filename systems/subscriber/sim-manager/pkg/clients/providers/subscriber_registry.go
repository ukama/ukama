package providers

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// SubscriberRegistryClientProvider creates a local client to interact with
// a remote instance of  Subscriber Registry service.
type SubscriberRegistryClientProvider interface {
	GetClient() (pb.RegistryServiceClient, error)
}

type subscriberRegistryClientProvider struct {
	subscriberRegistryService pb.RegistryServiceClient
	subscriberRegistryHost    string
}

func NewSubscriberRegistryClientProvider(subscriberRegistryHost string) SubscriberRegistryClientProvider {
	return &subscriberRegistryClientProvider{subscriberRegistryHost: subscriberRegistryHost}
}

func (p *subscriberRegistryClientProvider) GetClient() (pb.RegistryServiceClient, error) {
	if p.subscriberRegistryService == nil {
		var conn *grpc.ClientConn

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		log.Infoln("Connecting to Subscriber Registry service ", p.subscriberRegistryHost)

		conn, err := grpc.DialContext(ctx, p.subscriberRegistryHost, grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Errorf("Failed to connect to Subscriber Registry service %s. Error: %v", p.subscriberRegistryHost, err)

			return nil, fmt.Errorf("failed to connect to remote Subscriber Registry service: %w", err)
		}

		p.subscriberRegistryService = pb.NewRegistryServiceClient(conn)
	}

	return p.subscriberRegistryService, nil
}
