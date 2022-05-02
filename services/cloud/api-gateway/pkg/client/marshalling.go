package client

import (
	"errors"

	"github.com/ukama/ukama/services/common/rest"

	grpcGate "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
)

func marshalError(err error) (grpcError *rest.HttpError, isItAnError bool) {
	if err != nil {
		var customStatus *grpcGate.HTTPStatusError
		if errors.As(err, &customStatus) {
			err = customStatus.Err
		}
		s := status.Convert(err)
		pb := s.Proto()
		logrus.Errorf("Error response: %v\n", pb)

		st := grpcGate.HTTPStatusFromCode(s.Code())
		if customStatus != nil {
			st = customStatus.HTTPStatus
		}
		return &rest.HttpError{
			HttpCode: st,
			Message:  pb.Message,
		}, true
	}
	return nil, false
}
