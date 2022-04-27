package db_test

import (
	extsql "database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/ukama/ukama/services/factory/nmr/internal/db"
	intDb "github.com/ukama/ukama/services/factory/nmr/internal/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_moduleRepo_GetModuleMfgStatus(t *testing.T) {

	t.Run("ModuleMfgStatus", func(t *testing.T) {

		module := intDb.Module{
			ModuleID:   "1001",
			Type:       "TRX",
			PartNumber: "a1",
			HwVersion:  "h1",
			Mac:        "00:01:02:03:04:05",
			SwVersion:  "1.1",
			PSwVersion: "0.1",
			MfgDate:    time.Now(),
			MfgName:    "ukama",
			Status:     "StatusLabelGenrated",
		}

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"module_id", "type", "part_number", "hw_version", "mac", "sw_version", "p_sw_version", "mfg_date", "mfg_name", "status"}).
			AddRow(module.ModuleID, module.Type, module.PartNumber, module.HwVersion, module.Mac, module.SwVersion, module.PSwVersion, module.MfgDate, module.MfgName, module.Status)

		mock.ExpectQuery(`^SELECT.*modules.*`).
			WithArgs(module.ModuleID).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := intDb.NewModuleRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		modStatus, err := r.GetModuleMfgStatus(module.ModuleID)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		if assert.NotNil(t, modStatus) {
			assert.Equal(t, module.Status, (*modStatus).String())
		}
	})

}

func Test_moduleRepo_UpdateModuleMfgStatus(t *testing.T) {

	t.Run("UpdateModuleMfgStatus", func(t *testing.T) {

		status := db.MfgStatus("ModuleTest")
		moduleId := "1001"

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE")).WithArgs(status, "1001").
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

		r := intDb.NewModuleRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.UpdateModuleMfgStatus(moduleId, status)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})

}

func Test_moduleRepo_GetModule(t *testing.T) {

	t.Run("GetModule", func(t *testing.T) {

		module := intDb.Module{
			ModuleID:   "1001",
			Type:       "TRX",
			PartNumber: "a1",
			HwVersion:  "h1",
			Mac:        "00:01:02:03:04:05",
			SwVersion:  "1.1",
			PSwVersion: "0.1",
			MfgDate:    time.Now(),
			MfgName:    "ukama",
			Status:     "StatusLabelGenrated",
		}

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"module_id", "type", "part_number", "hw_version", "mac", "sw_version", "p_sw_version", "mfg_date", "mfg_name", "status"}).
			AddRow(module.ModuleID, module.Type, module.PartNumber, module.HwVersion, module.Mac, module.SwVersion, module.PSwVersion, module.MfgDate, module.MfgName, module.Status)

		mock.ExpectQuery(`^SELECT.*modules.*`).
			WithArgs(module.ModuleID).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := intDb.NewModuleRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		mod, err := r.GetModule(module.ModuleID)

		// Assert
		assert.NoError(t, err)

		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}

		assert.NoError(t, err)
		assert.NotNil(t, mod)
	})

}

func Test_moduleRepo_GetModuleFailure(t *testing.T) {
	t.Run("ModuleDoesNotExist", func(t *testing.T) {

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

		r := intDb.NewModuleRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		modData, err := r.GetModule("1002")

		// Assert
		if assert.Error(t, err) {
			assert.Equal(t, fmt.Errorf("no matching id"), err)
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}

		assert.NoError(t, err)
		assert.Nil(t, modData)
	})

}

func Test_moduleRepo_DeleteModule(t *testing.T) {
	t.Run("DeleteModule", func(t *testing.T) {

		var db *extsql.DB
		var err error
		modId := "1001"
		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE")).WithArgs(modId).
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

		r := intDb.NewModuleRepo(&UkamaDbMock{
			GormDb: gdb,
		})
		assert.NoError(t, err)

		// Act
		err = r.DeleteModule(modId)

		// Assert
		assert.NoError(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}

		assert.NoError(t, err)
	})

}

func Test_moduleRepo_GetModuleList(t *testing.T) {

	t.Run("ModuleList", func(t *testing.T) {

		module := intDb.Module{
			ModuleID:   "1001",
			Type:       "TRX",
			PartNumber: "a1",
			HwVersion:  "h1",
			Mac:        "00:01:02:03:04:05",
			SwVersion:  "1.1",
			PSwVersion: "0.1",
			MfgDate:    time.Now(),
			MfgName:    "ukama",
			Status:     "LabelGenrated",
		}

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"module_id", "type", "part_number", "hw_version", "mac", "sw_version", "p_sw_version", "mfg_date", "mfg_name", "status"}).
			AddRow(module.ModuleID, module.Type, module.PartNumber, module.HwVersion, module.Mac, module.SwVersion, module.PSwVersion, module.MfgDate, module.MfgName, module.Status).
			AddRow("1002", module.Type, module.PartNumber, module.HwVersion, module.Mac, module.SwVersion, module.PSwVersion, module.MfgDate, module.MfgName, module.Status)

		mock.ExpectQuery(`^SELECT.*modules.*`).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := intDb.NewModuleRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		modList, err := r.ListModules()

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		if assert.NotNil(t, modList) {
			assert.Equal(t, 2, len(*modList))
		}
	})

}
