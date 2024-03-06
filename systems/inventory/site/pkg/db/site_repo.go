/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"github.com/ukama/ukama/systems/common/sql"
)

type SiteRepo interface {
	Get() (*Site, error)
}

type siteRepo struct {
	Db sql.Db
}

func NewSiteRepo(db sql.Db) SiteRepo {
	return &siteRepo{
		Db: db,
	}
}

func (s siteRepo) Get() (*Site, error) {
	var site Site

	result := s.Db.GetGormDb().First(&site)
	if result.Error != nil {
		return nil, result.Error
	}

	return &site, nil
}
