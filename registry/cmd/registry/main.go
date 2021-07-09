package main

import (
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukamaX/common/config"
	"github.com/ukama/ukamaX/common/sql"
	"time"
	"ukamaX/registry/internal"
	"ukamaX/registry/internal/db"
)

// Logging
// Configuration
// Traceability
// Validation

var svcConf *internal.Config

func main() {
	initConfig()

	log.Infof("Starting the registry")
	workWithDb()

	for {
		log.Infof("Just sitting here, doing nothing...")
		time.Sleep(30 * time.Second)
	}
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	svcConf = internal.NewConfig()
	config.LoadConfig("", svcConf)
}

func workWithDb() {
	d := sql.NewDb(svcConf.DB, svcConf.DebugMode)
	err := d.Init(&db.Node{}, &db.Org{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	orgRepo := db.NewOrgRepo(d)
	org := &db.Org{Name: "AbcOrg"}
	err = orgRepo.Add(org)
	if err != nil {
		panic(err)
	}

	repo := db.NewNodeRepo(d)
	u := uuid.NewV1()
	err = repo.Add(&db.Node{
		UUID:        u,
		DeviceIP:    10,
		Certificate: "cert",
		Org:         org,
	})
	if err != nil {
		log.Fatalf("Error adding node. Error: %v", err)
	}
	log.Infof("Record added")

	addedNode, err := repo.Get(u)
	if err != nil {
		log.Fatalf("Error getting record. Error: %v", err)
	}
	if err != nil {
		log.Fatalf("Error adding record. Error: %v", err)
	}

	log.Infof("Node: %v", addedNode)
	log.Infof("Node's org: %s", addedNode.Org.Name)
}
