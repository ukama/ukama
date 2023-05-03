package main

import (
	"os"

	"github.com/ukama/ukama/systems/common/metrics"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gopkg.in/yaml.v2"

	"github.com/num30/config"
	"github.com/ukama/ukama/systems/billing/collector/cmd/version"
	"github.com/ukama/ukama/systems/billing/collector/internal"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"

	egenerated "github.com/ukama/ukama/systems/common/pb/gen/events"

	client "github.com/ukama/ukama/systems/billing/collector/internal/clients"
	"github.com/ukama/ukama/systems/billing/collector/internal/server"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"google.golang.org/grpc"
)

var serviceConfig = internal.NewConfig(internal.ServiceName)

func main() {
	ccmd.ProcessVersionArgument(internal.ServiceName, os.Args, version.Version)

	/* Log level */
	log.SetLevel(log.TraceLevel)
	log.Infof("Starting %s service", internal.ServiceName)

	initConfig()

	metrics.StartMetricsServer(serviceConfig.Metrics)

	runGrpcServer()

	log.Infof("Exiting service %s", internal.ServiceName)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	err := config.NewConfReader(internal.ServiceName).Read(serviceConfig)
	if err != nil {
		log.Fatalf("Error reading config file. Error: %v", err)
	} else if serviceConfig.DebugMode {
		// output config in debug mode
		b, err := yaml.Marshal(serviceConfig)
		if err != nil {
			log.Infof("Config:\n%s", string(b))
		}
	}

	log.Debugf("\nService: %s DB Config: %+v Service: %+v MsgClient Config %+v",
		internal.ServiceName, serviceConfig.DB, serviceConfig.Service, serviceConfig.MsgClient)

	internal.IsDebugMode = serviceConfig.DebugMode
}

func runGrpcServer() {
	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		/* used on local machines */
		inst := uuid.NewV4()
		instanceId = inst.String()
	}

	mbClient := mb.NewMsgBusClient(serviceConfig.MsgClient.Timeout, internal.SystemName,
		internal.ServiceName, instanceId, serviceConfig.Queue.Uri,
		serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
		serviceConfig.MsgClient.RetryCount,
		serviceConfig.MsgClient.ListenerRoutes)

	log.Debugf("MessageBus Client is %+v", mbClient)

	lagoClient := client.NewLagoClient(serviceConfig.LagoAPIKey,
		serviceConfig.LagoHost, serviceConfig.LagoPort)

	eSrv := server.NewBillingCollectorEventServer(lagoClient)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		egenerated.RegisterEventNotificationServiceServer(s, eSrv)
	})

	go msgBusListener(mbClient)

	grpcServer.StartServer()
}

func msgBusListener(m mb.MsgBusServiceClient) {
	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service. Error %s", err.Error())
	}
	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s",
			internal.ServiceName, err.Error())
	}
}
