package rest

type SendEmailReq struct {
	To      []string          `json:"to" validate:"required"`
	Subject string            `json:"subject" validate:"required"`
	Body    string            `json:"body" validate:"required"`
	Values  map[string]string `json:"values"`
}

type SendEmailRes struct {
	Message string `json:"message"`
}

type GetEmailByIdReq struct {
	MailerId string `json:"mailer_id" validate:"required" path:"mailer_id" binding:"required"`
}

