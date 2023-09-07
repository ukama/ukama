package validations

func IsEmpty(ss ...string) bool {
	for _, s := range ss {
		if s == "" {
			return true
		}
	}
	return false
}

func IsValidUploadReqArgs(fileUrl, effectiveAt, simType string) bool {
	return !IsEmpty(fileUrl, effectiveAt, simType)
}
