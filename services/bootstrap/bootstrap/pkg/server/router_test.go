package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/stretchr/testify/assert"

	"github.com/ukama/ukama/services/bootstrap/bootstrap/pkg"
	"github.com/ukama/ukama/services/common/rest"

	sr "github.com/ukama/ukama/services/common/srvcrouter"
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
	rs := sr.ServiceRouter{}

	r := NewRouter(defaultCongif, &rs).fizz.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}
