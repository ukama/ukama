package rest

type GetNodeMetricsInput struct {
	FilterBase
	NodeID string `path:"node" validate:"required"`
}

type GetOrgMetricsInput struct {
	Metric string `path:"metric" validate:"required"`
	Org    string `path:"org" validate:"required"`
}

type GetNetworkMetricsInput struct {
	Metric  string `path:"metric" validate:"required"`
	Network string `path:"network" validate:"required"`
}

type FilterBase struct {
	Metric string `path:"metric" validate:"required"`
	From   int64  `query:"from" validate:"required"`
	To     int64  `query:"to" default:"0"`      // can be omitted, if omitted the Now() is used
	Step   uint   `query:"step" default:"3600"` // default 1 hour
}

type GetSubscriberMetricsInput struct {
	FilterBase
	Metric     string `path:"metric" validate:"required"`
	Org        string `path:"org" validate:"required"`
	Network    string `path:"network" validate:"required"`
	Subscriber string `path:"subscriber" validate:"required"`
}

type GetMetricsInput struct {
	Metric string `path:"metric" validate:"required"`
}

type GetSimMetricsInput struct {
	FilterBase
	Metric     string `path:"metric" validate:"required"`
	Org        string `path:"org" validate:"required"`
	Network    string `path:"network" validate:"required"`
	Subscriber string `path:"subscriber" validate:"required"`
	Sim        string `path:"sim" validate:"required"`
}

type MetricIntput struct {
	Metric string `query:"metric" validate:"required"`
}
type DummyParameters struct {
}
