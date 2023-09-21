package db_test

import (
	extsql "database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/ukama/ukama/systems/node/health/pkg/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/tj/assert"

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

// func TestHealthRepo_StoreRunningAppsInfo(t *testing.T) {
// 	health := db.Health{
// 		Id:        uuid.NewV4(),
// 		NodeID:    uuid.NewV4(),
// 		Timestamp: time.Now().String(),
// 		System: []db.System{
// 			{
// 				Id:    uuid.NewV4(),
// 				Name:  "test",
// 				Value: "test",
// 			},
// 		},
// 		Capps: []db.Capp{
// 			{
// 				Id:     uuid.NewV4(),
// 				Name:   "test",
// 				Tag:    "test",
// 				Status: db.Running, // Assuming db.Running represents the Running status
// 				Resources: []db.Resource{
// 					{
// 						Id:    uuid.NewV4(),
// 						Name:  "test",
// 						Value: "test",
// 					},
// 				},
// 			},
// 		},
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 		DeletedAt: nil,
// 	}

// 	_db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)

// 	dialector := postgres.New(postgres.Config{
// 		DSN:                  "sqlmock_db_0",
// 		DriverName:           "postgres",
// 		Conn:                 _db,
// 		PreferSimpleProtocol: true,
// 	})

// 	gdb, err := gorm.Open(dialector, &gorm.Config{})
// 	assert.NoError(t, err)

// 	r := db.NewHealthRepo(&UkamaDbMock{
// 		GormDb: gdb,
// 	})

// 	t.Run("AddHealth", func(t *testing.T) {
// 		mock.ExpectBegin()

// 		// Insert System data
// 		for _, sys := range health.System {
// 			mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "system"`)).
// 				WithArgs(sys.Id, health.Id, sys.Name, sys.Value).
// 				WillReturnResult(sqlmock.NewResult(1, 1))
// 		}

// 		// Insert Capps data
// 		for _, capp := range health.Capps {
// 			mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "capps"`)).
// 				WithArgs(capp.Id, health.Id, capp.Name, capp.Tag, capp.Status).
// 				WillReturnResult(sqlmock.NewResult(1, 1))

// 			// Insert Resource data
// 			for _, resource := range capp.Resources {
// 				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "resource"`)).
// 					WithArgs(resource.Id, capp.Id, resource.Name, resource.Value).
// 					WillReturnResult(sqlmock.NewResult(1, 1))
// 			}
// 		}
// 		// systemJSON, _ := json.Marshal(health.System)
// 		// cappsJSON, _ := json.Marshal(health.Capps)

// 		// Insert Health data
// 		// Insert Health data

// 		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "healths"`)).
// 			WithArgs(
// 				health.Id,
// 				health.NodeID,
// 				health.Timestamp,
// 				health.CreatedAt,
// 				health.UpdatedAt,
// 				health.DeletedAt,
// 			).
// 			WillReturnResult(sqlmock.NewResult(1, 1))
// 		mock.ExpectCommit()

// 		// Act
// 		err = r.StoreRunningAppsInfo(&health, nil)

// 		// Assert
// 		assert.NoError(t, err)

// 		err = mock.ExpectationsWereMet()
// 		assert.NoError(t, err)
// 	})
// }

func TestHealthRepo_GetRunningAppsInfo(t *testing.T) {
	t.Run("HealthExist", func(t *testing.T) {
		// Arrange
		health := db.Health{
			Id:        uuid.NewV4(),
			NodeID:    uuid.NewV4(),
			Timestamp: time.Now().String(),
			System: []db.System{
				{
					Id:    uuid.NewV4(),
					Name:  "test",
					Value: "test",
				},
			},
			Capps: []db.Capp{
				{
					Id:     uuid.NewV4(),
					Name:   "test",
					Tag:    "test",
					Status: db.Running, // Assuming db.Running represents the Running status
					Resources: []db.Resource{
						{
							Id:    uuid.NewV4(),
							Name:  "test",
							Value: "test",
						},
					},
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		}

		var _db *extsql.DB

		_db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		// Convert System and Capps slices to JSON strings
		systemJSON, _ := json.Marshal(health.System)
		cappsJSON, _ := json.Marshal(health.Capps)

		rows := sqlmock.NewRows([]string{"id", "node_id", "time_stamp", "system", "capps"}).
			AddRow(health.Id, health.NodeID, health.Timestamp, systemJSON, cappsJSON)

		mock.ExpectQuery(`^SELECT.*healths.*`).
			WithArgs(health.NodeID).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 _db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := db.NewHealthRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		rm, err := r.GetRunningAppsInfo(health.NodeID)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, rm)
	})
}
