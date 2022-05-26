package pkg

import (
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	pb "github.com/ukama/ukama/services/cloud/registry/pb/gen"
	regMock "github.com/ukama/ukama/services/cloud/registry/pb/gen/mocks"
	"testing"
	"time"
)

func Test_Collect(t *testing.T) {
	regM := &regMock.RegistryServiceClient{}
	regM.On("List", mock.Anything, mock.Anything).Return(&pb.ListResponse{
		Orgs: []*pb.ListResponse_Org{
			&pb.ListResponse_Org{
				Name: "a",
				Networks: []*pb.ListResponse_Network{
					{Name: "n1", NumberOfNodes: 1},
					{Name: "n2", NumberOfNodes: 2},
				},
			},
			&pb.ListResponse_Org{
				Name: "b",
				Networks: []*pb.ListResponse_Network{
					{Name: "n3", NumberOfNodes: 5},
				},
			},
		},
	}, nil)

	coll := NewMetricsCollector(regM, 1*time.Second, 10*time.Microsecond)

	mch := make(chan prometheus.Metric)

	stop := false
	// collecting metrics in separate goroutine
	go func() {
		for {
			coll.Collect(mch)
			if stop {
				return
			}
		}
	}()

	for i := 0; i < 4; i++ {
		actual := <-mch

		dtoM := dto.Metric{}
		err := actual.Write(&dtoM)
		assert.NoError(t, err)

		org, net, nodes := parseMetric(dtoM)

		switch org + net {
		case "a" + "n1":
			assert.Equal(t, 1, nodes)
		case "a" + "n2":
			assert.Equal(t, 2, nodes)

		case "b" + "n3":
			assert.Equal(t, 5, nodes)
		default:
			assert.Fail(t, "unexpected metric")
		}
	}

	stop = true
}

func parseMetric(m dto.Metric) (org string, net string, nodes int) {
	for _, l := range m.GetLabel() {
		if l.GetName() == "network" {
			net = l.GetValue()
		} else if l.GetName() == "org" {
			org = l.GetValue()
		}
	}

	return org, net, int(m.Gauge.GetValue())
}
