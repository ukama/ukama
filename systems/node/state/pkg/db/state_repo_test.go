/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
package db

import (
	"log"
	"testing"
	"time"

	extsql "database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UkamaDbMock struct {
	GormDb *gorm.DB
}

func (u UkamaDbMock) Init(model ...interface{}) error {
	panic("implement me: Init()")
}

func (u UkamaDbMock) Connect() error {
	panic("implement me: Connect()")
}

func (u UkamaDbMock) GetGormDb() *gorm.DB {
	return u.GormDb
}

func (u UkamaDbMock) InitDB() error {
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func() error) error {
	log.Fatal("implement me: ExecuteInTransaction()")
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func(tx *gorm.DB) error) error {
	log.Fatal("implement me: ExecuteInTransaction2()")
	return nil
}

func TestState_GetLatestState(t *testing.T) {
	t.Run("state exists for node", func(t *testing.T) {
		nid := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() 
		assert.NoError(t, err)

		latestState := State{
			NodeId:    nid.String(),
			CreatedAt: time.Now(), 
		}

		rows := sqlmock.NewRows([]string{"node_id", "created_at"}).
			AddRow(latestState.NodeId, latestState.CreatedAt)

		mock.ExpectQuery(`^SELECT \* FROM "states" WHERE node_id = \$1 AND "states"."deleted_at" IS NULL ORDER BY created_at DESC,"states"."id" LIMIT \$2`).
			WithArgs(nid.String(), 1).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewStateRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		c, err := r.GetLatestState(nid.String())

		assert.NoError(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, nid.String(), c.NodeId)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}



func TestState_GetStateHistory(t *testing.T) {
	
	t.Run("GetStateHistory", func(t *testing.T) {
		nid := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE).String()
		stateId := uuid.NewV4()

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() 
		assert.NoError(t, err)

		history := State{
			Id:        stateId,
			NodeId:    nid,
			CreatedAt: time.Now(), 
		}
   rows := sqlmock.NewRows([]string{"id", "node_id", "created_at"}).
		AddRow(history.Id, history.NodeId, history.CreatedAt)

		mock.ExpectQuery(`^SELECT \* FROM "states" WHERE node_id = \$1 AND "states"."deleted_at" IS NULL ORDER BY created_at DESC,"states"."id" LIMIT \$2`).
		WithArgs(nid, 1).
		WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

	gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewStateRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		c, err := r.GetStateHistory(nid)

		assert.NoError(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, stateId, c[0].Id)
		
	})
}


