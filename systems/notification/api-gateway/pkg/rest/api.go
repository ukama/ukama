package rest

type SendEmailReq struct {
	To      []string          `json:"to" validate:"required"`
	Subject string            `json:"subject" validate:"required"`
	Body    string            `json:"body" validate:"required"`
	Values  map[string]string `json:"values"`
}

type SendEmailRes struct {
	Message  string `json:"message"`
	MailerId string `json:"mailer_id"`
}

type GetEmailByIdReq struct {
	MailerId string `json:"mailer_id" validate:"required" path:"mailer_id" binding:"required"`
}

type AddNotificationReq struct {
	NodeId      string `json:"node_id"`
	Severity    string `json:"severity,omitempty" type:"string"`
	Type        string `json:"notification_type,omitempty" validate:"eq=alert|eq=event"`
	ServiceName string `json:"service_name,omitempty"`
	Time        uint32 `json:"time,omitempty"`
	Description string `json:"description,omitempty"`
	Details     string `json:"details,omitempty"`
}

type GetNotificationReq struct {
	NotificationId string `json:"notification_id" path:"notification_id" validate:"required"`
}

type GetNotificationsReq struct {
	NodeId      string `form:"node_id" json:"node_id" query:"node_id" binding:"required"`
	ServiceName string `form:"service_name" json:"service_name" query:"service_name" binding:"required"`
	Type        string `form:"notification_type" json:"notification_type" query:"notification_type" binding:"required"`
	Count       uint32 `form:"count" json:"count" query:"count" binding:"required"`
	Sort        bool   `form:"sort" json:"sort" query:"sort" binding:"required"`
}

type DelNotificationsReq struct {
	NodeId      string `form:"node_id" json:"node_id" query:"node_id" binding:"required"`
	ServiceName string `form:"service_name" json:"service_name" query:"service_name" binding:"required"`
	Type        string `form:"notification_type" json:"notification_type" query:"notification_type" binding:"required"`
}
