package emailTemplate



const (
    EmailTemplateSimAllocation     = "sim-allocation"
    EmailTemplateMemberInvite      = "member-invite"
    EmailTemplateOrgInvite         = "org-invite"
    EmailTemplatePackageAddition   = "topup-plan" 
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
	EmailTemplatePackageAddition: { 
        TemplateName: EmailTemplatePackageAddition,
        Keys: []string{
            "SUBSCRIBER",
            "NETWORK",
            "NAME",
            "VOLUME",
            "UNIT",
            "ORG",
			"AMOUNT",
			"DURATION",
			"PACKAGE",
			"EXPIRATION_DATE",
			"PACKAGES_COUNT",
			"PACKAGES_DETAILS",
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
	EmailKeyPackage    = "PACKAGE"
	EmailKeyAmount     = "AMOUNT"
	EmailKeyDuration   = "DURATION"
	EmailKeyExpiration = "EXPIRATION_DATE"
	EmailKeyPackagesCount = "PACKAGES_COUNT"
	EmailKeyPackagesDetails = "PACKAGES_DETAILS"
)