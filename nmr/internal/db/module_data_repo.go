package db

import (
	"github.com/ukama/openIoR/services/common/ukama"
	"github.com/ukama/ukamaX/common/sql"
	"gorm.io/gorm/clause"
)

type ModuleDataRepo interface {
	AddModuleData(module *ModuleData) error
	GetModuleData(moduleId ukama.NodeID) (*ModuleData, error)
	DeleteModuleData(moduleId ukama.NodeID) error
	UpdateModuleProdStatusData(moduleData *ModuleData) error
	UpdateModuleBootstrapCert(moduleData *ModuleData) error
	UpdateModuleUserConfig(moduleData *ModuleData) error
	UpdateModuleFactoryConfig(moduleData *ModuleData) error
	UpdateModuleUserCalib(moduleData *ModuleData) error
	UpdateModuleFactoryCalib(moduleData *ModuleData) error
}

type moduleDataRepo struct {
	Db sql.Db
}

func NewModuleDataRepo(db sql.Db) *moduleDataRepo {
	return &moduleDataRepo{
		Db: db,
	}
}

func (r *moduleDataRepo) AddModuleData(Module *ModuleData) error {
	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "module_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"prod_test_status", "prod_report", "boostrap_cert", "user_calib", "factory_calib", "user_config", "factory_config", "invenotry_data"}),
	}).Create(Module)
	return d.Error
}

func (r *moduleDataRepo) GetModuleData(ModuleId ukama.NodeID) (*ModuleData, error) {
	var Module ModuleData
	result := r.Db.GetGormDb().Preload(clause.Associations).First(&Module, "module_id = ?", ModuleId.StringLowercase())
	if result.Error != nil {
		return nil, result.Error
	}
	return &Module, nil
}

/* Delete Module  */
func (r *moduleDataRepo) DeleteModuleData(ModuleId ukama.NodeID) error {
	result := r.Db.GetGormDb().Where("module_id = ?", ModuleId).Delete(&Module{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

/* Update Production Status  Data */
func (r *moduleDataRepo) UpdateModuleProdStatusData(moduleData *ModuleData) error {
	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "node_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"prod_test_status", "prod_report"}),
	}).Create(moduleData)
	return d.Error
}

/* Update Production Bootstrap Cert*/
func (r *moduleDataRepo) UpdateModuleBootstrapCert(moduleData *ModuleData) error {
	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "module_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"boostrap_cert"}),
	}).Create(moduleData)
	return d.Error
}

/* Update Production Bootstrap Cert*/
func (r *moduleDataRepo) DeleteBootstrapCert(ModuleId *ukama.NodeID) error {
	result := r.Db.GetGormDb().Where("module_id = ?", ModuleId).UpdateColumn("boostrap_cert", nil)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

/* Update Production user config  Data */
func (r *moduleDataRepo) UpdateModuleUserConfig(moduleData *ModuleData) error {
	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "module_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"user_config"}),
	}).Create(moduleData)
	return d.Error
}

/* Update Production factory config Data */
func (r *moduleDataRepo) UpdateModuleFactoryConfig(moduleData *ModuleData) error {
	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "module_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"factory_config"}),
	}).Create(moduleData)
	return d.Error
}

/* Update Production user calib Data */
func (r *moduleDataRepo) UpdateModuleUserCalib(moduleData *ModuleData) error {
	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "module_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"user_calib"}),
	}).Create(moduleData)
	return d.Error
}

/* Update Production factory calib Data */
func (r *moduleDataRepo) UpdateModuleFactoryCalib(moduleData *ModuleData) error {
	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "module_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"factory_calib"}),
	}).Create(moduleData)
	return d.Error
}
