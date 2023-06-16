package rest

type SendEmailReq struct {
	To      []string          `json:"to"`
	Message string            `json:"message"`
	Subject string            `json:"subject"`
	Body    string            `json:"body"`
	Values  map[string]string `json:"values"`
}

type SendEmailRes struct {
	Message string `json:"message"`
}
