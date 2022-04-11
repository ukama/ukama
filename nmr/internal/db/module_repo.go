package db

import (
	"fmt"

	"github.com/ukama/openIoR/services/common/ukama"
	"github.com/ukama/ukamaX/common/sql"
	"gorm.io/gorm/clause"
)

type ModuleRepo interface {
	AddOrUpdateModule(module *Module) error
	GetModule(moduleId ukama.NodeID) (*Module, error)
	DeleteModule(moduleId ukama.NodeID) error
	ListModules() (*[]Module, error)
	UpdateModuleProdStatus(moduleId ukama.NodeID, status string, data []byte) error
	UpdateModuleData(moduleId ukama.NodeID, field string, data []byte) error
	GetModuleData(moduleId ukama.NodeID, field string) ([]byte, error)
	DeleteBootstrapCert(ModuleId *ukama.NodeID) error
}

type moduleRepo struct {
	Db sql.Db
}

func GetModuleDataFieldName(field string) (*string, error) {
	var columnName string
	switch field {
	case "boostrap_cert":
		columnName = "boostrap_cert"
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

func (r *moduleRepo) AddOrUpdateModule(Module *Module) error {
	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "module_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"type", "part_number", "hw_version", "mac", "sw_version", "p_sw_version", "mfg_date", "mfg_name", "prod_test_status", "prod_report", "bootstrap_cert", "user_calib", "factory_calib", "user_config", "factory_config", "inventory_data"}),
	}).Create(Module)
	return d.Error
}

func (r *moduleRepo) GetModule(ModuleId ukama.NodeID) (*Module, error) {
	var Module Module
	result := r.Db.GetGormDb().Preload(clause.Associations).First(&Module, "module_id = ?", ModuleId.StringLowercase())
	if result.Error != nil {
		return nil, result.Error
	}
	return &Module, nil
}

/* Delete Module  */
func (r *moduleRepo) DeleteModule(ModuleId ukama.NodeID) error {
	result := r.Db.GetGormDb().Where("module_id = ?", ModuleId).Delete(&Module{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

/* List all Modules */
func (r *moduleRepo) ListModules() (*[]Module, error) {
	var Modules []Module

	result := r.Db.GetGormDb().Find(&Modules)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	} else {
		return &Modules, nil
	}
}

/* Update Production Status  Data */
func (r *moduleRepo) UpdateModuleProdStatus(moduleId ukama.NodeID, status string, data []byte) error {
	module := Module{
		ModuleID:       moduleId,
		ProdTestStatus: status,
		ProdReport:     data,
	}
	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "module_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"prod_test_status", "prod_report"}),
	}).Create(module)
	return d.Error
}

/* Update module data */
func (r *moduleRepo) UpdateModuleData(moduleId ukama.NodeID, field string, data []byte) error {

	columnName, err := GetModuleDataFieldName(field)
	if err != nil {
		return err
	}

	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "module_id"}},
		DoUpdates: clause.AssignmentColumns([]string{*columnName}),
	}).Create(data)
	return d.Error
}

/* Get Module data */
func (r *moduleRepo) GetModuleData(moduleId ukama.NodeID, field string) ([]byte, error) {
	var data []byte

	columnName, err := GetModuleDataFieldName(field)
	if err != nil {
		return nil, err
	}

	result := r.Db.GetGormDb().Preload(clause.Associations).Select(*columnName).First(data, "module_id = ?", moduleId.StringLowercase())
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
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
func (r *moduleRepo) DeleteBootstrapCert(ModuleId *ukama.NodeID) error {
	result := r.Db.GetGormDb().Where("module_id = ?", ModuleId).UpdateColumn("boostrap_cert", nil)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// /* Update Production user config  Data */
// func (r *moduleRepo) UpdateModuleUserConfig(moduleData *Module) error {
// 	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
// 		Columns:   []clause.Column{{Name: "module_id"}},
// 		DoUpdates: clause.AssignmentColumns([]string{"user_config"}),
// 	}).Create(moduleData)
// 	return d.Error
// }
//
// /* Update Production factory config Data */
// func (r *moduleRepo) UpdateModuleFactoryConfig(moduleData *Module) error {
// 	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
// 		Columns:   []clause.Column{{Name: "module_id"}},
// 		DoUpdates: clause.AssignmentColumns([]string{"factory_config"}),
// 	}).Create(moduleData)
// 	return d.Error
// }

// /* Update Production user calib Data */
// func (r *moduleRepo) UpdateModuleUserCalib(moduleData *Module) error {
// 	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
// 		Columns:   []clause.Column{{Name: "module_id"}},
// 		DoUpdates: clause.AssignmentColumns([]string{"user_calib"}),
// 	}).Create(moduleData)
// 	return d.Error
// }

// /* Update Production factory calib Data */
// func (r *moduleRepo) UpdateModuleFactoryCalib(moduleData *Module) error {
// 	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
// 		Columns:   []clause.Column{{Name: "module_id"}},
// 		DoUpdates: clause.AssignmentColumns([]string{"factory_calib"}),
// 	}).Create(moduleData)
// 	return d.Error
// }
