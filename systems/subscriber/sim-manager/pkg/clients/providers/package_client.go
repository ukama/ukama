package providers

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/data-plan/package/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// PackageClientProvider creates a local client to interact with
// a remote instance of  Package service.
type PackageClientProvider interface {
	GetClient() (pb.PackagesServiceClient, error)
}

type packageClientProvider struct {
	packageService pb.PackagesServiceClient
	packageHost    string
}

func NewPackageClientProvider(packageHost string) PackageClientProvider {
	return &packageClientProvider{packageHost: packageHost}
}

func (p *packageClientProvider) GetClient() (pb.PackagesServiceClient, error) {
	if p.packageService == nil {
		var conn *grpc.ClientConn

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		log.Infoln("Connecting to Package service ", p.packageHost)

		conn, err := grpc.DialContext(ctx, p.packageHost, grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Errorf("Failed to connect to Package service %s. Error: %v", p.packageHost, err)

			return nil, fmt.Errorf("failed to connect to remote package service: %w", err)
		}

		p.packageService = pb.NewPackagesServiceClient(conn)
	}

	return p.packageService, nil
}
