/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import "gorm.io/gorm"

type gormHandle interface {
	GetGormDb() *gorm.DB
}

type gormOnly struct {
	db *gorm.DB
}

func (g gormOnly) GetGormDb() *gorm.DB {
	return g.db
}
