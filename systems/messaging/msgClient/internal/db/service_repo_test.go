package db_test

import (
	extsql "database/sql"
	"testing"

	int_db "github.com/ukama/ukama/systems/messaging/msgClient/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Commenting below lines because of linting errors
// var route1 = int_db.Route{
// 	Key: "event.cloud.lookup.organization.create",
// }

var service1 = int_db.Service{
	Name:        "test",
	InstanceId:  "1",
	ServiceUuid: "1ce2fa2f-2997-422c-83bf-92cf2e7334dd",
	MsgBusUri:   "",
	ListQueue:   "",
	PublQueue:   "",
	Exchange:    "test-exchange",
	ServiceUri:  "test-service.host",
	GrpcTimeout: 5,
}

func Test_serviceRepo_Get(t *testing.T) {

	t.Run("ServiceExist", func(t *testing.T) {
		// Arrange

		var db *extsql.DB
		var err error

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

		gdbx := gdb.Debug()

		r := int_db.NewServiceRepo(&UkamaDbMock{
			GormDb: gdbx,
		})

		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"name", "instance_id", "service_uuid", "msg_bus_uri", "list_queue", "publ_queue", "exchange", "service_uri", "grpc_timeout"}).
			AddRow(service1.Name, service1.InstanceId, service1.ServiceUuid, service1.MsgBusUri, service1.ListQueue, service1.PublQueue, service1.Exchange, service1.ServiceUri, service1.GrpcTimeout)

		mock.ExpectQuery(`^SELECT.*services.*`).
			WithArgs(service1.ServiceUuid).
			WillReturnRows(rows)

		// Act
		svc, err := r.Get(service1.ServiceUuid)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		if assert.NotNil(t, svc) {
			assert.Equal(t, svc.Name, service1.Name)
		}

	})

}

func Test_serviceRepo_List(t *testing.T) {

	t.Run("ListService", func(t *testing.T) {
		// Arrange

		var db *extsql.DB
		var err error

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

		gdbx := gdb.Debug()

		r := int_db.NewServiceRepo(&UkamaDbMock{
			GormDb: gdbx,
		})

		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"name", "instance_id", "service_uuid", "msg_bus_uri", "list_queue", "publ_queue", "exchange", "service_uri", "grpc_timeout"}).
			AddRow(service1.Name, service1.InstanceId, service1.ServiceUuid, service1.MsgBusUri, service1.ListQueue, service1.PublQueue, service1.Exchange, service1.ServiceUri, service1.GrpcTimeout)

		mock.ExpectQuery(`^SELECT.*services.*`).
			WillReturnRows(rows)

		// Act
		svc, err := r.List()

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		if assert.NotNil(t, svc) {
			assert.Equal(t, svc[0].Name, service1.Name)
		}

	})

}
