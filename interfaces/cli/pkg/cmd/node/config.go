package node

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/ukama/ukama/services/cli/pkg"
	"github.com/ukama/ukama/services/cli/pkg/clients"
	"github.com/ukama/ukama/services/cli/pkg/cmd/config"
	"os"
)

// Yep, that's config for node config command
type nodeConfigConfig struct {
	config.GlobalConfig `mapstructure:",squash"`
	Node                nodeConfig
}

type nodeConfig struct {
	Ip       string `flag:"ip"`
	Cert     string `flag:"cert" validate:"file"`
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
			confReader.ReadConfig("node", cmd.Flags(), nc)

			if nc.Verbose {
				fmt.Fprintf(cmd.OutOrStdout(), "Node Config: '%+v'\n", nc)
			}

			clt := clients.NewNodeClient(pkg.NewLogger(cmd.OutOrStdout(), cmd.ErrOrStderr(), nc.Verbose))
			conf := cmd.InOrStdin()
			if nc.Node.ConfPath != "" {
				f, err := os.Open(nc.Node.ConfPath)
				if err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "Failed to open config file: %s\n", err)
					os.Exit(1)
				}
				conf = f
			}
			err := clt.SendFile(nc.Node.Ip, nc.Node.Cert, "", "", conf)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Failed to send config: %s\n", err)
				os.Exit(1)
			}
		},
	}

	getCmd.Flags().StringP("ip", "i", "", "Node ip or hostname")
	getCmd.Flags().StringP("cert", "c", "", "Node certificate")
	getCmd.Flags().StringP("conf", "k", "", "Path to config")

	return &getCmd
}
