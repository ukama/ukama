package client

import (
	"fmt"
	grpcGate "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
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
