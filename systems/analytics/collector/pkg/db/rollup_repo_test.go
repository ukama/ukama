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
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/uuid"

	col_db "github.com/ukama/ukama/systems/analytics/collector/pkg/db"
)

func Test_RollupRepo_UpsertBusinessSalesDaily(t *testing.T) {
	t.Run("RowIsUpserted", func(t *testing.T) {
		// Arrange
		mock, gdb := setupMockDB(t)
		repo := col_db.NewRollupRepo(&UkamaDbMock{GormDb: gdb})

		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "analytics_business_sales_rollup_daily".*ON CONFLICT.*`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		// Act
		err := repo.UpsertBusinessSalesDaily(&col_db.BusinessSalesRollupDaily{
			Day:       time.Now().Truncate(24 * time.Hour),
			NetworkId: uuid.NewV4(),
			SiteId:    uuid.NewV4(),
			Revenue:   100.0,
			Purchases: 5,
		})

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_RollupRepo_RebuildSalesDaily(t *testing.T) {
	t.Run("RebuildExecutesAggregate", func(t *testing.T) {
		// Arrange
		mock, gdb := setupMockDB(t)
		repo := col_db.NewRollupRepo(&UkamaDbMock{GormDb: gdb})

		from := time.Now().AddDate(0, 0, -7)
		to := time.Now()

		mock.ExpectExec(`INSERT INTO analytics_business_sales_rollup_daily.*FROM analytics_payment_events.*ON CONFLICT.*`).
			WithArgs(from, to).
			WillReturnResult(sqlmock.NewResult(0, 3))

		// Act
		err := repo.RebuildSalesDaily(from, to)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("RebuildPropagatesError", func(t *testing.T) {
		// Arrange
		mock, gdb := setupMockDB(t)
		repo := col_db.NewRollupRepo(&UkamaDbMock{GormDb: gdb})

		from := time.Now().AddDate(0, 0, -7)
		to := time.Now()

		mock.ExpectExec(`INSERT INTO analytics_business_sales_rollup_daily.*`).
			WithArgs(from, to).
			WillReturnError(errors.New("db failure"))

		// Act
		err := repo.RebuildSalesDaily(from, to)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
