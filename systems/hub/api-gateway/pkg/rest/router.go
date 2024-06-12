/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

package server

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/hub/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/hub/api-gateway/pkg"
	"github.com/ukama/ukama/systems/hub/api-gateway/pkg/client"
	apb "github.com/ukama/ukama/systems/hub/artifactManager/pb/gen"
	dpb "github.com/ukama/ukama/systems/hub/distributor/pb/gen"

	"github.com/Masterminds/semver/v3"
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
}

type Clients struct {
	a artifactManager
	d distributor
}

type artifactManager interface {
	StoreArtifact(in *apb.StoreArtifactRequest) (*apb.StoreArtifactResponse, error)
	GetArtifactLocation(in *apb.GetArtifactLocationRequest) (*apb.GetArtifactLocationResponse, error)
	GetArtifact(in *apb.GetArtifactRequest) (*apb.GetArtifactResponse, error)
	GetArtifactVersionList(in *apb.GetArtifactVersionListRequest) (*apb.GetArtifactVersionListResponse, error)
	ListArtifacts(in *apb.ListArtifactRequest) (*apb.ListArtifactResponse, error)
}

type distributor interface {
	CreateChunk(in *dpb.CreateChunkRequest) (*dpb.CreateChunkResponse, error)
	Chunk(in *dpb.GetChunkRequest) (*dpb.GetChunkResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.a = client.NewArtifactManager(endpoints.ArtifactManager, endpoints.Timeout)
	c.d = client.NewDistributor(endpoints.Distributor, endpoints.Timeout)

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
		artifact.PUT("/:type/:filename/:version", formatDoc("Upload artifact", "Upload a artifact contents"), tonic.Handler(r.artifactPutHandler, http.StatusCreated))
		artifact.GET("/:type/:name", formatDoc("List of versions for artifcat", "List all the available version and location info for artifact"), tonic.Handler(r.artifactListVersionsHandler, http.StatusOK))
		artifact.GET("/:type", formatDoc("List all artifact", "List all artifact of the matching type"), tonic.Handler(r.listArtifactsHandler, http.StatusOK))
		//capps.GET("/location/:type/:name/:filename/:version", formatDoc("Location", "Provide a location of the artifact to download from"), tonic.Handler(r.artifactLocationHandler, http.StatusOK))

		distr := auth.Group("/ditributor", "Get Artifacts in chunks", "Download Artifact in chunk")
		distr.GET("/*proxypath", formatDoc("Get Artifact contents", "Get artifact contents or its index files"), tonic.Handler(proxy, http.StatusOK))

	}
}

func proxy(c *gin.Context) error {
	remote, err := url.Parse("http://distributor:8099")
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = c.Param("proxyPath")
	}

	proxy.ServeHTTP(c.Writer, c.Request)

	return nil
}

func (r *Router) artifactGetHandler(c *gin.Context, req *ArtifactRequest) error {
	log.Infof("Getting artifact %s of type %s with filname %s", req.Name, req.ArtifactType, req.FileName)

	resp, err := r.clients.a.GetArtifact(&apb.GetArtifactRequest{
		Name:     req.Name,
		Type:     apb.ArtifactType(apb.ArtifactType_value[req.ArtifactType]),
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

// curl --request PUT   http://0.0.0.0:8000/v1/hub/cappart/test-capp/0.0.1  -F "file=@/cdrive/handson/desync/ukaa.caidx"
func (r *Router) artifactPutHandler(c *gin.Context) (*apb.StoreArtifactResponse, error) {

	req := &ArtifactUploadRequest{}
	req.ArtifactName = c.Param("filename")
	req.ArtifactType = c.Param("type")
	req.Version = c.Param("version")
	log.Infof("Adding artifact %s with version %s of type : %s", req.ArtifactName, req.Version, req.ArtifactType)

	_, err := r.parseVersion(req.Version)
	if err != nil {
		return nil, err
	}

	bufReader := NewBufReader(c.Request.Body)
	defer c.Request.Body.Close()

	uncompressedStream, err := gzip.NewReader(bufReader)
	if err != nil {
		log.Infof("Failed to read gz file: %v", err)

		return nil, rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  "Not a tar.gz file",
		}
	}

	tr := tar.NewReader(uncompressedStream)

	_, err = tr.Next()
	if err != nil {
		log.Infof("Failed to read tar file: %v", err)

		return nil, rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  "Not a tar.gz file",
		}
	}

	bufReader.Reset()

	// Call gRPresp, err :C client to send the file
	return r.clients.a.StoreArtifact(&apb.StoreArtifactRequest{
		Name:    req.ArtifactName,
		Type:    apb.ArtifactType(apb.ArtifactType_value[req.ArtifactType]),
		Version: req.Version,
		Data:    bufReader.buff,
	})

	// loc, err := r.storage.PutFile(ctx, name, v, pkg.TarGzExtension, bufReader)
	// if err != nil {
	// 	log.Errorf("Error adding artifact: %s %s", name, ver)

	// 	return err
	// }

	/* Call chunker */
	// go func() {
	// 	err = r.chunker.Chunk(name, v, loc)
	// 	if err != nil {
	// 		log.Errorf("Error chunking artifact: %s %s. Error: %+v", name, ver, err)
	// 	}
	// }()

}

func (r *Router) artifactListVersionsHandler(c *gin.Context, req *ArtifactVersionListRequest) (*apb.GetArtifactVersionListResponse, error) {
	log.Infof("Getting version list: %s of type %s", req.Name, req.ArtifactType)

	return r.clients.a.GetArtifactVersionList(&apb.GetArtifactVersionListRequest{
		Name: req.Name,
		Type: apb.ArtifactType(apb.ArtifactType_value[req.ArtifactType]),
	})
}

func (r *Router) artifactLocationHandler(c *gin.Context, req *ArtifactLocationRequest) (*apb.GetArtifactLocationResponse, error) {
	log.Infof("Getting location for %s version %s  of type %s", req.Name, req.Version, req.ArtifactType)

	return r.clients.a.GetArtifactLocation(&apb.GetArtifactLocationRequest{
		Name:    req.Name,
		Type:    apb.ArtifactType(apb.ArtifactType_value[req.ArtifactType]),
		Version: req.Version,
	})

}

func (r *Router) listArtifactsHandler(c *gin.Context, req *ArtifactListRequest) (*apb.ListArtifactResponse, error) {
	log.Infof("Getting list of artifacts of type %s", req.ArtifactType)

	return r.clients.a.ListArtifacts(&apb.ListArtifactRequest{
		Type: apb.ArtifactType(apb.ArtifactType_value[req.ArtifactType]),
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

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
