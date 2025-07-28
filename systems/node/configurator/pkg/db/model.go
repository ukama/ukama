/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"database/sql/driver"
	"strings"

	"gorm.io/gorm"
)

type Commit struct {
	gorm.Model
	Hash string `gorm:"type:string;uniqueIndex:idx_hash_id_case_insensitive,not null"`
}

type Configuration struct {
	gorm.Model
	NodeId          string      `gorm:"type:string;uniqueIndex:idx_node_id_case_insensitive,where:deleted_at is null;size:23;not null"`
	State           CommitState `gorm:"type:uint;not null"`
	Commit          Commit      `gorm:"foreignKey:CommitId"` /* Should be updated by health event after receiving update from node */
	CommitId        int
	LastCommit      Commit `gorm:"foreignKey:LastCommitId"` /* Should be updated by config store after pushing config to msgclient */
	LastCommitId    int
	LastCommitState CommitState `gorm:"type:uint;not null"`
}

type CommitState uint8

const (
	Undefined CommitState = iota
	Default   CommitState = 1 /* First time when node connects */
	Success   CommitState = 2 /* After first successful commit */
	Failed    CommitState = 3 /* After failed  commit */
	Published CommitState = 4 /* After commit is pushed to NodeFeeder but still waiting for confirmation from node */
	Partial   CommitState = 5 /* After partial commits */
	Rollback  CommitState = 6 /* After commit is pushed to msgclient.*/
)

func (e *CommitState) Scan(value interface{}) error {
	*e = CommitState(uint8(value.(int64)))

	return nil
}

func (e CommitState) Value() (driver.Value, error) {
	return int64(e), nil
}

func (e CommitState) String() string {
	ns := map[CommitState]string{
		Undefined: "undefined",
		Default:   "default",
		Success:   "success",
		Failed:    "failed",
		Rollback:  "rollback",
		Partial:   "partial",
		Published: "published",
	}

	return ns[e]
}

func ParseCommitState(s string) CommitState {
	switch strings.ToLower(s) {
	case "partial":
		return Partial
	case "default":
		return Default
	case "success":
		return Success
	case "failed":
		return Failed
	case "committed":
		return Rollback
	case "rollback":
		return Published
	default:
		return Undefined
	}
}
