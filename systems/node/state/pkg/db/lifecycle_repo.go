/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"context"
	"errors"
	"time"

	"github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/state/pkg/lifecycle"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// One assignment decision per node. Completion and the observed event cursor
// are committed with node state, so a crash cannot acknowledge half a change.
type LifecycleRecord struct {
	NodeID              string `gorm:"primaryKey"`
	Network             string
	Site                string
	RequestID           string
	Completed           bool
	AwaitingOperational bool             `gorm:"index:idx_lifecycle_retry,priority:1"`
	RetryAt             time.Time        `gorm:"index:idx_lifecycle_retry,priority:2"`
	Cursor              lifecycle.Cursor `gorm:"serializer:json;type:jsonb"`
}

type LifecyclePublication struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement"`
	NodeID   string
	State    string
	Substate string
	Event    string
}

type LifecycleRepo struct {
	db *gorm.DB
}

func NewLifecycleRepo(db *gorm.DB) *LifecycleRepo {
	return &LifecycleRepo{db: db}
}

// Update serializes every event for a node, including online/offline events,
// across backend replicas. The callback changes a detached state snapshot.
func (r *LifecycleRepo) Update(ctx context.Context, nodeID, event string,
	change func(*LifecycleRecord, *State) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		record := LifecycleRecord{NodeID: nodeID}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&record).Error; err != nil {
			return err
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&record, "node_id = ?", nodeID).Error; err != nil {
			return err
		}
		previous := State{}
		err := tx.Preload("Config").Where("node_id = ?", nodeID).Order("created_at DESC").First(&previous).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		current := previous
		current.NodeId = nodeID
		current.SubState = append(StringArray{}, previous.SubState...)
		current.Events = append(StringArray{}, previous.Events...)
		if err := change(&record, &current); err != nil {
			return err
		}
		if err := tx.Save(&record).Error; err != nil {
			return err
		}
		return saveLifecycleState(tx, &previous, &current, event)
	})
}

func saveLifecycleState(tx *gorm.DB, previous, current *State, event string) error {
	if previous.Id == uuid.Nil && current.Config == nil {
		current.Config = &NodeConfig{Id: uuid.NewV4(), NodeId: current.NodeId}
	}
	changed := previous.Id == uuid.Nil || previous.CurrentState != current.CurrentState ||
		lastSubstate(previous) != lastSubstate(current)
	if current.Config != nil && current.Config.Id != previous.ConfigId {
		if err := tx.Create(current.Config).Error; err != nil {
			return err
		}
		current.ConfigId = current.Config.Id
	}
	if !changed {
		if current.ConfigId != previous.ConfigId {
			return tx.Model(previous).Update("config_id", current.ConfigId).Error
		}
		return nil
	}
	current.Id = uuid.NewV4()
	if previous.Id != uuid.Nil {
		current.PreviousStateId = &previous.Id
	}
	current.CreatedAt = time.Time{}
	current.UpdatedAt = time.Time{}
	current.Events = StringArray{event}
	current.SubState = StringArray{lastSubstate(current)}
	// Reuse connection metadata; inserting another state is not a reconnect.
	if err := tx.Omit("Config").Create(current).Error; err != nil {
		return err
	}
	publication := LifecyclePublication{
		NodeID: current.NodeId, State: current.CurrentState.String(),
		Substate: lastSubstate(current), Event: event,
	}
	return tx.Create(&publication).Error
}

func lastSubstate(state *State) string {
	if len(state.SubState) == 0 {
		return "off"
	}
	return state.SubState[len(state.SubState)-1]
}

func (r *LifecycleRepo) RetryAssignments(ctx context.Context, now time.Time,
	send func(*LifecycleRecord)) error {
	var nodes []string
	if err := r.db.WithContext(ctx).Model(&LifecycleRecord{}).
		Where("request_id <> '' AND awaiting_operational = ? AND retry_at <= ?", true, now).
		Order("retry_at").Limit(100).Pluck("node_id", &nodes).Error; err != nil {
		return err
	}
	for _, nodeID := range nodes {
		err := r.Update(ctx, nodeID, "assignment-retry", func(record *LifecycleRecord, state *State) error {
			if state.CurrentState == ukama.NodeState_Offboarded {
				record.AwaitingOperational = false
				return nil
			}
			if !record.AwaitingOperational || record.RetryAt.After(now) {
				return nil
			}
			// Keep the same request ID even when publishing fails. A crash
			// after send but before commit safely sends the decision again.
			send(record)
			record.RetryAt = now.Add(60 * time.Second)
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Publications are sent in commit order. A crash before delete may replay an
// event, but cannot lose a committed transition or publish before persistence.
func (r *LifecycleRepo) PublishPending(ctx context.Context,
	publish func(*LifecyclePublication) error) error {
	for count := 0; count < 100; count++ {
		empty := false
		err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			var item LifecyclePublication
			err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Order("id").Limit(1).Find(&item).Error
			if err == nil && item.ID == 0 {
				empty = true
				return nil
			}
			if err != nil {
				return err
			}
			if err := publish(&item); err != nil {
				return err
			}
			return tx.Delete(&item).Error
		})
		if err != nil || empty {
			return err
		}
	}
	return nil
}
