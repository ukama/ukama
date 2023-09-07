package node

import (
	"github.com/spf13/cobra"
	"github.com/ukama/ukama/interfaces/cli/pkg/cmd/config"
)

func NewNodeCommand(confReader config.ConfigReader) *cobra.Command {
	nodeCmd := &cobra.Command{
		Use:   "node",
		Short: "Access node",
	}

	nodeCmd.AddCommand(configCommand(confReader))
	return nodeCmd
}
