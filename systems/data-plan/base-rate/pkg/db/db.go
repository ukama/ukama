package db

import (
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Handler struct {
	*gorm.DB
}

func Init(url string) Handler {
	loggerConf := logger.Config{
		SlowThreshold:             time.Second, // Slow SQL threshold
		LogLevel:                  logger.Warn, // Log level
		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
		Colorful:                  true,        // Disable color
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		loggerConf,
	)
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger:                 newLogger,
		SkipDefaultTransaction: true,
	})

	if err != nil {
		log.Fatalln(err)
		panic("Database is not connected. Make sure you call Init() first")
	}
	logrus.Infof("Connected to rate db")

	// db.AutoMigrate(&models.Rate{})
	// db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&models.Rate{})

	return Handler{db}
}
