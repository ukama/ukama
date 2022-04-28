package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ukama/ukama/services/cli/cmd/version"
	"github.com/ukama/ukama/services/cli/pkg"
	"github.com/ukama/ukama/services/cli/pkg/cmd/config"
	"github.com/ukama/ukama/services/cli/pkg/cmd/node"
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

	return rootCmd
}
