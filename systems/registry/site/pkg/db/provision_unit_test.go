/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 * Copyright (c) 2026-present, Ukama Inc.
 */
package db

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func provisionDatabase(t *testing.T) (*ProvisionRepo, sqlmock.Sqlmock) {
	t.Helper()
	sql, expect, err := sqlmock.New()
	require.NoError(t, err)
	database, err := gorm.Open(postgres.New(postgres.Config{Conn: sql}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, expect.ExpectationsWereMet())
		expect.ExpectClose()
		require.NoError(t, sql.Close())
	})
	return NewProvisionRepo(database), expect
}
func provisionRow(t *testing.T, op *SiteProvision) *sqlmock.Rows {
	t.Helper()
	site, err := json.Marshal(op.Site)
	require.NoError(t, err)
	nodes, err := json.Marshal(op.Nodes)
	require.NoError(t, err)
	return sqlmock.NewRows([]string{"id", "key", "site", "nodes", "phase", "attempt", "revision", "deadline"}).AddRow(op.ID, op.Key, string(site), string(nodes), op.Phase, op.Attempt, op.Revision, op.Deadline)
}
func storedProvision() *SiteProvision {
	site := Site{Id: uuid.NewV4(), NetworkId: uuid.NewV4(), Name: "site"}
	return &SiteProvision{ID: site.Id.String(), Key: site.NetworkId.String() + ":site", Site: site, Nodes: []string{"tower", "amplifier", "controller"}, Phase: "configuring", Attempt: 1, Revision: 1, Deadline: time.Now().UTC().Add(time.Minute)}
}
func expectProvisionLookup(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
	mock.ExpectExec("SELECT pg_advisory_xact_lock").WillReturnResult(sqlmock.NewResult(0, 1))
}

func TestProvisionCreateAtomicReservations(t *testing.T) {
	for _, fail := range []string{"", "lock", "read", "insert", "reservation"} {
		t.Run("failure="+fail, func(t *testing.T) {
			repo, m := provisionDatabase(t)
			op := storedProvision()
			cause := errors.New("database failure")
			m.ExpectBegin()
			lock := m.ExpectExec("SELECT pg_advisory_xact_lock")
			if fail == "lock" {
				lock.WillReturnError(cause)
			} else {
				lock.WillReturnResult(sqlmock.NewResult(0, 1))
				query := m.ExpectQuery(`SELECT .* FROM "site_provisions"`)
				if fail == "read" {
					query.WillReturnError(cause)
				} else {
					query.WillReturnRows(sqlmock.NewRows([]string{"id"}))
					insert := m.ExpectExec(`INSERT INTO "site_provisions"`)
					if fail == "insert" {
						insert.WillReturnError(cause)
					} else {
						insert.WillReturnResult(sqlmock.NewResult(0, 1))
						for i := 0; i < 3; i++ {
							reserve := m.ExpectExec(`INSERT INTO "provision_reservations"`).WithArgs(op.Nodes[i], sqlmock.AnyArg())
							if fail == "reservation" && i == 1 {
								reserve.WillReturnError(cause)
								break
							}
							reserve.WillReturnResult(sqlmock.NewResult(0, 1))
						}
					}
				}
			}
			if fail == "" {
				m.ExpectCommit()
			} else {
				m.ExpectRollback()
			}
			result, err := repo.Create(context.Background(), &op.Site, op.Nodes)
			if fail != "" {
				require.ErrorIs(t, err, cause)
				return
			}
			require.NoError(t, err)
			require.Equal(t, "configuring", result.Phase)
			require.Equal(t, 1, result.Attempt)
			require.Equal(t, op.Nodes, result.Nodes)
			require.WithinDuration(t, time.Now().Add(time.Minute), result.Deadline, time.Second)
		})
	}
}

