package server

type GetNodeMetricsInput struct {
	FilterBase
	NodeID string `path:"node" validate:"required"`
}

type GetLatestMetricInput struct {
	Metric string `path:"metric" validate:"required"`
	NodeID string `path:"node" validate:"required"`
}

type GetOrgMetricsInput struct {
	Metric string `path:"metric" validate:"required"`
	Org    string `path:"org" validate:"required"`
}

type GetNetMetricsInput struct {
	Metric string `path:"metric" validate:"required"`
	Org    string `path:"org" validate:"required"`
	Net    string `path:"net" validate:"required"`
}

type FilterBase struct {
	Metric string `path:"metric" validate:"required"`
	From   int64  `query:"from" validate:"required"`
	To     int64  `query:"to" default:"0"`      // can be omitted, if omitted the Now() is used
	Step   uint   `query:"step" default:"3600"` // default 1 hour
}
