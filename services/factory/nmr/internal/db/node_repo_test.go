package db_test

import (
	"database/sql"
	extsql "database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	intDb "github.com/ukama/openIoR/services/factory/nmr/internal/db"
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

func Test_nodeRepo_GetNodeStatus(t *testing.T) {

	t.Run("NodeExist", func(t *testing.T) {

		node := intDb.Node{
			NodeID:        "1001",
			Type:          "hnode",
			PartNumber:    "a1",
			Skew:          "s1",
			Mac:           "00:01:02:03:04:05",
			SwVersion:     "1.1",
			PSwVersion:    "0.1",
			AssemblyDate:  time.Now(),
			OemName:       "ukama",
			MfgTestStatus: "Pending",
			Status:        "StatusLabelGenrated",
		}

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"node_id", "type", "part_number", "skew", "mac", "sw_version", "p_sw_version", "assembly_date", "oem_name", "mfg_test_status", "status"}).
			AddRow(node.NodeID, node.Type, node.PartNumber, node.Skew, node.Mac, node.SwVersion, node.PSwVersion, node.AssemblyDate, node.OemName, node.MfgTestStatus, node.Status)

		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(node.NodeID).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := intDb.NewNodeRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		nodeStatus, err := r.GetNodeStatus(node.NodeID)

		// Assert
		assert.NoError(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
		assert.NoError(t, err)

		if assert.NotNil(t, nodeStatus) {
			assert.Equal(t, node.Status, (*nodeStatus).String())
		}
	})

}

func Test_nodeRepo_GetNode(t *testing.T) {

	t.Run("NodeExist", func(t *testing.T) {

		node := intDb.Node{
			NodeID:        "1001",
			Type:          "hnode",
			PartNumber:    "a1",
			Skew:          "s1",
			Mac:           "00:01:02:03:04:05",
			SwVersion:     "1.1",
			PSwVersion:    "0.1",
			AssemblyDate:  time.Now(),
			OemName:       "ukama",
			MfgTestStatus: "MfgTestStatusPending",
			Status:        "StatusLabelGenrated",
		}

		module := intDb.Module{
			ModuleID:   "M1001",
			Type:       "TRX",
			PartNumber: "a1",
			HwVersion:  "h1",
			Mac:        "00:01:02:03:04:05",
			SwVersion:  "1.1",
			PSwVersion: "0.1",
			MfgDate:    time.Now(),
			MfgName:    "ukama",
			Status:     "StatusLabelGenrated",
			UnitID:     sql.NullString{String: node.NodeID, Valid: true},
		}

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"node_id", "type", "part_number", "skew", "mac", "sw_version", "p_sw_version", "assembly_date", "oem_name", "mfg_test_status", "status"}).
			AddRow(node.NodeID, node.Type, node.PartNumber, node.Skew, node.Mac, node.SwVersion, node.PSwVersion, node.AssemblyDate, node.OemName, node.MfgTestStatus, node.Status)

		mrows := sqlmock.NewRows([]string{"module_id", "type", "part_number", "hw_version", "mac", "sw_version", "p_sw_version", "mfg_date", "mfg_name", "status", "unit_id"}).
			AddRow(module.ModuleID, module.Type, module.PartNumber, module.HwVersion, module.Mac, module.SwVersion, module.PSwVersion, module.MfgDate, module.MfgName, module.Status, module.UnitID)

		mock.ExpectQuery(`SELECT`).
			WithArgs(node.NodeID).
			WillReturnRows(rows)

		mock.ExpectQuery(`^SELECT.*modules.*`).
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(mrows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := intDb.NewNodeRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		nodeData, err := r.GetNode(node.NodeID)

		// Assert
		assert.NoError(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
		assert.NoError(t, err)

		if assert.NotNil(t, nodeData) {
			assert.Equal(t, node.NodeID, nodeData.NodeID)
		}
	})

}

func Test_nodeRepo_GetNodeFailure(t *testing.T) {
	t.Run("NodeDoesNotExist", func(t *testing.T) {

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectQuery(`SELECT`).
			WithArgs("1002").
			WillReturnError(fmt.Errorf("no matching id"))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := intDb.NewNodeRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		nodeData, err := r.GetNode("1002")

		// Assert
		if assert.Error(t, err) {
			assert.Equal(t, fmt.Errorf("no matching id"), err)
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}

		assert.NoError(t, err)
		assert.Nil(t, nodeData)
	})

}

func Test_nodeRepo_DeleteNode(t *testing.T) {
	t.Run("DeleteNode", func(t *testing.T) {

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE")).WithArgs("1001").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := intDb.NewNodeRepo(&UkamaDbMock{
			GormDb: gdb,
		})
		assert.NoError(t, err)

		// Act
		err = r.DeleteNode("1001")

		// Assert
		assert.NoError(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}

		assert.NoError(t, err)
	})

}

func Test_nodeRepo_GetNodeList(t *testing.T) {

	t.Run("NodeList", func(t *testing.T) {

		node := intDb.Node{
			NodeID:        "1001",
			Type:          "hnode",
			PartNumber:    "a1",
			Skew:          "s1",
			Mac:           "00:01:02:03:04:05",
			SwVersion:     "1.1",
			PSwVersion:    "0.1",
			AssemblyDate:  time.Now(),
			OemName:       "ukama",
			MfgTestStatus: "MfgTestStatusPending",
			Status:        "StatusLabelGenrated",
		}

		module := intDb.Module{
			ModuleID:   "M1001",
			Type:       "TRX",
			PartNumber: "a1",
			HwVersion:  "h1",
			Mac:        "00:01:02:03:04:05",
			SwVersion:  "1.1",
			PSwVersion: "0.1",
			MfgDate:    time.Now(),
			MfgName:    "ukama",
			Status:     "StatusLabelGenrated",
			UnitID:     sql.NullString{String: node.NodeID, Valid: true},
		}

		// nodeID2 := "1002"
		// unitID2 := sql.NullString{String: nodeId2, Valid: true}

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"node_id", "type", "part_number", "skew", "mac", "sw_version", "p_sw_version", "assembly_date", "oem_name", "status", "status"}).
			AddRow(node.NodeID, node.Type, node.PartNumber, node.Skew, node.Mac, node.SwVersion, node.PSwVersion, node.AssemblyDate, node.OemName, node.Status, node.Status).
			AddRow("1002", node.Type, node.PartNumber, node.Skew, node.Mac, node.SwVersion, node.PSwVersion, node.AssemblyDate, node.OemName, node.Status, node.Status)

		mrow := sqlmock.NewRows([]string{"module_id", "type", "part_number", "hw_version", "mac", "sw_version", "p_sw_version", "mfg_date", "mfg_name", "status", "unit_id"}).
			AddRow(module.ModuleID, module.Type, module.PartNumber, module.HwVersion, module.Mac, module.SwVersion, module.PSwVersion, module.MfgDate, module.MfgName, module.Status, module.UnitID)

		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WillReturnRows(rows)

		mock.ExpectQuery(`^SELECT.*modules.*`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(mrow)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := intDb.NewNodeRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		nodeList, err := r.ListNodes()

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		if assert.NotNil(t, nodeList) {
			assert.Equal(t, 2, len(*nodeList))
		}
	})

}
