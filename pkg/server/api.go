package server

type BootstrapRequest struct {
	nodeid string `nodeid:"nodeid" validate:"required"`
	role   string `role:"role" validate:"required"`
}

type BootstrapResponse struct {
}
