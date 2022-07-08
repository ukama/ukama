package deploy

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/ukama/ukama/interfaces/cli/pkg/config"
)

type deployConfig struct {
	config.GlobalConfig `mapstructure:",squash"`
	Cloud               cloudConfig
	Service             []string `flag:"service"`
	BaseDomain          string   `flag:"baseDomain" validate:"required"`
}

type cloudConfig struct {
	Type      string `flag:"cloud" validate:"oneof:aws"`
	AccessKey string `flag:"access-key"`
	SecretKey string `flag:"secret-key"`
}

func NewDeployCommand(confReader config.ConfigReader) *cobra.Command {
	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy Ukama service",
		Run: func(cmd *cobra.Command, args []string) {
			nc := &deployConfig{}
			confReader.ReadConfig("deploy", cmd.Flags(), nc)

			if nc.Verbose {
				fmt.Fprintf(cmd.OutOrStdout(), "Deploy Config: '%+v'\n", nc)
			}
		},
	}

	return deployCmd
}

//ukama deploy --cloud AWS  --accessKeyId AKIAJXQZQZQZQZQZQZQ --secretAccessKey SECRET --baseDomain ukama.com  // deploy all services and provision AWS cluster
//ukama deploy --service ukama@v1.0.1  --cloud AWS --accessKeyId AKIAJXQZQZQZQZQZQZQ --secretAccessKey SECRET --baseDomain ukama.com  // deploys ukamax helm v1.0.1 chart and provision AWS cluster
//ukama deploy --service ukama  --clusterName ukama-dev --smtpRelayHost mail.ukama.com --smtpRelayHostUsername user --smtpRelayHostPassword pass --baseDomain ukama.com // deploys ukamax helm chart to ukama-dev cluster
//ukama deploy --service ukama  --service metrics --service hub --cloud AWS  --accessKeyId AKIAJXQZQZQZQZQZQZQ --secretAccessKey SECRET --baseDomain ukama.com   // deploy latest versions of metrics , hub, ukama and provision cluster
//ukama deploy --cloud AWS_EKS --accessKeyId AKIAJXQZQZQZQZQZQZQ --secretAccessKey SECRET --baseDomain ukama.com  // deploy latest version of all services and provision AWS cluster
