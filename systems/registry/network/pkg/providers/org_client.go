package providers

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/registry/org/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// OrgClientProvider creates a local client to interact with
// a remote instance of  Org service.
type OrgClientProvider interface {
	GetClient() (pb.OrgServiceClient, error)
}

type orgClientProvider struct {
	orgService pb.OrgServiceClient
	orgHost    string
}

func NewOrgClientProvider(orgHost string) OrgClientProvider {
	return &orgClientProvider{orgHost: orgHost}
}

func (o *orgClientProvider) GetClient() (pb.OrgServiceClient, error) {
	if o.orgService == nil {
		var conn *grpc.ClientConn

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		log.Infoln("Connecting to Org service ", o.orgHost)

		conn, err := grpc.DialContext(ctx, o.orgHost, grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect to Org service %s. Error: %v", o.orgHost, err)
		}

		o.orgService = pb.NewOrgServiceClient(conn)
	}

	return o.orgService, nil
}
