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
	"time" // Make sure to import time package for CreatedAt

	extsql "database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
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
		// Arrange
		nid := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		// Define a sample State object to return
		latestState := State{
			NodeId:    nid.String(),
			CreatedAt: time.Now(), // Assume you have this field in your State struct
		}

		rows := sqlmock.NewRows([]string{"node_id", "created_at"}).
			AddRow(latestState.NodeId, latestState.CreatedAt)

		// Update the mock query to match the actual query
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

		// Act
		c, err := r.GetLatestState(nid.String())

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, nid.String(), c.NodeId)

		// Ensure all expectations were met
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

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		// Define a sample State object to return
		history := State{
			Id:        stateId,
			NodeId:    nid,
			CreatedAt: time.Now(), // Assume you have this field in your State struct
		}
   rows := sqlmock.NewRows([]string{"id", "node_id", "created_at"}).
		AddRow(history.Id, history.NodeId, history.CreatedAt)

		// Update the mock query to match the actual query
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

		// Act
		c, err := r.GetStateHistory(nid)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, stateId, c[0].Id)
		
	})
}


func TestState_AddState(t *testing.T) {
	t.Run("AddState", func(t *testing.T) {
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		newState := &State{
			Id:               uuid.NewV4(), 
			NodeId:          ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE).String(),
			PreviousStateId: nil, 
			CurrentState:    cpb.NodeState_Unknown ,
			SubState:        StringArray([]string{"on"}), 
			Events:          StringArray([]string{"online"}), 
			Version:         "1.0.0", 
			NodeType:        "example", 
			NodeIp:         "192.168.1.1", 
			NodePort:       8080, 
			MeshIp:         "192.168.1.1", 
			MeshPort:       8081, 
			MeshHostName:   "example-host", 
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(), 
			DeletedAt:      gorm.DeletedAt{}, 
		}

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "states"`).
			WithArgs(
				newState.Id, newState.NodeId, newState.PreviousStateId, newState.CurrentState,
				newState.SubState, newState.Events, newState.Version, newState.NodeType,
				newState.NodeIp, newState.NodePort, newState.MeshIp, newState.MeshPort,
				newState.MeshHostName, newState.CreatedAt, newState.UpdatedAt,
				newState.DeletedAt, 
			).
			WillReturnResult(sqlmock.NewResult(1, 1)) // Expect the result
		mock.ExpectCommit()

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

		// Act
		err = r.AddState(newState, nil)

		// Assert
		assert.NoError(t, err)

		// Ensure all expectations were met
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
