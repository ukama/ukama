package rest



type SendEmailReq struct {
	To      []string `json:"to" validate:"required"`
	TemplateName string `json:"template_name" validate:"required"`
	Values  map[string]interface{}
}


type SendEmailRes struct {
	Message string `json:"message"`
	MailerId string `json:"mailer_id"`
}

type GetEmailByIdReq struct {
	MailerId string `json:"mailer_id" validate:"required" path:"mailer_id" binding:"required"`
}

