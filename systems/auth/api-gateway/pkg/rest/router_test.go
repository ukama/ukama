package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ory/client-go"
	mauth "github.com/ukama/ukama/systems/auth/api-gateway/mocks"
	cconfig "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"

	"github.com/stretchr/testify/assert"
)

var token = "test-token"
var mockEmail = "test@ukama.com"
var mockPassword = "@Test123"

var defaultCors = cors.Config{
	AllowAllOrigins: true,
}

var routerConfig = &RouterConfig{
	serverConf: &rest.HttpConfig{
		Cors: defaultCors,
	},
	auth: &cconfig.Auth{
		AuthServerUrl: "http://localhost",
	},
	k: cconfig.LoadAuthKey(),
}

func init() {
	gin.SetMode(gin.TestMode)
}

func TestPingRoute(t *testing.T) {
	w := httptest.NewRecorder()
	cma := &mauth.AuthManager{}
	req, _ := http.NewRequest("GET", "/ping", nil)
	r := NewRouter(&Clients{au: cma}, routerConfig).f.Engine()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestLogin(t *testing.T) {
	cma := &mauth.AuthManager{}
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
	r := NewRouter(&Clients{au: cma}, routerConfig).f.Engine()

	tm := time.Now()
	e := tm.Add(time.Hour * 24 * 7)
	st := "session-token"
	cma.On("LoginUser", mock.Anything, mock.Anything).Return(&client.SuccessfulNativeLogin{
		SessionToken: &st,
		Session: client.Session{
			ExpiresAt:       &e,
			AuthenticatedAt: &tm,
		},
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
	cma := &mauth.AuthManager{}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/auth", nil)

// 	req.Header.Set("X-Session-Token", token)

	r := NewRouter(&Clients{au: cma}, routerConfig).f.Engine()

	cma.On("ValidateSession", mock.Anything, mock.Anything).Return(&client.Session{}, nil)

// 	r.ServeHTTP(w, req)

// 	assert.Equal(t, 200, w.Code)
// }

func TestWhoami(t *testing.T) {
	cma := &mauth.AuthManager{}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/whoami", nil)

// 	req.Header.Set("X-Session-Token", token)

	r := NewRouter(&Clients{au: cma}, routerConfig).f.Engine()

	cma.On("ValidateSession", mock.Anything, mock.Anything).Return(&client.Session{}, nil)

// 	r.ServeHTTP(w, req)

// 	assert.Equal(t, 200, w.Code)
// }
