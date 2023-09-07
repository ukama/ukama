package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/common/ukama"
)

func TestGetNodeIdFromPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	nodeIdKey := "nodeId"
	nodeId := ukama.NewVirtualNodeId("homenode")
	httpRecorder := httptest.NewRecorder()
	t.Run("validNodeId", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httpRecorder)
		c.Params = []gin.Param{
			gin.Param{Key: nodeIdKey, Value: nodeId.String()},
		}
		actualNodeId, isValid := GetNodeIdFromPath(c, nodeIdKey)
		assert.True(t, isValid)
		assert.Equal(t, nodeId, actualNodeId)
	})

	t.Run("invalidNodeId", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httpRecorder)
		c.Params = []gin.Param{
			gin.Param{Key: "448cca7e-b893-456c-8bf9-4c531ac85db8", Value: nodeId.String()},
		}
		actualNodeId, isValid := GetNodeIdFromPath(c, nodeIdKey)
		assert.False(t, isValid)
		assert.Equal(t, ukama.NodeID(""), actualNodeId)
	})

	t.Run("noNodeIdKey", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httpRecorder)
		c.Params = []gin.Param{}
		actualNodeId, isValid := GetNodeIdFromPath(c, nodeIdKey)
		assert.False(t, isValid)
		assert.Equal(t, ukama.NodeID(""), actualNodeId)
		assert.Equal(t, http.StatusBadRequest, httpRecorder.Result().StatusCode)
	})
}
