package rest

import (
	"github.com/ukama/ukama/systems/common/rest"
)

type GetNodeMetricsInput struct {
	rest.BaseRequest
	FilterBase
	NodeID string `path:"node" validate:"required"`
}

type GetOrgMetricsInput struct {
	rest.BaseRequest
	Metric string `path:"metric" validate:"required"`
	Org    string `path:"org" validate:"required"`
}

type GetNetworkMetricsInput struct {
	rest.BaseRequest
	Metric  string `path:"metric" validate:"required"`
	Network string `path:"network" validate:"required"`
}

type FilterBase struct {
	rest.BaseRequest
	Metric string `path:"metric" validate:"required"`
	From   int64  `query:"from" validate:"required"`
	To     int64  `query:"to" default:"0"`      // can be omitted, if omitted the Now() is used
	Step   uint   `query:"step" default:"3600"` // default 1 hour
}

type GetSubscriberMetricsInput struct {
	rest.BaseRequest
	FilterBase
	Metric     string `path:"metric" validate:"required"`
	Org        string `path:"org" validate:"required"`
	Network    string `path:"network" validate:"required"`
	Subscriber string `path:"subscriber" validate:"required"`
}

type GetMetricsRangeInput struct {
	rest.BaseRequest
	FilterBase
	Org        string `query:"org"`
	Network    string `query:"network"`
	Subscriber string `query:"subscriber" `
	Sim        string `query:"sim"`
	User       string `query:"user"`
	Site       string `query:"site"`
	NodeID     string `query:"node"`
}
type GetMetricsInput struct {
	rest.BaseRequest
	Metric string `path:"metric" validate:"required"`
}

type GetSimMetricsInput struct {
	rest.BaseRequest
	FilterBase
	Metric     string `path:"metric" validate:"required"`
	Org        string `path:"org" validate:"required"`
	Network    string `path:"network" validate:"required"`
	Subscriber string `path:"subscriber" validate:"required"`
	Sim        string `path:"sim" validate:"required"`
}

type GetWsMetricIntput struct {
	rest.BaseRequest
	Metric     string `query:"metric" validate:"required"`
	Interval   int    `query:"interval" validate:"required"` //Node/Network/Organization/Site/Sub/User
	Org        string `query:"org"`
	Network    string `query:"network"`
	Subscriber string `query:"subscriber" `
	Sim        string `query:"sim"`
	User       string `query:"user"`
	Site       string `query:"site"`
	NodeID     string `query:"node"`
}

type DummyParameters struct {
}
