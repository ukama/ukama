package validationErrors

func IsRequestEmpty(ss ...string) bool {
	for _, s := range ss {
		if s == "" {
			return true
		}
	}
	return false
}

func ReqSimTypeToPb(simType string) string {
	switch simType {
	case "INTER_MNO_ALL":
		return "inter_mno_all"
	case "INTER_MNO_DATA":
		return "inter_mno_data"
	case "INTER_NONE":
		return "inter_none"
	case "INTER_UKAMA_ALL":
		return "inter_ukama_all"
	default:
		return "inter_none"
	}

}
