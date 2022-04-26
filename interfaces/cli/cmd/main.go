package main

import (
	"github.com/spf13/cobra"
	"github.com/ukama/ukamaX/cli/pkg/cmd"
)

func main() {
	cobra.CheckErr(cmd.RootCommand().Execute())
}
