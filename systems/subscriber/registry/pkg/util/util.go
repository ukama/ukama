package server

import (
	"errors"
	"time"
)

func ValidateDOB(dob string) (string, error) {
	t, err := time.Parse(time.RFC3339, dob)
	if err != nil {
		return "", errors.New("invalid date format, must be RFC3339 standard")
	}
	return t.Format(time.RFC3339), nil
}
