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
	err := AddSimUsage(key, msg, s)
	if err != nil {
		return err
	}

	err = AddSimUsageDuration("event.cloud.simmanager.sim.duration", msg, s)
	if err != nil {
		return err
	}

	return nil
}

func AddSimUsage(key string, msg *pb.SimUsage, s *ExporterEventServer) error {

	cfg, err := s.mc.GetConfigForEvent(key)
	if err != nil {
		log.Errorf("Event %s not implemented.", key)
		return err
	}

	n := cfg.Name

	/* Check if metric exist */
	m, err := s.mc.GetMetric(n)
	if err == nil {
		/* Update value */
		return m.SetMetric(float64(msg.BytesUsed), SetUpDynamicLabelsForSim(cfg.DynamicLabels, msg))

	} else {
		l := cfg.Labels

		m, err := collector.SetUpMetric(key, s.mc, l, n, cfg.DynamicLabels)
		if err != nil {
			return err
		}

		err = m.SetMetric(float64(msg.BytesUsed), SetUpDynamicLabelsForSim(cfg.DynamicLabels, msg))
		if err != nil {
			return err
		}

	}
	return nil
}

func AddSimUsageDuration(key string, msg *pb.SimUsage, s *ExporterEventServer) error {
	cfg, err := s.mc.GetConfigForEvent(key)
	if err != nil {
		log.Errorf("Event %s not implemented.", key)
		return err
	}

	n := cfg.Name

	/* Check if metric exist */
	m, err := s.mc.GetMetric(n)
	if err == nil {
		/* Update value */
		return m.SetMetric(float64(msg.EndTime-msg.StartTime), SetUpDynamicLabelsForSim(cfg.DynamicLabels, msg))

	} else {
		l := cfg.Labels

		m, err := collector.SetUpMetric(key, s.mc, l, n, cfg.DynamicLabels)
		if err != nil {
			return err
		}

		err = m.SetMetric(float64(msg.EndTime-msg.StartTime), SetUpDynamicLabelsForSim(cfg.DynamicLabels, msg))
		if err != nil {
			return err
		}

	}
	return nil
}

func SetUpDynamicLabelsForSim(keys []string, msg *pb.SimUsage) prometheus.Labels {
	l := make(prometheus.Labels)
	for _, k := range keys {
		switch k {
		case "sim":
			l[k] = msg.Id
		case "org":
			l[k] = msg.OrgID
		case "network":
			l[k] = msg.NetworkID
		case "subscriber":
			l[k] = msg.SubscriberID
		case "sim_type":
			l[k] = msg.Type
		}
	}

	return l
}

/*

sum_over_time(sim_usage_b20c61f11c5a4559bfffcd00f746697d_sum[5m])

- increase(sim_usage_b20c61f11c5a4559bfffcd00f746697d_sum[10m])
{instance="localhost:10251", job="ukama-org", name="usage", network="9fd07299-2826-4f8b-aea9-69da56440bec", org="75ec112a-8745-49f9-ab64-1a37edade794", sim="b20c61f1-1c5a-4559-bfff-cd00f746697d", sim_type="test_simple", subscriber="c214f255-0ed6-4aa1-93e7-e333658c7318", system="dev"}
	194.87179487179486
- sum(sim_usage_b20c61f11c5a4559bfffcd00f746697d_sum)
{} 360

- sum_over_time(sim_usage_b20c61f11c5a4559bfffcd00f746697d_sum[2m])
{instance="localhost:10251", job="ukama-org", name="usage", network="9fd07299-2826-4f8b-aea9-69da56440bec", org="75ec112a-8745-49f9-ab64-1a37edade794", sim="b20c61f1-1c5a-4559-bfff-cd00f746697d", sim_type="test_simple", subscriber="c214f255-0ed6-4aa1-93e7-e333658c7318", system="dev"} 2880

sum_over_time(sim_usage_b20c61f11c5a4559bfffcd00f746697d_sum[15s])
{instance="localhost:10251", job="ukama-org", name="usage", network="9fd07299-2826-4f8b-aea9-69da56440bec", org="75ec112a-8745-49f9-ab64-1a37edade794", sim="b20c61f1-1c5a-4559-bfff-cd00f746697d", sim_type="test_simple", subscriber="c214f255-0ed6-4aa1-93e7-e333658c7318", system="dev"}
	420

*/
