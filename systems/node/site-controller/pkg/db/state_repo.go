package db

import (
	"errors"
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
	"time"
)

type StateRepo interface {
	Get(siteID string) (*SiteState, error)
	Upsert(state *SiteState) error
}
type stateRepo struct{ db sql.Db }

func NewStateRepo(db sql.Db) StateRepo { return &stateRepo{db: db} }
func (r *stateRepo) Get(siteID string) (*SiteState, error) {
	var m SiteState
	err := r.db.GetGormDb().First(&m, "site_id = ?", siteID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}
func (r *stateRepo) Upsert(m *SiteState) error {
	m.UpdatedAt = time.Now().UTC()
	return r.db.GetGormDb().Save(m).Error
}
