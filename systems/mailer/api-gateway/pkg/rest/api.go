package rest

type SendEmailReq struct {
	To           string         `json:"to"`
	TemplateName string         `json:"templateName"`
	Values       map[string]any `json:"values"`
}

type SendEmailRes struct {
	Success bool `json:"success"`
}
