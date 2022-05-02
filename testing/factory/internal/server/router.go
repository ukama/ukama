package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/rest"
	sr "github.com/ukama/ukama/services/common/srvcrouter"
	"github.com/ukama/ukama/testing/factory/cmd/version"
	"github.com/ukama/ukama/testing/factory/internal"
	"github.com/ukama/ukama/testing/factory/internal/worker"
	"github.com/wI2L/fizz"
)

const (
	NodePath = "/node"
)

type Router struct {
	fizz *fizz.Fizz
	port int
	w    *worker.Worker
}

func (r *Router) Run(close chan error) {
	logrus.Info("Listening on port ", r.port)
	err := r.fizz.Engine().Run(fmt.Sprint(":", r.port))
	if err != nil {
		close <- err
	}
	close <- nil
}

func NewRouter(config *internal.Config, svcR *sr.ServiceRouter) *Router {

	f := rest.NewFizzRouter(&config.Server, internal.ServiceName, version.Version, internal.IsDebugMode)

	worker := worker.NewWorker(svcR)
	r := &Router{fizz: f,
		port: config.Server.Port,
		w:    worker,
	}

	r.init()

	return r
}

func (r *Router) init() {
	node := r.fizz.Group(NodePath, "Node", "Node related operations")
	node.PUT("/", nil, tonic.Handler(r.PostBuildNode, http.StatusAccepted))
	node.DELETE("/", nil, tonic.Handler(r.DeleteNode, http.StatusOK))
}

func (r *Router) PostBuildNode(c *gin.Context, req *ReqBuildNode) (*RespBuildNode, error) {
	logrus.Debugf("Handling buid new node request %+v.", req)
	list, err := r.w.WorkOnBuildOrder(req.Type, req.Count)
	resp := &RespBuildNode{
		NodeIDList: list,
	}

	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  err.Error(),
		}
	}

	return resp, nil
}

func (r *Router) DeleteNode(c *gin.Context, req *ReqDeleteNode) error {
	logrus.Debugf("Handling delete node %+v.", req)

	return nil
}
