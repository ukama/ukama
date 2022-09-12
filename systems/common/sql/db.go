package sql

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/jackc/pgconn"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	wrp "github.com/ukama/ukama/systems/common/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const PGERROR_CODE_UNIQUE_VIOLATION = "23505"

type db struct {
	gorm      *gorm.DB
	DebugMode bool
	dbConfig  DbConfig
}

// would be better to seaprate migration logic from actuall ORM
type Db interface {
	GetGormDb() *gorm.DB
	Init(model ...interface{}) error
	Connect() error
	ExecuteInTransaction(dbOperation func(tx *gorm.DB) *gorm.DB, nestedFuncs ...func() error) (err error)
	// version 2 of execute in transaction to pass transaction object to nested functions
	ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB, nestedFuncs ...func(tx *gorm.DB) error) (err error)
}

type DbConfig interface {
	GetConnString() string
	ChangeDbName(name string) DbConfig
	GetDbName() string
}

func NewDb(dbConfig DbConfig, debugMode bool) Db {
	return &db{
		dbConfig:  dbConfig,
		DebugMode: debugMode,
	}
}

func NewDbFromGorm(gormDb *gorm.DB, debugMode bool) Db {
	return &db{
		DebugMode: debugMode,
		gorm:      gormDb,
	}
}

func (d *db) GetGormDb() *gorm.DB {
	if d.gorm == nil {
		panic("Database is not connected. Make sure you call Connect() first")
	}
	return d.gorm
}

func (d *db) initDbConn() error {
	err := d.Connect()
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "3D000" {
			logrus.Info("Database does not exist")
			err = d.createDb()
			if err != nil {
				return wrp.Wrap(err, "error creating database")
			}
			logrus.Info("Connecting to newly created database")
			return d.Connect()
		}
		return wrp.Wrap(err, "error setting default database")
	}
	return nil
}

func (d *db) Connect() error {
	dsn := d.dbConfig.GetConnString()
	loggerConf := logger.Config{
		SlowThreshold:             time.Second, // Slow SQL threshold
		LogLevel:                  logger.Warn, // Log level
		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
		Colorful:                  true,        // Disable color
	}

	if d.DebugMode {
		loggerConf.SlowThreshold = 0
		loggerConf.LogLevel = logger.Info
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		loggerConf,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 newLogger,
		SkipDefaultTransaction: true,
	})
	d.gorm = db
	return err
}

func (d *db) migrateDb(dst ...interface{}) error {
	logrus.Info("Migrating DB")
	err := d.GetGormDb().AutoMigrate(dst...)
	if err != nil {
		return err
	}
	return nil
}

func (d *db) CloseDb() {
	pgSql, err := d.GetGormDb().DB()
	if err != nil {
		panic(err)
	}

	err = pgSql.Close()
	if err != nil {
		panic(err)
	}
}

func (d *db) Init(model ...interface{}) error {
	err := d.initDbConn()
	if err != nil {
		return err
	}

	err = d.migrateDb(model...)
	if err != nil {
		return err
	}

	return nil
}

func (d *db) createDb() error {

	dbInfo := d.dbConfig.ChangeDbName("postgres").GetConnString()
	logrus.Info("Creating database ", d.dbConfig.GetDbName())
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	_, err = db.Exec("create database " + d.dbConfig.GetDbName())
	if err != nil {
		return err
	}
	return nil
}

func IsNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func IsDuplicateKeyError(err error) bool {
	var pge *pgconn.PgError
	if errors.As(err, &pge) {
		return pge.Code == PGERROR_CODE_UNIQUE_VIOLATION
	}
	return false
}

// ExecuteInTransaction executes dbOperation in transaction with all nested functions
// if any of nested function returns error then transaction is rolled back
func (d *db) ExecuteInTransaction(dbOperation func(tx *gorm.DB) *gorm.DB, nestedFuncs ...func() error) error {
	return d.gorm.Transaction(func(tx *gorm.DB) error {
		d := dbOperation(tx)

		if d.Error != nil {
			return d.Error
		}

		if len(nestedFuncs) > 0 {
			for _, n := range nestedFuncs {
				if n != nil {
					nestErr := n()
					if nestErr != nil {
						return nestErr
					}
				}
			}
		}

		return nil
	})
}

// ExecuteInTransaction executes dbOperation in transaction with all nested functions
// if any of nested function returns error then transaction is rolled back
// all nested functions receive transaction as parameter
func (d *db) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB, nestedFuncs ...func(tx *gorm.DB) error) (err error) {
	return d.gorm.Transaction(func(tx *gorm.DB) error {
		d := dbOperation(tx)

		if d.Error != nil {
			return d.Error
		}

		if len(nestedFuncs) > 0 {
			for _, n := range nestedFuncs {
				if n != nil {
					nestErr := n(tx)
					if nestErr != nil {
						return nestErr
					}
				}
			}
		}

		return nil
	})
}
