package grpc

import (
	"github.com/sirupsen/logrus"
	"github.com/ukama/openIoR/services/common/sql"
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
