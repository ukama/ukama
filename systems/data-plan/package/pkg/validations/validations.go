package validationErrors

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

func IsReqEmpty(id uint64) bool {
	return id == 0
}

func ValidateDOB(dob string) (time.Time, error) {
	t, err := time.Parse("02-01-2006", dob)
	if err != nil {
		return time.Time{}, errors.New("invalid date format, must be dd-mm-yyyy")
	}
	return t, nil
}
