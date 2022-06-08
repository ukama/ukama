package pkg

import (
	"context"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	reg "github.com/ukama/ukama/services/cloud/network/pb/gen"
)

func NewMetricsCollector(reg reg.RegistryServiceClient, timeout time.Duration, requestInterval time.Duration) prometheus.Collector {
	mx := sync.RWMutex{}

	c := &OrgCollector{
		mx:              &mx,
		reg:             reg,
		timeout:         timeout,
		requestInterval: requestInterval,
		metrics:         map[string]map[string]map[string]uint32{},
	}
	c.StartMetricsUpdate()
	return c
}

type OrgCollector struct {
	mx              *sync.RWMutex
	reg             reg.RegistryServiceClient
	timeout         time.Duration
	requestInterval time.Duration
	// map[org][network][nodeType] = count
	metrics map[string]map[string]map[string]uint32
}

func (c *OrgCollector) StartMetricsUpdate() {
	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), c.timeout)

			logrus.Infof("Getting data from network")
			resp, err := c.reg.List(ctx, &reg.ListRequest{})
			if err != nil {
				logrus.Errorf("Error while getting network list: %v", err)
				time.Sleep(c.requestInterval)
				continue
			}
			c.mx.Lock()
			for _, o := range resp.Orgs {
				nl := map[string]map[string]uint32{}
				for _, n := range o.GetNetworks() {
					nl[n.GetName()] = n.GetNumberOfNodes()
				}

				c.metrics[o.Name] = nl
			}
			logrus.Infof("Orgs count %d", len(c.metrics))
			c.mx.Unlock()
			cancel()

			logrus.Infof("Data retreival sleeps for %v", c.requestInterval)
			time.Sleep(c.requestInterval)
		}
	}()
}

func (o *OrgCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- prometheus.NewDesc("node_count", "org metrics", []string{"org", "network"}, nil)
}

func (o *OrgCollector) Collect(c chan<- prometheus.Metric) {
	o.mx.RLock()
	defer o.mx.RUnlock()

	for org, networks := range o.metrics {
		for network, nodes := range networks {
			for nodeType, count := range nodes {
				c <- prometheus.MustNewConstMetric(
					prometheus.NewDesc("node_count", "org metrics", []string{"org", "network", "node_type"}, nil),
					prometheus.GaugeValue,
					float64(count),
					org, network, nodeType,
				)
			}
		}
	}

}
