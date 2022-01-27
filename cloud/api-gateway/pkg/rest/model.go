package rest

type UserRequest struct {
	Org       string `path:"org" validate:"required"`
	Imsi      string `json:"imsi" validate:"required"`
	FirstName string `json:"firstName,omitempty" validate:"required"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
}

type GetNodeMetricsInput struct {
	Org    string `path:"org" validate:"required"`
	NodeID string `path:"node" validate:"required"`
	Metric string `path:"metric" validate:"required"`
	From   int64  `query:"from" validate:"required"`
	To     int64  `query:"to" validate:"required"`
	Step   uint   `query:"step" default:"3600"` // default 1 hour
}
