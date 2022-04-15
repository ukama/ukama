package server

import (
	"bytes"
	ds "database/sql"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/openIoR/services/common/rest"
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

func NewRouter(config *pkg.Config, rs *rs.RouterServer, nodeRepo db.NodeRepo, moduleRepo db.ModuleRepo) *Router {

	f := rest.NewFizzRouter(&config.Server, pkg.ServiceName, version.Version, pkg.IsDebugMode)

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
	node.DELETE("/", nil, tonic.Handler(r.DeleteNodeHandler, http.StatusOK))
	node.GET("/status", nil, tonic.Handler(r.GetNodeStatusHandler, http.StatusOK))
	node.PUT("/status", nil, tonic.Handler(r.PutNodeStatusHandler, http.StatusOK))
	node.PUT("/mfgstatus", nil, tonic.Handler(r.PutNodeMfgTestStatusHandler, http.StatusOK))
	node.GET("/mfgstatus", nil, tonic.Handler(r.GetNodeMfgTestStatusHandler, http.StatusOK))
	node.GET("/all", nil, tonic.Handler(r.GetNodeListHandler, http.StatusOK))

	module := r.fizz.Group(ModulePath, "Module", "Module related operations")
	module.GET("/", nil, tonic.Handler(r.GetModuleHandler, http.StatusOK))
	module.PUT("/", nil, tonic.Handler(r.PutModuleHandler, http.StatusOK))
	module.DELETE("/", nil, tonic.Handler(r.DeleteModuleHandler, http.StatusOK))
	module.GET("/all", nil, tonic.Handler(r.GetModuleListHandler, http.StatusOK))
	module.PUT("/assign", nil, tonic.Handler(r.PutAssignModuleToNode, http.StatusOK))
	module.GET("/status", nil, tonic.Handler(r.GetModuleMfgStatusHandler, http.StatusOK))
	module.PUT("/status", nil, tonic.Handler(r.PutModuleMfgStatusHandler, http.StatusOK))
	module.GET("/field", nil, tonic.Handler(r.GetModuleMfgFieldHandler, http.StatusOK))
	module.PUT("/field", nil, tonic.Handler(r.PutModuleMfgFieldHandler, http.StatusOK))
	module.GET("/data", nil, tonic.Handler(r.GetModuleMfgDataHandler, http.StatusOK))
	module.DELETE("/bootstrapcerts", nil, tonic.Handler(r.DeleteBootstrapCertsHandler, http.StatusOK))

}

func (r *Router) GetNodeHandler(c *gin.Context, req *ReqGetNode) (*RespGetNode, error) {
	logrus.Debugf("Handling NMR node info request %+v.", req)

	node, err := r.nodeRepo.GetNode(req.NodeID)
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

	_, err := db.MfgState(req.Status)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  err.Error(),
		}
	}

	_, err = db.MfgTestState(req.MfgTestStatus)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  err.Error(),
		}
	}

	node := &db.Node{
		NodeID:        req.NodeID,
		Type:          req.Type,
		PartNumber:    req.PartNumber,
		Skew:          req.Skew,
		SwVersion:     req.SwVersion,
		PSwVersion:    req.PSwVersion,
		Mac:           req.Mac,
		AssemblyDate:  req.AssemblyDate,
		Status:        req.Status,
		OemName:       req.OemName,
		MfgTestStatus: req.MfgTestStatus,
	}

	err = r.nodeRepo.AddOrUpdateNode(node)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	return nil
}

func (r *Router) DeleteNodeHandler(c *gin.Context, req *ReqDeleteNode) error {
	logrus.Debugf("Handling NMR delete node %+v.", req)

	err := r.nodeRepo.DeleteNode(req.NodeID)
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

	status, err := r.nodeRepo.GetNodeStatus(req.NodeID)
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	resp := &RespGetNodeStatus{}

	if status != nil {
		resp.Status = string(*status)
	}

	return resp, nil
}

