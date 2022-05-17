package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/rest"
	sr "github.com/ukama/ukama/services/common/srvcrouter"
	"github.com/ukama/ukama/services/common/ukama"
	"github.com/ukama/ukama/testing/services/network/cmd/version"
	"github.com/ukama/ukama/testing/services/network/internal"
	"github.com/ukama/ukama/testing/services/network/internal/controller"
	"github.com/ukama/ukama/testing/services/network/internal/db"
	"github.com/ukama/ukama/testing/services/network/internal/nmr"

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
	nmr  *nmr.NMR
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
		r.nmr = nmr.NewNMR(svcR)
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
	node.PUT("", nil, tonic.Handler(r.PutNode, http.StatusAccepted))
	node.DELETE("", nil, tonic.Handler(r.DeleteNode, http.StatusOK))

	list := r.fizz.Group(ListPath, "List", "Virtual Node list")
	list.GET("", nil, tonic.Handler(r.GetList, http.StatusOK))

}

func (r *Router) PutNode(c *gin.Context, req *ReqActionOnNode) error {
	logrus.Debugf("Handling Node opertaion on %+v.", req)

	/* validate nodeid */
	_, err := ukama.ValidateNodeId(req.NodeID)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  "Invalid node:" + err.Error(),
		}
	}

	switch req.LookingTo {
	case "vnode_power_on":
		return r.PutPowerOn(&ReqPowerOnNode{*req})
	case "vnode_power_off":
		return r.PutPowerOff(&ReqPowerOffNode{*req})
	default:
		return rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  "invalid opertaion on node",
		}
	}
}

func (r *Router) PutPowerOn(req *ReqPowerOnNode) error {
	logrus.Debugf("Handling node power on %+v.", req)

	/* Validation node from NMR */
	err := r.nmr.NmrLookForNode(req.NodeID)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  err.Error(),
		}
	}

	/* Add to db*/
	err = r.repo.Insert(req.NodeID, db.VNodePreCheck.String())
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  "Adding Node:" + err.Error(),
		}
	}

	/* PowerOn Node */
	err = r.c.PowerOnNode(req.NodeID)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  "PowerOn Node:" + err.Error(),
		}
	}

	return nil
}

func (r *Router) PutPowerOff(req *ReqPowerOffNode) error {
	logrus.Debugf("Handling node power off %+v.", req)

	node, err := r.repo.GetInfo(req.NodeID)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  "Reading Node Info:" + err.Error(),
		}
	}

	if node.Status != db.VNodeOn.String() {
		return rest.HttpError{
			HttpCode: http.StatusForbidden,
			Message:  fmt.Sprintf("Node state: Node %s expected state %s but found in %s state.", req.NodeID, db.VNodeOn, node.Status),
		}
	}

	err = r.c.PowerOffNode(req.NodeID)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  "Powering off:" + err.Error(),
		}
	}

	/* Add to db*/
	err = r.repo.Update(req.NodeID, db.VNodeOff.String())
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  "Updating Node State:" + err.Error(),
		}
	}

	return nil
}

func (r *Router) GetInfo(c *gin.Context, req *ReqGetNode) (*RespGetNode, error) {
	logrus.Debugf("Handling get node info %+v.", req)

	resp := &RespGetNode{
		Runtime: controller.VNodeUnkown,
	}

	/* validate nodeid */
	_, err := ukama.ValidateNodeId(req.NodeID)
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  "Invalid node:" + err.Error(),
		}
	}

	var rstate *string
	node, err := r.repo.GetInfo(req.NodeID)
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  "Reading Node Info:" + err.Error(),
		}
	}

	if node.Status == db.VNodeOn.String() {
		rstate, err = r.c.GetNodeRuntimeStatus(req.NodeID)
		if err != nil {
			return nil, rest.HttpError{
				HttpCode: http.StatusInternalServerError,
				Message:  "Reading Node runtime Info:" + err.Error(),
			}
		}
	}

	if node != nil {
		resp.Node = *node
		if rstate != nil {
			resp.Runtime = *rstate
		}
	}

	return resp, nil
}

func (r *Router) GetList(c *gin.Context, req *ReqGetNodeList) (*RespGetNodeList, error) {
	logrus.Debugf("Handling get nodes info %+v.", req)

	resp := &RespGetNodeList{}

	nodes, err := r.repo.List()
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  "Reading Node Info:" + err.Error(),
		}
	}

	if nodes != nil {
		resp.NodeList = *nodes
	}

	return resp, nil
}

func (r *Router) DeleteNode(c *gin.Context, req *ReqDeleteNode) error {
	logrus.Debugf("Handling delete node info %+v.", req)

	/* validate nodeid */
	_, err := ukama.ValidateNodeId(req.NodeID)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  "Invalid node:" + err.Error(),
		}
	}

	node, err := r.repo.GetInfo(req.NodeID)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  "Failed to find node:" + err.Error(),
		}
	}

	/* If Node is ON */
	if node.Status == db.VNodeOn.String() {
		err = r.c.PowerOffNode(req.NodeID)
		if err != nil {
			return rest.HttpError{
				HttpCode: http.StatusInternalServerError,
				Message:  "Failed to power off node :" + err.Error(),
			}
		}

		err = r.repo.Delete(req.NodeID)
		if err != nil {
			return rest.HttpError{
				HttpCode: http.StatusInternalServerError,
				Message:  "Failed to remove node from db:" + err.Error(),
			}
		}

	}
	return nil
}
