package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/stretchr/testify/assert"

	"github.com/ukama/ukama/services/common/rest"
	"github.com/ukama/ukama/testing/services/network/internal/db"
	"github.com/ukama/ukama/testing/services/network/mocks"

	"github.com/ukama/ukama/testing/services/network/internal"
)

func init() {
	internal.IsDebugMode = true
}

var defaultConfig = &internal.Config{
	Server: rest.HttpConfig{
		Cors: cors.Config{
			AllowAllOrigins: true,
		},
	},
	ServiceRouter: "http://localhost",
}

func Test_RouterPing(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	r := NewRouter(defaultConfig, nil, nil).fizz.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func Test_GetList(t *testing.T) {
	node := db.VNode{
		NodeID: "1001",
		Status: "PowerOn",
	}
	list := []db.VNode{node}
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/list?node=1001&looking_for=vnode_list", nil)

	vNodeRepo := mocks.VNodeRepo{}
	r := NewRouter(defaultConfig, nil, &vNodeRepo).fizz.Engine()

	vNodeRepo.On("List").Return(&list, nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), node.NodeID)
}