func (r *Router) PutNodeStatusHandler(c *gin.Context, req *ReqUpdateNodeStatus) error {
	logrus.Debugf("Handling NMR update node request %+v.", req)

	status, err := db.MfgState(req.Status)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  err.Error(),
		}
	}

	err = r.nodeRepo.UpdateNodeStatus(req.NodeID, *status)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	return nil
}

func (r *Router) PutNodeMfgTestStatusHandler(c *gin.Context, req *ReqUpdateNodeMfgStatus) error {
	logrus.Debugf("Handling NMR update node request %+v.", req)

	_, err := db.MfgState(req.Status)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  err.Error(),
		}
	}

	_, err = db.MfgTestState(req.MfgTestStatus)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  err.Error(),
		}
	}

	node := &db.Node{
		NodeID:        req.NodeID,
		Status:        req.Status,
		MfgTestStatus: req.MfgTestStatus,
		MfgReport:     (*[]byte)(req.MfgReport),
	}

	err = r.nodeRepo.UpdateNodeMfgTestStatus(node)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	return nil
}

func (r *Router) GetNodeMfgTestStatusHandler(c *gin.Context, req *ReqGetNodeMfgStatus) (*RespGetNodeMfgStatus, error) {
	logrus.Debugf("Handling NMR get node status request %+v.", req)

	status, mfg, err := r.nodeRepo.GetNodeMfgTestStatus(req.NodeID)
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	resp := &RespGetNodeMfgStatus{}

	if status != nil {
		resp.Status = string(*status)
	}

	if mfg != nil {
		resp.ProdReport = (*json.RawMessage)(mfg)
	}

	return resp, nil
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

	resp := &RespGetNodeList{}
	if list != nil {
		resp.NodeList = *list
	}

	return resp, nil

}

func (r *Router) GetModuleHandler(c *gin.Context, req *ReqGetModule) (*RespGetModule, error) {
	logrus.Debugf("Handling NMR module info request %+v.", req)

	module, err := r.moduleRepo.GetModule(req.ModuleID)
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

	sqlUnitId := ds.NullString{
		Valid:  true,
		String: req.UnitID,
	}

	if req.UnitID == "" {
		sqlUnitId.Valid = false
	}

	_, err := db.MfgState(req.Status)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  err.Error(),
		}
	}

	module := &db.Module{
		ModuleID:   req.ModuleID,
		Type:       req.Type,
		PartNumber: req.PartNumber,
		SwVersion:  req.SwVersion,
		PSwVersion: req.PSwVersion,
		Mac:        req.Mac,
		MfgDate:    req.MfgDate,
		MfgName:    req.MfgName,
		Status:     req.Status,
		UnitID:     sqlUnitId,
	}

	err = r.moduleRepo.UpsertModule(module)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

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

	resp := &RespGetModuleList{}
	if modules != nil {
		resp.Modules = *modules
	}

	return resp, nil
}

func (r *Router) PutAssignModuleToNode(c *gin.Context, req *ReqAssignModuleToNode) error {
	logrus.Debugf("Handling NMR assign Module to node request %+v.", req)

	err := r.moduleRepo.UpdateNodeId(req.ModuleID, req.NodeID)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}

	}

	return nil
}

/* Delete module from the db */
func (r *Router) DeleteModuleHandler(c *gin.Context, req *ReqDeleteModule) error {
	logrus.Debugf("Handling NMR delete module %+v.", req)

	err := r.moduleRepo.DeleteModule(req.ModuleID)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}
	return nil
}

/* Read module mfg status */
func (r *Router) GetModuleMfgStatusHandler(c *gin.Context, req *ReqGetModuleMfgStatusData) (*RespGetModuleMfgStatusData, error) {
	logrus.Debugf("Handling NMR get module mfg status request %+v.", req)

	status, err := r.moduleRepo.GetModuleMfgStatus(string(req.ModuleID))
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	resp := &RespGetModuleMfgStatusData{}

	if status != nil {
		resp.Status = string(*status)
	}

	return resp, nil
}

