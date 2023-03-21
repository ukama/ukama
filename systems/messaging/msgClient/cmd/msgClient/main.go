package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/num30/config"
	uconf "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/messaging/msgClient/cmd/version"
	"github.com/ukama/ukama/systems/messaging/msgClient/internal"
	"github.com/ukama/ukama/systems/messaging/msgClient/internal/db"
	"github.com/ukama/ukama/systems/messaging/msgClient/internal/queue"
	"github.com/ukama/ukama/systems/messaging/msgClient/internal/server"
	"gopkg.in/yaml.v3"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/sql"
	generated "github.com/ukama/ukama/systems/messaging/msgClient/pb/gen"

	"google.golang.org/grpc"
)

var serviceConfig = internal.NewConfig()

func main() {
	ccmd.ProcessVersionArgument("msgClient", os.Args, version.Version)

	/* Log level */
	log.SetLevel(log.TraceLevel)
	log.Infof("Starting the msgClient service")

	initConfig()

	metrics.StartMetricsServer(serviceConfig.Metrics)

	db := initDb()

	runGrpcServer(db)

	log.Infof("Exiting service %s", internal.ServiceName)

}

func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(serviceConfig.DB, internal.IsDebugMode)
	err := d.Init(&db.Service{}, &db.Route{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func initConfig() {
	log.Infof("Initializing config")
	serviceConfig = &internal.Config{
		DB: &uconf.Database{
			DbName: internal.ServiceName,
		},
		Grpc: &uconf.Grpc{
			Port: 9095,
		},
	}

	err := config.NewConfReader(internal.ServiceName).Read(serviceConfig)
	if err != nil {
		log.Fatal("Error reading config ", err)
	} else if internal.IsDebugMode {
		b, err := yaml.Marshal(serviceConfig)
		if err != nil {
			log.Infof("Config:\n%s", string(b))
		}
	}

	log.Debugf("Service: %s Config: %+v", internal.ServiceName, serviceConfig.Grpc)

}

func runGrpcServer(d sql.Db) {

	serviceRepo, routeRepo := db.NewServiceRepo(d), db.NewRouteRepo(d)
	handler := queue.NewMessageBusHandler(serviceRepo, routeRepo, serviceConfig.HeathCheck.AllowedMiss, serviceConfig.HeathCheck.Period)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		srv := server.NewMsgClientServer(serviceRepo, routeRepo, handler, serviceConfig.System)
		generated.RegisterMsgClientServiceServer(s, srv)
	})

	signalHandler(handler, grpcServer)
	log.Infof("Message Bus Handler is %+v", handler)
	err := handler.CreateServiceMsgBusHandler()
	if err != nil {
		log.Fatalf("Failed to start message bus queue listener. Error: %s", err.Error())
	}

	grpcServer.StartServer()
}

func signalHandler(handler *queue.MsgBusHandler, server *ugrpc.UkamaGrpcServer) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		handler.StopQueueListener()
		server.StopServer()
	}()
}
