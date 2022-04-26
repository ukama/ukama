//go:build tools
// +build tools

package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint@v1.44.0"
	_ "github.com/vektra/mockery/v2"
)
