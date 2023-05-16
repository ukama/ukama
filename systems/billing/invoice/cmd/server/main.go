package main

import (
	_ "embed"
	"os"

	"github.com/num30/config"
	"github.com/ukama/ukama/systems/billing/invoice/cmd/version"
	"github.com/ukama/ukama/systems/billing/invoice/internal"
	"github.com/ukama/ukama/systems/billing/invoice/internal/db"
	"github.com/ukama/ukama/systems/billing/invoice/internal/server"
	"github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
	fs "github.com/ukama/ukama/systems/billing/invoice/internal/pdf/server"
	generated "github.com/ukama/ukama/systems/billing/invoice/pb/gen"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
)

var serviceConfig = internal.NewConfig(internal.ServiceName)

func main() {
	ccmd.ProcessVersionArgument(internal.ServiceName, os.Args, version.Version)

	/* Log level */
	log.SetLevel(log.TraceLevel)
	log.Infof("Starting %s service", internal.ServiceName)

	initConfig()

	metrics.StartMetricsServer(serviceConfig.Metrics)

	invoiceDB := initDb()

	runGrpcServer(invoiceDB)

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

func initDb() sql.Db {
	log.Infof("Initializing Database")

	d := sql.NewDb(serviceConfig.DB, serviceConfig.DebugMode)

	err := d.Init(&db.Invoice{})
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

	mbClient := mb.NewMsgBusClient(serviceConfig.MsgClient.Timeout, internal.SystemName,
		internal.ServiceName, instanceId, serviceConfig.Queue.Uri,
		serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
		serviceConfig.MsgClient.RetryCount,
		serviceConfig.MsgClient.ListenerRoutes)

	log.Debugf("MessageBus Client is %+v", mbClient)

	invoiceServer := server.NewInvoiceServer(
		db.NewInvoiceRepo(gormDB),
		mbClient,
	)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		generated.RegisterInvoiceServiceServer(s, invoiceServer)
	})

	pdfServer := fs.NewPDFServer(serviceConfig.PdfHost, serviceConfig.PdfFolder, serviceConfig.PdfPrefix, serviceConfig.PdfPort)

	go pdfServer.Start()

	grpcServer.StartServer()
}
