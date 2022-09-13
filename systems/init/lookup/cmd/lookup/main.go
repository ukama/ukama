package main

import (
	"os"

	sr "github.com/ukama/ukama/systems/common/srvcrouter"
	"github.com/ukama/ukama/systems/init/lookup/cmd/version"
	"github.com/ukama/ukama/systems/init/lookup/internal"
	"github.com/ukama/ukama/systems/init/lookup/internal/db"
	"github.com/ukama/ukama/systems/init/lookup/internal/rest"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/sql"
)

func main() {
	ccmd.ProcessVersionArgument("lookup", os.Args, version.Version)

	/* Log level */
	logrus.SetLevel(logrus.TraceLevel)

	log.Infof("Starting the lookup service")

	initConfig()
	if internal.ServiceConf.DebugMode {
		log.Infof("Service running in debug mode")
	}
	log.Infof("")

	rs := sr.NewServiceRouter(internal.ServiceConf.ServiceRouter)

	d := initDb()

	ext := make(chan error)

	r := rest.NewRouter(internal.ServiceConf, rs, db.NewNodeRepo(d), db.NewOrgRepo(d), internal.ServiceConf.DebugMode)
	go r.Run(ext)

	/* Register service */
	if err := rs.RegisterService(internal.ServiceConf.ApiIf); err != nil {
		logrus.Errorf("Exiting the bootstarp service.")
		return
	}

	perr := <-ext
	if perr != nil {
		panic(perr)
	}

	logrus.Infof("Exiting service %s", internal.ServiceName)

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
