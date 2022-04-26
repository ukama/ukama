package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukamaX/common/rest"
	"github.com/ukama/ukamaX/templates/rest-service/pkg"
)

func init() {
	pkg.IsDebugMode = true
}

var defaultCongif = &rest.HttpConfig{
	Cors: cors.Config{
		AllowAllOrigins: true,
	},
}

func Test_RouterPing(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	r := NewRouter(defaultCongif).fizz.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func Test_RouterGet(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/foos/bar", nil)

	r := NewRouter(defaultCongif).fizz.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
}
