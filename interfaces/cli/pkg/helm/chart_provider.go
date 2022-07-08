package helm

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/ukama/ukama/interfaces/cli/pkg"
	"github.com/ukama/ukama/services/common/errors"
	"path"
)

type ChartProvider struct {
	log       pkg.Logger
	repoToken any
}

func NewChartProvider(log pkg.Logger, repoToken string) *ChartProvider {
	return &ChartProvider{log: log, repoToken: repoToken}
}

func (c ChartProvider) DownloadChart(name string, version string) (string, error) {
	chartUrl := fmt.Sprintf("https://%s@raw.githubusercontent.com/ukama/helm-charts/repo-index/%s-%s.tgz", c.repoToken, name, version)
	dir := "tmp/"
	chartPath := path.Join(dir, "chart.tgz")

	client := resty.New()
	_, err := client.R().
		SetOutput(chartPath).
		Get(chartUrl)

	if err != nil {
		return "", errors.Wrap(err, "error downloading chart")
	}

	return chartPath, err
}
