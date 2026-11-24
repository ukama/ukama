/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/uuid"
)

func Test_AuditAdd(t *testing.T) {
	t.Run("FillsIdAndTime", func(t *testing.T) {
		mock, udb := setupTestDB(t)
		repo := NewAuditRepo(udb)
		opId := uuid.NewV4()
		audit := &OperationAudit{OperationId: opId, ResourceKey: "node:a",
			Event: "rejected_running", Actor: "controller", Reason: "lease expired"}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "operation_audits"`)).
			WithArgs(sqlmock.AnyArg(), opId, "node:a", "rejected_running", "controller",
				"lease expired", sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Add(ctx, audit)

		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, audit.Id)
		assert.False(t, audit.At.IsZero())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("KeepsGivenIdAndTimeInTx", func(t *testing.T) {
		mock, udb := setupTestDB(t)
		repo := NewAuditRepo(udb)
		id := uuid.NewV4()
		at := time.Now().Add(-time.Minute).UTC()
		audit := &OperationAudit{Id: id, OperationId: uuid.NewV4(), Event: "completed",
			At: at}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "operation_audits"`)).
			WithArgs(id, audit.OperationId, "", "completed", "", "", at).
			WillReturnResult(sqlmock.NewResult(1, 1))

		assert.NoError(t, repo.AddInTx(udb.GetGormDb().Begin(), audit))
		assert.Equal(t, id, audit.Id)
		assert.Equal(t, at, audit.At)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ErrorIsReturned", func(t *testing.T) {
		mock, udb := setupTestDB(t)
		repo := NewAuditRepo(udb)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "operation_audits"`)).
			WillReturnError(errors.New("connection reset"))
		mock.ExpectRollback()

		assert.Error(t, repo.Add(ctx, &OperationAudit{Event: "deleted"}))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
