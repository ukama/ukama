package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/hub/distributor/cmd/version"
	"github.com/ukama/ukama/systems/hub/distributor/pkg"
	"github.com/ukama/ukama/systems/hub/distributor/pkg/chunk"

	"github.com/Masterminds/semver/v3"
	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"

	log "github.com/sirupsen/logrus"
)

const ChunksPath = "/v1/chunks"

type Router struct {
	fizz  *fizz.Fizz
	port  int
	Store pkg.StoreConfig
	Chunk pkg.ChunkConfig
}

func (r *Router) Run() {
	log.Info("Listening on port ", r.port)

	err := r.fizz.Engine().Run(fmt.Sprint(":", r.port))
	if err != nil {
		panic(err)
	}
}

func NewRouter(config *pkg.Config) *Router {
	f := rest.NewFizzRouter(&config.Server, pkg.ServiceName, version.Version, pkg.IsDebugMode, "")

	r := &Router{fizz: f,
		port:  config.Server.Port,
		Store: config.Distribution.StoreCfg,
		Chunk: config.Distribution.Chunk,
	}

	chunks := f.Group(ChunksPath, "ChunksServer", "Chunks Server for content to be distributed")
	chunks.PUT("/:name/:version", nil, tonic.Handler(r.chunkPutHandler, http.StatusOK))
	chunks.PUT("/", nil, tonic.Handler(chunkRootHandler, http.StatusOK))

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

	log.Debugf("Handling chunking request %+v.", req)

	index, err := chunk.CreateChunks(ctx, &r.Store, &r.Chunk, fname, ver, req.Store)
	if err != nil {
		log.Errorf("Error while chunking the file %s: %s", req.Name, err.Error())

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
			log.Errorf("Error while creating index file.")

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
