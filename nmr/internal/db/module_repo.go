package db

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukamaX/common/sql"
	"gorm.io/gorm/clause"
)

type ModuleRepo interface {
	AddModule(module *Module) error
	UpsertModule(Module *Module) error
	UpdateNodeId(moduleId string, nodeId string) error
	GetModule(moduleId string) (*Module, error)
	DeleteModule(moduleId string) error
	ListModules() (*[]Module, error)
	GetModuleMfgStatus(moduleId string) (*string, *[]byte, error)
	UpdateModuleMfgStatus(moduleId string, status string, data []byte) error
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

/* Only done when in mfg line when node Id is yet to be decided.
TODO: Get better solutuion */
func (r *moduleRepo) AddModule(module *Module) error {

	//result := r.Db.GetGormDb().Select("module_id", "type", "part_number", "hw_version", "mac", "sw_version", "p_sw_version", "mfg_date", "mfg_name", "mfg_test_status", "mfg_report", "bootstrap_certs", "user_calibration", "factory_calibration", "user_config", "factory_config", "inventory_data", "unit_id").Create(module)
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
	result := r.Db.GetGormDb().Preload(clause.Associations).First(&Module, "module_id = ?", ModuleId)
	if result.Error != nil {
		return nil, result.Error
	}
	return &Module, nil
}

func (r *moduleRepo) UpdateNodeId(moduleId string, nodeId string) error {
	module := Module{}
	result := r.Db.GetGormDb().Model(&module).Where("module_id = ?", moduleId).UpdateColumn("unit_id", nodeId)
	if result.Error != nil {
		logrus.Errorf("This error ss %+v", result)
		return result.Error
	}
	return nil
}

/* Delete Module with module Ip permanently  */
func (r *moduleRepo) DeleteModule(moduleId string) error {
	result := r.Db.GetGormDb().Unscoped().Where("module_id = ?", moduleId).Delete(&Module{})
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
func (r *moduleRepo) GetModuleMfgStatus(moduleId string) (*string, *[]byte, error) {
	var module Module
	result := r.Db.GetGormDb().Select("mfg_test_status", "mfg_report").First(&module, "module_id = ?", moduleId)
	if result.Error != nil {
		return nil, nil, result.Error
	}
	return &module.MfgTestStatus, module.MfgReport, nil
}

/* Update Mfg Status  Data */
func (r *moduleRepo) UpdateModuleMfgStatus(moduleId string, status string, data []byte) error {
	module := Module{
		ModuleID:      moduleId,
		MfgTestStatus: status,
	}

	result := r.Db.GetGormDb().Model(&Module{}).Where("module_id = ?", moduleId).UpdateColumns(module)
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

	// columnName, err := GetModuleDataFieldName(field)
	// if err != nil {
	// 	return nil, err
	// }

	//result := r.Db.GetGormDb().Model(&Module{}).Where("module_id = ?", moduleId).Pluck(*columnName, &data)
	result := r.Db.GetGormDb().Select(columnName).First(&module, "module_id = ?", moduleId)
	if result.Error != nil {
		return nil, result.Error
	}
	logrus.Tracef("Read module mfg field  for %s with %v. result %+v", moduleId, module, result)
	return &module, nil
}

/* Update Module Mfg field data */
func (r *moduleRepo) UpdateModuleMfgField(moduleId string, field string, module Module) error {

	result := r.Db.GetGormDb().Model(&Module{}).Where("module_id = ?", moduleId).UpdateColumns(&module)
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
	result := r.Db.GetGormDb().Where("module_id = ?", moduleId).UpdateColumn("boostrap_cert", nil)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
