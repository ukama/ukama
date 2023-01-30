package main

import (
	"os"

	"github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/common/sql"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gopkg.in/yaml.v2"

	"github.com/num30/config"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/cmd/version"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"

	generated "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"

	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/clients/adapters"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/clients/providers"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/interceptor"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/server"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"google.golang.org/grpc"
)

var svcConf = pkg.NewConfig(pkg.ServiceName)

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	/* Log level */
	log.SetLevel(log.TraceLevel)
	log.Infof("Starting %s service", pkg.ServiceName)

	initConfig()

	metrics.StartMetricsServer(svcConf.Metrics)

	simDB := initDb()

	runGrpcServer(simDB)

	log.Infof("Exiting service %s", pkg.ServiceName)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	err := config.NewConfReader(pkg.ServiceName).Read(svcConf)
	if err != nil {
		log.Fatalf("Error reading config file. Error: %v", err)
	} else if svcConf.DebugMode {
		// output config in debug mode
		b, err := yaml.Marshal(svcConf)
		if err != nil {
			logrus.Infof("Config:\n%s", string(b))
		}
	}

	log.Debugf("\nService: %s DB Config: %+v Service: %+v MsgClient Config %+v", pkg.ServiceName, svcConf.DB, svcConf.Service, svcConf.MsgClient)

	pkg.IsDebugMode = svcConf.DebugMode
}

func initDb() sql.Db {
	log.Infof("Initializing Database")

	d := sql.NewDb(svcConf.DB, svcConf.DebugMode)

	err := d.Init(&db.Sim{}, &db.Package{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}

	return d
}

func runGrpcServer(gormDB sql.Db) {
	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		/* used on local machines */
		inst := uuid.NewV4()
		instanceId = inst.String()
	}

	mbClient := mb.NewMsgBusClient(svcConf.MsgClient.Timeout, pkg.SystemName,
		pkg.ServiceName, instanceId, svcConf.Queue.Uri,
		svcConf.Service.Uri, svcConf.MsgClient.Host, svcConf.MsgClient.Exchange,
		svcConf.MsgClient.ListenQueue, svcConf.MsgClient.PublishQueue,
		svcConf.MsgClient.RetryCount,
		svcConf.MsgClient.ListenerRoutes)

	log.Debugf("MessageBus Client is %+v", mbClient)

	simManagerServer := server.NewSimManagerServer(
		db.NewSimRepo(gormDB),
		db.NewPackageRepo(gormDB),
		adapters.NewAgentFactory(svcConf.TestAgentHost, svcConf.Timeout),
		providers.NewPackageClientProvider(svcConf.PackageHost),
		providers.NewSubscriberRegistryClientProvider(svcConf.SubscriberRegistryHost),
		providers.NewSimPoolClientProvider(svcConf.SimPoolHost),
		svcConf.Key,
		mbClient,
	)

	tsInterceptor := interceptor.NewTestSimInterceptor(svcConf.TestAgentHost, svcConf.Timeout)

	grpcServer := ugrpc.NewGrpcServer(*svcConf.Grpc, func(s *grpc.Server) {
		generated.RegisterSimManagerServiceServer(s, simManagerServer)
	})

	grpcServer.ExtraUnaryInterceptors = []grpc.UnaryServerInterceptor{
		tsInterceptor.UnaryServerInterceptor}

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
