package server

import (
	"context"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/metrics/exporter/pb/gen"
	"github.com/ukama/ukama/systems/metrics/exporter/pkg/collector"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var customLabelsSimUsage = []string{"session", "start", "end"}

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
	n := strings.ReplaceAll("sim_usage_"+msg.Id, "-", "")

	/* Check if metric exist */
	m, err := s.mc.GetMetric(n)
	if err == nil {
		/* Update value */
		return m.SetMetric(float64(msg.BytesUsed), nil)

	} else {
		l := SetUpStaticLabelsForSimUsage(msg)

		m, err := collector.SetUpMetric(key, s.mc, l, n, nil)
		if err != nil {
			return err
		}

		m.SetMetric(float64(msg.BytesUsed), nil)

	}
	return nil
}

func AddSimUsageDuration(key string, msg *pb.SimUsage, s *ExporterEventServer) error {
	n := strings.ReplaceAll("sim_duration_"+msg.Id, "-", "")

	/* Check if metric exist */
	m, err := s.mc.GetMetric(n)
	if err == nil {
		/* Update value */
		return m.SetMetric(float64(msg.EndTime-msg.StartTime), nil)

	} else {
		l := SetUpStaticLabelsForSimUsageDuration(msg)

		m, err := collector.SetUpMetric(key, s.mc, l, n, nil)
		if err != nil {
			return err
		}

		m.SetMetric(float64(msg.EndTime-msg.StartTime), nil)

	}
	return nil
}

func SetUpStaticLabelsForSimUsage(msg *pb.SimUsage) map[string]string {
	labels := make(map[string]string)
	labels["sim"] = msg.Id
	labels["org"] = msg.OrgID
	labels["network"] = msg.NetworkID
	labels["subscriber"] = msg.SubscriberID
	labels["sim_type"] = msg.Type
	return labels
}

func SetUpDynamicLabelsForSimUsage(msg *pb.SimUsage) map[string]string {
	labels := make(map[string]string)
	labels["session"] = msg.SessionId
	labels["start"] = strconv.FormatInt(msg.StartTime, 10)
	labels["end"] = strconv.FormatInt(msg.EndTime, 10)
	return labels
}

func SetUpStaticLabelsForSimUsageDuration(msg *pb.SimUsage) map[string]string {
	labels := make(map[string]string)
	labels["sim"] = msg.Id
	labels["org"] = msg.OrgID
	labels["network"] = msg.NetworkID
	labels["subscriber"] = msg.SubscriberID
	labels["sim_type"] = msg.Type
	return labels
}

func SetUpDynamicLabelsForSimUsageDuration(msg *pb.SimUsage) map[string]string {
	labels := make(map[string]string)
	labels["session"] = msg.SessionId
	labels["start"] = strconv.FormatInt(msg.StartTime, 10)
	labels["end"] = strconv.FormatInt(msg.EndTime, 10)
	return labels
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

	