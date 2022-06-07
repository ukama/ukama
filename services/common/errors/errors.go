// This package is used to replace usage of github.com/pkg/errors package that used all over the code base
// with more conventional wrapping
package errors

import (
	"fmt"
)

func Wrap(err error, message string) error {
	return fmt.Errorf("%s %w", message, err)
}
