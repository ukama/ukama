package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/rest"
	sr "github.com/ukama/ukama/services/common/srvcrouter"
	"github.com/ukama/ukama/testing/network/cmd/version"
	"github.com/ukama/ukama/testing/network/internal"
	"github.com/ukama/ukama/testing/network/internal/controller"
	"github.com/ukama/ukama/testing/network/internal/db"

	"github.com/wI2L/fizz"
)

const (
	NodePath = "/node"
	ListPath = "/list"
)

type Router struct {
	fizz *fizz.Fizz
	port int
	repo db.VNodeRepo
	c    *controller.Controller
}

func (r *Router) Run(close chan error) {
	logrus.Info("Listening on port ", r.port)
	err := r.fizz.Engine().Run(fmt.Sprint(":", r.port))
	if err != nil {
		close <- err
	}
	close <- nil
}

func NewRouter(config *internal.Config, svcR *sr.ServiceRouter, vNodeRepo db.VNodeRepo) *Router {

	f := rest.NewFizzRouter(&config.Server, internal.ServiceName, version.Version, internal.IsDebugMode)

	r := &Router{fizz: f,
		port: config.Server.Port,
		repo: vNodeRepo,
	}

	if svcR != nil {
		r.c = controller.NewController(r.repo)
		if err := r.c.ControllerInit(); err != nil {
			logrus.Errorf("Controller init failed to start watcher for virtual nodes.")
		}
	}

	r.init()

	return r
}

func (r *Router) init() {
	node := r.fizz.Group(NodePath, "Node", "Node related operations")
	node.GET("", nil, tonic.Handler(r.GetInfo, http.StatusOK))
	node.PUT("/poweron", nil, tonic.Handler(r.PutPowerOn, http.StatusOK))
	node.PUT("/poweroff", nil, tonic.Handler(r.PutPowerOff, http.StatusOK))

	list := r.fizz.Group(ListPath, "List", "Virtual Node list")
	list.GET("", nil, tonic.Handler(r.GetList, http.StatusOK))

}

func (r *Router) PutPowerOn(c *gin.Context, req *ReqPowerOnNode) error {
	logrus.Debugf("Handling node power on %+v.", req)
	err := r.c.PowerOnNode(req.NodeID)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  err.Error(),
		}
	}
	return nil
}

func (r *Router) PutPowerOff(c *gin.Context, req *ReqPowerOffNode) error {
	logrus.Debugf("Handling node power off %+v.", req)

	return nil
}

func (r *Router) GetInfo(c *gin.Context, req *ReqGetNode) (*RespGetNode, error) {
	logrus.Debugf("Handling get node info %+v.", req)

	return nil, nil
}

func (r *Router) GetList(c *gin.Context, req *ReqGetNodeList) (*RespGetNodeList, error) {
	logrus.Debugf("Handling get nodes info %+v.", req)

	return nil, nil
}
