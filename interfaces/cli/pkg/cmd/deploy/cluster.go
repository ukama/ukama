package deploy

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ukama/ukama/interfaces/cli/pkg"
	"github.com/ukama/ukama/interfaces/cli/pkg/cluster"
	"github.com/ukama/ukama/interfaces/cli/pkg/config"
	"gopkg.in/yaml.v3"
)

type clusterDeployConf struct {
	pkg.GlobalConfig `mapstructure:",squash"`
	Name             string `flag:"name"`
	Region           string `flag:"region" default:"us-east-1"`
	DnsName          string `flag:"dns" default:"ukama.k8s.local"`
	Bucket           string `flag:"bucket"`
}

func NewDeployClusterCommand(confReader config.ConfigReader) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Provision Kubernetes cluster",
		Run: func(cmd *cobra.Command, args []string) {
			nc := &clusterDeployConf{}
			confReader.ReadConfig("deploy", cmd.Flags(), nc)

			if nc.Verbose {
				b, _ := yaml.Marshal(nc)
				fmt.Fprintf(cmd.OutOrStdout(), "Deploy Config:\n '%s'\n", string(b))
			}
			logger := pkg.NewLogger(cmd.OutOrStdout(), cmd.ErrOrStderr(), nc.Verbose)

			kops := cluster.NewKopsWrapper(logger, nc.Bucket, nc.Verbose)
			err := kops.ProvisionAwsCluster(nc.DnsName, nc.Region)
			if err != nil {
				logger.Errorf(err.Error() + "\n")
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringP("name", "n", "ukama-cluster", "Cluster name")
	cmd.Flags().StringP("region", "r", "us-east-1", "Region to deploy cluster. The format depends on the cloud provider")
	cmd.Flags().StringP("dns", "d", "ukama.k8s.local", "Dns name to access the cluster APi. Use `<subbdomain>.k8s.local` domain name in case you don't want to configure dns")
	cmd.Flags().StringP("bucket", "b", "", "Bucket name to store cluster state. Refer to https://kops.sigs.k8s.io/getting_started/aws/#cluster-state-storage")

	return cmd
}
