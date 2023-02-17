package server

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/metrics/exporter/pb/gen"
	"github.com/ukama/ukama/systems/metrics/exporter/pkg/collector"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var customLabels = []string{""}

type ExporterEventServer struct {
	mc *collector.MetricsCollector
	epb.UnimplementedEventNotificationServiceServer
}

func NewExporterEventServer(m *collector.MetricsCollector) *ExporterEventServer {
	return &ExporterEventServer{
		mc: m,
	}
}

func (s *ExporterEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)

	switch e.RoutingKey {
	case "event.cloud.simmanager.sim.usage":
		msg, err := unmarshalEventSimUsage(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleEventSimUsage(e.RoutingKey, msg, s)
		if err != nil {
			return nil, err
		}
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func unmarshalEventSimUsage(msg *anypb.Any) (*pb.SimUsage, error) {
	p := &pb.SimUsage{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddOrgRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func handleEventSimUsage(key string, msg *pb.SimUsage, s *ExporterEventServer) error {
	n := "usage" + msg.Id

	/* Check if metric exist */
	m, err := s.mc.GetMetric(n)
	if err != nil {
		/* Update value */
		m.SetMetric(m.Type, float64(msg.BytesUsed), prometheus.Labels{})
	} else {
		/* Initialize metric first */
		c, err := s.mc.GetConfigForEvent(key)
		if err != nil {
			return err
		}

		nm := collector.NewMetrics(n, c.Type)

		labels := make(map[string]string)
		labels["org"] = msg.OrgID
		labels["network"] = msg.NetworkID
		labels["subscriber"] = msg.SubscriberID

		nm.MergeLabels(c.Labels, labels)

		nm.InitializeMetric(n, *c, customLabels)

		/* Add a metric */
		err = s.mc.AddMetrics(n, *nm)
		if err != nil {
			return err
		}
	}
	return nil
}
