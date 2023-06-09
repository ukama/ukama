package rest

type SendEmailReq struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

type SendEmailRes struct {
	Message string `json:"message"`
}
