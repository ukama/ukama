package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ory "github.com/ory/client-go"
	mauth "github.com/ukama/ukama/systems/auth/api-gateway/mocks"
	"github.com/ukama/ukama/systems/auth/api-gateway/pkg/client"
	cconfig "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var token = ""
var mockEmail = "test@ukama.com"
var mockPassword = "Pass2021"

var defaultCors = cors.Config{
	AllowAllOrigins: true,
}

var routerConfig = &RouterConfig{
	serverConf: &rest.HttpConfig{
		Cors: defaultCors,
	},
	auth: &cconfig.Auth{
		AuthServerUrl: "http://localhost:8080",
	},
	k: cconfig.LoadAuthKey(),
}

func init() {
	gin.SetMode(gin.TestMode)
}

func TestPingRoute(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r := NewRouter(client.NewAuthManager(routerConfig.auth.AuthServerUrl, 3*time.Second), routerConfig).f.Engine()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestLogin(t *testing.T) {
	cma := &mauth.Auth{}
	w := httptest.NewRecorder()
	payload := &LoginReq{
		Email:    mockEmail,
		Password: mockPassword,
	}
	jp, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	reqBody := []byte(jp)
	req, _ := http.NewRequest("POST", "/v1/login", bytes.NewBuffer(reqBody))
	r := NewRouter(client.NewAuthManager(routerConfig.auth.AuthServerUrl, 3*time.Second), routerConfig).f.Engine()

	cma.On("LoginUser", mock.Anything, mock.Anything).Return(&LoginRes{
		Token: "some-token",
	}, nil)

	r.ServeHTTP(w, req)

	lg := LoginRes{}
	err = json.Unmarshal(w.Body.Bytes(), &lg)
	if err != nil {
		t.Error(err)
	}
	token = lg.Token
	assert.Equal(t, 200, w.Code)
}

func TestAuthenticate(t *testing.T) {
	cma := &mauth.Auth{}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/auth", nil)

	req.Header.Set("X-Session-Token", token)

	r := NewRouter(client.NewAuthManager(routerConfig.auth.AuthServerUrl, 3*time.Second), routerConfig).f.Engine()

	cma.On("ValidateSession", mock.Anything, mock.Anything).Return(&ory.Session{}, nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestWhoami(t *testing.T) {
	cma := &mauth.Auth{}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/whoami", nil)

	req.Header.Set("X-Session-Token", token)

	r := NewRouter(client.NewAuthManager(routerConfig.auth.AuthServerUrl, 3*time.Second), routerConfig).f.Engine()

	cma.On("ValidateSession", mock.Anything, mock.Anything).Return(&ory.Session{}, nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
