package main

import (
	"errors"
	"os"

	"github.com/num30/config"
	"github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/users/pkg/server"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/registry/users/pkg"

	"github.com/ukama/ukama/systems/registry/users/cmd/version"

	"github.com/ukama/ukama/systems/registry/users/pkg/db"

	provider "github.com/ukama/ukama/systems/registry/users/pkg/providers"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"

	uconf "github.com/ukama/ukama/systems/common/config"

	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/registry/users/pb/gen"
	"google.golang.org/grpc"
)

var svcConf *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	pkg.InstanceId = os.Getenv("POD_NAME")

	initConfig()

	usersDb := initDb()

	metrics.StartMetricsServer(svcConf.Metrics)

	runGrpcServer(usersDb)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	svcConf = &pkg.Config{
		DB: &uconf.Database{
			DbName: pkg.ServiceName,
		},
		Grpc: &uconf.Grpc{
			Port: 9090,
		},
		Metrics: &uconf.Metrics{
			Port: 10250,
		},
	}

	err := config.NewConfReader(pkg.ServiceName).Read(svcConf)
	if err != nil {
		log.Fatalf("Error reading config file. Error: %v", err)
	} else if svcConf.DebugMode {
		// output config in debug mode
		b, err := yaml.Marshal(svcConf)
		if err != nil {
			log.Infof("Config:\n%s", string(b))
		}
	}

	pkg.IsDebugMode = svcConf.DebugMode
}

func initDb() sql.Db {
	log.Infof("Initializing Database")

	d := sql.NewDb(svcConf.DB, svcConf.DebugMode)

	err := d.Init(&db.User{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}

	usersDB := d.GetGormDb()

	if usersDB.Migrator().HasTable(&db.User{}) {
		if err := usersDB.First(&db.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("Iniiialzing users table")
			var ukamaUUID uuid.UUID
			var err error

			if ukamaUUID, err = uuid.FromString(os.Getenv("UKAMA_UUID")); err != nil {
				log.Fatalf("Database initialization failed, need valid UKAMA UUID env var. Error: %v", err)
			}

			usersDB.Create(&db.User{
				Uuid:  ukamaUUID,
				Name:  "Ukama Root",
				Email: "hello@ukama.com",
				Phone: "0000000000",
			})
		}
	}

	return d
}

func runGrpcServer(gormdb sql.Db) {
	userService := server.NewUserService(db.NewUserRepo(gormdb),

		provider.NewOrgClientProvider(svcConf.OrgHost),
	)
	grpcServer := ugrpc.NewGrpcServer(*svcConf.Grpc, func(s *grpc.Server) {
		gen.RegisterUserServiceServer(s, userService)
	})

	grpcServer.StartServer()
}
