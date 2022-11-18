package validationErrors

import (
	"time"
)

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

func IsFutureDate(date string) bool {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return false
	}
	return time.Now().Before(t)
}
