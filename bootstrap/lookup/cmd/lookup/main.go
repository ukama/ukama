package main

import (
	"ukamaX/bootstrap/lookup/internal"
	"ukamaX/bootstrap/lookup/internal/db"
	"ukamaX/bootstrap/lookup/internal/rest"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukamaX/common/config"
	"github.com/ukama/ukamaX/common/sql"
)

func main() {
	log.Infof("Starting the lookub service")
	initConfig()
	if internal.ServiceConf.DebugMode {
		log.Infof("Service running in debug mode")
	}
	log.Infof("")
	d := initDb()
	runHttpServer(d)
}

func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(internal.ServiceConf.DB, internal.ServiceConf.DebugMode)
	err := d.Init(&db.Org{}, &db.Node{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func initConfig() {
	log.Infof("Initializing config")
	internal.ServiceConf = internal.NewConfig()
	config.LoadConfig(internal.ServiceName, internal.ServiceConf)
}

func runHttpServer(d sql.Db) {
	r := rest.NewRouter(db.NewNodeRepo(d), db.NewOrgRepo(d), internal.ServiceConf.DebugMode)
	r.Run()
}
