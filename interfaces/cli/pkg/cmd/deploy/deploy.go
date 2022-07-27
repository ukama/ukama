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

var servicesDefaults = map[string]map[string]string{
	"ukamax": {
		"baseDomain":     "example.com",
		"amqppass":       "testPass",
		"nodeMetricsUrl": "http://localhost:9091/metrics",
		"postgresPass":   "pass",
		"smtpRelayHost":  "smtp.example.com",
		"smtpUsername":   "user",
		"smtpPassword":   "pass",
	},
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

			if len(nc.Service) == 1 && strings.HasPrefix(nc.Service[0], "ukama") {
				chartName, chartVer := parsName(nc.Service[0])
				namesapce := nc.K8s.Namespace
				if len(namesapce) == 0 {
					namesapce = chartName
				}
				params, err := svcParamsToMap(chartName, nc.Helm.ServiceParams)
				if err != nil {
					logger.Errorf(err.Error())
					os.Exit(1)
				}

				if !nc.Helm.SkipDeps {
					haveNginx, err := helm.IsNginxInstalled()
					if err != nil {
						logger.Errorf("Error checking for NGinx in cluster. Error: %s\n", err.Error())
						os.Exit(1)
					}

					if !haveNginx {
						err = installIngress(logger, nc.Verbose)
						if err != nil {
							logger.Errorf("Error installing ingress: %s\n", err.Error())
							os.Exit(1)
						}
					}
				}

				chartProvider := helm.NewChartProvider(logger, nc.Helm.RepoUrl, nc.Helm.Token)
				helmClient := helm.NewHelmClient(chartProvider, logger, nc.Verbose)
				err = helmClient.InstallChart(chartName, chartVer, namesapce, valueOpts, nc.Helm.IsUpgrade, params)
				if err != nil {
					logger.Errorf("Failed to install chart: %s\n", err)
					os.Exit(1)
				}
			}
		},
	}

	// Warning! Make sure that flags do not collide with Helm flags
	cmd.Flags().StringP("service", "s", "", "Service name")
	cmd.Flags().StringP("baseDomain", "d", "", "Base domain")

	cmd.Flags().StringP("cloud", "c", "", "Cloud type")
	cmd.Flags().StringP("aws.accessKey", "", "", "access key to access AWS account")
	cmd.Flags().StringP("aws.secret", "", "", "AWS secret access key to access the AWS account")

	// Helm flags
	cmd.Flags().StringP("token", "t", "", "Helm repository token")

	cmd.Flags().StringP("helmRepo", "r", "https://raw.githubusercontent.com/ukama/helm-charts/repo-index", "Helm repository url")
	cmd.Flags().BoolP("upgrade", "u", false, "Specify if upgrading already deployed chart")

	cmd.Flags().StringP("k8s.namespace", "", "", "Target Kubernetes namespace")
	addValueOptionsFlags(cmd.Flags(), valueOpts)

	cmd.Flags().StringArrayP("svcParams", "p", []string{}, "Comma separated key=value pairs for service parameters")

	cmd.Flags().BoolP("skipDeps", "", false, "Skip validation and installation of dependencies")

	cmd.AddCommand(NewDeployClusterCommand(confReader))

	return cmd
}

func installIngress(log pkg.Logger, verbose bool) error {
	chartProvider := helm.NewChartProvider(log, "https://github.com/kubernetes/ingress-nginx/releases/download/helm-chart-4.2.0/", "")
	helmClient := helm.NewHelmClient(chartProvider, log, verbose)
	err := helmClient.InstallChart("ingress-nginx", "4.2.0", "kube-system", &values.Options{}, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func svcParamsToMap(chartName string, params []string) (res map[string]string, err error) {
	res = servicesDefaults[chartName]

	for _, p := range params {
		kv := strings.Split(p, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid service parameter: %s, Should have key=value format", p)
		}
		res[kv[0]] = kv[1]
	}
	return res, nil
}

func parsName(chartName string) (name string, version string) {

	i := strings.LastIndex(chartName, "@v")
	if i == -1 {
		name = chartName
		version = ""
	} else {
		name = chartName[:i]
		version = chartName[i+2:]
	}
	if strings.EqualFold(name, "ukama") {
		name = "ukamax"
	}

	return
}

func addValueOptionsFlags(f *pflag.FlagSet, v *values.Options) {
	f.StringSliceVarP(&v.ValueFiles, "values", "", []string{}, "specify values in a YAML file or a URL (can specify multiple)")
	f.StringArrayVar(&v.Values, "set", []string{}, "set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	f.StringArrayVar(&v.StringValues, "set-string", []string{}, "set STRING values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	f.StringArrayVar(&v.FileValues, "set-file", []string{}, "set values from respective files specified via the command line (can specify multiple or separate values with commas: key1=path1,key2=path2)")
}
