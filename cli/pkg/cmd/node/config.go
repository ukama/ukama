package node

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/ukama/ukamaX/cli/pkg/cmd/config"
	"os"
)

// Yep, that's config for node config command
type nodeConfigConfig struct {
	config.GlobalConfig `mapstructure:",squash"`
	Node                nodeConfig
}

type nodeConfig struct {
	Ip       string `flag:"ip"`
	Cert     string `flag:"cert"`
	ConfPath string `flag:"conf"`
}

// getCmd is a test command to demonstrate configurability
func configCommand(confReader config.ConfigReader) *cobra.Command {

	getCmd := cobra.Command{
		Use:   "config",
		Short: "Configure the node",
		Long:  `Sends a configuration to a node`,
		Run: func(cmd *cobra.Command, args []string) {

			nc := &nodeConfigConfig{}
			err := confReader.ReadConfig("node", cmd.Flags(), nc)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error reading config: '%+v'\n", err)
				os.Exit(1)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Node Config: '%+v'\n", nc)
		},
	}

	getCmd.Flags().StringP("ip", "i", "", "Node ip or hostname")
	getCmd.Flags().StringP("cert", "cr", "", "Node certificate")
	getCmd.Flags().StringP("conf", "c", "", "Path to config")

	return &getCmd
}
