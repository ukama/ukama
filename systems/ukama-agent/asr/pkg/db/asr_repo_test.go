package db_test

import (
	extsql "database/sql"
	"regexp"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/client"
	int_db "github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Iccid = "0123456789012345678912"

var sub = int_db.Asr{
	Iccid:          Iccid,
	Imsi:           Imsi,
	Op:             []byte("0123456789012345"),
	Key:            []byte("0123456789012345"),
	Amf:            []byte("800"),
	AlgoType:       1,
	UeDlAmbrBps:    2000000,
	UeUlAmbrBps:    2000000,
	Sqn:            1,
	CsgIdPrsent:    false,
	CsgId:          0,
	DefaultApnName: "ukama",
}

var sim = client.SimCardInfo{
	Iccid:          Iccid,
	Imsi:           Imsi,
	Op:             []byte("0123456789012345"),
	Key:            []byte("0123456789012345"),
	Amf:            []byte("800"),
	AlgoType:       1,
	UeDlAmbrBps:    2000000,
	UeUlAmbrBps:    2000000,
	Sqn:            1,
	CsgIdPrsent:    false,
	CsgId:          0,
	DefaultApnName: "ukama",
}

var tai = int_db.Tai{
	PlmnId:          "00101",
	Tac:             101,
	DeviceUpdatedAt: time.Now(),
}

func TestAsrRecordRepo_Add(t *testing.T) {

	t.Run("Add", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		//row := sqlmock.NewRows([]string{"iccid", "imsi", "op", "amf", "key", "algo_type", "ue_dl_ambr_bps", "ue_ul_ambr_bps", "sqn", "csg_id_prsent", "csg_id", "default_apn_name", "network_id", "package_id"}).

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sub.Iccid, sub.Imsi, sub.Op, sub.Amf, sub.Key, sub.AlgoType, sub.UeDlAmbrBps, sub.UeUlAmbrBps, sub.Sqn, sub.CsgIdPrsent, sub.CsgId, sub.DefaultApnName, sub.NetworkID, sub.PackageId).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewAsrRecordRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Add(&sub)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})

}

func TestAsrRecordRepo_Update(t *testing.T) {

	t.Run("UpdatePackage", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		var err error
		PackageId := uuid.NewV4()
		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
			WithArgs(sqlmock.AnyArg(), PackageId, Imsi).
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

		r := int_db.NewAsrRecordRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.UpdatePackage(Imsi, PackageId)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})

}

