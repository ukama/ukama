package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/stretchr/testify/assert"

	"github.com/ukama/openIoR/services/common/rest"
	"github.com/ukama/openIoR/services/factory/nmr/mocks"
	"github.com/ukama/openIoR/services/factory/nmr/pkg"

	"github.com/ukama/openIoR/services/factory/nmr/internal/db"
	"github.com/ukama/openIoR/services/factory/nmr/pkg/router"
)

func init() {
	pkg.IsDebugMode = true
}

var defaultCongif = &pkg.Config{
	Server: rest.HttpConfig{
		Cors: cors.Config{
			AllowAllOrigins: true,
		},
	},
}

func Test_RouterPing(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	nodeRepo := mocks.NodeRepo{}
	moduleRepo := mocks.ModuleRepo{}
	rs := router.RouterServer{}

	r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func NewNode(id string) *db.Node {
	return &db.Node{
		NodeID:        id,
		Type:          "hnode",
		PartNumber:    "a1",
		Skew:          "s1",
		Mac:           "00:01:02:03:04:05",
		SwVersion:     "1.1",
		OemName:       "ukama",
		MfgTestStatus: "",
		Status:        "Assembly",
	}
}

func NewModule(id string) *db.Module {
	return &db.Module{
		ModuleID:      id,
		Type:          "hnode",
		PartNumber:    "a1",
		HwVersion:     "s1",
		Mac:           "00:01:02:03:04:05",
		SwVersion:     "1.1",
		MfgName:       "ukama",
		MfgTestStatus: "UnderTest",
		UnitID: sql.NullString{
			String: "1001",
			Valid:  true,
		},
	}
}

func Test_PutNode(t *testing.T) {
	// Arrange
	nodeId := "1001"
	node := NewNode(nodeId)

	body, _ := json.Marshal(node)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/node/?node=1001&looking_to=update", bytes.NewReader(body))

	nodeRepo := mocks.NodeRepo{}
	moduleRepo := mocks.ModuleRepo{}
	rs := router.RouterServer{}

	r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

	nodeRepo.On("AddOrUpdateNode", node).Return(nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)

}

func Test_GetNode(t *testing.T) {
	t.Run("Read node", func(t *testing.T) {
		// Arrange
		nodeId := "1001"
		node := NewNode(nodeId)

		//body, _ := json.Marshal(node)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/node/?node=1001&looking_for=*", nil)

		nodeRepo := mocks.NodeRepo{}
		moduleRepo := mocks.ModuleRepo{}
		rs := router.RouterServer{}

		r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

		nodeRepo.On("GetNode", nodeId).Return(node, nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), nodeId)
	})

}

func Test_DeleteNode(t *testing.T) {
	t.Run("Read node", func(t *testing.T) {
		// Arrange
		nodeId := "1001"
		//body, _ := json.Marshal(node)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/node/?node=1001&looking_to=*", nil)

		nodeRepo := mocks.NodeRepo{}
		moduleRepo := mocks.ModuleRepo{}
		rs := router.RouterServer{}

		r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

		nodeRepo.On("DeleteNode", nodeId).Return(nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 200, w.Code)

	})

}

func Test_PutNodeStatus(t *testing.T) {
	// Arrange
	nodeId := "1001"
	status := "testing"

	url := "/node/status?node=" + nodeId + "&looking_to=update&status=" + status
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", url, nil)

	nodeRepo := mocks.NodeRepo{}
	moduleRepo := mocks.ModuleRepo{}
	rs := router.RouterServer{}

	r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

	nodeRepo.On("UpdateNodeStatus", nodeId, status).Return(nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)

}

func Test_GetNodeStatus(t *testing.T) {
	t.Run("Read node status", func(t *testing.T) {
		// Arrange
		nodeId := "1001"
		status := "testing"

		url := "/node/status?node=" + nodeId + "&looking_for=*"

		//body, _ := json.Marshal(node)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)

		nodeRepo := mocks.NodeRepo{}
		moduleRepo := mocks.ModuleRepo{}
		rs := router.RouterServer{}

		r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

		nodeRepo.On("GetNodeStatus", nodeId).Return(&status, nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "testing")

	})

}

func Test_PutNodeMfgStatus(t *testing.T) {
	// Arrange
	nodeId := "1001"
	status := "testing"

	url := "/node/mfgstatus?node=" + nodeId + "&looking_to=update&status=" + status
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", url, nil)

	nodeRepo := mocks.NodeRepo{}
	moduleRepo := mocks.ModuleRepo{}
	rs := router.RouterServer{}

	r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

	nodeRepo.On("UpdateNodeStatus", nodeId, status).Return(nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)

}

func Test_GetNodeMfgStatus(t *testing.T) {
	t.Run("Read node mfg status", func(t *testing.T) {
		// Arrange
		nodeId := "1001"
		status := "Tested"
		mfg := []byte("\"report: passed\"")

		url := "/node/mfgstatus?node=" + nodeId + "&looking_for=*"

		//body, _ := json.Marshal(node)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)

		nodeRepo := mocks.NodeRepo{}
		moduleRepo := mocks.ModuleRepo{}
		rs := router.RouterServer{}

		r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

		nodeRepo.On("GetNodeMfgStatus", nodeId).Return(&status, &mfg, nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "Tested")

	})

}
