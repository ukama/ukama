package main

import (
	"errors"
	"os"

	"github.com/num30/config"
	"github.com/ukama/ukama/systems/common/msgBusServiceClient"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/users/pkg/server"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/registry/users/pkg"

	"github.com/ukama/ukama/systems/registry/users/cmd/version"

	"github.com/ukama/ukama/systems/registry/users/pkg/db"

	provider "github.com/ukama/ukama/systems/registry/users/pkg/providers"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"

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

	runGrpcServer(usersDb)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	svcConf = pkg.NewConfig(pkg.ServiceName)
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

	initUsersDB(usersDB)

	return d
}

func runGrpcServer(gormdb sql.Db) {
	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		/* used on local machines */
		inst := uuid.NewV4()
		instanceId = inst.String()
	}

	mbClient := msgBusServiceClient.NewMsgBusClient(svcConf.MsgClient.Timeout, pkg.SystemName, pkg.ServiceName, instanceId, svcConf.Queue.Uri, svcConf.Service.Uri, svcConf.MsgClient.Host, svcConf.MsgClient.Exchange, svcConf.MsgClient.ListenQueue, svcConf.MsgClient.PublishQueue, svcConf.MsgClient.RetryCount, svcConf.MsgClient.ListenerRoutes)

	logrus.Debugf("MessageBus Client is %+v", mbClient)
	userService := server.NewUserService(db.NewUserRepo(gormdb),

		provider.NewOrgClientProvider(svcConf.Org), mbClient,
	)
	grpcServer := ugrpc.NewGrpcServer(*svcConf.Grpc, func(s *grpc.Server) {
		gen.RegisterUserServiceServer(s, userService)
	})

	go msgBusListener(mbClient)
	grpcServer.StartServer()

}

func initUsersDB(usersDB *gorm.DB) {
	if usersDB.Migrator().HasTable(&db.User{}) {
		if err := usersDB.First(&db.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("Iniiialzing users table")
			var ownerUUID uuid.UUID
			var err error

			if ownerUUID, err = uuid.FromString(svcConf.OrgOWnerUUID); err != nil {
				log.Fatalf("Database initialization failed, need valid %s envronment variable. Error: %v", "ORGOWNERUUID", err)
			}

			usersDB.Create(&db.User{
				Uuid:  ownerUUID,
				Name:  svcConf.OrgOWnerName,
				Email: svcConf.OrgOWnerEmail,
				Phone: svcConf.OrgOWnerPhone,
			})
		}
	}
}
func msgBusListener(m mb.MsgBusServiceClient) {

	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service. Error %s", err.Error())
	}

	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
	}
}
