package rest

import (
	"encoding/json"
	"fmt"
	grpcGate "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"testing"
)

func Test_errorHook(t *testing.T) {
	dummyError := "dummy error"
	dummyErrResp := fmt.Sprint("{\"error\":\"", dummyError, "\"}")

	tests := []struct {
		name    string
		err     error
		expCode int
		expMsg  string
	}{
		{
			name:    "HttpError",
			err:     HttpError{HttpCode: 401, Message: dummyError},
			expCode: 401,
			expMsg:  dummyErrResp,
		},
		{
			name:    "tonic.BindError",
			err:     tonic.BindError{},
			expCode: 400,
			expMsg:  "",
		},
		{
			name:    "no_error",
			err:     nil,
			expCode: 0,
			expMsg:  "",
		},
		{
			name:    "no_error",
			err:     &grpcGate.HTTPStatusError{HTTPStatus: 404, Err: fmt.Errorf(dummyError)},
			expCode: 404,
			expMsg:  dummyErrResp,
		},
		{
			name:    "grpc_error_extract_description",
			err:     status.Errorf(codes.NotFound, "rpc error: code = Aborted desc = limit of sim cards reached for org"),
			expCode: http.StatusNotFound,
			expMsg:  "{\"error\":\"limit of sim cards reached for org\"}",
		},

		{
			name:    "grpc_error_plain_description",
			err:     status.Errorf(codes.NotFound, dummyError),
			expCode: http.StatusNotFound,
			expMsg:  dummyErrResp,
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(tt *testing.T) {
			c, resp := errorHook(nil, ts.err)
			b, err := json.Marshal(resp)
			if assert.NoError(tt, err) {
				assert.Equal(tt, ts.expCode, c)
				assert.Contains(tt, string(b), ts.expMsg)
			}
		})
	}
}
