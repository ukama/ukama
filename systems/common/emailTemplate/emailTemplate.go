package emailTemplate



const (
    EmailTemplateSimAllocation     = "sim-allocation"
    EmailTemplateMemberInvite      = "member-invite"
    EmailTemplateOrgInvite         = "org-invite"
    EmailTemplatePackageAddition   = "topup-plan" 
    EmailTemplatePaymentReceipt    = "payment-receipt"
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
			"PACKAGE",
			"VOLUME",
			"UNIT",
			"ENDDATE",
			"ORG",
			"DURATION",
			"AMOUNT",
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
	EmailTemplatePaymentReceipt: {
		TemplateName: EmailTemplatePaymentReceipt,
		Keys: []string{
			"NAME",
			"ORG",
			"AMOUNT",
			"RECEIPT_NUMBER",
			"PAYMENT_DATE",
			"PAYMENT_METHOD",
			"DESCRIPTION",
		},
	},
}



const (
	EmailKeyInvitation = "INVITATION"
	EmailKeyLink       = "LINK"
	EmailKeySubscriber = "SUBSCRIBER"
	EmailKeyNetwork    = "NETWORK"
	EmailKeyName       = "NAME"
	EmailKeyQRCode     = "QRCODE"
	EmailKeyVolume     = "VOLUME"
	EmailKeyUnit       = "UNIT"
	EmailKeyOrg        = "ORG"
	EmailKeyOwner      = "OWNER"
	EmailKeyRole       = "ROLE"
	EmailKeyExpiration = "EXPIRATION_DATE"
	EmailKeyPackagesCount = "PACKAGES_COUNT"
	EmailKeyPackagesDetails = "PACKAGES_DETAILS"
	EmailKeyEndDate    = "ENDDATE"
	EmailKeyPackage    = "PACKAGE"
	EmailKeyAmount     = "AMOUNT"
	EmailKeyDescription   = "DESCRIPTION"
	EmailKeyReceiptNumber = "RECEIPT_NUMBER"
	EmailKeyPaymentDate   = "PAYMENT_DATE"
	EmailKeyPaymentMethod = "PAYMENT_METHOD"
	EmailKeyDuration   = "DURATION"

)