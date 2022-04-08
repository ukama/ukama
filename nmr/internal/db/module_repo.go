package db

import (
	"github.com/ukama/openIoR/services/common/ukama"
	"github.com/ukama/ukamaX/common/sql"
	"gorm.io/gorm/clause"
)

type ModuleRepo interface {
	AddOrUpdateModule(module *Module) error
	GetModule(moduleId ukama.NodeID) (*Module, error)
	DeleteModule(moduleId ukama.NodeID) error
	ListModules() (*[]Module, error)
}

type moduleRepo struct {
	Db sql.Db
}

func NewModuleRepo(db sql.Db) *moduleRepo {
	return &moduleRepo{
		Db: db,
	}
}

func (r *moduleRepo) AddOrUpdateModule(Module *Module) error {
	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "module_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"prod_test_status", "prod_report", "bootstrap_cert", "user_calib", "factory_calib", "user_config", "factory_config", "inventory_data"}),
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
