/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 * Copyright (c) 2026-present, Ukama Inc.
 */
package db

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestProvisionRecoveryDatabase(t *testing.T) {
	dsn := os.Getenv("PROVISION_TEST_DSN")
	if dsn == "" {
		t.Skip("PROVISION_TEST_DSN is not set")
	}
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	pool, poolErr := database.DB()
	require.NoError(t, poolErr)
	if os.Getenv("PROVISION_TEST_SINGLE_CONNECTION") != "" {
		pool.SetMaxOpenConns(1)
	}
	require.NoError(t, database.AutoMigrate(&Site{}, &SiteProvision{}, &ProvisionReservation{}))
	repo := NewProvisionRepo(database)
	ctx := context.Background()
	site := &Site{NetworkId: uuid.NewV4(), Name: "test-" + uuid.NewV4().String(), AccessId: uuid.NewV4()}
	nodes := []string{uuid.NewV4().String(), uuid.NewV4().String(), uuid.NewV4().String()}
	op, err := repo.Create(ctx, site, nodes)
	require.NoError(t, err)
	duplicate, err := repo.Create(ctx, site, nodes)
	require.NoError(t, err)
	require.Equal(t, op.ID, duplicate.ID)
	op.Phase = "cancelling"
	require.NoError(t, repo.Save(ctx, op))
	restored, err := NewProvisionRepo(database).Get(ctx, op.ID)
	require.NoError(t, err)
	require.Equal(t, "cancelling", restored.Phase)
	require.Equal(t, 1, restored.Attempt)
	require.WithinDuration(t, op.Deadline, restored.Deadline, time.Millisecond)
	duplicate.Phase = "creating"
	require.Error(t, repo.Save(ctx, duplicate), "stale worker must not undo cancellation")
	require.NoError(t, repo.Run(ctx, op.ID, func(work context.Context, current *SiteProvision) error {
		require.Equal(t, "cancelling", current.Phase)
		current.Phase = "failed"
		return repo.Save(work, current)
	}))
	var count int64
	require.NoError(t, database.Model(&ProvisionReservation{}).Where("operation_id = ?", op.ID).Count(&count).Error)
	require.Zero(t, count)
	retry, err := repo.Create(ctx, site, nodes)
	require.NoError(t, err)
	require.NotEqual(t, op.ID, retry.ID)
	retry.LeaseID = "crashed-worker"
	retry.LeaseUntil = time.Now().Add(time.Minute)
	require.NoError(t, repo.Save(ctx, retry))
	called := false
	require.NoError(t, repo.Run(ctx, retry.ID, func(context.Context, *SiteProvision) error { called = true; return nil }))
	require.False(t, called, "an unexpired worker lease must be respected")
	retry.LeaseUntil = time.Now().Add(-time.Second)
	require.NoError(t, repo.Save(ctx, retry))
	require.NoError(t, repo.Run(ctx, retry.ID, func(context.Context, *SiteProvision) error { called = true; return nil }))
	require.True(t, called, "an expired worker lease must be recoverable")
	retry, err = repo.Get(ctx, retry.ID)
	require.NoError(t, err)
	retry.Phase = "creating"
	require.NoError(t, repo.Save(ctx, retry))
	require.NoError(t, database.Transaction(func(tx *gorm.DB) error {
		if err := MarkSitePublishing(tx, retry.ID, retry.Revision); err != nil {
			return err
		}
		return tx.Create(&retry.Site).Error
	}))
	saved, err := repo.Get(ctx, retry.ID)
	require.NoError(t, err)
	require.Equal(t, "publishing", saved.Phase)
	_, err = repo.Create(ctx, &Site{NetworkId: site.NetworkId, Name: "conflict"}, nodes)
	require.Error(t, err)

}
