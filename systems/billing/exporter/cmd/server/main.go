package main

import (
	"os"

	"github.com/ukama/ukama/systems/common/metrics"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gopkg.in/yaml.v2"

	"github.com/num30/config"
	"github.com/ukama/ukama/systems/billing/exporter/cmd/version"
	"github.com/ukama/ukama/systems/billing/exporter/pkg"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"

	egenerated "github.com/ukama/ukama/systems/common/pb/gen/events"

	"github.com/ukama/ukama/systems/billing/exporter/pkg/server"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"google.golang.org/grpc"
)

var serviceConfig = pkg.NewConfig(pkg.ServiceName)

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	/* Log level */
	log.SetLevel(log.TraceLevel)
	log.Infof("Starting %s service", pkg.ServiceName)

	initConfig()

	metrics.StartMetricsServer(serviceConfig.Metrics)

	// simDB := initDb()

	// runGrpcServer(simDB)
	runGrpcServer()

	log.Infof("Exiting service %s", pkg.ServiceName)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	err := config.NewConfReader(pkg.ServiceName).Read(serviceConfig)
	if err != nil {
		log.Fatalf("Error reading config file. Error: %v", err)
	} else if serviceConfig.DebugMode {
		// output config in debug mode
		b, err := yaml.Marshal(serviceConfig)
		if err != nil {
			log.Infof("Config:\n%s", string(b))
		}
	}

	log.Debugf("\nService: %s DB Config: %+v Service: %+v MsgClient Config %+v", pkg.ServiceName, serviceConfig.DB, serviceConfig.Service, serviceConfig.MsgClient)

	pkg.IsDebugMode = serviceConfig.DebugMode
}

// func initDb() sql.Db {
// log.Infof("Initializing Database")

// d := sql.NewDb(serviceConfig.DB, serviceConfig.DebugMode)

// err := d.Init(&db.Sim{}, &db.Package{})
// if err != nil {
// log.Fatalf("Database initialization failed. Error: %v", err)
// }

// return d
// }

func runGrpcServer() {
	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		/* used on local machines */
		inst := uuid.NewV4()
		instanceId = inst.String()
	}

	mbClient := mb.NewMsgBusClient(serviceConfig.MsgClient.Timeout, pkg.SystemName,
		pkg.ServiceName, instanceId, serviceConfig.Queue.Uri,
		serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
		serviceConfig.MsgClient.RetryCount,
		serviceConfig.MsgClient.ListenerRoutes)

	log.Debugf("MessageBus Client is %+v", mbClient)

	// pckgClient, err := providers.NewPackageClient(serviceConfig.DataPlan, pkg.IsDebugMode)
	// if err != nil {
	// log.Fatalf("Failed to connect to Data Plan API Gateway service for retriving packages %s. Error: %v",
	// serviceConfig.DataPlan, err)
	// }

	simManagerEventServer := server.NewBillingExporterEventServer(serviceConfig.LagoHost, serviceConfig.LagoAPIKey, serviceConfig.LagoPort)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		egenerated.RegisterEventNotificationServiceServer(s, simManagerEventServer)
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
