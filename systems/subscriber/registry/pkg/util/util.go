package server

import (
	"errors"
	"time"
)

func ValidateDOB(dob string) (time.Time, error) {
	t, err := time.Parse("02-01-2006", dob)
	if err != nil {
		return time.Time{}, errors.New("invalid date format, must be dd-mm-yyyy")
	}
	return t, nil
}
