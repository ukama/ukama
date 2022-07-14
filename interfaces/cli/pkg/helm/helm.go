package helm

import (
	"fmt"
	"os"

	"github.com/ukama/ukama/interfaces/cli/pkg"
	"github.com/ukama/ukama/services/common/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/kube"
)

type HelmClient struct {
	log           pkg.Logger
	chartProvider *ChartProvider
	verbose       bool
}

func NewHelmClient(chartProvider *ChartProvider, log pkg.Logger, verbose bool) *HelmClient {
	return &HelmClient{
		log:           log,
		chartProvider: chartProvider,
		verbose:       verbose,
	}
}

func (h *HelmClient) InstallChart(chartName string, chartVersion string, namespace string, valueOpts *values.Options) error {
	chartPath, err := h.chartProvider.DownloadChart(chartName, chartVersion)
	if err != nil {
		return err
	}

	settings := cli.New()

	chart, err := loader.Load(chartPath)
	if err != nil {
		h.log.Errorf("Error loading chart %+v\n", err)
		return errors.Wrap(err, "error loading chart")
	}

	actionConfig := new(action.Configuration)

	k8sClientGetter := settings.RESTClientGetter()
	if err := actionConfig.Init(k8sClientGetter, namespace, os.Getenv("HELM_DRIVER"), h.debug); err != nil {
		h.log.Errorf("Error loading kube config %+v\n", err)
		return errors.Wrap(err, "error loading kube config")
	}

	kubeClient := kube.New(k8sClientGetter)
	kubeClient.Namespace = namespace
	actionConfig.KubeClient = kubeClient

	//
	path, err := h.chartProvider.RenderDefaultValues(chartName, map[string]string{
		"baseDomain":     "exmple.com",
		"amqppass":       "testPass",
		"nodeMetricsUrl": "http://localhost:9091/metrics",
		"postgresPass":   "pass",
		"smtpRelayHost":  "smtp.example.com",
		"smtpUsername":   "user",
		"smtpPassword":   "pass",
	})
	valueOpts.ValueFiles = []string{path}

	// prepare values
	p := getter.All(settings)
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		return err
	}

	iCli := action.NewInstall(actionConfig)
	iCli.Namespace = namespace
	iCli.ReleaseName = "my-release"
	iCli.CreateNamespace = true
	iCli.UseReleaseName = true
	iCli.Devel = h.verbose
	rel, err := iCli.Run(chart, vals)
	if err != nil {
		h.log.Errorf("Error applying chart: %s\n", err)
		return errors.Wrap(err, "error applying chart")
	}
	h.log.Printf("Successfully installed release: ", rel.Name)
	return nil
}

func (h *HelmClient) debug(format string, v ...interface{}) {
	format = fmt.Sprintf("[debug] %s\n", format)
	h.log.Printf(format, v...)
}
