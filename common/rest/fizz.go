package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukamaX/common/rest/swagger"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
)

type PingResponse struct {
	Message string `json:"message"`
	Service string `json:"service"`
}

type HttpConfig struct {
	Port int
	Cors cors.Config
}

type HttpError struct {
	HttpCode int
	Message  string
}

func (g HttpError) Error() string {
	return g.Message
}

func NewFizzRouter(httpConfig *HttpConfig, srvName string, srvVersion string, isDebug bool) *fizz.Fizz {

	gin.SetMode(gin.ReleaseMode)
	if isDebug {
		gin.SetMode(gin.DebugMode)
	}

	g := gin.Default()
	g.Use(gin.Logger())
	g.Use(cors.New(httpConfig.Cors))

	m := ginmetrics.GetMonitor()
	m.UseWithoutExposingEndpoint(g)

	tonic.SetErrorHook(errorHook)

	f := fizz.NewFromEngine(g)
	f.GET("/ping", nil, tonic.Handler(func(c *gin.Context) (*PingResponse, error) {
		return &PingResponse{Message: "pong", Service: fmt.Sprintf("%s@%s", srvName, srvVersion)}, nil
	}, http.StatusOK))

	infos := &openapi.Info{
		Title:       srvName,
		Description: "Rest API for " + srvName,
		Version:     srvVersion,
	}

	f.GET("/openapi.json", nil, f.OpenAPI(infos, "json"))
	swagger.AddOpenApiUIHandler(g, "swagger", "/openapi.json")

	return f
}

func errorHook(c *gin.Context, e error) (int, interface{}) {
	if e == nil {
		logrus.Errorf("This erro means that something is broken but it's no clear what. Usually something bad with serialization")
		return 0, nil
	}
	errcode, errpl := 500, e.Error()
	if _, ok := e.(tonic.BindError); ok {
		errcode = 400
		errpl = e.Error()
	} else if gErr, ok := e.(HttpError); ok {
		errcode = gErr.HttpCode
		errpl = gErr.Message
	}

	return errcode, gin.H{`error`: errpl}
}
