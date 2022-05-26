package pkg

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	reg "github.com/ukama/ukama/services/cloud/registry/pb/gen"
	"sync"
	"time"
)

func NewMetricsCollector(reg reg.RegistryServiceClient, timeout time.Duration, requestInterval time.Duration) prometheus.Collector {
	mx := sync.Mutex{}

	c := &OrgCollector{
		mx:              &mx,
		reg:             reg,
		timeout:         timeout,
		requestInterval: requestInterval,
	}
	c.StartMetricsUpdate()
	return c
}

type OrgCollector struct {
	mx              *sync.Mutex
	reg             reg.RegistryServiceClient
	timeout         time.Duration
	requestInterval time.Duration
	metrics         map[string]map[string]int32
}

func (c *OrgCollector) StartMetricsUpdate() {
	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
			defer cancel()

			resp, err := c.reg.List(ctx, &reg.ListRequest{})
			if err != nil {
				logrus.Errorf("Error while getting registry list: %v", err)
			}
			c.mx.Lock()
			for _, o := range resp.Orgs {
				nl := map[string]int32{}
				for _, n := range o.GetNetworks() {
					nl[n.GetName()] = n.GetNumberOfNodes()
				}

				c.metrics[o.Name] = nl
			}
			c.mx.Unlock()

			time.Sleep(c.requestInterval)
		}
	}()
}

func (o *OrgCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- prometheus.NewDesc("org_metrics", "org metrics", []string{"org", "network"}, nil)
}

func (o *OrgCollector) Collect(c chan<- prometheus.Metric) {
	o.mx.Lock()
	defer o.mx.Unlock()

	for org, networks := range o.metrics {
		for network, nodes := range networks {
			c <- prometheus.MustNewConstMetric(
				prometheus.NewDesc("nodes_count", "org metrics", []string{"org", "network"}, nil),
				prometheus.GaugeValue,
				float64(nodes),
				org,
				network,
			)
		}
	}

}
