package pkg

import (
	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukamaX/cloud/hss/pb/gen"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type ImsiClientProvider interface {
	GetClient() (pb.ImsiServiceClient, error)
}

type imsiClientProvider struct {
	imsiService pb.ImsiServiceClient
	hssHost     string
}

func NewImsiClientProvider(hssHost string) ImsiClientProvider {
	return &imsiClientProvider{hssHost: hssHost}
}

func (i *imsiClientProvider) GetClient() (pb.ImsiServiceClient, error) {
	if i.imsiService == nil {
		var conn *grpc.ClientConn

		// connect to Grpc service
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		log.Infoln("Connecting to hss service ", i.hssHost)

		conn, err := grpc.DialContext(ctx, i.hssHost, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect to service %s. Error: %v", i.hssHost, err)
		}

		i.imsiService = pb.NewImsiServiceClient(conn)
	}
	return i.imsiService, nil
}
