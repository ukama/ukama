/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
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
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	softwaredb "github.com/ukama/ukama/systems/node/software/pkg/db"
)

const (
	releaseName    = "ukamaos"
	releaseVersion = "1.2.0"
)

func setupReleaseTestDB(t *testing.T) (sqlmock.Sqlmock, softwaredb.ReleaseRepo) {
	t.Helper()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	require.NoError(t, err)
	return mock, softwaredb.NewReleaseRepo(ukamaDbMock{gormDb: gdb})
}

func catalogRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "name", "type", "version", "chunked"})
}

func desiredRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "name", "type", "desired_version"})
}

func returning() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"created_at", "updated_at"}).
		AddRow(time.Now(), time.Now())
}

func TestReleaseRepoUpsert(t *testing.T) {
	upsert := `INSERT INTO "release_catalogs" .* ON CONFLICT \("name","type","version"\) ` +
		`DO UPDATE SET "digest"="excluded"."digest","size_bytes"="excluded"."size_bytes",` +
		`"available"="excluded"."available","uploaded_at"="excluded"."uploaded_at",` +
		`"updated_at"="excluded"."updated_at"`

	t.Run("AssignsIdAndDefaultTypeAndLeavesChunkedAlone", func(t *testing.T) {
		mock, repo := setupReleaseTestDB(t)
		rel := &softwaredb.ReleaseCatalog{Name: releaseName, Version: releaseVersion}

		mock.ExpectBegin()
		mock.ExpectQuery(upsert).WillReturnRows(returning())
		mock.ExpectCommit()

		require.NoError(t, repo.Upsert(rel))
		assert.NotEqual(t, uuid.Nil, rel.Id)
		assert.Equal(t, "app", rel.Type)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("KeepsGivenIdAndType", func(t *testing.T) {
		mock, repo := setupReleaseTestDB(t)
		id := uuid.NewV4()
		rel := &softwaredb.ReleaseCatalog{Id: id, Name: releaseName, Type: "chunk",
			Version: releaseVersion}

		mock.ExpectBegin()
		mock.ExpectQuery(upsert).WillReturnRows(returning())
		mock.ExpectCommit()

		require.NoError(t, repo.Upsert(rel))
		assert.Equal(t, id, rel.Id)
		assert.Equal(t, "chunk", rel.Type)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("PassesThroughDbError", func(t *testing.T) {
		mock, repo := setupReleaseTestDB(t)

		mock.ExpectBegin()
		mock.ExpectQuery(upsert).WillReturnError(errors.New("insert failed"))
		mock.ExpectRollback()

		assert.Error(t, repo.Upsert(&softwaredb.ReleaseCatalog{Name: releaseName,
			Version: releaseVersion}))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestReleaseRepoSetChunked(t *testing.T) {
	update := regexp.QuoteMeta(`UPDATE "release_catalogs" SET "chunked"=$1,` +
		`"updated_at"=$2 WHERE name = $3 AND type = $4 AND version = $5`)

	t.Run("DefaultsTypeToApp", func(t *testing.T) {
		mock, repo := setupReleaseTestDB(t)

		mock.ExpectBegin()
		mock.ExpectExec(update).
			WithArgs(true, sqlmock.AnyArg(), releaseName, "app", releaseVersion).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		assert.NoError(t, repo.SetChunked(releaseName, "", releaseVersion))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("PassesThroughDbError", func(t *testing.T) {
		mock, repo := setupReleaseTestDB(t)

		mock.ExpectBegin()
		mock.ExpectExec(update).WillReturnError(errors.New("update failed"))
		mock.ExpectRollback()

		assert.Error(t, repo.SetChunked(releaseName, "chunk", releaseVersion))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestReleaseRepoGet(t *testing.T) {
	get := regexp.QuoteMeta(`SELECT * FROM "release_catalogs" WHERE name = $1 AND ` +
		`type = $2 AND version = $3 ORDER BY "release_catalogs"."id" LIMIT $4`)

	t.Run("Found", func(t *testing.T) {
		mock, repo := setupReleaseTestDB(t)
		id := uuid.NewV4()

		mock.ExpectQuery(get).WithArgs(releaseName, "app", releaseVersion, 1).
			WillReturnRows(catalogRows().AddRow(id, releaseName, "app", releaseVersion,
				true))

		rel, err := repo.Get(releaseName, "", releaseVersion)

		require.NoError(t, err)
		assert.Equal(t, id, rel.Id)
		assert.True(t, rel.Chunked)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		mock, repo := setupReleaseTestDB(t)

		mock.ExpectQuery(get).WithArgs(releaseName, "app", releaseVersion, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		rel, err := repo.Get(releaseName, "app", releaseVersion)

		assert.Nil(t, rel)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestReleaseRepoExists(t *testing.T) {
	count := regexp.QuoteMeta(`SELECT count(*) FROM "release_catalogs" WHERE name = $1 ` +
		`AND type = $2 AND version = $3`)

	for name, tc := range map[string]struct {
		count  int64
		exists bool
	}{
		"Present": {1, true},
		"Absent":  {0, false},
	} {
		t.Run(name, func(t *testing.T) {
			mock, repo := setupReleaseTestDB(t)

			mock.ExpectQuery(count).WithArgs(releaseName, "app", releaseVersion).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(tc.count))

			exists, err := repo.Exists(releaseName, "", releaseVersion)

			require.NoError(t, err)
			assert.Equal(t, tc.exists, exists)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

	t.Run("PassesThroughDbError", func(t *testing.T) {
		mock, repo := setupReleaseTestDB(t)

		mock.ExpectQuery(count).WillReturnError(errors.New("count failed"))

		exists, err := repo.Exists(releaseName, "app", releaseVersion)

		assert.Error(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestReleaseRepoList(t *testing.T) {
	t.Run("NoFilters", func(t *testing.T) {
		mock, repo := setupReleaseTestDB(t)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "release_catalogs" ORDER BY name, version`)).
			WillReturnRows(catalogRows().
				AddRow(uuid.NewV4(), releaseName, "app", "1.0.0", false).
				AddRow(uuid.NewV4(), releaseName, "app", releaseVersion, true))

		out, err := repo.List("", "")

		require.NoError(t, err)
		assert.Len(t, out, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NameAndTypeFilters", func(t *testing.T) {
		mock, repo := setupReleaseTestDB(t)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "release_catalogs" WHERE name = $1 `+
			`AND type = $2 ORDER BY name, version`)).
			WithArgs(releaseName, "chunk").
			WillReturnRows(catalogRows())

		out, err := repo.List(releaseName, "chunk")

		require.NoError(t, err)
		assert.Empty(t, out)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestReleaseRepoGetDesired(t *testing.T) {
	get := regexp.QuoteMeta(`SELECT * FROM "app_desired_releases" WHERE name = $1 AND ` +
		`type = $2 ORDER BY "app_desired_releases"."id" LIMIT $3`)

	t.Run("Found", func(t *testing.T) {
		mock, repo := setupReleaseTestDB(t)

		mock.ExpectQuery(get).WithArgs(releaseName, "app", 1).
			WillReturnRows(desiredRows().AddRow(uuid.NewV4(), releaseName, "app",
				releaseVersion))

		d, err := repo.GetDesired(releaseName, "")

		require.NoError(t, err)
		assert.Equal(t, releaseVersion, d.DesiredVersion)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotSetIsNotAnError", func(t *testing.T) {
		mock, repo := setupReleaseTestDB(t)

		mock.ExpectQuery(get).WithArgs(releaseName, "app", 1).WillReturnRows(desiredRows())

		d, err := repo.GetDesired(releaseName, "app")

		assert.NoError(t, err)
		assert.Nil(t, d)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("PassesThroughDbError", func(t *testing.T) {
		mock, repo := setupReleaseTestDB(t)

		mock.ExpectQuery(get).WillReturnError(errors.New("select failed"))

		d, err := repo.GetDesired(releaseName, "app")

		assert.Error(t, err)
		assert.Nil(t, d)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestReleaseRepoSetDesired(t *testing.T) {
	upsert := `INSERT INTO "app_desired_releases" .* ON CONFLICT \("name","type"\) ` +
		`DO UPDATE SET "desired_version"="excluded"."desired_version",` +
		`"promoted_at"="excluded"."promoted_at","promoted_by"="excluded"."promoted_by",` +
		`"updated_at"="excluded"."updated_at"`

	t.Run("AssignsIdAndDefaultType", func(t *testing.T) {
		mock, repo := setupReleaseTestDB(t)
		d := &softwaredb.AppDesiredRelease{Name: releaseName,
			DesiredVersion: releaseVersion}

		mock.ExpectBegin()
		mock.ExpectQuery(upsert).WillReturnRows(returning())
		mock.ExpectCommit()

		require.NoError(t, repo.SetDesired(d))
		assert.NotEqual(t, uuid.Nil, d.Id)
		assert.Equal(t, "app", d.Type)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("PassesThroughDbError", func(t *testing.T) {
		mock, repo := setupReleaseTestDB(t)

		mock.ExpectBegin()
		mock.ExpectQuery(upsert).WillReturnError(errors.New("insert failed"))
		mock.ExpectRollback()

		assert.Error(t, repo.SetDesired(&softwaredb.AppDesiredRelease{Name: releaseName,
			DesiredVersion: releaseVersion}))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestReleaseRepoListDesired(t *testing.T) {
	mock, repo := setupReleaseTestDB(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "app_desired_releases"`)).
		WillReturnRows(desiredRows().AddRow(uuid.NewV4(), releaseName, "app",
			releaseVersion))

	out, err := repo.ListDesired()

	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, releaseVersion, out[0].DesiredVersion)
	assert.NoError(t, mock.ExpectationsWereMet())
}
