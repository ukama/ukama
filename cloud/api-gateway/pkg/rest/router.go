package rest

import (
	"fmt"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/ukama/ukamaX/cloud/api-gateway/cmd/version"
	"github.com/ukama/ukamaX/cloud/api-gateway/pkg/swagger"
	pb "github.com/ukama/ukamaX/cloud/registry/pb/gen"
	"github.com/ukama/ukamaX/common/config"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"net/http"

	"github.com/gin-contrib/cors"

	"github.com/ukama/ukamaX/cloud/api-gateway/pkg"
	"github.com/ukama/ukamaX/cloud/api-gateway/pkg/client"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	hsspb "github.com/ukama/ukamaX/cloud/hss/pb/gen"
)

const ORG_URL_PARAMETER = "org"

type Router struct {
	gin            *gin.Engine
	port           int
	authMiddleware AuthMiddleware
	cors           cors.Config
	clients        *Clients
	metricsConfig  config.Metrics
}

type Clients struct {
	Registry *client.Registry
	Hss      *client.Hss
}

type AuthMiddleware interface {
	IsAuthenticated(c *gin.Context)
	IsAuthorized(c *gin.Context)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.Registry = client.NewRegistry(endpoints.Registry, endpoints.TimeoutSeconds)
	c.Hss = client.NewHss(endpoints.Hss, endpoints.TimeoutSeconds)
	return c
}

func NewRouter(port int,
	debugMode bool,
	authMiddleware AuthMiddleware,
	cors cors.Config,
	metrics config.Metrics,
	clients *Clients) *Router {
	r := &Router{
		authMiddleware: authMiddleware,
		clients:        clients,
		cors:           cors,
		metricsConfig:  metrics,
	}
	if !debugMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r.init(port)
	return r
}

func (rt *Router) Run() {
	logrus.Info("Listening on port ", rt.port)
	err := rt.gin.Run(fmt.Sprint(":", rt.port))
	if err != nil {
		panic(err)
	}
}

func (r *Router) init(port int) {
	r.gin = gin.Default()
	r.gin.Use(gin.Logger())
	r.gin.Use(cors.New(r.cors))
	r.port = port

	if r.metricsConfig.Enabled {
		prometheus := ginprometheus.NewPrometheus("api_gateway")
		prometheus.SetListenAddress(fmt.Sprint(":", r.metricsConfig.Port))
		prometheus.Use(r.gin)
	}

	tonic.SetErrorHook(errorHook)

	f := fizz.NewFromEngine(r.gin)

	authorized := f.Group("/", "Authorization", "Requires authorization", r.authMiddleware.IsAuthenticated,
		r.authMiddleware.IsAuthorized)

	authorized.Use()
	{
		const org = "/orgs/" + ":" + ORG_URL_PARAMETER

		authorized.GET(org, []fizz.OperationOption{}, tonic.Handler(r.orgHandler, http.StatusOK))

		// registry
		nodes := authorized.Group(org+"/nodes", "Nodes", "Nodes operations")
		nodes.GET("", nil, tonic.Handler(r.nodesHandler, http.StatusOK))

		// hss
		hss := authorized.Group(org+"/users", "Network Users", "Operations on network users and SIM cards"+
			"Do not confuse with organization users")
		hss.GET("", nil, tonic.Handler(r.getUsersHandler, http.StatusOK))
		authorized.POST(org+"/users", []fizz.OperationOption{}, tonic.Handler(r.postUsersHandler, http.StatusCreated))
		hss.DELETE("/:user", nil, tonic.Handler(r.deleteUserHandler, http.StatusOK))
	}

	f.GET("/ping", nil, tonic.Handler(r.pingHandler, http.StatusOK))

	infos := &openapi.Info{
		Title:       "Ukama API Gateway",
		Description: `Ukam API Gateway server`,
		Version:     version.Version,
	}
	f.GET("/openapi.json", nil, f.OpenAPI(infos, "json"))
	swagger.AddOpenApiUIHandler(r.gin, "swagger-ui", "/openapi.json")
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
	} else {
		if gErr, ok := e.(client.GrpcClientError); ok {
			errcode = gErr.HttpCode
			errpl = gErr.Message
		}
	}

	return errcode, gin.H{`error`: errpl}
}

func (r *Router) getOrgNameFromRoute(c *gin.Context) string {
	return c.Param("org")
}

func (r *Router) orgHandler(c *gin.Context) (*pb.Organization, error) {
	orgName := r.getOrgNameFromRoute(c)
	return r.clients.Registry.GetOrg(orgName)
}

func (r *Router) nodesHandler(c *gin.Context) (*pb.NodesList, error) {
	orgName := r.getOrgNameFromRoute(c)

	return r.clients.Registry.GetNodes(orgName)
}

type PingResponse struct {
	Message string `json:"message"`
}

func (rt *Router) pingHandler(c *gin.Context) (*PingResponse, error) {
	return &PingResponse{Message: "pong"}, nil
}

func (r *Router) getUsersHandler(c *gin.Context) (*hsspb.ListUsersResponse, error) {
	orgName := r.getOrgNameFromRoute(c)
	return r.clients.Hss.GetUsers(orgName)
}

func (r *Router) postUsersHandler(c *gin.Context, req *UserRequest) (*hsspb.AddUserResponse, error) {
	return r.clients.Hss.AddUser(req.Org, &hsspb.User{
		Imsi:      req.Imsi,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
	})
}

func (r *Router) deleteUserHandler(c *gin.Context) error {
	orgName := r.getOrgNameFromRoute(c)
	userId := c.Param("user")
	return r.clients.Hss.Delete(orgName, userId)
}
