package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ory "github.com/ory/client-go"
	cconfig "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"

	"github.com/stretchr/testify/assert"
)

var token string
var mockEmail = "test@ukama.com"
var mockPassword = "@Pass2021"

var defaultCors = cors.Config{
	AllowAllOrigins: true,
}

var routerConfig = &RouterConfig{
	serverConf: &rest.HttpConfig{
		Cors: defaultCors,
	},
	auth: cconfig.LoadAuthHostConfig("auth"),
	o:    nil,
	k:    cconfig.LoadAuthKey(),
}

func init() {
	jar, _ := cookiejar.New(nil)
	gin.SetMode(gin.TestMode)
	routerConfig.o = ory.NewAPIClient(&ory.Configuration{
		Servers: []ory.ServerConfiguration{
			{
				URL: routerConfig.auth.AuthServerUrl,
			},
		},
	})
	routerConfig.o.GetConfig().HTTPClient = &http.Client{
		Jar: jar,
	}
	routerConfig.o.GetConfig().DefaultHeader = map[string]string{}
}

func TestPingRoute(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	r := NewRouter(routerConfig).f.Engine()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestLogin(t *testing.T) {
	r := gin.Default()
	router := NewRouter(routerConfig)
	payload := &LoginReq{
		Email:    mockEmail,
		Password: mockPassword,
	}
	router.init()

	r.POST("/v1/login", func(ctx *gin.Context) {
		res, err := router.login(ctx, payload)

		if err != nil {
			return
		}
		token = res.Token
		ctx.JSON(http.StatusOK, res)
	})

	jp, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	reqBody := []byte(jp)
	req, _ := http.NewRequest("POST", "/v1/login", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var res LoginRes
	err = json.Unmarshal(w.Body.Bytes(), &res)
	assert.Nil(t, err)
	assert.NotEmpty(t, res.Token)
	assert.Equal(t, res.Token, token)
}

func TestAuthenticate(t *testing.T) {
	r := NewRouter(routerConfig)
	g := gin.New()
	g.GET("/v1/auth", func(c *gin.Context) {
		err := r.authenticate(c, &OptionalReqHeader{
			XSessionToken: token,
		})
		if err != nil {
			return
		}
		c.JSON(http.StatusOK, "")
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/auth", nil)
	req.Header.Add("X-Session-Token", token)
	g.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", w.Result().StatusCode)
	}
}

func TestWhoami(t *testing.T) {
	r := NewRouter(routerConfig)
	g := gin.New()
	g.GET("/v1/whoami", func(c *gin.Context) {
		res, err := r.getUserInfo(c, &OptionalReqHeader{
			XSessionToken: token,
		})
		if err != nil {
			return
		}
		c.JSON(http.StatusOK, res)
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/whoami", nil)
	req.Header.Add("X-Session-Token", token)
	g.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", w.Result().StatusCode)
	}
	assert.Equal(t, http.StatusOK, w.Code)
	var res GetUserInfo
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.Nil(t, err)
	assert.NotEmpty(t, res.Email)
	assert.Equal(t, res.Email, mockEmail)
}
