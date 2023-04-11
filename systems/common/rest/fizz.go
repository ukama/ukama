package rest

import (
	"fmt"
	"net/http"
	"strings"

	grpcGate "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest/swagger"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
)

var SESSION_KEY = "ukama_session"

type PingResponse struct {
	Message string `json:"message"`
	Service string `json:"service"`
}

type HttpError struct {
	HttpCode int
	Message  string
}

func (g HttpError) Error() string {
	return g.Message
}

func NewFizzRouter(httpConfig *HttpConfig, srvName string, srvVersion string, isDebug bool, redirectUrl string) *fizz.Fizz {

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
	tonic.SetRenderHook(renderHook, jsonContentType[0])

	f := fizz.NewFromEngine(g)
	f.GET("/ping", nil, tonic.Handler(func(c *gin.Context) (*PingResponse, error) {
		return &PingResponse{Message: "pong", Service: fmt.Sprintf("%s@%s", srvName, srvVersion)}, nil
	}, http.StatusOK))

	f.Generator().SetSecuritySchemes(map[string]*openapi.SecuritySchemeOrRef{
		"cookie": {
			SecurityScheme: &openapi.SecurityScheme{
				Type: "oauth2",
				In:   "header",
				Name: SESSION_KEY,
				Flows: &openapi.OAuthFlows{
					Implicit: &openapi.OAuthFlow{
						AuthorizationURL: redirectUrl,
					},
				},
			},
		},
	})

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
		logrus.Errorf("This error means that something is broken but it's no clear what. Usually something bad with serialization")
		return 0, nil
	}
	errcode, errpl := 500, e.Error()

	switch et := e.(type) {
	case HttpError:
		errcode = et.HttpCode
		errpl = et.Message

	case tonic.BindError:
		errcode = http.StatusBadRequest
		errpl = e.Error()

	case *grpcGate.HTTPStatusError:
		errcode = et.HTTPStatus
		errpl = e.Error()
	}

	if stat, ok := status.FromError(e); ok {
		errcode = grpcGate.HTTPStatusFromCode(stat.Code())
		pb := stat.Proto()
		// Get rid of extra info we have in message like code, grpc error etc
		const desc = " desc = "
		idx := strings.Index(pb.Message, desc)
		if idx > 0 {
			errpl = pb.Message[idx+len(desc):]
		} else {
			errpl = pb.GetMessage()
		}

	}

	return errcode, gin.H{`error`: errpl}
}

// renderHook is identical to default renderHook from gin except it renders filds that are ommited
func renderHook(c *gin.Context, statusCode int, payload interface{}) {
	var status int

	// This is a tricky part.
	// We need to be able to set the status in toni.Handeler for cases when it's not default
	// Here is how it done in defaul gin renderHook https://github.com/loopfz/gadgeto/blob/c4f8b2f64586099b9b281cbe99aa2f8b05e7d8b0/tonic/tonic.go#L111
	// but this does not work because here c.Writer.Written() is always false
	// We have to realy on default status from Gin taht is always 200 for no reason
	if c.Writer.Status() != 200 {
		status = c.Writer.Status()
	} else {
		status = statusCode
	}
	if payload != nil {
		if gin.IsDebugging() {
			c.Render(status, ExtJson{Data: payload, Indent: true})
		} else {
			c.Render(status, ExtJson{Data: payload})
		}
	} else {
		c.String(status, "")
	}
}
