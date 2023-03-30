package validations

import (
	"errors"
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

func IsFutureDate(date string) error {
	t, err := time.Parse(time.RFC1123, date)
	if err != nil {
		return err
	}
	if t.After(time.Now()) {
		return nil
	}
	return errors.New("Date is not in the future")
}

func ValidateDate(date string) (string, error) {
	t, err := time.Parse(time.RFC1123, date)
	if err != nil {
		return "", errors.New("invalid date format, must be RFC1123 standard")
	}
	return t.Format(time.RFC1123), nil
}
