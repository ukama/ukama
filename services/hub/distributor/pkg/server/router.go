package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Masterminds/semver/v3"
	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/rest"
	"github.com/ukama/ukama/services/hub/distributor/cmd/version"
	"github.com/ukama/ukama/services/hub/distributor/pkg"
	"github.com/ukama/ukama/services/hub/distributor/pkg/chunk"
	"github.com/wI2L/fizz"
)

type Router struct {
	fizz  *fizz.Fizz
	port  int
	Store pkg.StoreConfig
	Chunk pkg.ChunkConfig
}

func (r *Router) Run() {
	logrus.Info("Listening on port ", r.port)
	err := r.fizz.Engine().Run(fmt.Sprint(":", r.port))
	if err != nil {
		panic(err)
	}
}

func NewRouter(config *pkg.Config) *Router {

	f := rest.NewFizzRouter(&config.Server, pkg.ServiceName, version.Version, pkg.IsDebugMode)

	r := &Router{fizz: f,
		port:  config.Server.Port,
		Store: config.Distribution.StoreCfg,
		Chunk: config.Distribution.Chunk,
	}

	chunk := f.Group("/chunk", "ChunkServer", "Chunk Server for content to be distributed")
	chunk.PUT("/:name/:version", nil, tonic.Handler(r.chunkPutHandler, http.StatusOK))
	chunk.PUT("/", nil, tonic.Handler(chunkRootHandler, http.StatusOK))

	return r
}

func (r *Router) chunkPutHandler(c *gin.Context, req *ChunkRequest) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fname := req.Name

	ver, err := r.parseVersion(req.Version)
	if err != nil {
		return err
	}

	logrus.Debugf("Handling chunking request %+v.", req)
	index, err := chunk.CreateChunks(ctx, &r.Store, &r.Chunk, fname, ver, req.Store)
	if err != nil {
		logrus.Errorf("Error while chunking the file %s: %s", req.Name, err.Error())
		return rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  err.Error(),
		}
	} else {

		if index != nil {
			c.Header("Content-Type", "application/octet-stream")
			_, err = index.WriteTo(c.Writer)

		}

		if err != nil {
			logrus.Errorf("Error while creating index file.")
			return rest.HttpError{
				HttpCode: http.StatusInternalServerError,
				Message:  err.Error(),
			}
		}

	}

	return nil

}

func chunkRootHandler(c *gin.Context, r *ChunkRequest) error {
	c.JSON(http.StatusOK, gin.H{
		"message": "Chunk Server",
	})
	return nil
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
