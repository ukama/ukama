package db_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/state/pkg/db"
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
	panic("implement me: ExecuteInTransaction()")
}

func (u UkamaDbMock) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func(tx *gorm.DB) error) error {
	panic("implement me: ExecuteInTransaction2()")
}

func TestState_Create(t *testing.T) {
	sqlDb, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDb.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := db.NewStateRepo(&UkamaDbMock{GormDb: gormDb})

	state := &db.State{
		Id:              uuid.NewV4(),
		NodeId:          "node1",
		CurrentState:    db.StateConfigure,
		LastHeartbeat:   time.Now(),
		LastStateChange: time.Now(),
		Type:            "someType",
		Version:         "1.0",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"states\"").WithArgs(
		state.Id,
		state.NodeId,
		state.CurrentState,
		state.LastHeartbeat,
		state.LastStateChange,
		state.Type,
		state.Version,
		state.CreatedAt,
		state.UpdatedAt,
		sqlmock.AnyArg(),
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.Create(state, nil)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestState_GetByNodeId(t *testing.T) {
	sqlDb, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDb.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := db.NewStateRepo(&UkamaDbMock{GormDb: gormDb})
	var nodeId = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

	expectedState := &db.State{
		Id:           uuid.NewV4(),
		NodeId:       nodeId.String(),
		CurrentState: db.StateOperational,
	}

	rows := sqlmock.NewRows([]string{"id", "node_id", "current_state"}).
		AddRow(expectedState.Id, expectedState.NodeId, expectedState.CurrentState)

	mock.ExpectQuery(`^SELECT.*states.*`).
		WithArgs(nodeId, sqlmock.AnyArg()).
		WillReturnRows(rows)
	state, err := repo.GetByNodeId(nodeId)
	assert.NoError(t, err)
	assert.Equal(t, expectedState.NodeId, state.NodeId)
	assert.Equal(t, expectedState.CurrentState, state.CurrentState)
	assert.NoError(t, mock.ExpectationsWereMet())
}


// func TestState_GetStateHistoryByTimeRange(t *testing.T) {
// 	sqlDb, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer sqlDb.Close()

// 	gormDb, err := gorm.Open(postgres.New(postgres.Config{
// 		Conn: sqlDb,
// 	}), &gorm.Config{})
// 	assert.NoError(t, err)

// 	repo := db.NewStateRepo(&UkamaDbMock{GormDb: gormDb})
// 	nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
// 	from := time.Now().Add(-24 * time.Hour)
// 	to := time.Now()

// 	expectedHistory := []db.StateHistory{
// 		{
// 			Id:            uuid.NewV4(),
// 			NodeStateId:   nodeId.String(),
// 			PreviousState: db.StateFaulty,
// 			NewState:      db.StateActive,
// 			Timestamp:     time.Now().Add(-12 * time.Hour),
// 		},
// 		{
// 			Id:            uuid.NewV4(),
// 			NodeStateId:   nodeId.String(),
// 			PreviousState: db.StateFaulty,
// 			NewState:      db.StateMaintenance,
// 			Timestamp:     time.Now().Add(-18 * time.Hour),
// 		},
// 	}

// 	rows := sqlmock.NewRows([]string{"id", "node_state_id", "previous_state", "new_state", "timestamp"})
// 	for _, history := range expectedHistory {
// 		rows.AddRow(history.Id, history.NodeStateId, history.PreviousState, history.NewState, history.Timestamp)
// 	}

// 	mock.ExpectQuery("SELECT \\* FROM \"state_histories\"").
// 		WithArgs(nodeId.String(), from, to).
// 		WillReturnRows(rows)

// 	history, err := repo.GetStateHistoryByTimeRange(nodeId, from, to)
// 	assert.NoError(t, err)
// 	assert.Equal(t, len(expectedHistory), len(history))
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }


func TestState_ListAll(t *testing.T) {
	sqlDb, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDb.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := db.NewStateRepo(&UkamaDbMock{GormDb: gormDb})

	expectedStates := []db.State{
		{
			Id:           uuid.NewV4(),
			NodeId:       ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE).String(),
			CurrentState: db.StateConfigure,
		},
		{
			Id:           uuid.NewV4(),
			NodeId:       ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE).String(),
			CurrentState: db.StateOperational,
		},
	}

	rows := sqlmock.NewRows([]string{"id", "node_id", "current_state", "last_heartbeat", "last_state_change", "type", "version", "created_at", "updated_at"})
	for _, state := range expectedStates {
		rows.AddRow(state.Id, state.NodeId, state.CurrentState, time.Now(), time.Now(), "type", "version", time.Now(), time.Now())
	}

	// Expect the SELECT query
	mock.ExpectQuery(`SELECT \* FROM "states" WHERE "states"."deleted_at" IS NULL`).
		WillReturnRows(rows)

	states, err := repo.ListAll()
	assert.NoError(t, err)
	assert.Equal(t, len(expectedStates), len(states))
	assert.NoError(t, mock.ExpectationsWereMet())

	// Check if the returned states match the expected states
	for i, state := range states {
		assert.Equal(t, expectedStates[i].Id, state.Id)
		assert.Equal(t, expectedStates[i].NodeId, state.NodeId)
		assert.Equal(t, expectedStates[i].CurrentState, state.CurrentState)
	}
}
