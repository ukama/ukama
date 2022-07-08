package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ukama/ukama/interfaces/cli/cmd/version"
	"github.com/ukama/ukama/interfaces/cli/pkg"
	"github.com/ukama/ukama/interfaces/cli/pkg/cmd/deploy"
	"github.com/ukama/ukama/interfaces/cli/pkg/cmd/node"
	"github.com/ukama/ukama/interfaces/cli/pkg/config"
)

func RootCommand() *cobra.Command {
	var cfgFile string

	var rootCmd = &cobra.Command{
		Use:   pkg.CliName,
		Short: "Ukama CLI",
	}

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "f", "", "config file (default is $HOME/.ukama.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose mode")
	rootCmd.Version = version.Version
	cfMgr := config.NewConfMgr(cfgFile, rootCmd.OutOrStdout(), rootCmd.ErrOrStderr())

	// top level commands

	rootCmd.AddCommand(node.NewNodeCommand(cfMgr))
	rootCmd.AddCommand(deploy.NewDeployCommand(cfMgr))

	return rootCmd
}
