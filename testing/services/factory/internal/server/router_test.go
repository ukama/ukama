package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/stretchr/testify/assert"

	"github.com/ukama/ukama/services/common/rest"

	"github.com/ukama/ukama/testing/factory/internal"
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

	r := NewRouter(defaultConfig, nil).fizz.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func Test_PostBuildNodeQueryParamValidationFailure(t *testing.T) {

	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/node", nil)

	r := NewRouter(defaultConfig, nil).fizz.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	if assert.Equal(t, 400, w.Code) {
		assert.Contains(t, w.Body.String(), " Error:Field validation for 'LookingTo' ")
	}

}

func Test_PostBuildNodeQueryParamTypeValidationFailure(t *testing.T) {

	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/node?looking_to=create_node&type=xnode&count=1", nil)

	r := NewRouter(defaultConfig, nil).fizz.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	if assert.Equal(t, 400, w.Code) {
		assert.Contains(t, w.Body.String(), "Error:Field validation for 'Type' failed on the 'eq=HNODE|eq=TNODE|eq=ANODE|eq=hnode|eq=tnode|eq=anode'")
	}

}

func Test_PostBuildNodeWorkerInitFail(t *testing.T) {

	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/node?looking_to=create_node&type=hnode&count=1", nil)

	r := NewRouter(defaultConfig, nil).fizz.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	if assert.Equal(t, 500, w.Code) {
		assert.Contains(t, w.Body.String(), "factory worker not initialized")
	}

}
