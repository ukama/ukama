package db

import (
	"fmt"
	"testing"

	extsql "database/sql"

	log "github.com/sirupsen/logrus"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/DATA-DOG/go-sqlmock"
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

func TestHealthRepo_GetRunningAppsInfo(t *testing.T) {
	var nodeId = ukama.NewVirtualNodeId("HomeNode")
	id := uuid.NewV4()

	health := Health{
		Id:        id,
		NodeId:    nodeId.String(),
		TimeStamp: "test",
		System: []System{
			{
				Id:    id,
				Name:  "test",
				Value: "test",
			},
		},
		Capps: []Capp{
			{
				Id:     id,
				Name:   "test",
				Tag:    "test",
				Status: Status(1),
				Resources: []Resource{
					{
						Id:    id,
						Name:  "test",
						Value: "test",
					},
				},
			},
		},
	}
	t.Run("RunningAppFound", func(t *testing.T) {
		var db *extsql.DB

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

		r := NewHealthRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		healthRows := sqlmock.NewRows([]string{"id", "node_id", "time_stamp"}).
			AddRow(health.Id, nodeId.String(), "12-12-2024")

		systemRows := sqlmock.NewRows([]string{"id", "health_id", "name", "value"}).
			AddRow(health.Id, health.Id, health.System[0].Name, health.System[0].Value)

		cappRows := sqlmock.NewRows([]string{"id", "health_id", "name", "tag", "status"}).
			AddRow(health.Id, health.Id, health.Capps[0].Name, health.Capps[0].Tag, health.Capps[0].Status)

		resourceRows := sqlmock.NewRows([]string{"id", "capp_id", "name", "value"}).
			AddRow(health.Capps[0].Resources[0].Id, health.Capps[0].Id, health.Capps[0].Resources[0].Name, health.Capps[0].Resources[0].Value)

		mock.ExpectQuery(`^SELECT.*healths.*`).
			WithArgs(nodeId).
			WillReturnRows(healthRows)

		mock.ExpectQuery(`^SELECT.*systems.*`).
			WithArgs(health.System[0].Id).
			WillReturnRows(systemRows)

		mock.ExpectQuery(`^SELECT.*capps.*`).
			WithArgs(health.Capps[0].Id).
			WillReturnRows(cappRows, resourceRows)

		mock.ExpectQuery(`^SELECT.*resources.*`).
			WithArgs(health.Capps[0].Resources[0].CappID).
			WillReturnRows(resourceRows)

		mock.ExpectQuery(`^SELECT \* FROM "capps" WHERE "capps"."health_id" = \$1`).
			WithArgs(health.NodeId). // Update this to match the actual parameter used
			WillReturnRows(cappRows)

		apps, err := r.GetRunningAppsInfo(nodeId)
		fmt.Println("BRACKLEY", apps)
		// Assert
		assert.NotNil(t, apps)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})

}
