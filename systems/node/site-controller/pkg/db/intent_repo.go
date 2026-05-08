package db

import (
	"errors"
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
	"time"
)

type IntentRepo interface {
	Get(siteID string) (*SiteIntent, error)
	Upsert(intent *SiteIntent) error
}
type intentRepo struct{ db sql.Db }

func NewIntentRepo(db sql.Db) IntentRepo { return &intentRepo{db: db} }
func (r *intentRepo) Get(siteID string) (*SiteIntent, error) {
	var m SiteIntent
	err := r.db.GetGormDb().First(&m, "site_id = ?", siteID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}
func (r *intentRepo) Upsert(m *SiteIntent) error {
	now := time.Now().UTC()
	m.UpdatedAt = now
	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
	return r.db.GetGormDb().Save(m).Error
}
