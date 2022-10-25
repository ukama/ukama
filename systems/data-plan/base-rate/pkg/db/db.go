package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func Init(url string) Handler {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connecte to db")
	// db.AutoMigrate(&models.Rate{})
	// db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&models.Rate{})

	return Handler{db}
}
