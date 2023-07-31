package rest

type SendEmailReq struct {
	To      []string `json:"to" validate:"required"`
	TemplateName string `json:"template_name" validate:"required"`
	Values  map[string]interface{}
}


type GetEmailByIdReq struct {
	MailerId string `json:"mailer_id" validate:"required" path:"mailer_id" binding:"required"`
}

