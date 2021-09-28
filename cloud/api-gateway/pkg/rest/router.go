package rest

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"net/http"

	"github.com/ukama/ukamaX/cloud/api-gateway/pkg"
	"github.com/ukama/ukamaX/cloud/api-gateway/pkg/client"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	urest "github.com/ukama/ukamaX/common/rest"
)

type Router struct {
	gin            *gin.Engine
	port           int
	authMiddleware AuthMiddleware
	cors           cors.Config
	clients        *Clients
}

type Clients struct {
	Registry *client.Registry
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.Registry = client.NewRegistry(endpoints.Registry, endpoints.TimeoutSeconds)
	return c
}

type AuthMiddleware interface {
	IsAuthenticated() gin.HandlerFunc
}

func NewRouter(port int,
	debugMode bool,
	authMiddlware AuthMiddleware,
	cors cors.Config,
	clients *Clients) *Router {
	r := &Router{
		authMiddleware: authMiddlware,
		clients:        clients,
		cors:           cors,
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

	authorized := r.gin.Group("/")

	authorized.Use(r.authMiddleware.IsAuthenticated())
	{
		authorized.GET("/orgs/:name", r.orgHandler)
		authorized.GET("/nodes", r.nodesHandler)
	}

	r.gin.GET("/ping", r.pingHandler)
}

func (r *Router) orgHandler(c *gin.Context) {
	orgName := c.Param("name")
	resp, err := r.clients.Registry.GetOrg(orgName)

	if err != nil {
		urest.ThrowError(c, err.HttpCode, err.Message, "", nil)
		return
	}
	c.String(http.StatusOK, resp)
}

func (r *Router) nodesHandler(c *gin.Context) {
	userId := c.GetString(USER_ID_KEY)
	if len(userId) == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	resp, err := r.clients.Registry.GetNodes(userId, "")

	if err != nil {
		urest.ThrowError(c, err.HttpCode, err.Message, "", nil)
		return
	}
	c.String(http.StatusOK, resp)
}

func (rt *Router) pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
