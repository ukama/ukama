package main

import (
	"errors"
	"os"

	"github.com/ukama/ukama/systems/common/msgBusServiceClient"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"

	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/registry/org/pb/gen"
	"github.com/ukama/ukama/systems/registry/org/pkg/server"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/registry/org/pkg"
	"gopkg.in/yaml.v2"

	"github.com/ukama/ukama/systems/registry/org/cmd/version"

	"github.com/ukama/ukama/systems/registry/org/pkg/db"

	"github.com/num30/config"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	uconf "github.com/ukama/ukama/systems/common/config"

	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/sql"
	"google.golang.org/grpc"
)

var svcConf *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	pkg.InstanceId = os.Getenv("POD_NAME")

	initConfig()

	orgDb := initDb()

	runGrpcServer(orgDb)
}

func initConfig() {
	svcConf = &pkg.Config{
		DB: &uconf.Database{
			DbName: pkg.ServiceName,
		},
		Grpc: &uconf.Grpc{
			Port: 9090,
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

	err := d.Init(&db.Org{}, &db.User{}, &db.OrgUser{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}

	orgDB := d.GetGormDb()

	initOrgDB(orgDB)

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
	
	log.Debugf("MessageBus Client is %+v", mbClient)
	regServer := server.NewOrgServer(db.NewOrgRepo(gormdb),
		db.NewUserRepo(gormdb),
		svcConf.OrgName,mbClient,
	)

	grpcServer := ugrpc.NewGrpcServer(*svcConf.Grpc, func(s *grpc.Server) {
		pb.RegisterOrgServiceServer(s, regServer)
	})
	go msgBusListener(mbClient)

	grpcServer.StartServer()
}

func initOrgDB(orgDB *gorm.DB) {
	if orgDB.Migrator().HasTable(&db.Org{}) {
		if err := orgDB.First(&db.Org{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("Iniiialzing orgs table")

			var OwnerUUID uuid.UUID
			var err error

			if OwnerUUID, err = uuid.FromString(svcConf.OrgOwnerUUID); err != nil {
				log.Fatalf("Database initialization failed, need valid %v environment variable. Error: %v", "ORGOWNERUUID", err)
			}

			org := &db.Org{
				ID:    uuid.NewV4(),
				Name:  svcConf.OrgName,
				Owner: OwnerUUID,
			}

			usr := &db.User{
				Uuid: OwnerUUID,
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
					Uuid:   usr.Uuid,
				}).Error; err != nil {
					return err
				}

				return nil
			}); err != nil {
				log.Fatalf("Database initialization failed, invalid initial state. Error: %v", err)
			}
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