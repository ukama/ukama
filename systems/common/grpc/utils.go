package grpc

import (
	"context"

	"google.golang.org/grpc/credentials/insecure"
	"github.com/ukama/ukama/systems/common/config"
	"google.golang.org/grpc"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Returns grpc error with code based on Sql error.
// Handles cases such as "not found" or "duplicate key"
func SqlErrorToGrpc(err error, entity string) error {
	logrus.Error(err)
	if err != nil {
		if sql.IsNotFoundError(err) {
			return status.Errorf(codes.NotFound, entity+" record not found")
		}

		if sql.IsDuplicateKeyError(err) {
			return status.Errorf(codes.AlreadyExists, entity+" already exist")
		}
	}

	return status.Errorf(codes.Internal, err.Error())
}

func CreateGrpcConn(conf config.GrpcService) *grpc.ClientConn {
	var conn *grpc.ClientConn

	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), conf.Timeout)
	defer cancel()

	logrus.Infoln("Connecting to service ", conf.Host)

	conn, err := grpc.DialContext(ctx, conf.Host, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(
			grpc.ConnectParams{
				MinConnectTimeout: conf.Timeout,
			}))
	if err != nil {
		log.Fatalf("Failed to connect to service %s. Error: %v", conf.Host, err)
	}
	return conn
}
