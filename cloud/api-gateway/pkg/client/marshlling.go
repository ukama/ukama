package client

import (
	"errors"
	grpcGate "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
	jsonpb "google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"net/http"
)

func marshallResponse(err error, res proto.Message) (string, *GrpcClientError) {

	clientError, done := marshalError(err)
	if done {
		return "", clientError
	}

	b, err := jsonpb.Marshal(res)
	if err != nil {
		return "", &GrpcClientError{
			HttpCode: http.StatusInternalServerError,
			Message:  err.Error(),
		}
	}

	return string(b), nil
}

func marshalError(err error) (*GrpcClientError, bool) {
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
		return &GrpcClientError{
			HttpCode: st,
			Message:  pb.Message,
		}, true
	}
	return nil, false
}
