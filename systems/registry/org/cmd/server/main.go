package main

import (
	"errors"
	"os"

	"github.com/google/uuid"
	"github.com/ukama/ukama/systems/common/msgBusServiceClient"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"

	generated "github.com/ukama/ukama/systems/registry/org/pb/gen"
	"gorm.io/gorm"

	"gopkg.in/yaml.v2"

	"github.com/num30/config"
	"github.com/ukama/ukama/systems/registry/org/cmd/version"
	"github.com/ukama/ukama/systems/registry/org/pkg"
	"github.com/ukama/ukama/systems/registry/org/pkg/db"
	"github.com/ukama/ukama/systems/registry/org/pkg/server"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/sql"
	"google.golang.org/grpc"
)

var serviceConfig *pkg.Config
 const defaultOrgName = "ukama"

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()
	orgDb := initDb()
	runGrpcServer(orgDb)
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
	err := d.Init(&db.Org{}, &db.User{}, &db.OrgUser{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}

	orgDB := d.GetGormDb()

	if orgDB.Migrator().HasTable(&db.Org{}) {
		if err := orgDB.First(&db.Org{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Info("Iniiialzing orgs table")
			var ukamaUUID uuid.UUID
			var err error

			if ukamaUUID, err = uuid.Parse(os.Getenv("UKAMA_UUID")); err != nil {
				log.Fatalf("Database initialization failed, need valid UKAMA UUID env var. Error: %v", err)
			}

			org := &db.Org{
				Name:  defaultOrgName,
				Owner: ukamaUUID,
			}

			usr := &db.User{
				Uuid: ukamaUUID,
			}

			if err := orgDB.Transaction(func(tx *gorm.DB) error {
				if err := tx.Create(org).Error; err != nil {
					return err
				}

				if err := tx.Create(usr).Error; err != nil {
					return err
				}

				if err := tx.Create(&db.OrgUser{
					OrgID:  org.ID,
					UserID: usr.ID,
					Uuid:   ukamaUUID,
				}).Error; err != nil {
					return err
				}

				return nil
			}); err != nil {
				log.Fatalf("Database initialization failed, invalid initial state. Error: %v", err)
			}
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
	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		srv := server.NewOrgServer(db.NewOrgRepo(gormdb),db.NewUserRepo(gormdb),mbClient)
		generated.RegisterOrgServiceServer(s, srv)
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
