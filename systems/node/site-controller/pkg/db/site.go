/*
* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at https://mozilla.org/MPL/2.0/.
*
* Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"errors"

	"gorm.io/gorm"
)

// ensureSite ensures a registry row exists for siteID before inserting FK children.
func ensureSite(tx *gorm.DB, siteID string) error {
	if siteID == "" {
		return errors.New("site_id is required")
	}
	s := Site{SiteID: siteID}
	return tx.Where(&Site{SiteID: siteID}).FirstOrCreate(&s).Error
}
