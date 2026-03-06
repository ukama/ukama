/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ukama/ukama/systems/common/ukama"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	softwaredb "github.com/ukama/ukama/systems/node/software/pkg/db"
)

// Software test fixture constants (shared to avoid hardcoding in each test).
// defaultAppName is defined in app_repo_test.go and reused here.
const (
	defaultNodeId         = "node-1"
	defaultCurrentVersion = "1.0.0"
	defaultDesiredVersion = "1.0.1"
	defaultChangeLogsJSON = `["fix bug","new feature"]`
)

// softwareTestDbResult holds mock and SoftwareRepo for tests.
type softwareTestDbResult struct {
	Mock sqlmock.Sqlmock
	Repo softwaredb.SoftwareRepo
}

// setupSoftwareTestDB creates a sqlmock-backed GORM DB and SoftwareRepo.
func setupSoftwareTestDB(t *testing.T) *softwareTestDbResult {
	t.Helper()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gdb, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	repo := softwaredb.NewSoftwareRepo(ukamaDbMock{gormDb: gdb})
	return &softwareTestDbResult{Mock: mock, Repo: repo}
}

// defaultTestSoftware returns a shared Software fixture. Override fields in tests as needed.
func defaultTestSoftware() *softwaredb.Software {
	releaseDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	return &softwaredb.Software{
		Id:             uuid.NewV4(),
		NodeId:         defaultNodeId,
		AppName:        defaultAppName, // from app_repo_test.go
		ChangeLogs:     []string{"fix bug", "new feature"},
		CurrentVersion: defaultCurrentVersion,
		DesiredVersion: defaultDesiredVersion,
		ReleaseDate:    releaseDate,
		Status:         ukama.UpToDate,
	}
}

// softwareWithStatus returns a copy of the default fixture with the given status.
func softwareWithStatus(status ukama.SoftwareStatusType) *softwaredb.Software {
	s := defaultTestSoftware()
	s.Status = status
	return s
}

func TestSoftwareRepoCreate(t *testing.T) {
	sw := defaultTestSoftware()

	t.Run("creates software successfully", func(t *testing.T) {
		res := setupSoftwareTestDB(t)
		mock, repo := res.Mock, res.Repo

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "softwares" ("id","node_id","app_name","change_logs","current_version","desired_version","release_date","status") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "created_at","updated_at","deleted_at"`)).
			WithArgs(sw.Id, sw.NodeId, sw.AppName, defaultChangeLogsJSON, sw.CurrentVersion, sw.DesiredVersion, sw.ReleaseDate, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at"}).AddRow(time.Now(), time.Now(), nil))
		mock.ExpectCommit()

		err := repo.Create(sw)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSoftwareRepoGet(t *testing.T) {
	sw := defaultTestSoftware()

	t.Run("returns software when found", func(t *testing.T) {
		res := setupSoftwareTestDB(t)
		mock, repo := res.Mock, res.Repo

		// Main query: Get by id
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "softwares" WHERE id = $1 ORDER BY "softwares"."id" LIMIT $2`)).
			WithArgs(sw.Id, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "node_id", "app_name", "change_logs", "current_version", "desired_version", "release_date", "created_at", "updated_at", "deleted_at", "status"}).
				AddRow(sw.Id, sw.NodeId, sw.AppName, defaultChangeLogsJSON, sw.CurrentVersion, sw.DesiredVersion, sw.ReleaseDate, time.Now(), time.Now(), nil, sw.Status))
		// Preload App
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "apps" WHERE "apps"."name" = $1`)).
			WithArgs(sw.AppName).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "space", "notes", "metrics_keys"}).
				AddRow(uuid.NewV4(), sw.AppName, "default-space", "notes", defaultMetricsKeysJSON))

		got, err := repo.Get(sw.Id)
		assert.NoError(t, err)
		assert.Equal(t, sw.Id, got.Id)
		assert.Equal(t, sw.NodeId, got.NodeId)
		assert.Equal(t, sw.AppName, got.AppName)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns ErrRecordNotFound when software does not exist", func(t *testing.T) {
		res := setupSoftwareTestDB(t)
		mock, repo := res.Mock, res.Repo
		missingID := uuid.NewV4()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "softwares" WHERE id = $1 ORDER BY "softwares"."id" LIMIT $2`)).
			WithArgs(missingID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		got, err := repo.Get(missingID)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Equal(t, softwaredb.Software{}, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSoftwareRepoList(t *testing.T) {
	sw := defaultTestSoftware()

	t.Run("returns empty list when no software", func(t *testing.T) {
		res := setupSoftwareTestDB(t)
		mock, repo := res.Mock, res.Repo

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "softwares"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "node_id", "app_name", "change_logs", "current_version", "desired_version", "release_date", "created_at", "updated_at", "deleted_at", "status"}))

		list, err := repo.List("", ukama.SoftwareStatusType(0), "")
		assert.NoError(t, err)
		assert.Empty(t, list)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns list with filters applied", func(t *testing.T) {
		res := setupSoftwareTestDB(t)
		mock, repo := res.Mock, res.Repo

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "softwares" WHERE app_name = $1 AND node_id = $2 AND status = $3`)).
			WithArgs(defaultAppName, defaultNodeId, ukama.UpToDate).
			WillReturnRows(sqlmock.NewRows([]string{"id", "node_id", "app_name", "change_logs", "current_version", "desired_version", "release_date", "created_at", "updated_at", "deleted_at", "status"}).
				AddRow(sw.Id, sw.NodeId, sw.AppName, defaultChangeLogsJSON, sw.CurrentVersion, sw.DesiredVersion, sw.ReleaseDate, time.Now(), time.Now(), nil, sw.Status))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "apps" WHERE "apps"."name" = $1`)).
			WithArgs(defaultAppName).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "space", "notes", "metrics_keys"}).
				AddRow(uuid.NewV4(), sw.AppName, "default-space", "notes", defaultMetricsKeysJSON))

		list, err := repo.List(defaultNodeId, ukama.UpToDate, defaultAppName)
		assert.NoError(t, err)
		require.Len(t, list, 1)
		assert.Equal(t, sw.Id, list[0].Id)
		assert.Equal(t, defaultAppName, list[0].AppName)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns error when query fails", func(t *testing.T) {
		res := setupSoftwareTestDB(t)
		mock, repo := res.Mock, res.Repo

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "softwares"`)).
			WillReturnError(errors.New("db error"))

		list, err := repo.List("", ukama.SoftwareStatusType(0), "")
		assert.Error(t, err)
		assert.Nil(t, list)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSoftwareRepoUpdate(t *testing.T) {
	sw := softwareWithStatus(ukama.UpdateInProgress)
	sw.DesiredVersion = "1.0.2"

	t.Run("updates software successfully", func(t *testing.T) {
		res := setupSoftwareTestDB(t)
		mock, repo := res.Mock, res.Repo

		// GORM Save() sends SET col1=$1, col2=$2, ... WHERE id=$last
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "softwares" SET "node_id"=$1,"app_name"=$2,"change_logs"=$3,"current_version"=$4,"desired_version"=$5,"release_date"=$6,"created_at"=$7,"updated_at"=$8,"deleted_at"=$9,"status"=$10 WHERE "id" = $11`)).
			WithArgs(sw.NodeId, sw.AppName, defaultChangeLogsJSON, sw.CurrentVersion, sw.DesiredVersion,
				sw.ReleaseDate, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, sw.Status, sw.Id).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.Update(sw)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
