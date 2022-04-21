package main

import (
	"os"

	"github.com/ukama/openIoR/services/bootstrap/lookup/cmd/version"
	"github.com/ukama/openIoR/services/bootstrap/lookup/internal"
	"github.com/ukama/openIoR/services/bootstrap/lookup/internal/db"
	"github.com/ukama/openIoR/services/bootstrap/lookup/internal/rest"
	sr "github.com/ukama/openIoR/services/common/srvcrouter"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/openIoR/services/common/cmd"
	"github.com/ukama/openIoR/services/common/config"
	"github.com/ukama/openIoR/services/common/sql"
)

func main() {
	ccmd.ProcessVersionArgument("lookup", os.Args, version.Version)

	log.Infof("Starting the lookup service")
	initConfig()
	if internal.ServiceConf.DebugMode {
		log.Infof("Service running in debug mode")
	}
	log.Infof("")

	rs := sr.NewServiceRouter(internal.ServiceConf.ServiceRouter)

	/* Register service */
	if err := rs.RegisterService(internal.ServiceConf.ApiIf); err != nil {
		logrus.Errorf("Exiting the bootstarp service.")
		//return
	}

	d := initDb()
	runHttpServer(d, rs)
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

func runHttpServer(d sql.Db, rs *sr.ServiceRouter) {
	r := rest.NewRouter(rs, db.NewNodeRepo(d), db.NewOrgRepo(d), internal.ServiceConf.DebugMode)
	r.Run()
}