/* Update module mfg status */
func (r *Router) PutModuleMfgStatusHandler(c *gin.Context, req *ReqUpdateModuleMfgStatusData) error {
	logrus.Debugf("Handling NMR update Module Mfg Status Data request %+v.", req)

	status, err := db.MfgState(req.Status)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  err.Error(),
		}
	}
	err = r.moduleRepo.UpdateModuleMfgStatus(req.ModuleID, *status)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	return nil
}

/* Convert mfg data into base64 encided json string */
func bytesToJsonCompatibleRawMsg(data *[]byte) json.RawMessage {
	es := []byte{34} //""
	if data == nil {
		return []byte{}
	}
	base := b64.StdEncoding.EncodeToString(*data)
	if bytes.HasPrefix(*data, es) && bytes.HasSuffix(*data, es) {
		return (json.RawMessage)(base)
	} else {
		jdata := append(es, base...)
		msg := append(jdata, es...)
		return (json.RawMessage)(msg)
	}
}

/* Read all the mfg data */
func (r *Router) GetModuleMfgDataHandler(c *gin.Context, req *ReqGetModuleMfgData) (*RespGetModuleMfgData, error) {
	logrus.Debugf("Handling NMR get module mfg data request %+v.", req)

	module, err := r.moduleRepo.GetModule(req.ModuleID)
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	if _, err := db.MfgState(module.Status); err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  err.Error(),
		}
	}

	resp := &RespGetModuleMfgData{
		MfgTestStatus:      module.Status,
		MfgReport:          bytesToJsonCompatibleRawMsg(module.MfgReport),
		BootstrapCerts:     bytesToJsonCompatibleRawMsg(module.BootstrapCerts),
		UserCalibration:    bytesToJsonCompatibleRawMsg(module.UserCalibration),
		FactoryCalibration: bytesToJsonCompatibleRawMsg(module.FactoryCalibration),
		UserConfig:         bytesToJsonCompatibleRawMsg(module.UserConfig),
		FactoryConfig:      bytesToJsonCompatibleRawMsg(module.FactoryConfig),
		InventoryData:      bytesToJsonCompatibleRawMsg(module.InventoryData),
	}

	logrus.Tracef("Read data is %+v", resp)
	return resp, nil
}

