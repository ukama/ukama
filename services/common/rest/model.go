package rest

type ErrorMessage struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
