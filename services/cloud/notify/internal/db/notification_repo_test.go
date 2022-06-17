package db_test

import (
	extsql "database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
	"github.com/tj/assert"
	intdb "github.com/ukama/ukama/services/cloud/notify/internal/db"
	"github.com/ukama/ukama/services/common/ukama"
	jdb "gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UkamaDbMock struct {
	GormDb *gorm.DB
}

func (u UkamaDbMock) Init(model ...interface{}) error {
	panic("implement me")
}

func (u UkamaDbMock) Connect() error {
	panic("implement me")
}

func (u UkamaDbMock) GetGormDb() *gorm.DB {
	return u.GormDb
}

func (u UkamaDbMock) InitDB() error {
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction(dbOperation func(tx *gorm.DB) *gorm.DB, nestedFuncs ...func() error) error {
	panic("implement me")
}

func (u UkamaDbMock) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB, nestedFuncs ...func(tx *gorm.DB) error) (err error) {
	panic("implement me")
}

func NewTestDbNotification(nodeID string, ntype string) intdb.Notification {
	return intdb.Notification{
		NotificationID: uuid.NewV4(),
		NodeID:         nodeID,
		NodeType:       *ukama.GetNodeType(nodeID),
		Severity:       intdb.SeverityType("high"),
		Type:           intdb.NotificationType(ntype),
		ServiceName:    "noded",
		Time:           uint32(time.Now().Unix()),
		Description:    "Some random alert",
		Details:        jdb.JSON(`{"reason": "testing", "component":"router_test"}`),
	}
}

func prepare_db(t *testing.T) (sqlmock.Sqlmock, *gorm.DB) {
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

	return mock, gdb
}

func Test_Insert(t *testing.T) {
	node := ukama.NewVirtualHomeNodeId()
	nt := NewTestDbNotification(node.String(), "alert")
	var db *extsql.DB
	var err error

	db, mock, err := sqlmock.New() // mock sql.DB
	assert.NoError(t, err)

	// rows := sqlmock.NewRows([]string{"notification_id", "node_id", "node_type", "severity", "service_name", "time", "description", "details"}).
	// 	AddRow(nt.NotificationID, nt.NodeID, nt.NodeType, nt.Severity, nt.ServiceName, nt.Time, nt.Description, nt.Details)

	// mock.ExpectExec("^INSERT INTO \"groups\" (.+)$").WithArgs(AnyTime{}, AnyTime{}, nil, name, description).WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectQuery(regexp.QuoteMeta(
	// 	`INSERT INTO  ("id","created_at","updated_at","delted_at","notification_id", "node_id", "node_type", "severity", "service_name", "time", "description", "details")
	// 			 VALUES ($1,$2,$3,4,$5,$6,$7,$8) RETURNING "person"."id"`)).
	// 	WithArgs(nt.NotificationID, nt.NodeID, nt.NodeType, nt.Severity, nt.ServiceName, nt.Time, nt.Description, nt.Details).
	// 	WillReturnRows(
	// 		sqlmock.NewRows([]string{"id"}).AddRow(id.String()))

	// mock.ExpectQuery(`^SELECT.*nodes.*`).
	// 	WithArgs(node.NodeID).
	// 	WillReturnRows(rows)

	//mock.ExpectQuery("^INSERT INTO \"notifications\" (.+)$").WithArgs(nt.NotificationID, nt.NodeID, nt.NodeType, nt.Severity, nt.ServiceName, nt.Time, nt.Description, nt.Details).WillReturnRows(sqlmock.Rows(1, 1))

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	r := intdb.NewNotificationRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	assert.NoError(t, err)

	// Act
	err = r.Insert(&nt)

	// Assert
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
	assert.NoError(t, err)

}

func Test_List(t *testing.T) {

	node := ukama.NewVirtualHomeNodeId()
	nt := NewTestDbNotification(node.String(), "alert")
	var err error

	mock, gdb := prepare_db(t)

	rows := sqlmock.NewRows([]string{"notification_id", "node_id", "node_type", "severity", "service_name", "time", "description", "details"}).
		AddRow(nt.NotificationID, nt.NodeID, nt.NodeType, nt.Severity, nt.ServiceName, nt.Time, nt.Description, nt.Details)

	mock.ExpectQuery(`^SELECT.*notifications.*`).
		WithArgs().
		WillReturnRows(rows)

	r := intdb.NewNotificationRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	assert.NoError(t, err)

	// Act
	list, err := r.List()

	// Assert
	assert.NoError(t, err)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

	assert.NoError(t, err)
	assert.NotNil(t, list)

}

func Test_GetNotificationForService(t *testing.T) {

	node := ukama.NewVirtualHomeNodeId()
	nt := NewTestDbNotification(node.String(), "alert")
	var err error

	mock, gdb := prepare_db(t)

	rows := sqlmock.NewRows([]string{"notification_id", "node_id", "node_type", "severity", "service_name", "time", "description", "details"}).
		AddRow(nt.NotificationID, nt.NodeID, nt.NodeType, nt.Severity, nt.ServiceName, nt.Time, nt.Description, nt.Details)

	mock.ExpectQuery(`^SELECT.*notifications.*`).
		WithArgs().
		WillReturnRows(rows)

	r := intdb.NewNotificationRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	assert.NoError(t, err)

	// Act
	list, err := r.GetNotificationForService(nt.ServiceName, string(nt.Type))

	// Assert
	assert.NoError(t, err)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

	assert.NoError(t, err)
	assert.NotNil(t, list)

}
