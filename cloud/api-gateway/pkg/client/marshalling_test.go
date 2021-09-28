package client

import (
	"encoding/json"
	"fmt"
	grpcGate "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
	pb "github.com/ukama/ukamaX/cloud/registry/pb/gen"
	"net/http"
	"testing"
)

func Test_marshallStruct(t *testing.T) {
	msg := pb.GetNodeResponse{
		Org: &pb.Organization{
			Name: "org-1",
		},
		Node: &pb.Node{},
	}
	resp, err := marshallResponse(nil, &msg)

	assert.Nil(t, err)
	assert.True(t, json.Valid([]byte(resp)))
	assert.True(t, json.Valid([]byte(resp)))

	j := map[string]interface{}{}
	mErr := json.Unmarshal([]byte(resp), &j)
	assert.NoError(t, mErr)
	assert.Equal(t, msg.Org.Name, j["org"].(map[string]interface{})["name"])
	assert.Equal(t, nil, j["Node"])
}

func Test_marshallError(t *testing.T) {
	msg := "invalid request"
	resp, err := marshallResponse(fmt.Errorf(msg), nil)

	assert.Empty(t, resp)
	assert.Equal(t, msg, err.Message)
	assert.Equal(t, http.StatusInternalServerError, err.HttpCode)
}

func Test_marshallGrpcHttpError(t *testing.T) {
	grpcErr := &grpcGate.HTTPStatusError{
		HTTPStatus: http.StatusNotFound,
		Err:        fmt.Errorf("not found error"),
	}
	resp, err := marshallResponse(grpcErr, nil)

	assert.Empty(t, resp)
	assert.Equal(t, grpcErr.Error(), err.Message)
	assert.Equal(t, grpcErr.HTTPStatus, err.HttpCode)
}

func Test_marshallNodeStruct(t *testing.T) {
	msg := pb.GetNodeResponse{
		Org: &pb.Organization{
			Name: "org-1",
		},
		Node: &pb.Node{
			State:  pb.NodeState_UNDEFINED,
			NodeId: "node-id-1",
		},
	}
	resp, err := marshallResponse(nil, &msg)

	assert.Nil(t, err)

	j := map[string]interface{}{}
	mErr := json.Unmarshal([]byte(resp), &j)
	assert.NoError(t, mErr)
	assert.Equal(t, msg.Org.Name, j["org"].(map[string]interface{})["name"])
	node := j["node"].(map[string]interface{})
	assert.Equal(t, msg.Node.NodeId, node["nodeId"])
	assert.Equal(t, "UNDEFINED", node["state"])
}
