package provider

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/registry/users/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// UserClientProvider creates a local client to interact with
// a remote instance of  Users service.
type UserClientProvider interface {
	GetClient() (pb.UserServiceClient, error)
}

type userClientProvider struct {
	userService pb.UserServiceClient
	userHost    string
}

func NewUserClientProvider(userHost string) UserClientProvider {
	return &userClientProvider{userHost: userHost}
}

func (u *userClientProvider) GetClient() (pb.UserServiceClient, error) {
	if u.userService == nil {
		var conn *grpc.ClientConn

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		log.Infoln("Connecting to users service ", u.userHost)

		conn, err := grpc.DialContext(ctx, u.userHost, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect to users service %s. Error: %v", u.userHost, err)
		}

		u.userService = pb.NewUserServiceClient(conn)
	}

	return u.userService, nil
}
