package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

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
	logrus.Debugf("Handling NMR node info request %+v.", req)

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
	logrus.Debugf("Handling NMR adding node request %+v.", req)

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
	logrus.Debugf("Handling NMR get node status request %+v.", req)

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
	logrus.Debugf("Handling NMR update node request %+v.", req)

	err := r.nodeRepo.UpdateNodeStatus(ukama.NodeID(req.NodeID), req.Status)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	return nil
}

func (r *Router) GetNodeListHandler(c *gin.Context, req *ReqGetNodeList) (*RespGetNodeList, error) {
	logrus.Debugf("Handling NMR get list of nodes request %+v.", req)
	list, err := r.nodeRepo.ListNodes()
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	resp := &RespGetNodeList{
		NodeList: *list,
	}

	return resp, nil

}

func (r *Router) GetModuleHandler(c *gin.Context, req *ReqGetModule) (*RespGetModule, error) {
	logrus.Debugf("Handling NMR module info request %+v.", req)

	module, err := r.moduleRepo.GetModule(ukama.NodeID(req.ModuleID))
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	resp := &RespGetModule{
		Module: module,
	}

	return resp, nil
}

func (r *Router) PutModuleHandler(c *gin.Context, req *ReqAddOrUpdateModule) error {
	logrus.Debugf("Handling NMR adding module request %+v.", req)

	date, err := time.Parse("2006-01-02", req.MfgDate)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  err.Error(),
		}
	}

	module := &db.Module{
		ModuleID:       ukama.NodeID(req.ModuleID),
		Type:           req.Type,
		PartNumber:     req.PartNumber,
		SwVersion:      req.SwVersion,
		PSwVersion:     req.PSwVersion,
		Mac:            req.Mac,
		MfgDate:        date,
		MfgName:        req.MfgName,
		ProdTestStatus: req.ProdTestStatus,
		ProdReport:     req.ProdReport,
	}

	err = r.moduleRepo.AddOrUpdateModule(module)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	return nil
}

func (r *Router) GetModuleStatusHandler(c *gin.Context, req *ReqGetModuleStatusData) (*RespUpdateModuleStatusData, error) {
	logrus.Debugf("Handling NMR get module status request %+v.", req)

	module, err := r.moduleRepo.GetModule(ukama.NodeID(req.ModuleID))
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	resp := &RespUpdateModuleStatusData{
		ProdTestStatus:     module.ProdTestStatus,
		ProdReport:         module.ProdReport,
		BootstrapCerts:     module.BootstrapCerts,
		UserCalibrartion:   module.UserCalibrartion,
		FactoryCalibration: module.FactoryCalibration,
		UserConfig:         module.UserConfig,
		FactoryConfig:      module.FactoryConfig,
		InventoryData:      module.InventoryData,
	}

	return resp, nil
}

func (r *Router) PutModuleStatusHandler(c *gin.Context, req *ReqUpdateModuleStatusData) error {
	logrus.Debugf("Handling NMR update Module Status Data request %+v.", req)

	return nil
}

func (r *Router) GetModuleListHandler(c *gin.Context, req *ReqGetModuleList) (*RespGetModuleList, error) {
	logrus.Debugf("Handling NMR get module list request %+v.", req)

	modules, err := r.moduleRepo.ListModules()
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	resp := &RespGetModuleList{
		Modules: *modules,
	}

	return resp, nil
}

//408
func (r *Router) GetModuleDataHandler(c *gin.Context, req *ReqGetModuleStatusData) (*RespUpdateModuleStatusData, error) {
	logrus.Debugf("Handling NMR get module status data request %+v.", req)

	module, err := r.moduleRepo.GetModule(ukama.NodeID(req.ModuleID))
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	resp := &RespUpdateModuleStatusData{
		ProdTestStatus:     module.ProdTestStatus,
		ProdReport:         module.ProdReport,
		BootstrapCerts:     module.BootstrapCerts,
		UserCalibrartion:   module.UserCalibrartion,
		FactoryCalibration: module.FactoryCalibration,
		UserConfig:         module.UserConfig,
		FactoryConfig:      module.FactoryConfig,
		InventoryData:      module.InventoryData,
	}

	return resp, nil
}

func (r *Router) PutModuleDataHandler(c *gin.Context, req *ReqUpdateModuleData) error {
	logrus.Debugf("Handling NMR update Module Status Data request %+v.", req)

	byteBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  err.Error(),
		}
	}

	err = r.moduleRepo.UpdateModuleData(ukama.NodeID(req.ModuleID), req.Field, byteBody)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	return nil
}

func (r *Router) PutModuleProdStatusDataHandler(c *gin.Context, req *ReqUpdateModuleStatusData) error {
	logrus.Debugf("Handling NMR update Module Production Status Data request %+v.", req)

	byteBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  err.Error(),
		}
	}

	err = r.moduleRepo.UpdateModuleProdStatus(ukama.NodeID(req.ModuleID), req.Status, byteBody)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	return nil
}

func (r *Router) GetModuleProdStatusDataHandler(c *gin.Context, req *ReqGetModuleProdStatusData) (*RespGetModuleProdStatusData, error) {
	logrus.Debugf("Handling NMR get Module Production status data request %+v.", req)

	module, err := r.moduleRepo.GetModule(ukama.NodeID(req.ModuleID))
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	resp := &RespGetModuleProdStatusData{
		ProdTestStatus: module.ProdTestStatus,
		ProdReport:     module.ProdReport,
	}

	return resp, nil

}
