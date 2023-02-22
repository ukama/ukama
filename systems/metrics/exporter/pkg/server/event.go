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

var customLabels = []string{"test"}

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

	err = AddSimUsageDuration(key, msg, s)
	if err != nil {
		return err
	}

	err = AddSimUsageSession(key, msg, s)
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
		l := SetUpLabelsForSimUsage(msg)

		m, err := collector.SetUpMetric(key, s.mc, l, n)
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
		return m.SetMetric(float64(msg.BytesUsed), nil)

	} else {
		l := SetUpLabelsForSimDuration(msg)

		m, err := collector.SetUpMetric(key, s.mc, l, n)
		if err != nil {
			return err
		}

		m.SetMetric(float64(msg.EndTime-msg.StartTime), nil)

	}
	return nil
}

func AddSimUsageSession(key string, msg *pb.SimUsage, s *ExporterEventServer) error {
	n := strings.ReplaceAll("sim_usage_sessions_"+msg.Id, "-", "")

	/* Check if metric exist */
	m, err := s.mc.GetMetric(n)
	if err == nil {
		/* Update value */
		return m.SetMetric(float64(msg.BytesUsed), nil)

	} else {
		l := SetUpLabelsForSimUsageSession(msg)

		m, err := collector.SetUpMetric(key, s.mc, l, n)
		if err != nil {
			return err
		}

		m.SetMetric(0, nil)

	}
	return nil
}

func SetUpLabelsForSimUsage(msg *pb.SimUsage) map[string]string {
	labels := make(map[string]string)
	labels["org"] = msg.OrgID
	labels["network"] = msg.NetworkID
	labels["subscriber"] = msg.SubscriberID
	labels["sim_type"] = msg.Type
	labels["session"] = msg.SessionId
	labels["start"] = strconv.FormatInt(msg.StartTime, 10)
	labels["end"] = strconv.FormatInt(msg.StartTime, 10)
	return labels
}

func SetUpLabelsForSimDuration(msg *pb.SimUsage) map[string]string {
	labels := make(map[string]string)
	labels["org"] = msg.OrgID
	labels["network"] = msg.NetworkID
	labels["subscriber"] = msg.SubscriberID
	labels["sim_type"] = msg.Type
	labels["session"] = msg.SessionId
	labels["start"] = strconv.FormatInt(msg.StartTime, 10)
	labels["end"] = strconv.FormatInt(msg.StartTime, 10)

	return labels
}

func SetUpLabelsForSimUsageSession(msg *pb.SimUsage) map[string]string {
	labels := make(map[string]string)
	labels["org"] = msg.OrgID
	labels["network"] = msg.NetworkID
	labels["subscriber"] = msg.SubscriberID
	labels["sim_type"] = msg.Type
	labels["session"] = msg.SessionId
	labels["start"] = strconv.FormatInt(msg.StartTime, 10)
	labels["end"] = strconv.FormatInt(msg.StartTime, 10)

	return labels
}
