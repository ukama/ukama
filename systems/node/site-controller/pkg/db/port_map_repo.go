/*
* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at https://mozilla.org/MPL/2.0/.
*
* Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"time"

	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

type PortMapRepo interface {
	GetBySite(siteID string) ([]SitePortMap, error)
	Upsert(siteID string, cnodeID string, ports []SitePortMap) error
}
type portMapRepo struct{ db sql.Db }

func NewPortMapRepo(db sql.Db) PortMapRepo { return &portMapRepo{db: db} }
func (r *portMapRepo) GetBySite(siteID string) ([]SitePortMap, error) {
	var ports []SitePortMap
	err := r.db.GetGormDb().Where("site_id = ?", siteID).Order("port asc").Find(&ports).Error
	return ports, err
}
func (r *portMapRepo) Upsert(siteID string, cnodeID string, ports []SitePortMap) error {
	return r.db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if err := ensureSite(tx, siteID); err != nil {
			return err
		}
		if err := tx.Where("site_id = ?", siteID).Delete(&SitePortMap{}).Error; err != nil {
			return err
		}
		now := time.Now().UTC()
		for i := range ports {
			ports[i].SiteID = siteID
			// if ports[i].CNodeID == "" {
			// 	ports[i].CNodeID = cnodeID
			// }
			ports[i].CreatedAt = now
			ports[i].UpdatedAt = now
			if err := tx.Create(&ports[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
