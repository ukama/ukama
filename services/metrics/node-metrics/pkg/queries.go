package pkg

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type Query struct {
	Query string `json:"query"`
}

func (q Query) getQuery(filter *Filter) (string, error) {
	t, err := template.New("q").Parse(q.Query)
	if err != nil {
		return "", err
	}
	s := bytes.NewBufferString("")
	err = t.Execute(s, struct {
		Filter string
	}{
		Filter: filter.GetFilter(),
	})
	if err != nil {
		return "", err
	}

	return s.String(), nil
}

type Filter struct {
	nodeId  string
	org     string
	network string
}

func NewFilter() *Filter {
	return &Filter{}
}

func (f *Filter) WithNodeId(nodeId string) *Filter {
	f.nodeId = nodeId
	return f
}

func (f *Filter) WithOrg(org string) *Filter {
	f.org = org
	return f
}

func (f *Filter) HasNetwork() bool {
	return f.network != ""
}

func (f *Filter) WithNetwork(org string, network string) *Filter {
	f.org = org
	f.network = network
	return f
}

// GetFilter returns a prometheus filter
func (f *Filter) GetFilter() string {
	var filter []string
	if f.nodeId != "" {
		filter = append(filter, fmt.Sprintf("nodeid='%s'", f.nodeId))
	}
	if f.org != "" {
		filter = append(filter, fmt.Sprintf("org='%s'", f.org))
	}
	if f.network != "" {
		filter = append(filter, fmt.Sprintf("network='%s'", f.network))
	}
	return strings.Join(filter, ",")
}

func getExcludeStatements(labels ...string) string {
	el := []string{"job", "instance", "receive", "tenant_id"}
	el = append(el, labels...)
	return fmt.Sprintf("without (%s)", strings.Join(el, ","))
}

type Metric struct {
	NeedRate bool   `json:"needRate"`
	Metric   string `json:"metric"`
	// Range vector duration used in Rate func https://prometheus.io/docs/prometheus/latest/querying/basics/#time-durations
	// if NeedRate is false then this field is ignored
	// Example: 1d or 5h, or 30s
	RateInterval string `json:"rateInterval"`

	// Aggregate function used to aggregate metrics that involve multiple nodes
	AggregateFunc string `json:"aggregateFunc"`
}

func (m Metric) getQuery(metricFilter *Filter, defaultRateInterval string, aggregateFunc string) string {
	rateInterval := m.RateInterval
	if m.NeedRate && len(rateInterval) == 0 {
		rateInterval = defaultRateInterval
	}

	if m.NeedRate {
		return fmt.Sprintf("%s(rate(%s {%s}[%s])) %s", aggregateFunc, m.Metric,
			metricFilter.GetFilter(), rateInterval, getExcludeStatements())
	}

	return fmt.Sprintf("%s(%s {%s}) %s", aggregateFunc, m.Metric, metricFilter.GetFilter(), getExcludeStatements())
}

func (m Metric) getAggregateQuery(filter *Filter) string {
	exludSt := getExcludeStatements("nodeid")

	// org only filter
	if !filter.HasNetwork() {
		exludSt = getExcludeStatements("nodeid", "network")
	}
	af := m.AggregateFunc
	if len(af) == 0 {
		af = "sum"
	}

	return fmt.Sprintf("%s(%s {%s}) %s", af, m.Metric, filter.GetFilter(), exludSt)
}
