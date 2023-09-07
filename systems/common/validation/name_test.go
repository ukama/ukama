package validation

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestIsValidDnsLabelName(t *testing.T) {

	tests := []struct {
		name string
		want bool
	}{
		{"some-org", true},
		{"org-1", true},
		{"s", true},
		{strings.Repeat("a", 253), true},
		{"some-Org", false},
		{"sAasdfds", false},
		{"", false},
		{strings.Repeat("a", 254), false},
		{"asdfsd;", false},
		{"as_dfs;", false},
		{"asdfsd.", false},
		{"-asdfsd", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidDnsLabelName(tt.name)
			assert.Equal(t, tt.want, got, "Name: "+tt.name)
		})
	}
}
