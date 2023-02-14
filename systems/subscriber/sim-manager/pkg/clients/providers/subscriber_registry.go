package providers

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/subscriber/subscriber-registry/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// SubscriberRegistryClientProvider creates a local client to interact with
// a remote instance of  Subscriber Registry service.
type SubscriberRegistryClientProvider interface {
	GetClient() (pb.SubscriberRegistryServiceClient, error)
}

type subscriberRegistryClientProvider struct {
	subscriberRegistryService pb.SubscriberRegistryServiceClient
	subscriberRegistryHost    string
	timeout                   time.Duration
}

func NewSubscriberRegistryClientProvider(subscriberRegistryHost string, timeout time.Duration) SubscriberRegistryClientProvider {
	return &subscriberRegistryClientProvider{subscriberRegistryHost: subscriberRegistryHost, timeout: timeout}
}

func (sr *subscriberRegistryClientProvider) GetClient() (pb.SubscriberRegistryServiceClient, error) {
	if sr.subscriberRegistryService == nil {
		var conn *grpc.ClientConn

		ctx, cancel := context.WithTimeout(context.Background(), sr.timeout)
		defer cancel()

		log.Infoln("Connecting to Subscriber Registry service ", sr.subscriberRegistryHost)

		conn, err := grpc.DialContext(ctx, sr.subscriberRegistryHost, grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Errorf("Failed to connect to Subscriber Registry service %s. Error: %v", sr.subscriberRegistryHost, err)

			return nil, fmt.Errorf("failed to connect to remote Subscriber Registry service: %w", err)
		}

		sr.subscriberRegistryService = pb.NewSubscriberRegistryServiceClient(conn)
	}

	return sr.subscriberRegistryService, nil
}
