package main

import (
	"errors"
	"os"

	"github.com/ukama/ukama/systems/common/msgBusServiceClient"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/uuid"
	generated "github.com/ukama/ukama/systems/registry/member/pb/gen"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"

	"github.com/num30/config"
	"github.com/ukama/ukama/systems/registry/member/cmd/version"
	"github.com/ukama/ukama/systems/registry/member/pkg"

	"github.com/ukama/ukama/systems/registry/member/pkg/db"
	"github.com/ukama/ukama/systems/registry/member/pkg/providers"
	"github.com/ukama/ukama/systems/registry/member/pkg/server"

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
	mDb := initDb()
	runGrpcServer(mDb)
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
	err := d.Init(&db.Member{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {

	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		/* used on local machines */
		inst := uuid.NewV4()
		instanceId = inst.String()
	}

	id, err := uuid.FromString(serviceConfig.OrgId)
	if err != nil {
		log.Fatalf("invalid org uuid. Error %s", err.Error())
	}
	p := providers.NewNucleusClientProvider(serviceConfig.OrgRegistryHost, serviceConfig.DebugMode)
	mbClient := msgBusServiceClient.NewMsgBusClient(serviceConfig.MsgClient.Timeout, serviceConfig.OrgName, pkg.SystemName, pkg.ServiceName, instanceId, serviceConfig.Queue.Uri, serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange, serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue, serviceConfig.MsgClient.RetryCount, serviceConfig.MsgClient.ListenerRoutes)
	memberServer := server.NewMemberServer(serviceConfig.OrgName, db.NewMemberRepo(gormdb),
		p, mbClient, serviceConfig.PushGateway, id)

	logrus.Debugf("MessageBus Client is %+v", mbClient)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		generated.RegisterMemberServiceServer(s, memberServer)
	})

	go grpcServer.StartServer()

	go msgBusListener(mbClient)

	_ = memberServer.PushOrgMemberCountMetric(id)

	initMemberDB(gormdb, p)

	waitForExit()
}

func msgBusListener(m mb.MsgBusServiceClient) {

	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service. Error %s", err.Error())
	}

	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
	}
}

func waitForExit() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	go func() {

		sig := <-sigs
		log.Info(sig)
		done <- true
	}()

	log.Debug("awaiting terminate/interrrupt signal")
	<-done
	log.Infof("exiting service %s", pkg.ServiceName)
}

func initMemberDB(d sql.Db, p providers.NucleusClientProvider) {
	mDB := d.GetGormDb()
	if mDB.Migrator().HasTable(&db.Member{}) {
		if err := mDB.First(&db.Member{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("Initializing registry member table for org")

			var OwnerUUID uuid.UUID
			var err error

			if OwnerUUID, err = uuid.FromString(serviceConfig.OwnerId); err != nil {
				log.Fatalf("Database initialization failed, need valid %v environment variable. Error: %v", "OWNERID", err)
			}

			/* TODO: validate the user from user services */
			u, err := p.GetUserById(serviceConfig.OwnerId)
			if err != nil {
				log.Fatalf("Failed to connect to user service for validation of owner %s. Error: %v", serviceConfig.OwnerId, err)
			}

			o, err := p.GetOrgByName(serviceConfig.OrgName)
			if err != nil {
				log.Fatalf("Failed to connect to org service for validation of owner %s. Error: %v", serviceConfig.OrgName, err)
			}

			if u.User.Id != o.Org.Owner {
				log.Fatalf("Failed to validate user %s as owner of org %+v.", serviceConfig.OwnerId, o)
			}

			if u.User.IsDeactivated {
				log.Fatalf("User is %s is in %s state", serviceConfig.OwnerId, "deactivated")
			}

			if o.Org.IsDeactivated {
				log.Fatalf("Org is %s in %s state", serviceConfig.OwnerId, "deactivated")
			}

			member := &db.Member{
				UserId:      OwnerUUID,
				Deactivated: false,
				Role:        db.Owner,
			}

			if err := mDB.Transaction(func(tx *gorm.DB) error {

				err := p.UpdateOrgToUser(o.Org.Id, member.UserId.String())
				if err != nil {
					return err
				}

				if err := tx.Create(member).Error; err != nil {
					return err
				}
				return nil

			}); err != nil {
				log.Fatalf("Database initialization failed, invalid initial state. Error: %v", err)
			}
		}
	}
}
