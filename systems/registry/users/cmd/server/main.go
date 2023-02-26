package main

import (
	"errors"
	"os"

	"github.com/ukama/ukama/systems/common/msgBusServiceClient"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	generated "github.com/ukama/ukama/systems/registry/users/pb/gen"
	provider "github.com/ukama/ukama/systems/registry/users/pkg/providers"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"

	"github.com/num30/config"
	"github.com/ukama/ukama/systems/registry/users/cmd/version"
	"github.com/ukama/ukama/systems/registry/users/pkg"

	"github.com/google/uuid"
	"github.com/ukama/ukama/systems/registry/users/pkg/db"
	"github.com/ukama/ukama/systems/registry/users/pkg/server"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/sql"
	"google.golang.org/grpc"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()
	usersDb := initDb()
	runGrpcServer(usersDb)
}
func initConfig() {

	serviceConfig = pkg.NewConfig(pkg.ServiceName)
	err := config.NewConfReader(pkg.ServiceName).Read(serviceConfig)
	if err != nil {
		log.Fatal("Error reading config ", err)
	} else if serviceConfig.DebugMode {
		b, err := yaml.Marshal(serviceConfig)
		if err != nil {
			log.Infof("Config:\n%s", string(b))
		}
	}
	pkg.IsDebugMode = serviceConfig.DebugMode
}


func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(serviceConfig.DB, serviceConfig.DebugMode)
	err := d.Init(&db.User{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}

	usersDB := d.GetGormDb()

	if usersDB.Migrator().HasTable(&db.User{}) {
		if err := usersDB.First(&db.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Info("Iniiialzing users table")
			var ukamaUUID uuid.UUID
			var err error

			if ukamaUUID, err = uuid.Parse(os.Getenv("UKAMA_UUID")); err != nil {
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

	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		/* used on local machines */
		inst := uuid.New()
		instanceId = inst.String()
	}

	mbClient := msgBusServiceClient.NewMsgBusClient(serviceConfig.MsgClient.Timeout, pkg.SystemName, pkg.ServiceName, instanceId, serviceConfig.Queue.Uri, serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange, serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue, serviceConfig.MsgClient.RetryCount, serviceConfig.MsgClient.ListenerRoutes)
	
	logrus.Debugf("MessageBus Client is %+v", mbClient)
	userService := server.NewUserService(db.NewUserRepo(gormdb),

		provider.NewOrgClientProvider(serviceConfig.OrgHost),mbClient,
	)
	
	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		generated.RegisterUserServiceServer(s, userService)
	})

	go msgBusListener(mbClient)

	grpcServer.StartServer()
}

func msgBusListener(m mb.MsgBusServiceClient) {

	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service. Error %s", err.Error())
	}

	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
	}
}
