package server

import (
	"errors"
	"time"
)


func ValidateDOB(dob string) (string, error) {
    t, err := time.Parse("02-01-2006", dob)
    if err != nil {
        return "", errors.New("invalid date format, must be dd-mm-yyyy")
    }
    return t.Format(time.RFC1123), nil
}
