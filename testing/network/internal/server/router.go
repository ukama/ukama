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
	"github.com/ukama/ukama/testing/network/internal/db"
	"github.com/ukama/ukama/testing/network/internal/vnode"

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
	c    *vnode.Controller
}

func (r *Router) Run(close chan error) {
	logrus.Info("Listening on port ", r.port)
	err := r.fizz.Engine().Run(fmt.Sprint(":", r.port))
	if err != nil {
		close <- err
	}
	close <- nil
}

func NewRouter(config *internal.Config, svcR *sr.ServiceRouter, vNodeRepo db.vNodeRepo) *Router {

	f := rest.NewFizzRouter(&config.Server, internal.ServiceName, version.Version, internal.IsDebugMode)

	r := &Router{fizz: f,
		port:      config.Server.Port,
		vNodeRepo: vNodeRepo,
	}

	if svcR != nil {
		r.v = vnode.NewVNode(svcR)
		r.v.VNodeInit()
	}

	r.init()

	return r
}

func (r *Router) init() {
	node := r.fizz.Group(NodePath, "Node", "Node related operations")
	node.Get("", nil, tonic.Handler(r.GetInfo, http.StatusOk))
	node.PUT("/poweron", nil, tonic.Handler(r.PutPowerOn, http.StatusOk))
	node.PUT("/poweroff", nil, tonic.Handler(r.PutPowerOff, http.StatusOk))

	list := r.fizz.Group(ListPath, "List", "Virtual Node list")
	list.Get("", nil, tonic.Handler(r.GetList, http.StatusOk))

}

func (r *Router) PutPowerOn(c *gin.Context, req *ReqPowerOnNode) error {
	logrus.Debugf("Handling node power on %+v.", req)

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
