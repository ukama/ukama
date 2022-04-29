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
	"github.com/ukama/ukama/testing/factory/internal/db"
	"github.com/wI2L/fizz"
)

const (
	NodePath = "/node"
)

type Router struct {
	fizz     *fizz.Fizz
	port     int
	R        *sr.ServiceRouter
	nodeRepo db.NodeRepo
}

func (r *Router) Run(close chan error) {
	logrus.Info("Listening on port ", r.port)
	err := r.fizz.Engine().Run(fmt.Sprint(":", r.port))
	if err != nil {
		close <- err
	}
	close <- nil
}

func NewRouter(config *internal.Config, svcR *sr.ServiceRouter, nodeRepo db.NodeRepo) *Router {

	f := rest.NewFizzRouter(&config.Server, internal.ServiceName, version.Version, internal.IsDebugMode)

	r := &Router{fizz: f,
		port:     config.Server.Port,
		R:        svcR,
		nodeRepo: nodeRepo,
	}

	r.init()

	return r
}

func (r *Router) init() {
	node := r.fizz.Group(NodePath, "Node", "Node related operations")
	node.GET("/", nil, tonic.Handler(r.GetNode, http.StatusOK))
	node.PUT("/", nil, tonic.Handler(r.PostBuildNode, http.StatusAccepted))
	node.DELETE("/", nil, tonic.Handler(r.DeleteNode, http.StatusOK))
	node.GET("/listall", nil, tonic.Handler(r.GetNodeList, http.StatusOK))
}

func (r *Router) GetNode(c *gin.Context, req *ReqGetNode) (*RespGetNode, error) {
	logrus.Debugf("Handling node request %+v.", req)

	resp := &RespGetNode{}

	return resp, nil
}

func (r *Router) PostBuildNode(c *gin.Context, req *ReqBuildNode) error {
	logrus.Debugf("Handling buid new node request %+v.", req)

	return nil
}

func (r *Router) DeleteNode(c *gin.Context, req *ReqDeleteNode) error {
	logrus.Debugf("Handling delete node %+v.", req)

	return nil
}

func (r *Router) GetNodeList(c *gin.Context, req *ReqGetNodeList) (*RespGetNodeList, error) {
	logrus.Debugf("Handling node request %+v.", req)

	resp := &RespGetNodeList{}

	return resp, nil
}
