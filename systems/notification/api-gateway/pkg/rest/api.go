package rest

type SendEmailReq struct {
	To      []string          `json:"to" validate:"required"`
	Message string            `json:"message" validate:"required"`
	Subject string            `json:"subject" validate:"required"`
	Body    string            `json:"body" validate:"required"`
	Values  map[string]string `json:"values"`
}

type SendEmailRes struct {
	Message string `json:"message"`
}
