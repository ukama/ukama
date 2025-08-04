/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/nucleus/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/nucleus/api-gateway/pkg"
	"github.com/ukama/ukama/systems/nucleus/api-gateway/pkg/client"

	log "github.com/sirupsen/logrus"
	orgpb "github.com/ukama/ukama/systems/nucleus/org/pb/gen"
	userspb "github.com/ukama/ukama/systems/nucleus/user/pb/gen"
)

const (
	USER_ID_KEY       = "UserId"
	ORG_URL_PARAMETER = "org"
	DUMMY_AUTH_ID     = "4eaef5fa-6548-481c-87eb-2b2d85ce6141"
)

type Router struct {
	f       *fizz.Fizz
	clients *Clients
	config  *RouterConfig
}

type RouterConfig struct {
	metricsConfig config.Metrics
	httpEndpoints *pkg.HttpEndpoints
	debugMode     bool
	serverConf    *rest.HttpConfig
	auth          *config.Auth
}

type Clients struct {
	Organization organization
	User         *client.Users
}

type organization interface {
	AddOrg(orgName string, owner string, certificate string, country string, currency string) (*orgpb.AddResponse, error)
	GetOrg(orgName string) (*orgpb.GetByNameResponse, error)
	GetOrgs(ownerUUID string) (*orgpb.GetByUserResponse, error)
	UpdateOrgToUser(orgId string, userId string) (*orgpb.UpdateOrgForUserResponse, error)
	RemoveOrgForUser(orgId string, userId string) (*orgpb.RemoveOrgForUserResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.Organization = client.NewOrgRegistry(endpoints.Org, endpoints.Timeout)
	c.User = client.NewUsers(endpoints.User, endpoints.Timeout)
	return c
}

func NewRouter(clients *Clients, config *RouterConfig, authfunc func(*gin.Context, string) error) *Router {
	r := &Router{
		clients: clients,
		config:  config,
	}

	if !config.debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r.init(authfunc)
	return r
}

func NewRouterConfig(svcConf *pkg.Config) *RouterConfig {
	return &RouterConfig{
		metricsConfig: svcConf.Metrics,
		httpEndpoints: &svcConf.HttpServices,
		serverConf:    &svcConf.Server,
		debugMode:     svcConf.DebugMode,
		auth:          svcConf.Auth,
	}
}

func (rt *Router) Run() {
	log.Info("Listening on port ", rt.config.serverConf.Port)
	err := rt.f.Engine().Run(fmt.Sprint(":", rt.config.serverConf.Port))
	if err != nil {
		panic(err)
	}
}

func (r *Router) init(f func(*gin.Context, string) error) {
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect=true")
	auth := r.f.Group("/v1", "API gateway", "Registry system version v1", func(ctx *gin.Context) {
		if r.config.auth.BypassAuthMode {
			ctx.Set(USER_ID_KEY, DUMMY_AUTH_ID)
			log.Info("Bypassing auth")
			return
		}
		s := fmt.Sprintf("%s, %s, %s", pkg.SystemName, ctx.Request.Method, ctx.Request.URL.Path)
		ctx.Request.Header.Set("Meta", s)
		err := f(ctx, r.config.auth.AuthAPIGW)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		}
	})
	auth.Use()
	{
		// org routes
		const org = "/orgs"
		orgs := auth.Group(org, "Orgs", "Operations on Orgs")
		orgs.GET("", formatDoc("Get Orgs", "Get all organization owned by a user"), tonic.Handler(r.getOrgsHandler, http.StatusOK))
		orgs.GET("/:org", formatDoc("Get Org by name", "Get organization by name"), tonic.Handler(r.getOrgHandler, http.StatusOK))
		orgs.POST("", formatDoc("Add Org", "Add a new organization"), tonic.Handler(r.postOrgHandler, http.StatusCreated))
		orgs.PUT("/:org_id/users/:user_id", formatDoc("Add user to orgs", "set association between user and org"), tonic.Handler(r.updateOrgToUserHandler, http.StatusCreated))
		orgs.DELETE("/:org_id/users/:user_id", formatDoc("Remove user from org", "remove association between user and org"), tonic.Handler(r.removeUserFromOrgHandler, http.StatusCreated))

		// Users routes
		const user = "/users"
		users := auth.Group(user, "Users", "Operations on Users")
		users.POST("", formatDoc("Add User", "Add a new User to the registry"), tonic.Handler(r.postUserHandler, http.StatusCreated))
		users.GET("/:user_id", formatDoc("Get User", "Get a specific user"), tonic.Handler(r.getUserHandler, http.StatusOK))
		users.GET("/auth/:auth_id", formatDoc("Get User By AuthId", "Get a specific user by authId"), tonic.Handler(r.getUserByAuthIdHandler, http.StatusOK))
		users.GET("/whoami/:user_id", formatDoc("Get detailed User", "Get a specific user details with all linked orgs"), tonic.Handler(r.whoamiHandler, http.StatusOK))
		//add get user by email
		users.GET("/email/:email", formatDoc("Get User By Email", "Get a specific user by email"), tonic.Handler(r.getUserByEmailHandler, http.StatusOK))
		// user orgs-member
		// update user
		// Deactivate user
		// Delete user
		// users.DELETE("/:user_id", formatDoc("Remove User", "Remove a user from the registry"), tonic.Handler(r.removeUserHandler, http.StatusOK))
	}
}

