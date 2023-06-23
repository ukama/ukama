package main

import (
	"os"

	"github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/node/pkg/providers"
	"github.com/ukama/ukama/systems/registry/node/pkg/server"

	"github.com/num30/config"
	"gopkg.in/yaml.v3"

	"github.com/ukama/ukama/systems/registry/node/pkg"

	"github.com/ukama/ukama/systems/registry/node/cmd/version"

	"github.com/ukama/ukama/systems/registry/node/pkg/db"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	egenerated "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/sql"
	generated "github.com/ukama/ukama/systems/registry/node/pb/gen"
	"google.golang.org/grpc"
)

var serviceConfig = pkg.NewConfig(pkg.ServiceName)

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	initConfig()
	metrics.StartMetricsServer(serviceConfig.Metrics)

	nodeDb := initDb()
	runGrpcServer(nodeDb)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	err := config.NewConfReader(pkg.ServiceName).Read(serviceConfig)
	if err != nil {
		log.Fatal("Error reading config ", err)
	} else if serviceConfig.DebugMode {
		b, err := yaml.Marshal(serviceConfig)
		if err != nil {
			logrus.Infof("Config:\n%s", string(b))
		}
	}

	pkg.IsDebugMode = serviceConfig.DebugMode
}

func initDb() sql.Db {
	log.Infof("Initializing Database")

	d := sql.NewDb(serviceConfig.DB, serviceConfig.DebugMode)

	err := d.Init(&db.Node{})
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

	orgId, err := uuid.FromString(serviceConfig.Org)
	if err != nil {
		log.Fatalf("Invalid organization identifier %s. Error %s", serviceConfig.Org, err)
	}

	mbClient := mb.NewMsgBusClient(serviceConfig.MsgClient.Timeout, pkg.SystemName,
		pkg.ServiceName, instanceId, serviceConfig.Queue.Uri,
		serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
		serviceConfig.MsgClient.RetryCount,
		serviceConfig.MsgClient.ListenerRoutes)

	log.Debugf("MessageBus Client is %+v", mbClient)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		srv := server.NewNodeServer(db.NewNodeRepo(gormdb), db.NewSiteRepo(gormdb), db.NewNodeStatusRepo(gormdb),
			serviceConfig.PushGateway, mbClient,
			providers.NewOrgClientProvider(serviceConfig.OrgHost),
			providers.NewNetworkClientProvider(serviceConfig.NetworkHost),
			orgId)

		nSrv := server.NewNodeEventServer(srv)
		generated.RegisterNodeServiceServer(s, srv)
		egenerated.RegisterEventNotificationServiceServer(s, nSrv)
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
