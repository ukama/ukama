package db

import (
	"errors"
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
	"time"
)

type ComponentRepo interface {
	Get(siteID string) (*SiteComponent, error)
	Upsert(component *SiteComponent) error
}
type componentRepo struct{ db sql.Db }

func NewComponentRepo(db sql.Db) ComponentRepo { return &componentRepo{db: db} }
func (r *componentRepo) Get(siteID string) (*SiteComponent, error) {
	var m SiteComponent
	err := r.db.GetGormDb().First(&m, "site_id = ?", siteID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}
func (r *componentRepo) Upsert(m *SiteComponent) error {
	m.UpdatedAt = time.Now().UTC()
	return r.db.GetGormDb().Save(m).Error
}