/* Read specific mfg data */
func (r *Router) GetModuleMfgFieldHandler(c *gin.Context, req *ReqGetModuleMfgField) error {
	logrus.Debugf("Handling NMR get Module mfg field data request %+v.", req)
	var columnName string
	var module *db.Module
	var data *[]byte
	var err error
	switch req.Field {
	case "mfg_report":
		columnName = "mfg_report"
		module, err = r.moduleRepo.GetModuleMfgField(req.ModuleID, columnName)
		if err != nil {
			return rest.HttpError{
				HttpCode: http.StatusNotFound,
				Message:  err.Error(),
			}
		}
		data = module.MfgReport

	case "bootstrap_cert":
		columnName = "bootstrap_certs"
		module, err = r.moduleRepo.GetModuleMfgField(req.ModuleID, columnName)
		if err != nil {
			return rest.HttpError{
				HttpCode: http.StatusNotFound,
				Message:  err.Error(),
			}
		}
		data = module.BootstrapCerts

	case "user_config":
		columnName = "user_config"
		module, err = r.moduleRepo.GetModuleMfgField(req.ModuleID, columnName)
		if err != nil {
			return rest.HttpError{
				HttpCode: http.StatusNotFound,
				Message:  err.Error(),
			}
		}
		data = module.UserConfig

	case "factory_config":
		columnName = "factory_config"
		module, err = r.moduleRepo.GetModuleMfgField(req.ModuleID, columnName)
		if err != nil {
			return rest.HttpError{
				HttpCode: http.StatusNotFound,
				Message:  err.Error(),
			}
		}
		data = module.FactoryConfig

	case "user_calibration":
		columnName = "user_calibration"
		module, err = r.moduleRepo.GetModuleMfgField(req.ModuleID, columnName)
		if err != nil {
			return rest.HttpError{
				HttpCode: http.StatusNotFound,
				Message:  err.Error(),
			}
		}
		data = module.UserCalibration

	case "factory_calibration":
		columnName = "factory_calibration"
		module, err = r.moduleRepo.GetModuleMfgField(req.ModuleID, columnName)
		if err != nil {
			return rest.HttpError{
				HttpCode: http.StatusNotFound,
				Message:  err.Error(),
			}
		}
		data = module.FactoryCalibration

	// case "cloud_certs":
	// 	columnName = "cloud_certs"
	// 	module, err = r.moduleRepo.GetModuleMfgField(req.ModuleID, req.Field)
	// 	if err != nil {
	// 		return nil, rest.HttpError{
	// 			HttpCode: http.StatusNotFound,
	// 			Message:  err.Error(),
	// 		}
	// 	}
	// 	data = module.CloudCerts

	case "inventory_data":
		columnName = "inventory_data"
		module, err = r.moduleRepo.GetModuleMfgField(req.ModuleID, columnName)
		if err != nil {
			return rest.HttpError{
				HttpCode: http.StatusNotFound,
				Message:  err.Error(),
			}
		}
		data = module.InventoryData

	default:
		return rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  "field type not supported",
		}

	}

	if data == nil {
		return rest.HttpError{
			HttpCode: http.StatusNoContent,
			Message:  "No content found",
		}
	}

	_, err = io.Copy(c.Writer, bytes.NewReader(*data))
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  "Artifact not found",
		}
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s-%s%s", req.ModuleID, req.Field, ".file"))

	return nil

}

/* Update specific mfg data */
func (r *Router) PutModuleMfgFieldHandler(c *gin.Context) error {

	var req ReqUpdateModuleMfgField
	var columnName string
	var err error
	module := db.Module{}

	err = c.BindQuery(&req)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  err.Error(),
		}
	}

	byteBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  err.Error(),
		}
	}

	logrus.Debugf("Handling NMR update Module mfg field request %+v data %+v.", req, byteBody)

	switch req.Field {
	case "mfg_report":
		columnName = "mfg_report"
		module.MfgReport = &byteBody

	case "bootstrap_cert":
		columnName = "bootstrap_certs"
		module.BootstrapCerts = &byteBody

	case "user_config":
		columnName = "user_config"
		module.UserConfig = &byteBody
	case "factory_config":
		columnName = "factory_config"
		module.FactoryConfig = &byteBody

	case "user_calibration":
		columnName = "user_calibration"
		module.UserCalibration = &byteBody

	case "factory_calibration":
		columnName = "factory_calibration"
		module.FactoryCalibration = &byteBody

	// case "cloud_certs":
	// 	columnName = "cloud_certs"
	// 	module, err = r.moduleRepo.GetModuleMfgField(req.ModuleID, req.Field)
	// 	if err != nil {
	// 		return nil, rest.HttpError{
	// 			HttpCode: http.StatusNotFound,
	// 			Message:  err.Error(),
	// 		}
	// 	}
	// 	data = module.CloudCerts

	case "inventory_data":
		columnName = "inventory_data"
		module.InventoryData = &byteBody

	default:
		return rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  "field type not supported",
		}

	}

	err = r.moduleRepo.UpdateModuleMfgField(req.ModuleID, columnName, module)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}

	return nil
}

func (r *Router) DeleteBootstrapCertsHandler(c *gin.Context, req *ReqUpdateModuleBootStrapCerts) error {
	logrus.Debugf("Handling NMR delete bootstrap certs %+v.", req)

	err := r.moduleRepo.DeleteBootstrapCert(req.ModuleID)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  err.Error(),
		}
	}
	return nil
}
