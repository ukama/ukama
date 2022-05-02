//go:build tools
// +build tools

package tools

import (
	_ "github.com/vektra/mockery/v2"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
