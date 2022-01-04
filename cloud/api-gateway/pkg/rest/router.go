package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"

	"github.com/ukama/ukamaX/cloud/api-gateway/pkg"
	"github.com/ukama/ukamaX/cloud/api-gateway/pkg/client"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	hsspb "github.com/ukama/ukamaX/cloud/hss/pb/gen"
	urest "github.com/ukama/ukamaX/common/rest"
)

const ORG_URL_PARAMETER = "org"

type Router struct {
	gin            *gin.Engine
	port           int
	authMiddleware AuthMiddleware
	cors           cors.Config
	clients        *Clients
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
	clients *Clients) *Router {
	r := &Router{
		authMiddleware: authMiddleware,
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

	authorized.Use(r.authMiddleware.IsAuthenticated).Use(r.authMiddleware.IsAuthorized)
	{
		const org = "/orgs/" + ":" + ORG_URL_PARAMETER

		// registry
		authorized.GET(org, r.orgHandler)
		authorized.GET(org+"/nodes", r.nodesHandler)

		// hss
		// returns list of users

		authorized.GET(org+"/users", r.getUsersHandler)
		authorized.POST(org+"/users", r.postUsersHandler)
		authorized.DELETE(org+"/users/:user", r.deleteUserHandler)
	}

	r.gin.GET("/ping", r.pingHandler)
}

func (r *Router) getOrgNameFromRoute(c *gin.Context) string {
	return c.Param("org")
}

func (r *Router) orgHandler(c *gin.Context) {
	orgName := r.getOrgNameFromRoute(c)
	resp, err := r.clients.Registry.GetOrg(orgName)

	if err != nil {
		urest.ThrowError(c, err.HttpCode, err.Message, "", nil)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (r *Router) nodesHandler(c *gin.Context) {
	orgName := r.getOrgNameFromRoute(c)

	resp, err := r.clients.Registry.GetNodes(orgName)
	if err != nil {
		urest.ThrowError(c, err.HttpCode, "Registry request failed. Error:"+err.Message, "", nil)
		return
	}

	mResp, err := client.MarshallResponse(nil, resp)
	if err != nil {
		urest.ThrowError(c, err.HttpCode, "Failed marshaling response. Error:"+err.Message, "", nil)
		return
	}
	c.Header("Content-Type", "application/json")

	c.String(http.StatusOK, mResp)
}

func (rt *Router) pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (r *Router) getUsersHandler(c *gin.Context) {
	orgName := r.getOrgNameFromRoute(c)
	resp, err := r.clients.Hss.GetUsers(orgName)

	if err != nil {
		urest.ThrowError(c, err.HttpCode, err.Message, "", nil)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *Router) postUsersHandler(c *gin.Context) {
	var user hsspb.User
	orgName := r.getOrgNameFromRoute(c)
	err := c.ShouldBind(&user)
	if err != nil {
		urest.ThrowError(c, http.StatusInternalServerError, err.Error(), "", nil)
		return
	}

	resp, grpcErr := r.clients.Hss.AddUser(orgName, &user)
	if grpcErr != nil {
		urest.ThrowError(c, grpcErr.HttpCode, "Failed to add a user", grpcErr.Message, nil)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (r *Router) deleteUserHandler(c *gin.Context) {
	orgName := r.getOrgNameFromRoute(c)
	userId := c.Param("user")
	_, err := r.clients.Hss.Delete(orgName, userId)
	if err != nil {
		urest.ThrowError(c, http.StatusInternalServerError, err.Message, "", nil)
	}
}
