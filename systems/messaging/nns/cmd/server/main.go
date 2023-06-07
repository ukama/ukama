package main

import (
	"os"

	"github.com/num30/config"
	pb "github.com/ukama/ukama/systems/messaging/nns/pb/gen"
	"github.com/ukama/ukama/systems/messaging/nns/pkg/client"
	"github.com/ukama/ukama/systems/messaging/nns/pkg/server"

	"github.com/ukama/ukama/systems/messaging/nns/pkg"

	"github.com/ukama/ukama/systems/messaging/nns/cmd/version"

	dnspb "github.com/coredns/coredns/pb"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/metrics"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	egenerated "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/uuid"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	initConfig()

	nnsClient := pkg.NewNns(serviceConfig)
	nodeOrgMapping := pkg.NewNodeToOrgMap(serviceConfig)

	metrics.StartMetricsServer(serviceConfig.Metrics)
	go func() {
		srv := server.NewHttpServer(serviceConfig.Http, serviceConfig.Grpc, serviceConfig.NodeMetricsPort, nnsClient, nodeOrgMapping)
		srv.RunHttpServer()
	}()
	runGrpcServer(nnsClient, nodeOrgMapping)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	log.Infof("Initializing config")

	err := config.NewConfReader(pkg.ServiceName).Read(serviceConfig)
	if err != nil {
		log.Fatal("Error reading config ", err)
	} else if serviceConfig.DebugMode {
		b, err := yaml.Marshal(serviceConfig)
		if err != nil {
			log.Infof("Config:\n%s", string(b))
		}
	}

	log.Debugf("\nService: %s Service: %+v MsgClient Config %+v", pkg.ServiceName, serviceConfig.Service, serviceConfig.MsgClient)

	pkg.IsDebugMode = serviceConfig.DebugMode
}

func runGrpcServer(nns *pkg.Nns, nodeOrgMapping *pkg.NodeOrgMap) {

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

	registryClient, err := client.NewRegistryClient(serviceConfig.Registry, serviceConfig.DebugMode)
	if err != nil {
		log.Fatalf("Error creating registry client. Error: %v", err)
	}

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		srv := server.NewNnsServer(nns, nodeOrgMapping)
		pb.RegisterNnsServer(s, srv)

		dnspb.RegisterDnsServiceServer(s, server.NewDnsServer(nns, serviceConfig.Dns))

		eSrv := server.NewNnsEventServer(registryClient, serviceConfig.EtcdHost, serviceConfig.Timeout)
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
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
	}
}
