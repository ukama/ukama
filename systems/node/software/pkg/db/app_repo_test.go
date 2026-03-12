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

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	softwaredb "github.com/ukama/ukama/systems/node/software/pkg/db"
)

// Test fixture constants and helpers to avoid hardcoding in each test.
const (
	defaultAppName       = "test-app"
	defaultAppSpace      = "default-space"
	defaultAppNotes      = "test notes"
	defaultMetricsKeysJSON = `["key1","key2"]`
)

// ukamaDbMock implements sql.Db for tests.
type ukamaDbMock struct {
	gormDb *gorm.DB
}

func (u ukamaDbMock) Init(model ...interface{}) error { return nil }
func (u ukamaDbMock) Connect() error                  { return nil }
func (u ukamaDbMock) GetGormDb() *gorm.DB { return u.gormDb }
func (u ukamaDbMock) ExecuteInTransaction(dbOperation func(tx *gorm.DB) *gorm.DB, nestedFuncs ...func() error) error {
	return nil
}
func (u ukamaDbMock) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB, nestedFuncs ...func(tx *gorm.DB) error) error {
	return nil
}

// testDbResult holds the result of setting up a test DB and repo.
type testDbResult struct {
	Mock sqlmock.Sqlmock
	Repo softwaredb.AppRepo
}

// setupTestDB creates a sqlmock-backed GORM DB and AppRepo. Use this in tests
// to avoid repeating DB creation and mock setup.
func setupTestDB(t *testing.T) *testDbResult {
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

	var _ sql.Db = (*ukamaDbMock)(nil)

	repo := softwaredb.NewAppRepo(ukamaDbMock{gormDb: gdb})
	return &testDbResult{Mock: mock, Repo: repo}
}

// defaultTestApp returns a single shared fixture for an App. Override fields in tests as needed.
func defaultTestApp() softwaredb.App {
	return softwaredb.App{
		Id:          uuid.NewV4(),
		Name:        defaultAppName,
		Space:       defaultAppSpace,
		Notes:       defaultAppNotes,
		MetricsKeys: []string{"key1", "key2"}, // serialized as defaultMetricsKeysJSON in DB
	}
}

// testAppWithName returns an App fixture with the given name; other fields use defaults.
func testAppWithName(name string) softwaredb.App {
	app := defaultTestApp()
	app.Name = name
	return app
}

func TestAppRepoCreate(t *testing.T) {
	app := defaultTestApp()

	t.Run("creates app successfully", func(t *testing.T) {
		res := setupTestDB(t)
		mock, repo := res.Mock, res.Repo

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "apps" ("id","name","space","notes","metrics_keys") VALUES ($1,$2,$3,$4,$5)`)).
			WithArgs(app.Id, app.Name, app.Space, app.Notes, defaultMetricsKeysJSON).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.Create(app)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAppRepoGetAll(t *testing.T) {
	t.Run("returns empty list when no apps", func(t *testing.T) {
		res := setupTestDB(t)
		mock, repo := res.Mock, res.Repo

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "apps"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "space", "notes", "metrics_keys"}))

		apps, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Empty(t, apps)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns all apps", func(t *testing.T) {
		res := setupTestDB(t)
		mock, repo := res.Mock, res.Repo
		app := defaultTestApp()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "apps"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "space", "notes", "metrics_keys"}).
				AddRow(app.Id, app.Name, app.Space, app.Notes, defaultMetricsKeysJSON))

		apps, err := repo.GetAll()
		assert.NoError(t, err)
		require.Len(t, apps, 1)
		assert.Equal(t, app.Name, apps[0].Name)
		assert.Equal(t, app.Space, apps[0].Space)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns ErrRecordNotFound on db error", func(t *testing.T) {
		res := setupTestDB(t)
		mock, repo := res.Mock, res.Repo

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "apps"`)).
			WillReturnError(errors.New("db error"))

		apps, err := repo.GetAll()
		assert.Error(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Nil(t, apps)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAppRepoGet(t *testing.T) {
	appName := defaultAppName
	app := testAppWithName(appName)

	t.Run("returns app when found", func(t *testing.T) {
		res := setupTestDB(t)
		mock, repo := res.Mock, res.Repo

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "apps" WHERE name = $1 ORDER BY "apps"."id" LIMIT $2`)).
			WithArgs(appName, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "space", "notes", "metrics_keys"}).
				AddRow(app.Id, app.Name, app.Space, app.Notes, defaultMetricsKeysJSON))

		got, err := repo.Get(appName)
		assert.NoError(t, err)
		assert.Equal(t, app.Name, got.Name)
		assert.Equal(t, app.Space, got.Space)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns ErrRecordNotFound when app does not exist", func(t *testing.T) {
		res := setupTestDB(t)
		mock, repo := res.Mock, res.Repo

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "apps" WHERE name = $1 ORDER BY "apps"."id" LIMIT $2`)).
			WithArgs("nonexistent", 1).
			WillReturnError(gorm.ErrRecordNotFound)

		got, err := repo.Get("nonexistent")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.Equal(t, softwaredb.App{}, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
