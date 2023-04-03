package server

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/metrics/exporter/pkg"
	"github.com/ukama/ukama/systems/metrics/exporter/pkg/collector"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

//var customLabelsSimUsage = []string{"session", "start", "end"}

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
	case "event.cloud.cdr.sim.usage":
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

func unmarshalEventSimUsage(msg *anypb.Any) (*epb.SimUsage, error) {
	p := &epb.SimUsage{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddOrgRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func handleEventSimUsage(key string, msg *epb.SimUsage, s *ExporterEventServer) error {

	cfgs, err := s.mc.GetConfigForEvent(key)
	if err != nil {
		log.Errorf("Event %s not implemented.", key)
		return err
	}

	/* Iterating over metrics schema for event */
	for _, ms := range cfgs.Schema {

		switch ms.Name {
		case "sim_usage":
			err := AddSimUsage(msg, s, ms)
			if err != nil {
				return err
			}
		case "sim_usage_duration":
			err = AddSimUsageDuration(msg, s, ms)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func AddSimUsage(msg *epb.SimUsage, s *ExporterEventServer, ms pkg.MetricSchema) error {

	/* Check if metric exist */
	m, err := s.mc.GetMetric(ms.Name)
	if err == nil {
		/* Update value */
		return m.SetMetric(float64(msg.BytesUsed), SetUpDynamicLabelsForSim(ms.DynamicLabels, msg))

	} else {

		m, err := collector.SetUpMetric(s.mc, ms)
		if err != nil {
			return err
		}

		err = m.SetMetric(float64(msg.BytesUsed), SetUpDynamicLabelsForSim(ms.DynamicLabels, msg))
		if err != nil {
			return err
		}

	}
	return nil
}

func AddSimUsageDuration(msg *epb.SimUsage, s *ExporterEventServer, ms pkg.MetricSchema) error {

	/* Check if metric exist */
	m, err := s.mc.GetMetric(ms.Name)
	if err == nil {
		/* Update value */
		return m.SetMetric(float64(msg.EndTime-msg.StartTime), SetUpDynamicLabelsForSim(ms.DynamicLabels, msg))

	} else {
		m, err := collector.SetUpMetric(s.mc, ms)
		if err != nil {
			return err
		}

		err = m.SetMetric(float64(msg.EndTime-msg.StartTime), SetUpDynamicLabelsForSim(ms.DynamicLabels, msg))
		if err != nil {
			return err
		}

	}
	return nil
}

func SetUpDynamicLabelsForSim(keys []string, msg *epb.SimUsage) prometheus.Labels {
	l := make(prometheus.Labels)
	for _, k := range keys {
		switch k {
		case "sim":
			l[k] = msg.Id
		case "org":
			l[k] = msg.OrgId
		case "network":
			l[k] = msg.NetworkId
		case "subscriber":
			l[k] = msg.SubscriberId
		case "sim_type":
			l[k] = msg.Type
		}
	}

	return l
}
