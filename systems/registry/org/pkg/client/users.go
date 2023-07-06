package client

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/registry/users/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


type RegistryUsersClientProvider interface {
	GetClient() (pb.UserServiceClient, error)
}

type registryUsersClientProvider struct {
	registryUsersService pb.UserServiceClient
	registryUsersHost    string
	timeout                   time.Duration
}

func NewRegistryUsersClientProvider(registryUsersHost string, timeout time.Duration) RegistryUsersClientProvider {
	return &registryUsersClientProvider{registryUsersHost: registryUsersHost, timeout: timeout}
}

func (p *registryUsersClientProvider) GetClient() (pb.UserServiceClient, error) {
	if p.registryUsersService == nil {
		var conn *grpc.ClientConn

		ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
		defer cancel()

		log.Infoln("Connecting to Registry user service ", p.registryUsersHost)

		conn, err := grpc.DialContext(ctx, p.registryUsersHost, grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Errorf("Failed to connect to user Registry service %s. Error: %v", p.registryUsersHost, err)

			return nil, fmt.Errorf("failed to connect to remote user Registry service: %w", err)
		}

		p.registryUsersService = pb.NewUserServiceClient(conn)
	}

	return p.registryUsersService, nil
}