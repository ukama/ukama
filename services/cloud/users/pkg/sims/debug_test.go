package sims

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDubugIccid(t *testing.T) {
	got := GetDubugIccid()
	assert.Len(t, got, 18)

	imsi := GetDubugImsi(got)
	assert.Len(t, imsi, 15)

	assert.True(t, IsDebugIdentifier(got))
}
