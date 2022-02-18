package server

type GetNodeMetricsInput struct {
	NodeID string `path:"node" validate:"required"`
	Metric string `path:"metric" validate:"required"`
	From   int64  `query:"from" validate:"required"`
	To     int64  `query:"to" validate:"required"`
	Step   uint   `query:"step" default:"3600"` // default 1 hour
}
