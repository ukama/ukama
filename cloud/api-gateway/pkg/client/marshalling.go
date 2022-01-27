package client

import (
	"errors"
	grpcGate "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
)

type HttpError struct {
	HttpCode int
	Message  string
}

func (g HttpError) Error() string {
	return g.Message
}

func marshalError(err error) (grpcError *HttpError, isItAnError bool) {
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
		return &HttpError{
			HttpCode: st,
			Message:  pb.Message,
		}, true
	}
	return nil, false
}
