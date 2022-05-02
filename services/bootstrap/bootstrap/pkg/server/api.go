package server

type BootstrapRequest struct {
	Nodeid     string `query:"node" validate:"required"`
	LookingFor string `query:"looking_for" validate:"eq=validation,required"`
}

type BootstrapResponse struct {
}
