/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package storage

import (
	"errors"

	"github.com/ukama/ukama/systems/common/ukama"
)

var ErrNotFound = errors.New("sim not found")
var ErrInternal = errors.New("an unexpected error has occurred")

type Storage interface {
	Get(key string) (*SimInfo, error)
	Put(string, *SimInfo) error
	Delete(key string) error
}

type SimInfo struct {
	Iccid  string          `json:"iccid"`
	Imsi   string          `json:"imsi"`
	Status ukama.SimStatus `json:"status"`
}
