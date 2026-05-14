package db

import (
	"errors"
	"time"

	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

type IntentRepo interface {
	Get(siteID string) (*SiteIntent, error)
	Upsert(intent *SiteIntent) error
}

type intentRepo struct{ db sql.Db }

func NewIntentRepo(db sql.Db) IntentRepo { return &intentRepo{db: db} }

func (r *intentRepo) Get(siteID string) (*SiteIntent, error) {
	var m SiteIntent
	err := r.db.GetGormDb().
		Where("site_id = ?", siteID).
		Order("updated_at DESC, id DESC").
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *intentRepo) Upsert(m *SiteIntent) error {
	gdb := r.db.GetGormDb()
	if err := ensureSite(gdb, m.SiteID); err != nil {
		return err
	}
	now := time.Now().UTC()
	row := SiteIntent{
		SiteID:         m.SiteID,
		DesiredSite:    m.DesiredSite,
		DesiredService: m.DesiredService,
		DesiredRadio:   m.DesiredRadio,
		Reason:         m.Reason,
		RequestedBy:    m.RequestedBy,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	return gdb.Create(&row).Error
}
