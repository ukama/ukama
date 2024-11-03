package emailTemplate

const (
	EmailTemplateSimAllocation = "sim-allocation"
	EmailTemplateMemberInvite  = "member-invite"
	EmailTemplateOrgInvite     = "org-invite"
)

type EmailTemplateKeys struct {
	TemplateName string
	Keys         []string
}

var EmailTemplateConfig = map[string]EmailTemplateKeys{
	EmailTemplateSimAllocation: {
		TemplateName: EmailTemplateSimAllocation,
		Keys: []string{
			"SUBSCRIBER",
			"NETWORK",
			"NAME",
			"QRCODE",
			"VOLUME",
			"UNIT",
			"ORG",
		},
	},
	EmailTemplateMemberInvite: {
		TemplateName: EmailTemplateMemberInvite,
		Keys: []string{
			"ORG",
			"OWNER",
			"NAME",
			"ROLE",
		},
	},
}

const (
	EmailKeySubscriber = "SUBSCRIBER"
	EmailKeyNetwork    = "NETWORK"
	EmailKeyName       = "NAME"
	EmailKeyQRCode     = "QRCODE"
	EmailKeyVolume     = "VOLUME"
	EmailKeyUnit       = "UNIT"
	EmailKeyOrg        = "ORG"
	EmailKeyOwner      = "OWNER"
	EmailKeyRole       = "ROLE"
)
