/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

package server

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/minio/minio-go/v7"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/hub/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/hub/api-gateway/pkg"
	"github.com/ukama/ukama/systems/hub/api-gateway/pkg/client"

	log "github.com/sirupsen/logrus"
	apb "github.com/ukama/ukama/systems/hub/artifactmanager/pb/gen"
	dpb "github.com/ukama/ukama/systems/hub/distributor/pb/gen"
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
		debugMode:     svcConf.DebugMode,
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
	auth := r.f.Group("/v1", "API Gateway", "Hub system version v1", func(ctx *gin.Context) {
		if r.config.auth.BypassAuthMode {
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
		artifact := auth.Group("/hub", "Artifact store", "Artifact operations")
		artifact.GET("/:type/:name/:filename", formatDoc("Get Artifact contents", "Get artifact contents or its index files"), tonic.Handler(r.artifactGetHandler, http.StatusOK))
		artifact.PUT("/:type/:name/:version", formatDoc("Upload artifact", "Upload a artifact contents"), tonic.Handler(r.artifactPutHandler, http.StatusCreated))
		artifact.GET("/:type/:name", formatDoc("List of versions for artifact", "List all the available version and location info for artifact"), tonic.Handler(r.artifactListVersionsHandler, http.StatusOK))
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

func (r *Router) artifactPutHandler(c *gin.Context) (*apb.StoreArtifactResponse, error) {

	req := &ArtifactUploadRequest{}
	req.ArtifactName = c.Param("name")
	req.ArtifactType = c.Param("type")
	req.Version = c.Param("version")
	log.Infof("Adding artifact %s with version %s of type : %s", req.ArtifactName, req.Version, req.ArtifactType)

	_, err := r.parseVersion(req.Version)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	// Copy the uploaded data to the buffer
	size, err := io.Copy(&buf, c.Request.Body)
	if err != nil {
		log.Errorf("Failed to copy buffer: %v", err)
		return nil, rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  fmt.Sprintf("copy to buffer err: %s", err.Error()),
		}
	}
	log.Infof("Got tar file with size %d", size)

	err = IsValidGzip(buf)
	if err != nil {
		log.Errorf("Not a gzip format for file: %s", req.ArtifactName)
		return nil, err
	}

	return r.clients.a.StoreArtifact(&apb.StoreArtifactRequest{
		Name:    req.ArtifactName,
		Type:    apb.ArtifactType(apb.ArtifactType_value[strings.ToUpper(req.ArtifactType)]),
		Version: req.Version,
		Data:    buf.Bytes(),
	})

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

func (r *Router) parseVersion(version string) (*semver.Version, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  "Invalid version format. Refer to https://semver.org/ for more information",
		}
	}

	return v, err
}

func IsValidGzip(buf bytes.Buffer) error {
	// Validate if the buffer contains gzip-compressed data
	_, err := gzip.NewReader(&buf)
	if err != nil && err != io.EOF {
		return rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  "Invalid file format",
		}
	}
	return nil
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
