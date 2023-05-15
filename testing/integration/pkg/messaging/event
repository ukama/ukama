package messaging

import (
	"os"

	"github.com/ukama/ukama/testing/integration/pkg"
	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	egenerated "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/uuid"
)

var serviceConfig = pkg.NewConfig(pkg.ServiceName)
var mbClient mb.MsgBusServiceClient

func MessagingStartEventServer() {

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {

		Srv := NewEventServer()
		egenerated.RegisterEventNotificationServiceServer(s, Srv)
	})

	grpcServer.StartServer()

}

func StartBusListener() error {

	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		/* used on local machines */
		inst := uuid.NewV4()
		instanceId = inst.String()
	}

	mbClient = mb.NewMsgBusClient(serviceConfig.MsgClient.Timeout, pkg.SystemName,
		pkg.ServiceName, instanceId, serviceConfig.Queue.Uri,
		serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
		serviceConfig.MsgClient.RetryCount,
		serviceConfig.MsgClient.ListenerRoutes)

	log.Debugf("MessageBus Client is %+v", mbClient)

	if err := mbClient.Register(); err != nil {
		log.Errorf("Failed to register to Message Client Service. Error %s", err.Error())
		return err
	}

	if err := mbClient.Start(); err != nil {
		log.Errorf("Failed to start to Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
		return err
	}

	return nil
}

func StopMsgBusListner() error {
	if err := mbClient.Start(); err != nil {
		log.Errorf("Failed to start to Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
		return err
	}
	return nil
}
