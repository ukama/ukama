package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	_ "github.com/lib/pq"
	wrp "github.com/pkg/errors"
	"github.com/ukama/ukamaX/common/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
	"time"
)

type db struct {
	gorm      *gorm.DB
	DebugMode bool
	dbConfig  config.Database
}

type Db interface {
	GetGormDb() *gorm.DB
	Init(model ...interface{}) error
	Connect() error
}

func NewDb(dbConfig config.Database, debugMode bool) Db {
	return &db{
		dbConfig:  dbConfig,
		DebugMode: debugMode,
	}
}

func (d *db) GetGormDb() *gorm.DB {
	return d.gorm
}

func (d *db) initDbConn() error {
	err := d.Connect()
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "3D000" {
			err = d.createDb()
			if err != nil {
				return wrp.Wrap(err, "error creating database")
			}
			return d.initDbConn()
		}
		return wrp.Wrap(err, "error setting default database")
	}
	return nil
}

func (d *db) Connect() error {
	dsn := d.formatDbInfo(d.dbConfig.DbName)
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
		Logger: newLogger,
	})
	d.gorm = db
	return err
}

func (d *db) formatDbInfo(dbName string) string {
	sslMode := "disable"
	if d.dbConfig.SslEnabled {
		sslMode = "enable"
	}

	dsn := fmt.Sprintf("host=%s user=postgres password=%s database=%s port=%d sslmode=%s",
		d.dbConfig.Host, d.dbConfig.Password, dbName, d.dbConfig.Port, sslMode)
	return dsn
}

func (d *db) migrateDb(dst ...interface{}) error {
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
	dbInfo := d.formatDbInfo("postgres")
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	_, err = db.Exec("create database " + d.dbConfig.DbName)
	if err != nil {
		return err
	}
	return nil
}

func IsNotFoundError(err error) bool {
	return err.Error() == "record not found"
}

func IsDuplicateKeyError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value")
}
