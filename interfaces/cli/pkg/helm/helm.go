package helm

import (
	"fmt"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"os"
	"time"

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

type HelmRun interface {
	Run(chrt *chart.Chart, vals map[string]interface{}) (*release.Release, error)
}

func (h *HelmClient) InstallChart(chartName string, chartVersion string, namespace string, valueOpts *values.Options, isUpgrade bool,
	mustSetParams map[string]string) error {
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

	if mustSetParams != nil {
		path, err := h.chartProvider.RenderDefaultValues(chartName, mustSetParams)
		if err != nil {
			return errors.Wrap(err, "error rendering default values")
		}
		valueOpts.ValueFiles = []string{path}
	}

	// prepare values
	p := getter.All(settings)
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		return errors.Wrap(err, "error merging helm values")
	}
	releaseName := chartName + "-release"
	if isUpgrade { //
		iCli := action.NewUpgrade(actionConfig)
		iCli.Namespace = namespace
		iCli.Devel = h.verbose
		rel, err := iCli.Run(releaseName, chart, vals)
		if err != nil {
			h.log.Errorf("Error upgrading chart: %s\n", err)
			return err
		}
		h.log.Printf("Successfully upgraded release: ", rel.Name)

	} else { // fresh install
		iCli := action.NewInstall(actionConfig)
		iCli.Namespace = namespace
		iCli.ReleaseName = releaseName
		iCli.CreateNamespace = true
		iCli.UseReleaseName = true
		iCli.Devel = h.verbose
		iCli.Wait = true
		iCli.Timeout = 3 * time.Minute
		rel, err := iCli.Run(chart, vals)
		if err != nil {
			h.log.Errorf("Error applying chart: %s\n", err)
			return errors.Wrap(err, "error applying chart")
		}
		h.log.Printf("Successfully installed release: ", rel.Name)
	}

	return nil
}

func (h *HelmClient) debug(format string, v ...interface{}) {
	format = fmt.Sprintf("[debug] %s\n", format)
	h.log.Printf(format, v...)
}