func TestProvisionCreateDuplicateAndConflict(t *testing.T) {
	for _, mode := range []string{"duplicate", "conflict", "failed", "retire-error"} {
		t.Run(mode, func(t *testing.T) {
			repo, m := provisionDatabase(t)
			op := storedProvision()
			input := op.Site
			if mode == "conflict" {
				input.Location = "different"
			}
			if mode == "failed" || mode == "retire-error" {
				op.Phase = "failed"
			}
			expectProvisionLookup(m)
			m.ExpectQuery(`SELECT .* FROM "site_provisions"`).WillReturnRows(provisionRow(t, op))
			if mode == "failed" || mode == "retire-error" {
				update := m.ExpectExec(`UPDATE "site_provisions"`)
				if mode == "retire-error" {
					update.WillReturnError(errors.New("write failed"))
				} else {
					update.WillReturnResult(sqlmock.NewResult(0, 1))
					m.ExpectExec(`INSERT INTO "site_provisions"`).WillReturnResult(sqlmock.NewResult(0, 1))
					for range op.Nodes {
						m.ExpectExec(`INSERT INTO "provision_reservations"`).WillReturnResult(sqlmock.NewResult(0, 1))
					}
				}
			}
			if mode == "conflict" || mode == "retire-error" {
				m.ExpectRollback()
			} else {
				m.ExpectCommit()
			}
			result, err := repo.Create(context.Background(), &input, op.Nodes)
			if mode == "conflict" || mode == "retire-error" {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if mode == "duplicate" {
				require.Equal(t, op.ID, result.ID)
			} else {
				require.NotEqual(t, op.ID, result.ID)
			}
		})
	}
}

func TestProvisionCreateRequiresThreeDistinctNodes(t *testing.T) {
	repo, _ := provisionDatabase(t)
	for _, nodes := range [][]string{nil, {"a", "b"}, {"a", "a", "c"}, {"a", "b", "a"}, {"a", "b", "b"}} {
		_, err := repo.Create(context.Background(), &Site{}, nodes)
		require.Error(t, err)
	}
}

func TestProvisionSaveFencesWorkersAndReleasesReservations(t *testing.T) {
	for _, mode := range []string{"pending", "active", "failed", "stale", "write-error", "cleanup-error"} {
		t.Run(mode, func(t *testing.T) {
			repo, m := provisionDatabase(t)
			op := storedProvision()
			if mode == "active" || mode == "failed" {
				op.Phase = mode
			}
			if mode == "cleanup-error" {
				op.Phase = "failed"
			}
			m.ExpectBegin()
			update := m.ExpectExec(`UPDATE "site_provisions" .* WHERE id = .* AND revision =`)
			switch mode {
			case "write-error":
				update.WillReturnError(errors.New("write failed"))
			case "stale":
				update.WillReturnResult(sqlmock.NewResult(0, 0))
			default:
				update.WillReturnResult(sqlmock.NewResult(0, 1))
			}
			if op.Phase == "active" || op.Phase == "failed" {
				cleanup := m.ExpectExec(`DELETE FROM "provision_reservations"`).WithArgs(op.ID)
				if mode == "cleanup-error" {
					cleanup.WillReturnError(errors.New("cleanup failed"))
				} else {
					cleanup.WillReturnResult(sqlmock.NewResult(0, 3))
				}
			}
			failed := mode == "stale" || mode == "write-error" || mode == "cleanup-error"
			if failed {
				m.ExpectRollback()
			} else {
				m.ExpectCommit()
			}
			err := repo.Save(context.Background(), op)
			if failed {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, uint64(2), op.Revision)
			}
		})
	}
}

func TestProvisionReadAndPending(t *testing.T) {
	repo, m := provisionDatabase(t)
	op := storedProvision()
	m.ExpectQuery(`SELECT .* FROM "site_provisions"`).WithArgs(op.ID, 1).WillReturnRows(provisionRow(t, op))
	result, err := repo.Get(context.Background(), op.ID)
	require.NoError(t, err)
	require.Equal(t, op.Nodes, result.Nodes)
	require.Equal(t, op.Deadline, result.Deadline)
	m.ExpectQuery(`SELECT .* FROM "site_provisions" WHERE phase NOT IN`).WithArgs("active", "failed").WillReturnRows(provisionRow(t, op))
	pending, err := repo.Pending(context.Background())
	require.NoError(t, err)
	require.Len(t, pending, 1)
}

func TestProvisionLeaseAndRecovery(t *testing.T) {
	for _, mode := range []string{"run", "busy", "claim-error", "read-error", "callback-error"} {
		t.Run(mode, func(t *testing.T) {
			repo, m := provisionDatabase(t)
			op := storedProvision()
			cause := errors.New("worker failure")
			m.ExpectBegin()
			claim := m.ExpectExec(`UPDATE "site_provisions" .*lease_until IS NULL`)
			if mode == "claim-error" {
				claim.WillReturnError(cause)
				m.ExpectRollback()
			} else {
				rows := int64(1)
				if mode == "busy" {
					rows = 0
				}
				claim.WillReturnResult(sqlmock.NewResult(0, rows))
				m.ExpectCommit()
				if mode != "busy" {
					read := m.ExpectQuery(`SELECT .* FROM "site_provisions"`)
					if mode == "read-error" {
						read.WillReturnError(cause)
					} else {
						read.WillReturnRows(provisionRow(t, op))
					}
					m.ExpectBegin()
					m.ExpectExec(`UPDATE "site_provisions" .* WHERE id = .* AND lease_id =`).WillReturnResult(sqlmock.NewResult(0, 1))
					m.ExpectCommit()
				}
			}
			called := false
			err := repo.Run(context.Background(), op.ID, func(ctx context.Context, got *SiteProvision) error {
				called = true
				require.Equal(t, op.ID, got.ID)
				deadline, ok := ctx.Deadline()
				require.True(t, ok)
				require.WithinDuration(t, time.Now().Add(75*time.Second), deadline, time.Second)
				if mode == "callback-error" {
					return cause
				}
				return nil
			})
			require.Equal(t, mode == "run" || mode == "callback-error", called)
			if mode == "claim-error" || mode == "read-error" || mode == "callback-error" {
				require.ErrorIs(t, err, cause)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestProvisionPublicationRequiresCurrentRevision(t *testing.T) {
	for _, mode := range []string{"success", "stale", "error"} {
		t.Run(mode, func(t *testing.T) {
			repo, m := provisionDatabase(t)
			m.ExpectBegin()
			update := m.ExpectExec(`UPDATE "site_provisions" .* WHERE id = .* AND phase = .* AND revision =`)
			if mode == "error" {
				update.WillReturnError(errors.New("write failed"))
				m.ExpectRollback()
			} else {
				rows := int64(1)
				if mode == "stale" {
					rows = 0
				}
				update.WillReturnResult(sqlmock.NewResult(0, rows))
				m.ExpectCommit()
			}
			err := MarkSitePublishing(repo.db, "operation", 2)
			if mode == "success" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
