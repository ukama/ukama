package server

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/common/config"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/metrics/exporter/pb/gen"
	"github.com/ukama/ukama/systems/metrics/exporter/pkg"
	"github.com/ukama/ukama/systems/metrics/exporter/pkg/collector"
	"google.golang.org/protobuf/types/known/anypb"
)

type TestConfig struct {
	MetricConfig []pkg.MetricConfig
	Metrics      *config.Metrics
}

func InitTestConfig() *TestConfig {
	t := &TestConfig{}
	t.MetricConfig = []pkg.MetricConfig{
		{
			Name:    "subscriber_simusage",
			Event:   "event.cloud.simmanager.sim.usage", //"event.cloud.cdr.sim.usage"}
			Type:    "histogram",
			Units:   "bytes",
			Labels:  map[string]string{"name": "usage"},
			Details: "Data Usage of the sim",
			Buckets: []float64{1024, 10240, 102400, 1024000, 10240000, 102400000},
		},
		{
			Name:    "subscriber_simusage_duration",
			Event:   "event.cloud.simmanager.sim.duration", //"event.cloud.cdr.sim.usage"}
			Type:    "histogram",
			Units:   "seconds",
			Labels:  map[string]string{"name": "usage_duration"},
			Details: "Data Usage durations",
			Buckets: []float64{60, 300, 600, 1200, 1800, 2700, 3600, 7200, 18000},
		},
		{
			Name:    "subscriber_simcount",
			Event:   "event.cloud.simmanager.sim.count", //"event.cloud.cdr.sim.usage"}
			Type:    "counter",
			Units:   "",
			Labels:  map[string]string{"name": "simcount"},
			Details: "Counter test",
		},
		{
			Name:    "subscriber_activesim",
			Event:   "event.cloud.simmanager.sim.activecount", //"event.cloud.cdr.sim.usage"}
			Type:    "gauge",
			Units:   "",
			Labels:  map[string]string{"name": "active_simcount"},
			Details: "Gauge test",
		},
		{
			Name:    "subscriber_simusage_request",
			Event:   "event.cloud.simmanager.sim.request", //"event.cloud.cdr.sim.usage"}
			Type:    "summary",
			Units:   "seconds",
			Labels:  map[string]string{"name": "sim_usage"},
			Details: "Summary test",
			Buckets: []float64{60, 300, 600, 1200, 1800, 2700, 3600, 7200, 18000},
		},
		{
			Name:    "unkown",
			Event:   "event.cloud.simmanager.sim.unkown", //"event.cloud.cdr.sim.usage"}
			Type:    "unkown",
			Units:   "",
			Labels:  map[string]string{"name": "unkown"},
			Details: "Data Usage of the sim",
		},
	}

	t.Metrics = &config.Metrics{
		Port: 10251,
	}

	return t
}

func TestEvent_EventNotification(t *testing.T) {
	tC := InitTestConfig()
	mc := collector.NewMetricsCollector(tC.MetricConfig)
	s := NewExporterEventServer(mc)
	simUsage := pb.SimUsage{
		Id:           "b20c61f1-1c5a-4559-bfff-cd00f746697d",
		SubscriberID: "c214f255-0ed6-4aa1-93e7-e333658c7318",
		NetworkID:    "9fd07299-2826-4f8b-aea9-69da56440bec",
		OrgID:        "75ec112a-8745-49f9-ab64-1a37edade794",
		Type:         "test_simple",
		BytesUsed:    uint64(rand.Int63n(4096000)),
		SessionId:    uuid.NewV4().String(),
		StartTime:    time.Now().Unix() - int64(rand.Intn(30000)),
		EndTime:      time.Now().Unix(),
	}

	anyE, err := anypb.New(&simUsage)
	assert.NoError(t, err)

	msg := &epb.Event{
		RoutingKey: tC.MetricConfig[0].Event,
		Msg:        anyE,
	}
	_, err = s.EventNotification(context.TODO(), msg)
	assert.NoError(t, err)

}
