package db_test

import (
	extsql "database/sql"
	"regexp"
	"testing"
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"
	int_db "github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Iccid = "0123456789012345678912"

var sub = int_db.Asr{
	Model: gorm.Model{
		ID: 1,
	},
	Iccid:                   Iccid,
	Imsi:                    Imsi,
	Op:                      []byte("0123456789012345"),
	Key:                     []byte("0123456789012345"),
	Amf:                     []byte("800"),
	AlgoType:                1,
	UeDlAmbrBps:             2000000,
	UeUlAmbrBps:             2000000,
	Sqn:                     1,
	CsgIdPrsent:             false,
	CsgId:                   0,
	DefaultApnName:          "ukama",
	LastStatusChangeAt:      time.Now(),
	LastStatusChangeReasons: int_db.ACTIVATION,
	AllowedTimeOfService:    7200,
	Policy: int_db.Policy{
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Id:           uuid.NewV4(),
		Burst:        1500,
		TotalData:    1024000,
		ConsumedData: 0,
		Dlbr:         5000,
		Ulbr:         1000,
		StartTime:    1714008143,
		EndTime:      1914008143,
		AsrID:        1,
	},
}

// var sim = client.SimCardInfo{
// 	Iccid:          Iccid,
// 	Imsi:           Imsi,
// 	Op:             []byte("0123456789012345"),
// 	Key:            []byte("0123456789012345"),
// 	Amf:            []byte("800"),
// 	AlgoType:       1,
// 	UeDlAmbrBps:    2000000,
// 	UeUlAmbrBps:    2000000,
// 	Sqn:            1,
// 	CsgIdPrsent:    false,
// 	CsgId:          0,
// 	DefaultApnName: "ukama",
// }

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

		//row := sqlmock.NewRows(12[]string{"iccid", "imsi", "op", "amf", "key", "algo_type", "ue_dl_ambr_bps", "ue_ul_ambr_bps", "sqn", "csg_id_prsent", "csg_id", "default_apn_name", "network_id", "package_id"}).

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sub.Iccid, sub.Imsi, sub.Op, sub.Amf, sub.Key, sub.AlgoType, sub.UeDlAmbrBps, sub.UeUlAmbrBps, sub.Sqn, sub.CsgIdPrsent, sub.CsgId, sub.DefaultApnName, sub.NetworkId, sub.PackageId, sub.LastStatusChangeAt, sub.AllowedTimeOfService, sub.LastStatusChangeReasons).
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

		// Assert12
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
		hrow := sqlmock.NewRows([]string{"ID", "iccid", "imsi", "op", "amf", "key", "algo_type", "ue_dl_ambr_bps", "ue_ul_ambr_bps", "sqn", "csg_id_prsent", "csg_id", "default_apn_name", "network_id", "package_id", "last_status_chang_at", "allowed_time_of_service", "last_status_change_reasons"}).
			AddRow(sub.ID, sub.Iccid, sub.Imsi, sub.Op, sub.Amf, sub.Key, sub.AlgoType, sub.UeDlAmbrBps, sub.UeUlAmbrBps, sub.Sqn, sub.CsgIdPrsent, sub.CsgId, sub.DefaultApnName, sub.NetworkId, sub.PackageId, sub.LastStatusChangeAt, sub.AllowedTimeOfService, sub.LastStatusChangeReasons)

		mock.ExpectBegin()
		mock.ExpectQuery(`^SELECT.*asrs.*`).
			WithArgs(sub.Imsi).
			WillReturnRows(hrow)

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
			WithArgs(sqlmock.AnyArg(), sub.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sub.Policy.Id, sub.Policy.Burst, sub.Policy.TotalData, sub.Policy.ConsumedData, sub.Policy.Dlbr, sub.Policy.Ulbr, sub.Policy.StartTime, sub.Policy.EndTime, sub.Model.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
			WithArgs(sub.ID, sqlmock.AnyArg(), sub.Iccid, sub.Imsi, sub.Op, sub.Amf, sub.Key, sub.AlgoType, sub.UeDlAmbrBps, sub.UeUlAmbrBps, sub.Sqn, sub.DefaultApnName, PackageId.String(), sqlmock.AnyArg(), sub.AllowedTimeOfService, int_db.PACKAGE_UPDATE, sub.Imsi).
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
		err = r.UpdatePackage(Imsi, PackageId, &sub.Policy)

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
			AddRow(sub.ID, sub.Iccid, sub.Imsi, sub.Op, sub.Amf, sub.Key, sub.AlgoType, sub.UeDlAmbrBps, sub.UeDlAmbrBps, sub.Sqn, sub.CsgIdPrsent, sub.CsgId, sub.DefaultApnName, sub.NetworkId, sub.PackageId)

		trow := sqlmock.NewRows([]string{"asr_id", "plmn_id", "tac", "device_updated_at"})

		prow := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "id", "burst", "total_data", "consumed_data", "dlbr", "ulbr", "start_time", "end_time", "asr_id"}).
			AddRow(sub.Policy.CreatedAt, sub.Policy.UpdatedAt, sub.Policy.DeletedAt, sub.Policy.Id, sub.Policy.Burst, sub.Policy.TotalData, sub.Policy.ConsumedData, sub.Policy.Dlbr, sub.Policy.Ulbr, sub.Policy.StartTime, sub.Policy.EndTime, sub.Model.ID)

		mock.ExpectQuery(`^SELECT.*asrs.*`).
			WithArgs(sub.ID).
			WillReturnRows(hrow)

		mock.ExpectQuery(`^SELECT.*policies.*`).
			WithArgs(sub.ID).
			WillReturnRows(prow)

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

		hrow := sqlmock.NewRows([]string{"ID", "iccid", "imsi", "op", "amf", "key", "algo_type", "ue_dl_ambr_bps", "ue_ul_ambr_bps", "sqn", "csg_id_prsent", "csg_id", "default_apn_name", "network_id", "package_id"}).
			AddRow(sub.ID, sub.Iccid, sub.Imsi, sub.Op, sub.Amf, sub.Key, sub.AlgoType, sub.UeDlAmbrBps, sub.UeDlAmbrBps, sub.Sqn, sub.CsgIdPrsent, sub.CsgId, sub.DefaultApnName, sub.NetworkId, sub.PackageId)

		trow := sqlmock.NewRows([]string{"asr_id", "plmn_id", "tac", "device_updated_at"})

		prow := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "id", "burst", "total_data", "consumed_data", "dlbr", "ulbr", "start_time", "end_time", "asr_id"}).
			AddRow(sub.Policy.CreatedAt, sub.Policy.UpdatedAt, sub.Policy.DeletedAt, sub.Policy.Id, sub.Policy.Burst, sub.Policy.TotalData, sub.Policy.ConsumedData, sub.Policy.Dlbr, sub.Policy.Ulbr, sub.Policy.StartTime, sub.Policy.EndTime, sub.Model.ID)

		mock.ExpectQuery(`^SELECT.*asrs.*`).
			WithArgs(sub.Iccid).
			WillReturnRows(hrow)

		mock.ExpectQuery(`^SELECT.*policies.*`).
			WithArgs(sub.ID).
			WillReturnRows(prow)

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

		hrow := sqlmock.NewRows([]string{"ID", "iccid", "imsi", "op", "amf", "key", "algo_type", "ue_dl_ambr_bps", "ue_ul_ambr_bps", "sqn", "csg_id_prsent", "csg_id", "default_apn_name", "network_id", "package_id"}).
			AddRow(sub.ID, sub.Iccid, sub.Imsi, sub.Op, sub.Amf, sub.Key, sub.AlgoType, sub.UeDlAmbrBps, sub.UeDlAmbrBps, sub.Sqn, sub.CsgIdPrsent, sub.CsgId, sub.DefaultApnName, sub.NetworkId, sub.PackageId)

		trow := sqlmock.NewRows([]string{"asr_id", "plmn_id", "tac", "device_updated_at"})

		prow := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "id", "burst", "total_data", "consumed_data", "dlbr", "ulbr", "start_time", "end_time", "asr_id"}).
			AddRow(sub.Policy.CreatedAt, sub.Policy.UpdatedAt, sub.Policy.DeletedAt, sub.Policy.Id, sub.Policy.Burst, sub.Policy.TotalData, sub.Policy.ConsumedData, sub.Policy.Dlbr, sub.Policy.Ulbr, sub.Policy.StartTime, sub.Policy.EndTime, sub.Model.ID)

		mock.ExpectQuery(`^SELECT.*asrs.*`).
			WithArgs(Imsi).
			WillReturnRows(hrow)
		mock.ExpectQuery(`^SELECT.*policies.*`).
			WithArgs(sub.ID).
			WillReturnRows(prow)

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

func TestAsrRecordRepo_UpdateTai(t *testing.T) {
	t.Run("UpdateTai", func(t *testing.T) {

		sub.ID = 1

		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		hrow := sqlmock.NewRows([]string{"ID", "iccid", "imsi", "op", "amf", "key", "algo_type", "ue_dl_ambr_bps", "ue_ul_ambr_bps", "sqn", "csg_id_prsent", "csg_id", "default_apn_name", "network_id", "package_id"}).
			AddRow(sub.ID, sub.Iccid, sub.Imsi, sub.Op, sub.Amf, sub.Key, sub.AlgoType, sub.UeDlAmbrBps, sub.UeDlAmbrBps, sub.Sqn, sub.CsgIdPrsent, sub.CsgId, sub.DefaultApnName, sub.NetworkId, sub.PackageId)

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
