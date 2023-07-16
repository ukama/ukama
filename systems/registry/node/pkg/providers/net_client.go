package providers

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/registry/network/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NetworkClientProvider creates a local client to interact with
// a remote instance of  network service.
type NetworkClientProvider interface {
	GetClient() (pb.NetworkServiceClient, error)
}

type networkClientProvider struct {
	networkService pb.NetworkServiceClient
	networkHost    string
}

func NewNetworkClientProvider(networkHost string) NetworkClientProvider {
	return &networkClientProvider{networkHost: networkHost}
}

func (o *networkClientProvider) GetClient() (pb.NetworkServiceClient, error) {
	if o.networkService == nil {
		var conn *grpc.ClientConn

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		log.Infoln("Connecting to Network service ", o.networkHost)

		conn, err := grpc.DialContext(ctx, o.networkHost, grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Errorf("Failed to connect to Network service %s. Error: %v", o.networkHost, err)

			return nil, fmt.Errorf("failed to connect to remote network service: %w", err)
		}

		o.networkService = pb.NewNetworkServiceClient(conn)
	}

	return o.networkService, nil
}
