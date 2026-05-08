package db

import (
	"time"

	"github.com/ukama/ukama/systems/common/sql"
)

type SwitchPolicyRepo interface {
	Get(siteID string) (*SiteSwitchPolicy, error)
	Upsert(policy *SiteSwitchPolicy) error
}

type switchPolicyRepo struct{ db sql.Db }

func NewSwitchPolicyRepo(db sql.Db) SwitchPolicyRepo { return &switchPolicyRepo{db: db} }

func (r *switchPolicyRepo) Get(siteID string) (*SiteSwitchPolicy, error) {
	var policy SiteSwitchPolicy
	err := r.db.GetGormDb().Where("site_id = ?", siteID).First(&policy).Error
	if err != nil {
		return nil, nil
	}
	return &policy, nil
}

func (r *switchPolicyRepo) Upsert(policy *SiteSwitchPolicy) error {
	if policy == nil {
		return nil
	}
	now := time.Now().UTC()
	policy.ObservedAt = now
	policy.UpdatedAt = now
	if policy.CreatedAt.IsZero() {
		policy.CreatedAt = now
	}
	return r.db.GetGormDb().Save(policy).Error
}
