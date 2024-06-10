/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

package server

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/errors"
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
	GetArtifcatVersionList(in *apb.GetArtifactVersionListRequest) (*apb.GetArtifactVersionListResponse, error)
	ListArtifacts(in *apb.ListArtifactRequest) (*apb.ListArtifactResponse, error)
}

type distributor interface {
	Store(in *dpb.Request) (*dpb.Response, error)
	Get(in *dpb.ChunkRequest) (*dpb.ChunkResponse, error)
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
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect=true")
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
		capps := auth.Group("/hub", "Artifact store", "Artifact operations")
		capps.GET("/:type/:name/:filename", formatDoc("Get Artifact", "Get artifact of name and type"), tonic.Handler(r.artifactGetHandler, http.StatusOK))
		capps.PUT("/:type/:name/:version", formatDoc("Upload artifact", "Upload a artifact file"), tonic.Handler(r.artifactPutHandler, http.StatusCreated))
		capps.GET("/:type/:name", formatDoc("List of versions for artifcat", "List all the available version and location info for artifact"), tonic.Handler(r.artifactListVersionsHandler, http.StatusOK))
		capps.GET("/:type", formatDoc("List all artifact", "List all artifact of the matching type"), tonic.Handler(r.listArtifactsHandler, http.StatusOK))
		capps.GET("/location/:type/:name/:filename/:version", formatDoc("Location", "Provide a location of the artifact to download from"), tonic.Handler(r.artifactLocationHandler, http.StatusOK))
	}
}

func (r *Router) artifactGetHandler(c *gin.Context, req *ArtifactRequest) error {
	log.Infof("Getting artifact %s of type %s", req.Name, req.ArtifactType)

	v, ext, err := parseArtifactName(req.Name)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  "Artifact name is not valid",
		}
	}

	resp, err := r.clients.a.GetArtifact(&apb.GetArtifactRequest{
		Name:    req.Name,
		Type:    apb.ArtifactType(apb.ArtifactType_value[req.ArtifactType]),
		Version: v.String(),
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
		fmt.Sprintf("attachment; filename=%s-%s%s", resp.Name, resp.Version.Version, ext))

	return nil
}

func parseArtifactName(name string) (ver *semver.Version, ext string, err error) {
	if strings.HasSuffix(name, pkg.TarGzExtension) {
		name = strings.TrimSuffix(name, pkg.TarGzExtension)
		ext = pkg.TarGzExtension
	} else if strings.HasSuffix(name, pkg.ChunkIndexExtension) {
		name = strings.TrimSuffix(name, pkg.ChunkIndexExtension)
		ext = pkg.ChunkIndexExtension
	} else {
		return nil, "", fmt.Errorf("Unsupported extension")
	}

	ver, err = semver.NewVersion(name)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to parse version")
	}

	return ver, ext, nil
}

func (r *Router) artifactPutHandler(c *gin.Context, req *ArtifactUploadRequest) (*apb.StoreArtifactResponse, error) {
	log.Infof("Adding artifact %s with version %s of type : %s %s", req.ArtifactName, req.Version, req.ArtifactType)
	_, err := r.parseVersion(req.Version)
	if err != nil {
		return nil, err
	}

	file, err := c.FormFile("file")
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  fmt.Sprintf("get form err: %s", err.Error()),
		}
	}

	// Open the uploaded file
	uploadedFile, err := file.Open()
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  fmt.Sprintf("file open error: %s", err.Error()),
		}
	}
	defer uploadedFile.Close()

	// Validate if the file is proper gzip data
	if !r.isValidGzip(uploadedFile) {
		return nil, rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  fmt.Sprintf("The uploaded file is not valid gzip data", err.Error()),
		}
	}

	data, err := os.ReadFile(file.Filename)
	if err != nil {
		log.Fatalf("could not read file: %v", err)
	}

	// Call gRPresp, err :C client to send the file
	return r.clients.a.StoreArtifact(&apb.StoreArtifactRequest{
		Name:    req.ArtifactName,
		Type:    apb.ArtifactType(apb.ArtifactType_value[req.ArtifactType]),
		Version: req.Version,
		Data:    data,
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
	log.Infof("Getting version list: %s of type ", req.Name, req.ArtifactType)

	return r.clients.a.GetArtifcatVersionList(&apb.GetArtifactVersionListRequest{
		Name: req.Name,
		Type: apb.ArtifactType(apb.ArtifactType_value[req.ArtifactType]),
	})
}

func (r *Router) artifactLocationHandler(c *gin.Context, req *ArtifactLocationRequest) (*apb.GetArtifactLocationResponse, error) {
	log.Infof("Getting location for %s version %s  of type %s", req.Name, req.ArtifactType)

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

// isValidGzip checks if the file is valid gzip data
func (r *Router) isValidGzip(file io.Reader) bool {
	_, err := gzip.NewReader(file)
	return err == nil
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
