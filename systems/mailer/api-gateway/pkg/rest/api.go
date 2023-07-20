package rest

type SendEmailReq struct {
	To      string `json:"to"`
	Message string `json:"message"`
	Subject string `json:"subject"`
}

type SendEmailRes struct {
	Message string `json:"message"`
}
