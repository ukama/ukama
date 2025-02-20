/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

package server

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	apb "github.com/ukama/ukama/systems/hub/artifactmanager/pb/gen"
	dpb "github.com/ukama/ukama/systems/hub/distributor/pb/gen"
	"github.com/ukama/ukama/systems/hub/node-gateway/cmd/version"
	"github.com/ukama/ukama/systems/hub/node-gateway/pkg"
	"github.com/ukama/ukama/systems/hub/node-gateway/pkg/client"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/minio/minio-go/v7"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
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
	isGlobal      bool
	serverConf    *rest.HttpConfig
	auth          *config.Auth
	distributor   string
}

type Clients struct {
	a artifactManager
	c chunker
}

type artifactManager interface {
	StoreArtifact(in *apb.StoreArtifactRequest) (*apb.StoreArtifactResponse, error)
	GetArtifact(in *apb.GetArtifactRequest) (*apb.GetArtifactResponse, error)
	GetArtifactVersionList(in *apb.GetArtifactVersionListRequest) (*apb.GetArtifactVersionListResponse, error)
	ListArtifacts(in *apb.ListArtifactRequest) (*apb.ListArtifactResponse, error)
}

type chunker interface {
	CreateChunk(in *dpb.CreateChunkRequest) (*dpb.CreateChunkResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.a = client.NewArtifactManager(endpoints.ArtifactManager, endpoints.MaxMsgSize, endpoints.Timeout)
	c.c = client.NewChunker(endpoints.Distributor, endpoints.MaxMsgSize, endpoints.Timeout)

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
		debugMode:     svcConf.BaseConfig.DebugMode,
		auth:          svcConf.Auth,
		isGlobal:      svcConf.IsGlobal,
		distributor:   svcConf.HttpServices.Distributor,
	}
}

func (r *Router) Run() {
	log.Info("Listening on port ", r.config.serverConf.Port)
	err := r.f.Engine().Run(fmt.Sprint(":", r.config.serverConf.Port))
	if err != nil {
		log.Error(err)
	}
}

func (r *Router) init(f func(*gin.Context, string) error) {
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, true, r.config.auth.AuthAppUrl+"?redirect=true")
	auth := r.f.Group("/v1", "Node Gateway", "Hub system version v1", func(ctx *gin.Context) {
		if r.config.auth.BypassAuthMode {
			log.Info("Bypassing auth")
			return
		}
		s := fmt.Sprintf("%s, %s, %s", pkg.SystemName, ctx.Request.Method, ctx.Request.URL.Path)
		ctx.Request.Header.Set("Meta", s)
		err := f(ctx, r.config.auth.AuthAPIGW)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}
		if err == nil {
			return
		}
	})
	auth.Use()
	{
		artifact := auth.Group("/hub", "Artifact store", "Artifact operations")
		artifact.GET("/:type/:name/:filename", formatDoc("Get Artifact contents", "Get artifact contents or its index files"), tonic.Handler(r.artifactGetHandler, http.StatusOK))
		artifact.GET("/:type/:name", formatDoc("List of versions for artifcat", "List all the available version and location info for artifact"), tonic.Handler(r.artifactListVersionsHandler, http.StatusOK))
		artifact.GET("/:type", formatDoc("List all artifact", "List all artifact of the matching type"), tonic.Handler(r.listArtifactsHandler, http.StatusOK))

		distr := auth.Group("/distributor", "Get chunks", "Download Artifact in chunk")
		distr.GET("/*proxypath", formatDoc("Get chunks", "Get artifact chunks"), tonic.Handler(r.proxy, http.StatusOK))

	}
}

func (r *Router) proxy(c *gin.Context) error {
	remote, err := url.Parse(r.config.distributor)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = c.Param("proxypath")
	}

	proxy.ServeHTTP(c.Writer, c.Request)

	return nil
}

// curl --request GET http://0.0.0.0:8000/v1/hub/app/test-capp/0.0.34.tar.gz --output test-capp-v-0.0.34.tar.gz
// curl --request GET http://0.0.0.0:8000/v1/hub/app/test-capp/0.0.34.caibx --output test-capp-v0.0.34.caibx
func (r *Router) artifactGetHandler(c *gin.Context, req *ArtifactRequest) error {
	log.Infof("Getting artifact %s of type %s with filname %s", req.Name, req.ArtifactType, req.FileName)

	resp, err := r.clients.a.GetArtifact(&apb.GetArtifactRequest{
		Name:     req.Name,
		Type:     apb.ArtifactType(apb.ArtifactType_value[strings.ToUpper(req.ArtifactType)]),
		FileName: req.FileName,
	})
	if err != nil {
		return err
	}

	_, err = c.Writer.Write(resp.Data)
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return rest.HttpError{
				HttpCode: http.StatusNotFound,
				Message:  "Artifact not found",
			}
		}
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition",
		fmt.Sprintf("attachment; filename=%s", resp.FileName))

	return nil
}

func (r *Router) artifactListVersionsHandler(c *gin.Context, req *ArtifactVersionListRequest) (*apb.GetArtifactVersionListResponse, error) {
	log.Infof("Getting version list: %s of type %s", req.Name, req.ArtifactType)

	return r.clients.a.GetArtifactVersionList(&apb.GetArtifactVersionListRequest{
		Name: req.Name,
		Type: apb.ArtifactType(apb.ArtifactType_value[strings.ToUpper(req.ArtifactType)]),
	})
}

func (r *Router) listArtifactsHandler(c *gin.Context, req *ArtifactListRequest) (*apb.ListArtifactResponse, error) {
	log.Infof("Getting list of artifacts of type %s", req.ArtifactType)

	return r.clients.a.ListArtifacts(&apb.ListArtifactRequest{
		Type: apb.ArtifactType(apb.ArtifactType_value[strings.ToUpper(req.ArtifactType)]),
	})

}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
