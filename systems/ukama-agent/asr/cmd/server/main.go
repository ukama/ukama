package main

import (
	"os"

	"github.com/gofrs/uuid"
	"github.com/num30/config"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/client"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/server"
	"gopkg.in/yaml.v3"

	pkg "github.com/ukama/ukama/systems/ukama-agent/asr/pkg"

	"github.com/ukama/ukama/systems/ukama-agent/asr/cmd/version"

	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	egen "github.com/ukama/ukama/systems/common/pb/gen/events"

	"github.com/ukama/ukama/systems/common/sql"
	"google.golang.org/grpc"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	initConfig()
	hssDb := initDb()
	runGrpcServer(hssDb)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	log.Infof("Initializing config")
	serviceConfig = pkg.NewConfig(pkg.ServiceName)
	err := config.NewConfReader(pkg.ServiceName).Read(serviceConfig)
	if err != nil {
		log.Fatal("Error reading config ", err)
	} else if pkg.IsDebugMode {
		b, err := yaml.Marshal(serviceConfig)
		if err != nil {
			log.Infof("Config:\n%s", string(b))
		}
	}
	pkg.IsDebugMode = serviceConfig.DebugMode
	log.Infof("Config: %+v", serviceConfig)
}

func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(serviceConfig.DB, true)
	err := d.Init(&db.Asr{}, &db.Guti{}, &db.Tai{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {

	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		/* used on local machines */
		inst, err := uuid.NewV4()
		if err != nil {
			log.Fatalf("Failed to genrate instanceId. Error %s", err.Error())
		}
		instanceId = inst.String()
	}

	mbClient := mb.NewMsgBusClient(serviceConfig.MsgClient.Timeout, pkg.SystemName,
		pkg.ServiceName, instanceId, serviceConfig.Queue.Uri,
		serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
		serviceConfig.MsgClient.RetryCount,
		serviceConfig.MsgClient.ListenerRoutes)

	log.Debugf("MessageBus Client is %+v", mbClient)
	asr := db.NewAsrRecordRepo(gormdb)
	guti := db.NewGutiRepo(gormdb)

	factory, err := client.NewFactoryClient(serviceConfig.FactoryHost, pkg.IsDebugMode)
	if err != nil {
		log.Fatalf("Fcatory Client initilization failed. Error: %v", err)
	}

	network, err := client.NewNetworkClient(serviceConfig.NetworkHost, pkg.IsDebugMode)
	if err != nil {
		log.Fatalf("Network Client initilization failed. Error: %v", err)
	}

	pcrf, err := client.NewPolicyControlClient(serviceConfig.PCRFHost, pkg.IsDebugMode)
	if err != nil {
		log.Fatalf("PCRF Client initialization failed. Error: %v", err)
	}

	// asr service
	asrServer, err := server.NewAsrRecordServer(asr, guti,
		factory, network, pcrf, serviceConfig.Org, mbClient)

	if err != nil {
		log.Fatalf("asr server initialization failed. Error: %v", err)
	}
	nSrv := server.NewAsrEventServer(asr, guti)

	rpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		gen.RegisterAsrRecordServiceServer(s, asrServer)
		egen.RegisterEventNotificationServiceServer(s, nSrv)
	})

	if pkg.IsMessageBus {
		go msgBusListener(mbClient)
	}

	rpcServer.StartServer()

}

func msgBusListener(m mb.MsgBusServiceClient) {

	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service. Error %s", err.Error())
	}

	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
	}
}
