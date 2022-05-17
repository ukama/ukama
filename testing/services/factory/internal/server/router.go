package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/rest"
	sr "github.com/ukama/ukama/services/common/srvcrouter"
	"github.com/ukama/ukama/testing/services/factory/cmd/version"
	"github.com/ukama/ukama/testing/services/factory/internal"
	"github.com/ukama/ukama/testing/services/factory/internal/worker"
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

	r := &Router{fizz: f,
		port: config.Server.Port,
	}

	if svcR != nil {
		r.w = worker.NewWorker(svcR)
		r.w.WorkerInit()
	}

	r.init()

	return r
}

func (r *Router) init() {
	node := r.fizz.Group(NodePath, "Node", "Node related operations")
	node.PUT("", nil, tonic.Handler(r.PostBuildNode, http.StatusAccepted))
}

func (r *Router) PostBuildNode(c *gin.Context, req *ReqBuildNode) (*RespBuildNode, error) {
	logrus.Debugf("Handling buid new node request %+v.", req)

	var err error
	list := []string{}
	resp := &RespBuildNode{
		NodeIDList: list,
	}
	if r.w == nil {
		err = fmt.Errorf("factory worker not initialized")
	} else {
		list, err = r.w.WorkOnBuildOrder(req.Type, req.Count)
		resp.NodeIDList = list
	}

	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  err.Error(),
		}
	}

	return resp, nil
}
