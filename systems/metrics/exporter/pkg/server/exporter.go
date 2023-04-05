package server

import (
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"

	pb "github.com/ukama/ukama/systems/metrics/exporter/pb/gen"
	"github.com/ukama/ukama/systems/metrics/exporter/pkg"
)

type ExporterServer struct {
	pb.UnimplementedExporterServiceServer
	baseRoutingKey msgbus.RoutingKeyBuilder
	Org            string
	msgbus         mb.MsgBusServiceClient
}

func NewExporterServer(org string, msgBus mb.MsgBusServiceClient) (*ExporterServer, error) {

	exp := ExporterServer{
		Org:    org,
		msgbus: msgBus,
	}

	if msgBus != nil {
		exp.baseRoutingKey = msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName)
	}

	return &exp, nil
}
