/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 * Copyright (c) 2026-present, Ukama Inc.
 */
package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type SiteProvision struct {
	Revision   uint64
	LeaseID    string
	LeaseUntil time.Time
	ID         string   `gorm:"primaryKey"`
	Key        string   `gorm:"uniqueIndex"`
	Site       Site     `gorm:"serializer:json;type:jsonb"`
	Nodes      []string `gorm:"serializer:json;type:jsonb"`
	Phase      string   `gorm:"index"`
	Attempt    int
	Deadline   time.Time
	Failure    string
	Stop       bool
}

type ProvisionReservation struct {
	NodeID      string `gorm:"primaryKey"`
	OperationID string `gorm:"index"`
}

type ProvisionRepo struct{ db *gorm.DB }

func NewProvisionRepo(db *gorm.DB) *ProvisionRepo { return &ProvisionRepo{db: db} }

func (r *ProvisionRepo) Create(ctx context.Context, site *Site, nodes []string) (*SiteProvision, error) {
	op := &SiteProvision{}
	if len(nodes) != 3 || nodes[0] == nodes[1] || nodes[0] == nodes[2] || nodes[1] == nodes[2] {
		return nil, fmt.Errorf("site requires exactly three distinct nodes")
	}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		key := site.NetworkId.String() + ":" + site.Name
		// Serialize duplicate Add calls before reserving nodes.
		if err := tx.Exec("SELECT pg_advisory_xact_lock(hashtextextended(?, 0))", key).Error; err != nil {
			return err
		}
		err := tx.Where("key = ?", key).First(op).Error
		if err == nil {
			if !sameSiteRequest(op.Site, *site) {
				return fmt.Errorf("site name already reserved")
			}
			if op.Phase != "failed" {
				return nil
			}
			if err := tx.Model(op).Update("key", op.ID).Error; err != nil {
				return err
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		site.Id = uuid.NewV4()
		*op = SiteProvision{Revision: 1, ID: site.Id.String(), Key: key, Site: *site, Nodes: nodes, Phase: "configuring", Attempt: 1, Deadline: time.Now().UTC().Add(60 * time.Second)}
		if err := tx.Create(op).Error; err != nil {
			return err
		}
		for _, nodeID := range nodes {
			if err := tx.Create(&ProvisionReservation{NodeID: nodeID, OperationID: op.ID}).Error; err != nil {
				return fmt.Errorf("node %s is reserved: %w", nodeID, err)
			}
		}
		return nil
	})
	return op, err
}

func (r *ProvisionRepo) Get(ctx context.Context, id string) (*SiteProvision, error) {
	op := &SiteProvision{}
	err := r.db.WithContext(ctx).First(op, "id = ?", id).Error
	return op, err
}

func (r *ProvisionRepo) Pending(ctx context.Context) ([]SiteProvision, error) {
	var ops []SiteProvision
	err := r.db.WithContext(ctx).Where("phase NOT IN ?", []string{"active", "failed"}).Find(&ops).Error
	return ops, err
}

func (r *ProvisionRepo) Save(ctx context.Context, op *SiteProvision) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		version := op.Revision
		op.Revision++
		result := tx.Model(&SiteProvision{}).Where("id = ? AND revision = ?", op.ID, version).Select("*").Updates(op)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return fmt.Errorf("site operation changed; reload before proceeding")
		}
		if op.Phase == "failed" || op.Phase == "active" {
			return tx.Where("operation_id = ?", op.ID).Delete(&ProvisionReservation{}).Error
		}
		return nil
	})
}

// A bounded lease allows recovery without holding a database connection
// while waiting for nodes. Revision checks fence a worker whose lease expired.
func (r *ProvisionRepo) Run(ctx context.Context, id string, run func(context.Context, *SiteProvision) error) error {
	owner := uuid.NewV4().String()
	claim := r.db.WithContext(ctx).Model(&SiteProvision{}).
		Where("id = ? AND phase NOT IN ? AND (lease_until IS NULL OR lease_until < CURRENT_TIMESTAMP)",
			id, []string{"active", "failed"}).
		Updates(map[string]interface{}{
			"lease_id":    owner,
			"lease_until": gorm.Expr("CURRENT_TIMESTAMP + INTERVAL '90 seconds'"),
			"revision":    gorm.Expr("revision + 1"),
		})
	if claim.Error != nil {
		return claim.Error
	}
	if claim.RowsAffected == 0 {
		return nil
	}
	defer func() {
		release, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		r.db.WithContext(release).Model(&SiteProvision{}).Where("id = ? AND lease_id = ?", id, owner).
			Updates(map[string]interface{}{"lease_id": "", "lease_until": nil})
	}()
	op, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	work, cancel := context.WithTimeout(ctx, 75*time.Second)
	defer cancel()
	return run(work, op)
}

// Site creation and the move to publication commit together.
func MarkSitePublishing(tx *gorm.DB, id string, revision uint64) error {
	result := tx.Model(&SiteProvision{}).Where("id = ? AND phase = ? AND revision = ?", id, "creating", revision).
		Updates(map[string]interface{}{"phase": "publishing", "revision": gorm.Expr("revision + 1")})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return fmt.Errorf("site operation changed before creation")
	}
	return nil
}

func sameSiteRequest(left, right Site) bool {
	left.Id, right.Id = uuid.Nil, uuid.Nil
	left.CreatedAt, right.CreatedAt = time.Time{}, time.Time{}
	left.UpdatedAt, right.UpdatedAt = time.Time{}, time.Time{}
	left.DeletedAt, right.DeletedAt = gorm.DeletedAt{}, gorm.DeletedAt{}
	return left == right
}
