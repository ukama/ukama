package validationErrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsEmpty(t *testing.T) {
	// Success case
	check := IsEmpty("a", "b", "c")
	assert.False(t, check)

	_check := IsEmpty("", "", "c")
	assert.True(t, _check)

	__check := IsEmpty("", "", "")
	assert.True(t, __check)
}
func Test_IsInvalidId(t *testing.T) {
	// Success case
	check := IsInvalidId(0)
	assert.True(t, check)

	_check := IsInvalidId(1)
	assert.False(t, _check)
}
