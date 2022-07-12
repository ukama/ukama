package deploy

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ukama/ukama/interfaces/cli/pkg"
	"github.com/ukama/ukama/interfaces/cli/pkg/config"
	"github.com/ukama/ukama/interfaces/cli/pkg/helm"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/cli/values"
)

type deployConfig struct {
	pkg.GlobalConfig `mapstructure:",squash"`
	Cloud            string `flag:"cloud"`
	Aws              *awsConfig
	Service          []string `flag:"service"`
	BaseDomain       string   `flag:"baseDomain" validate:"required"`
	K8s              *k8sConfig
	Helm             *helmConfig
}

type k8sConfig struct {
	Namespace string `flag:"k8s.namespace"`
}

type helmConfig struct {
	RepoUrl string `flag:"helmRepo"`
	Token   string `flag:"token"`
}

type awsConfig struct {
	AccessKey string `flag:"aws.accessKey"`
	SecretKey string `flag:"aws.secret"`
}

func NewDeployCommand(confReader config.ConfigReader) *cobra.Command {
	valueOpts := &values.Options{}

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy Ukama service",
		Run: func(cmd *cobra.Command, args []string) {
			nc := &deployConfig{}
			confReader.ReadConfig("deploy", cmd.Flags(), nc)

			if nc.Verbose {
				b, _ := yaml.Marshal(nc)
				fmt.Fprintf(cmd.OutOrStdout(), "Deploy Config:\n '%s'\n", string(b))
			}
			logger := pkg.NewLogger(cmd.OutOrStdout(), cmd.ErrOrStderr(), nc.Verbose)
			chartProvider := helm.NewChartProvider(logger, nc.Helm.RepoUrl, nc.Helm.Token)

			helmClient := helm.NewHelmClient(chartProvider, logger)

			if len(nc.Service) == 1 && strings.HasPrefix(nc.Service[0], "ukama") {
				chartName, chartVer := parsName(nc.Service[0])
				namesapce := nc.K8s.Namespace
				if len(namesapce) == 0 {
					namesapce = chartName
				}

				err := helmClient.InstallChart("ukamax", chartVer, namesapce, valueOpts)
				if err != nil {
					logger.Errorf("Failed to install chart: %s", err)
					os.Exit(1)
				}
			}
		},
	}

	cmd.Flags().StringP("service", "s", "", "Service name")
	cmd.Flags().StringP("baseDomain", "d", "", "Base domain")

	cmd.Flags().StringP("cloud", "c", "", "Cloud type")
	cmd.Flags().StringP("aws.accessKey", "", "", "access key to access AWS account")
	cmd.Flags().StringP("aws.secret", "", "", "AWS secret access key to access the AWS account")

	// Helm flags
	cmd.Flags().StringP("token", "t", "", "Helm repository token")

	cmd.Flags().StringP("helmRepo", "r", "https://raw.githubusercontent.com/ukama/helm-charts/repo-index", "Helm repository url")

	cmd.Flags().StringP("k8s.namespace", "", "", "Target Kubernetes namespace")
	addValueOptionsFlags(cmd.Flags(), valueOpts)

	return cmd
}

func parsName(chartName string) (name string, version string) {
	i := strings.LastIndex(chartName, "@v")
	if i == -1 {
		return chartName, ""
	} else {
		return chartName[:i], chartName[i+2:]
	}

}

func addValueOptionsFlags(f *pflag.FlagSet, v *values.Options) {
	f.StringSliceVarP(&v.ValueFiles, "values", "", []string{}, "specify values in a YAML file or a URL (can specify multiple)")
	f.StringArrayVar(&v.Values, "set", []string{}, "set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	f.StringArrayVar(&v.StringValues, "set-string", []string{}, "set STRING values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	f.StringArrayVar(&v.FileValues, "set-file", []string{}, "set values from respective files specified via the command line (can specify multiple or separate values with commas: key1=path1,key2=path2)")
}

//ukama deploy --cloud AWS  --cloud.cloud.accessKeyId AKIAJXQZQZQZQZQZQZQ --cloud.secretAccessKey SECRET --baseDomain ukama.com --token UKAMA_ACCESS_KEY  // deploy all services and provision AWS cluster
//ukama deploy --service ukama@v1.0.1  --cloud AWS --cloud.accessKeyId AKIAJXQZQZQZQZQZQZQ --cloud.secretAccessKey SECRET --baseDomain ukama.com --token UKAMA_ACCESS_KEY // deploys ukamax helm v1.0.1 chart and provision AWS cluster
//ukama deploy --service ukama  --clusterName ukama-dev --smtpRelayHost mail.ukama.com --smtpRelayHostUsername user --smtpRelayHostPassword pass --baseDomain ukama.com --token UKAMA_ACCESS_KEY // deploys ukamax helm chart to ukama-dev cluster
//ukama deploy --service ukama  --service metrics --service hub --cloud AWS  --cloud.accessKeyId AKIAJXQZQZQZQZQZQZQ --cloud.secretAccessKey SECRET --baseDomain ukama.com --token UKAMA_ACCESS_KEY  // deploy latest versions of metrics , hub, ukama and provision cluster
//ukama deploy --cloud AWS_EKS --cloud.accessKeyId AKIAJXQZQZQZQZQZQZQ --cloud.secretAccessKey SECRET --baseDomain ukama.com --token UKAMA_ACCESS_KEY // deploy latest version of all services and provision AWS cluster