func TestAsrRecordRepo_Get(t *testing.T) {

	t.Run("ReadByID", func(t *testing.T) {
		sub.ID = 1
		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		hrow := sqlmock.NewRows([]string{"ID", "iccid", "imsi", "op", "amf", "key", "algo_type", "ue_dl_ambr_bps", "ue_ul_ambr_bps", "sqn", "csg_id_prsent", "csg_id", "default_apn_name", "network_id", "package_id"}).
			AddRow(sub.ID, sub.Iccid, sub.Imsi, sub.Op, sub.Amf, sub.Key, sub.AlgoType, sub.UeDlAmbrBps, sub.UeDlAmbrBps, sub.Sqn, sub.CsgIdPrsent, sub.CsgId, sub.DefaultApnName, sub.NetworkID, sub.PackageId)

		trow := sqlmock.NewRows([]string{"asr_id", "plmn_id", "tac", "device_updated_at"})

		mock.ExpectQuery(`^SELECT.*asrs.*`).
			WithArgs(sub.ID).
			WillReturnRows(hrow)

		mock.ExpectQuery(`^SELECT.*tais.*`).
			WithArgs(sub.ID).
			WillReturnRows(trow)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewAsrRecordRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		asr, err := r.Get(int(sub.ID))

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

		if assert.NotNil(t, asr) {
			assert.EqualValues(t, asr.ID, sub.ID)
		}

	})

	t.Run("ReadByICCID", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		hrow := sqlmock.NewRows([]string{"iccid", "imsi", "op", "amf", "key", "algo_type", "ue_dl_ambr_bps", "ue_ul_ambr_bps", "sqn", "csg_id_prsent", "csg_id", "default_apn_name", "network_id", "package_id"}).
			AddRow(sub.Iccid, sub.Imsi, sub.Op, sub.Amf, sub.Key, sub.AlgoType, sub.UeDlAmbrBps, sub.UeDlAmbrBps, sub.Sqn, sub.CsgIdPrsent, sub.CsgId, sub.DefaultApnName, sub.NetworkID, sub.PackageId)

		mock.ExpectQuery(`^SELECT.*asrs.*`).
			WithArgs(Iccid).
			WillReturnRows(hrow)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewAsrRecordRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		asr, err := r.GetByIccid(Iccid)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

		if assert.NotNil(t, asr) {
			assert.EqualValues(t, asr.Iccid, Iccid)
		}

	})

	t.Run("ReadByImsi", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		hrow := sqlmock.NewRows([]string{"iccid", "imsi", "op", "amf", "key", "algo_type", "ue_dl_ambr_bps", "ue_ul_ambr_bps", "sqn", "csg_id_prsent", "csg_id", "default_apn_name", "network_id", "package_id"}).
			AddRow(sub.Iccid, sub.Imsi, sub.Op, sub.Amf, sub.Key, sub.AlgoType, sub.UeDlAmbrBps, sub.UeDlAmbrBps, sub.Sqn, sub.CsgIdPrsent, sub.CsgId, sub.DefaultApnName, sub.NetworkID, sub.PackageId)

		mock.ExpectQuery(`^SELECT.*asrs.*`).
			WithArgs(Imsi).
			WillReturnRows(hrow)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewAsrRecordRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		asr, err := r.GetByImsi(Imsi)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

		if assert.NotNil(t, asr) {
			assert.EqualValues(t, asr.Imsi, Imsi)
		}

	})

}

func TestAsrRecordRepo_Delete(t *testing.T) {

	t.Run("DeleteByICCID", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta("UPDATE")).WithArgs(sqlmock.AnyArg(), Iccid).
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

		r := int_db.NewAsrRecordRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.DeleteByIccid(Iccid)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})

	t.Run("DeleteByImsi", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta("UPDATE")).WithArgs(sqlmock.AnyArg(), Imsi).
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

		r := int_db.NewAsrRecordRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Delete(Imsi)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})

}

func TestAsrRecordRepo_UpdateTai(t *testing.T) {
	t.Run("UpdateTai", func(t *testing.T) {

		sub.ID = 1

		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		hrow := sqlmock.NewRows([]string{"ID", "iccid", "imsi", "op", "amf", "key", "algo_type", "ue_dl_ambr_bps", "ue_ul_ambr_bps", "sqn", "csg_id_prsent", "csg_id", "default_apn_name", "network_id", "package_id"}).
			AddRow(sub.ID, sub.Iccid, sub.Imsi, sub.Op, sub.Amf, sub.Key, sub.AlgoType, sub.UeDlAmbrBps, sub.UeDlAmbrBps, sub.Sqn, sub.CsgIdPrsent, sub.CsgId, sub.DefaultApnName, sub.NetworkID, sub.PackageId)

		trow := sqlmock.NewRows([]string{"asr_id", "plmn_id", "tac", "device_updated_at"})

		mock.ExpectBegin()
		mock.ExpectQuery(`^SELECT.*asrs.*`).
			WithArgs(Imsi).
			WillReturnRows(hrow)

		mock.ExpectQuery(`^SELECT.*tais.*`).
			WithArgs(sub.ID, sqlmock.AnyArg()).
			WillReturnRows(trow)

		mock.ExpectExec(regexp.QuoteMeta("UPDATE")).
			WithArgs(sqlmock.AnyArg(), sub.ID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sub.ID, tai.PlmnId, tai.Tac, sqlmock.AnyArg()).
			WillReturnRows(trow)

		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewAsrRecordRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.UpdateTai(Imsi, tai)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})
}