// Org handlers
func (r *Router) getOrgHandler(c *gin.Context, req *GetOrgRequest) (*orgpb.GetByNameResponse, error) {
	return r.clients.Organization.GetOrg(req.OrgName)
}

func (r *Router) getOrgsHandler(c *gin.Context, req *GetOrgsRequest) (*orgpb.GetByUserResponse, error) {
	ownerUUID, ok := c.GetQuery("user_uuid")
	if !ok {
		return nil, &rest.HttpError{HttpCode: http.StatusBadRequest,
			Message: "user_uuid is a mandatory query parameter"}
	}

	return r.clients.Organization.GetOrgs(ownerUUID)
}

func (r *Router) postOrgHandler(c *gin.Context, req *AddOrgRequest) (*orgpb.AddResponse, error) {
	return r.clients.Organization.AddOrg(req.OrgName, req.Owner, req.Certificate, req.Country, req.Currency)
}

func (r *Router) updateOrgToUserHandler(c *gin.Context, req *UserOrgRequest) (*orgpb.UpdateOrgForUserResponse, error) {
	return r.clients.Organization.UpdateOrgToUser(req.OrgId, req.UserId)
}

func (r *Router) removeUserFromOrgHandler(c *gin.Context, req *UserOrgRequest) (*orgpb.RemoveOrgForUserResponse, error) {
	return r.clients.Organization.RemoveOrgForUser(req.OrgId, req.UserId)
}

func (r *Router) getUserByEmailHandler(c *gin.Context, req *GetByEmailRequest) (*userspb.GetResponse, error) {
	return r.clients.User.GetByEmail(strings.ToLower(req.Email))
}

// Users handlers
func (r *Router) getUserHandler(c *gin.Context, req *GetUserRequest) (*userspb.GetResponse, error) {
	return r.clients.User.Get(req.UserId)
}

func (r *Router) getUserByAuthIdHandler(c *gin.Context, req *GetUserByAuthIdRequest) (*userspb.GetResponse, error) {
	return r.clients.User.GetByAuthId(req.AuthId)
}

func (r *Router) whoamiHandler(c *gin.Context, req *GetUserRequest) (*userspb.WhoamiResponse, error) {
	return r.clients.User.Whoami(req.UserId)
}

func (r *Router) postUserHandler(c *gin.Context, req *AddUserRequest) (*userspb.AddResponse, error) {
	return r.clients.User.AddUser(&userspb.User{
		Name:   req.Name,
		Email:  strings.ToLower(req.Email),
		Phone:  req.Phone,
		AuthId: req.AuthId,
	})
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
