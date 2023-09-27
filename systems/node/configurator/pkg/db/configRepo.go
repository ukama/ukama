package db

import (
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/sql"

	"gorm.io/gorm/clause"
)

type ConfigRepo interface {
	Add(id string) error
	Get(id string) (*Configuration, error)
	GetAll() ([]Configuration, error)
	Delete(id string) error
	Update(c Configuration) error
	UpdateCurrentCommit(c Configuration, hash string) error
	UpdateLastCommit(c Configuration, hash string) error
}

type configRepo struct {
	Db sql.Db
}

func NewConfigRepo(db sql.Db) ConfigRepo {
	return &configRepo{
		Db: db,
	}
}

func (n *configRepo) Add(node string) error {
	config := Configuration{
		NodeId:     node,
		Status:     Default,
		LastStatus: Undefined,
	}

	r := n.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "node_id"}},
		DoNothing: true,
	}).Create(&config)

	return r.Error
}

func (n *configRepo) Get(id string) (*Configuration, error) {
	var config Configuration

	result := n.Db.GetGormDb().Preload("Commit").First(&config, "node_id=?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &config, nil
}

func (n *configRepo) GetAll() ([]Configuration, error) {
	var configs []Configuration

	result := n.Db.GetGormDb().Preload("Commit").Find(&configs)

	if result.Error != nil {
		return nil, result.Error
	}

	return configs, nil
}

func (n *configRepo) Delete(id string) error {
	var configs Configuration
	result := n.Db.GetGormDb().Where("node_id=?", id).Delete(&configs)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Update updated node with `id`. Only fields that are not nil are updated, eg name and state.
func (n *configRepo) Update(c Configuration) error {

	result := n.Db.GetGormDb().Where("node_id=?", c.NodeId).Updates(&c)
	if result.Error != nil {
		return result.Error
	}

	return result.Error
}

func (n *configRepo) UpdateLastCommit(c Configuration, hash string) error {
	err := n.Db.GetGormDb().Model(&c).Association("LastCommit").Replace(&Commit{Hash: hash})
	if err != nil {
		log.Errorf("Failed to cuurent commit: %v", err)
	}

	return err
}

func (n *configRepo) UpdateCurrentCommit(c Configuration, hash string) error {
	err := n.Db.GetGormDb().Model(&c).Association("Commit").Replace(&Commit{Hash: hash})
	if err != nil {
		log.Errorf("Failed to cuurent commit: %v", err)
	}

	return err
}
