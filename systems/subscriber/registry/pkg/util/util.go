package server

import (
	"errors"
	"time"
)

func ValidateDOB(dob string) (string, error) {
	t, err := time.Parse(time.RFC1123, dob)
	if err != nil {
		return "", errors.New("invalid date format, must be RFC1123 standard")
	}
	return t.Format(time.RFC1123), nil
}
