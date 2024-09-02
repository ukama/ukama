/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	extsql "database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	ukama "github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UkamaDbMock struct {
	GormDb *gorm.DB
}

func (u UkamaDbMock) Init(model ...interface{}) error {
	return nil
}

func (u UkamaDbMock) Connect() error {
	return nil
}

func (u UkamaDbMock) GetGormDb() *gorm.DB {
	return u.GormDb
}

func (u UkamaDbMock) InitDB() error {
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func() error) error {
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func(tx *gorm.DB) error) error {
	return nil
}

var nodeId = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

func TestCreate(t *testing.T) {
	stateId := uuid.NewV4()
	state := State{
		Id:              stateId,
		NodeId:          nodeId.String(),
		State:           ukama.StateConfigure,
		Type:            "hnode",
		LastHeartbeat:   time.Now(),
		LastStateChange: time.Now(),
		Connectivity:    Online,
		Version:         "v34",
	}

	db, mock, err := sqlmock.New() // mock sql.DB
	assert.NoError(t, err)

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

	t.Run("Create", func(t *testing.T) {
		// Arrange
		mock.ExpectBegin()

		// Adjust the expected arguments according to your State structure
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "states"`)).
			WithArgs(state.Id, state.NodeId, state.State, state.LastHeartbeat, state.LastStateChange, state.Connectivity, state.Type, state.Version, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Act
		err = r.Create(&state, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
func GetStateByNodeId(t *testing.T) {
	t.Run("StateExist", func(t *testing.T) {
		// Arrange
		stateId := uuid.NewV4()
		state := State{
			Id:              stateId,
			NodeId:          nodeId.String(),
			State:           ukama.StateConfigure,
			Type:            "hnode",
			LastHeartbeat:   time.Now(),
			LastStateChange: time.Now(),
			Connectivity:    Online,
			Version:         "v34",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			DeletedAt:       gorm.DeletedAt{},
		}

		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "node_id", "state", "type", "lastheartbeat", "laststatechange", "connectivity", "version"}).
			AddRow(stateId, state.NodeId, state.State, state.Type, state.LastHeartbeat, state.LastStateChange, state.Connectivity, state.Version)

		mock.ExpectQuery(`^SELECT.*states.*`).
			WithArgs(state.NodeId, 1). // Add '1' for the LIMIT clause
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

		assert.NoError(t, err)

		// Act
		rm, err := r.GetByNodeId(nodeId)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, rm)
	})
}

func TestGetStateHistory(t *testing.T) {
	t.Run("GetStateHistory", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		fromTime := time.Now().Add(-24 * time.Hour)
		toTime := time.Now()
		expectedStates := []State{
			{
				Id:              uuid.NewV4(),
				NodeId:          nodeId.String(),
				State:           ukama.StateConfigure,
				Type:            "hnode",
				LastHeartbeat:   time.Now(),
				LastStateChange: time.Now(),
				Connectivity:    Online,
				Version:         "v34",
			},
		}

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "states"`)).
			WithArgs(nodeId.String(), fromTime, toTime).
			WillReturnRows(sqlmock.NewRows([]string{"id", "node_id", "state", "last_heartbeat", "last_state_change", "connectivity", "type", "version"}).
				AddRow(expectedStates[0].Id, expectedStates[0].NodeId, expectedStates[0].State, expectedStates[0].LastHeartbeat, expectedStates[0].LastStateChange, expectedStates[0].Connectivity, expectedStates[0].Type, expectedStates[0].Version))

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
		states, err := r.GetStateHistory(nodeId, fromTime, toTime)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedStates, states)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
