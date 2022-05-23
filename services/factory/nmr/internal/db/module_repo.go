package db

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/sql"
	"gorm.io/gorm/clause"
)

const (
	MODULE_ID_CLAUSE = "module_id = ?"
)

type ModuleRepo interface {
	AddModule(module *Module) error
	UpsertModule(Module *Module) error
	UpdateNodeId(moduleId string, nodeId string) error
	GetModule(moduleId string) (*Module, error)
	DeleteModule(moduleId string) error
	ListModules() (*[]Module, error)
	GetModuleMfgStatus(moduleId string) (*MfgStatus, error)
	UpdateModuleMfgStatus(moduleId string, status MfgStatus) error
	GetModuleMfgData(moduleId string) (*Module, error)
	GetModuleMfgField(moduleId string, field string) (*Module, error)
	UpdateModuleMfgField(moduleId string, field string, module Module) error
	DeleteBootstrapCert(ModuleId string) error
}

type moduleRepo struct {
	Db sql.Db
}

func GetModuleDataFieldName(field string) (*string, error) {
	var columnName string
	switch field {
	case "bootstrap_cert":
		columnName = "bootstrap_certs"
	case "user_config":
		columnName = "user_config"
	case "factory_config":
		columnName = "factory_config"
	case "user_calibration":
		columnName = "user_calibration"
	case "factory_calibration":
		columnName = "factory_calibration"
	case "cloud_certs":
		columnName = "cloud_certs"
	case "inventory_data":
		columnName = "inventory_data"
	default:
		return nil, fmt.Errorf("not supported field %s", field)
	}
	return &columnName, nil
}

func NewModuleRepo(db sql.Db) *moduleRepo {
	return &moduleRepo{
		Db: db,
	}
}

func (r *moduleRepo) AddModule(module *Module) error {
	result := r.Db.GetGormDb().Create(module)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

/* Upsert is used when we know the node id */
func (r *moduleRepo) UpsertModule(Module *Module) error {
	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "module_id"}},
		UpdateAll: true,
	}).Create(Module)
	return d.Error
}

func (r *moduleRepo) GetModule(ModuleId string) (*Module, error) {
	var Module Module
	result := r.Db.GetGormDb().Preload(clause.Associations).First(&Module, MODULE_ID_CLAUSE, ModuleId)
	if result.Error != nil {
		return nil, result.Error
	}
	return &Module, nil
}

func (r *moduleRepo) UpdateNodeId(moduleId string, nodeId string) error {
	module := Module{}
	result := r.Db.GetGormDb().Model(&module).Where(MODULE_ID_CLAUSE, moduleId).UpdateColumn("unit_id", nodeId)
	if result.Error != nil {
		logrus.Errorf("This error ss %+v", result)
		return result.Error
	}
	return nil
}

/* Delete Module with module Ip permanently  */
func (r *moduleRepo) DeleteModule(moduleId string) error {
	result := r.Db.GetGormDb().Unscoped().Where(MODULE_ID_CLAUSE, moduleId).Delete(&Module{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

/* List all Modules */
func (r *moduleRepo) ListModules() (*[]Module, error) {
	var Modules []Module

	result := r.Db.GetGormDb().Preload(clause.Associations).Find(&Modules)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	} else {
		return &Modules, nil
	}
}

/* Read module mfg status */
func (r *moduleRepo) GetModuleMfgStatus(moduleId string) (*MfgStatus, error) {
	var module Module
	result := r.Db.GetGormDb().Select("status").First(&module, MODULE_ID_CLAUSE, moduleId)
	if result.Error != nil {
		return nil, result.Error
	}

	status, err := MfgState(module.Status)
	if err != nil {
		return nil, err
	}

	return status, nil
}

/* Update Mfg Status  Data */
func (r *moduleRepo) UpdateModuleMfgStatus(moduleId string, status MfgStatus) error {
	module := Module{
		Status: string(status),
	}

	result := r.Db.GetGormDb().Model(&Module{}).Where(MODULE_ID_CLAUSE, moduleId).UpdateColumns(module)
	if result.Error != nil {
		return result.Error
	}

	logrus.Tracef("Updated module mfg status for %s with %v. result %+v", moduleId, module, result)
	return nil
}

/* Read module mfg data */
func (r *moduleRepo) GetModuleMfgData(moduleId string) (*Module, error) {
	var module Module
	result := r.Db.GetGormDb().Select("mfg_test_status", "mfg_report", "bootstrap_certs", "user_calibration", "factory_calibration", "user_config", "factory_config", "inventory_data").First(&module, "module_id = ?", moduleId)
	if result.Error != nil {
		return nil, result.Error
	}
	return &module, nil
}

/* Get Module Mfg field */
func (r *moduleRepo) GetModuleMfgField(moduleId string, columnName string) (*Module, error) {
	var module Module

	result := r.Db.GetGormDb().Select(columnName).First(&module, MODULE_ID_CLAUSE, moduleId)
	if result.Error != nil {
		return nil, result.Error
	}
	logrus.Tracef("Read module mfg field  for %s with %v. result %+v", moduleId, module, result)
	return &module, nil
}

/* Update Module Mfg field data */
func (r *moduleRepo) UpdateModuleMfgField(moduleId string, field string, module Module) error {

	result := r.Db.GetGormDb().Model(&Module{}).Where(MODULE_ID_CLAUSE, moduleId).UpdateColumns(&module)
	if result.Error != nil {
		return result.Error
	}
	logrus.Tracef("Updated module mfg field  for %s with %v. result %+v", moduleId, module, result)
	return nil
}

/* Update Production Bootstrap Cert*/
func (r *moduleRepo) UpdateModuleBootstrapCert(moduleData *Module) error {
	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "module_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"boostrap_cert"}),
	}).Create(moduleData)
	return d.Error
}

/* Update Production Bootstrap Cert*/
func (r *moduleRepo) DeleteBootstrapCert(moduleId string) error {
	result := r.Db.GetGormDb().Model(&Module{}).Where(MODULE_ID_CLAUSE, moduleId).UpdateColumn("bootstrap_certs", nil)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
