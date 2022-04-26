package client

import (
	"fmt"
	grpcGate "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"testing"
)

func Test_marshallError(t *testing.T) {
	msg := "invalid request"
	resp, _ := marshalError(fmt.Errorf(msg))

	assert.Equal(t, msg, resp.Message)
	assert.Equal(t, http.StatusInternalServerError, resp.HttpCode)
}

func Test_marshallGrpcHttpError(t *testing.T) {
	grpcErr := &grpcGate.HTTPStatusError{
		HTTPStatus: http.StatusNotFound,
		Err:        fmt.Errorf("not found error"),
	}
	resp, isErr := marshalError(grpcErr)

	assert.True(t, isErr)
	assert.Equal(t, grpcErr.Error(), resp.Message)
	assert.Equal(t, grpcErr.HTTPStatus, resp.HttpCode)
}

func Test_marshallGrpcStatusCodes(t *testing.T) {

	tests := []struct {
		Code     codes.Code
		Expected int
	}{
		{Code: codes.Aborted, Expected: http.StatusConflict},
		{Code: codes.Canceled, Expected: http.StatusRequestTimeout},
		{Code: codes.Internal, Expected: http.StatusInternalServerError},
		{Code: codes.NotFound, Expected: http.StatusNotFound},
		{Code: codes.InvalidArgument, Expected: http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("Status_%s", tc.Code.String()), func(tt *testing.T) {
			gerr, isErr := marshalError(status.Errorf(tc.Code, "error"))
			assert.True(t, isErr)
			assert.Equal(t, tc.Expected, gerr.HttpCode, "Bad mapping for %s", tc.Code.String())
		})
	}
}
