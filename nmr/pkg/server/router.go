package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/openIoR/services/common/rest"
	"github.com/ukama/openIoR/services/common/sql"
	"github.com/ukama/openIoR/services/common/ukama"
	"github.com/ukama/openIoR/services/factory/nmr/cmd/version"
	"github.com/ukama/openIoR/services/factory/nmr/internal/db"
	"github.com/ukama/openIoR/services/factory/nmr/pkg"
	rs "github.com/ukama/openIoR/services/factory/nmr/pkg/router"
	"github.com/wI2L/fizz"
)

const (
	NodePath   = "/node"
	ModulePath = "/module"
)

type Router struct {
	fizz       *fizz.Fizz
	port       int
	R          *rs.RouterServer
	nodeRepo   db.NodeRepo
	moduleRepo db.ModuleRepo
}

func (r *Router) Run() {
	logrus.Info("Listening on port ", r.port)
	err := r.fizz.Engine().Run(fmt.Sprint(":", r.port))
	if err != nil {
		panic(err)
	}
}

func NewRouter(config *pkg.Config, rs *rs.RouterServer, d sql.Db) *Router {

	f := rest.NewFizzRouter(&config.Server, pkg.ServiceName, version.Version, pkg.IsDebugMode)
	nodeRepo := db.NewNodeRepo(d)

	moduleRepo := db.NewModuleRepo(d)

	r := &Router{fizz: f,
		port:       config.Server.Port,
		R:          rs,
		nodeRepo:   nodeRepo,
		moduleRepo: moduleRepo,
	}

	r.init()
	return r
}

func (r *Router) init() {
	node := r.fizz.Group(NodePath, "Node", "Node related operations")
	node.GET("/", nil, tonic.Handler(r.GetNodeHandler, http.StatusOK))
	node.PUT("/", nil, tonic.Handler(r.PutNodeHandler, http.StatusOK))
	node.GET("/status", nil, tonic.Handler(r.GetNodeStatusHandler, http.StatusOK))
	node.PUT("/status", nil, tonic.Handler(r.PutNodeStatusHandler, http.StatusOK))
	node.GET("/all", nil, tonic.Handler(r.GetNodeListHandler, http.StatusOK))

	module := r.fizz.Group(ModulePath, "Module", "Module related operations")
	module.GET("/", nil, tonic.Handler(r.GetModuleHandler, http.StatusOK))
	module.PUT("/", nil, tonic.Handler(r.PutModuleHandler, http.StatusOK))
	module.GET("/status", nil, tonic.Handler(r.GetModuleStatusHandler, http.StatusOK))
	module.PUT("/status", nil, tonic.Handler(r.PutModuleStatusHandler, http.StatusOK))
	module.GET("/all", nil, tonic.Handler(r.GetModuleListHandler, http.StatusOK))
	module.GET("/data", nil, tonic.Handler(r.GetModuleDataHandler, http.StatusOK))
	module.PUT("/data", nil, tonic.Handler(r.PutModuleDataHandler, http.StatusOK))
}

func (r *Router) GetNodeHandler(c *gin.Context, req *ReqGetNode) (*RespGetNode, error) {
	logrus.Debugf("Handling NMR get request %+v.", req)

	node, err := r.nodeRepo.GetNode(ukama.NodeID(req.NodeID))
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	resp := &RespGetNode{
		Node: node,
	}

	return resp, nil
}

func (r *Router) PutNodeHandler(c *gin.Context, req *ReqAddOrUpdateNode) error {
	logrus.Debugf("Handling NMR get request %+v.", req)

	node := &db.Node{
		NodeID:         ukama.NodeID(req.NodeID),
		Type:           req.Type,
		PartNumber:     req.PartNumber,
		Skew:           req.Skew,
		SwVersion:      req.SwVersion,
		PSwVersion:     req.PSwVersion,
		Mac:            req.Mac,
		AssemblyDate:   req.AssemblyDate,
		Status:         req.Status,
		OemName:        req.OemName,
		ProdTestStatus: req.ProdTestStatus,
		ProdReport:     req.ProdReport,
	}

	err := r.nodeRepo.AddOrUpdateNode(node)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	return nil
}

func (r *Router) GetNodeStatusHandler(c *gin.Context, req *ReqGetNodeStatus) (*RespGetNodeStatus, error) {
	logrus.Debugf("Handling NMR get request %+v.", req)

	status, err := r.nodeRepo.GetNodeStatus(ukama.NodeID(req.NodeID))
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	resp := &RespGetNodeStatus{
		Status: *status,
	}

	return resp, nil
}

func (r *Router) PutNodeStatusHandler(c *gin.Context, req *ReqUpdateNodeStatus) error {
	logrus.Debugf("Handling NMR get request %+v.", req)

	err := r.nodeRepo.UpdateNodeStatus(ukama.NodeID(req.NodeID), req.Status)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	return nil
}

func (r *Router) GetNodeListHandler(c *gin.Context, req *ReqGetNodeList) ([]RespGetNode, error) {
	logrus.Debugf("Handling NMR get request %+v.", req)

	return nil, nil
}

func (r *Router) GetModuleHandler(c *gin.Context, req *ReqGetModule) error {
	logrus.Debugf("Handling NMR get request %+v.", req)

	return nil
}

func (r *Router) PutModuleHandler(c *gin.Context, req *ReqAddOrUpdateModule) error {
	logrus.Debugf("Handling NMR get request %+v.", req)

	return nil
}

func (r *Router) GetModuleStatusHandler(c *gin.Context, req *ReqGetModuleStatusData) (*RespUpdateModuleStatusData, error) {
	logrus.Debugf("Handling NMR get request %+v.", req)

	return nil, nil
}

func (r *Router) PutModuleStatusHandler(c *gin.Context, req *ReqUpdateModuleStatusData) error {
	logrus.Debugf("Handling NMR get request %+v.", req)

	return nil
}

func (r *Router) GetModuleListHandler(c *gin.Context, req *ReqGetModuleList) (*RespGetModuleList, error) {
	logrus.Debugf("Handling NMR get request %+v.", req)

	return nil, nil
}

func (r *Router) GetModuleDataHandler(c *gin.Context, req *ReqGetNode) (*RespGetNode, error) {
	logrus.Debugf("Handling NMR get request %+v.", req)

	return nil, nil

}

func (r *Router) PutModuleDataHandler(c *gin.Context, req *ReqGetNode) error {
	logrus.Debugf("Handling NMR get request %+v.", req)

	return nil
}
