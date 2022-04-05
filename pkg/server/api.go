package server

type BootstrapRequest struct {
	Nodeid     string `query:"nodeid" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}

type BootstrapResponse struct {
}
