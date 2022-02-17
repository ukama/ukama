package server

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukamaX/common/rest"
	"github.com/ukama/ukamaX/hub/cmd/version"
	"github.com/ukama/ukamaX/hub/pkg"
	"github.com/wI2L/fizz"
	"io"
	"net/http"
	"time"
)

type Router struct {
	fizz                  *fizz.Fizz
	port                  int
	storage               pkg.Storage
	storageRequestTimeout time.Duration
	chunker               pkg.Chunker
}

func (r *Router) Run() {
	logrus.Info("Listening on port ", r.port)
	err := r.fizz.Engine().Run(fmt.Sprint(":", r.port))
	if err != nil {
		panic(err)
	}
}

func NewRouter(config *rest.HttpConfig, storage pkg.Storage, chunker pkg.Chunker, storageTimeout time.Duration) *Router {
	f := rest.NewFizzRouter(config, pkg.ServiceName, version.Version, pkg.IsDebugMode)

	r := &Router{
		fizz:                  f,
		port:                  config.Port,
		storage:               storage,
		storageRequestTimeout: storageTimeout,
		chunker:               chunker}
	r.init()
	return r
}

func (r *Router) init() {
	capps := r.fizz.Group("/capps", "CApps", "CApps operations")
	capps.GET("/:name/:version", nil, tonic.Handler(r.cappGetHandler, http.StatusOK))
	capps.PUT("/:name/:version", nil, tonic.Handler(r.cappPutHandler, http.StatusCreated))
	capps.GET("/:name/", nil, tonic.Handler(r.cappListHandler, http.StatusOK))
	capps.GET("/", nil, tonic.Handler(r.listAllAppsHandler, http.StatusOK))

}

func (r *Router) cappGetHandler(c *gin.Context, req *CAppRequest) error {
	logrus.Infof("Getting artifact: %s %s", req.Name, req.Version)
	ctx, cancel := context.WithTimeout(context.Background(), r.storageRequestTimeout)
	defer cancel()

	v, err := r.parseVersion(req.Version)
	if err != nil {
		return err
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s-%s.tar.gz", req.Name, v.String()))
	rd, err := r.storage.GetFile(ctx, req.Name, v, pkg.TarGzExtension)
	if err != nil {
		return err
	}
	defer rd.Close()

	_, err = io.Copy(c.Writer, rd)
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return rest.HttpError{
				HttpCode: http.StatusNotFound,
				Message:  "Artifact not found",
			}
		}
	}

	return nil
}

func (r *Router) cappPutHandler(c *gin.Context) error {
	name := c.Param("name")
	ver := c.Param("version")
	logrus.Infof("Adding artifact: %s %s", name, ver)
	ctx, cancel := context.WithTimeout(context.Background(), r.storageRequestTimeout)
	defer cancel()

	v, err := r.parseVersion(ver)
	if err != nil {
		return err
	}
	bufReader := NewBufReader(c.Request.Body)
	defer c.Request.Body.Close()

	uncompressedStream, err := gzip.NewReader(bufReader)
	if err != nil {
		logrus.Infof("Failed to read gz file: %v", err)
		return rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  "Not a tar.gz file",
		}
	}

	tr := tar.NewReader(uncompressedStream)
	_, err = tr.Next()
	if err != nil {
		logrus.Infof("Failed to read tar file: %v", err)
		return rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  "Not a tar.gz file",
		}
	}
	bufReader.Reset()

	loc, err := r.storage.PutFile(ctx, name, v, pkg.TarGzExtension, bufReader)
	if err != nil {
		logrus.Errorf("Error adding artifact: %s %s", name, ver)
		return err
	}

	go func() {
		err = r.chunker.Chunk(name, v, loc)
		if err != nil {
			logrus.Errorf("Error chunking artifact: %s %s. Error: %+v", name, ver, err)
		}
	}()

	return nil
}

func (r *Router) cappListHandler(c *gin.Context, req *CAppListRequest) (*CAppListResponse, error) {
	logrus.Infof("Getting version list: %s", req.Name)
	ctx, cancel := context.WithTimeout(context.Background(), r.storageRequestTimeout)
	defer cancel()

	ls, err := r.storage.List(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	return &CAppListResponse{
		Artifacts: ls,
	}, nil
}

func (r *Router) listAllAppsHandler(c *gin.Context) (*CAppsListResponse, error) {
	logrus.Infof("Getting list of apps")
	ctx, cancel := context.WithTimeout(context.Background(), r.storageRequestTimeout)
	defer cancel()

	ls, err := r.storage.ListApps(ctx)
	if err != nil {
		return nil, err
	}

	return &CAppsListResponse{
		Artifacts: ls,
	}, nil
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
