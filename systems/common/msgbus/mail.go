package msgbus

type MailMessage struct {
	To           string `json:"to"`
	TemplateName string `json:"templateName"`
	Values       map[string]any `json:"values"`
}
