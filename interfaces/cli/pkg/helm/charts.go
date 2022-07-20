package helm

import (
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/go-resty/resty/v2"
	"github.com/ukama/ukama/interfaces/cli/pkg"
	"github.com/ukama/ukama/services/common/errors"
	"golang.org/x/mod/semver"
	"gopkg.in/yaml.v3"
)

type ChartProvider struct {
	log        pkg.Logger
	repoToken  string
	repo       string
	restClient *resty.Client
}

var valueFiles = map[string]string{
	"ukamax": "ukamax-cli-values.gotmpl",
}
var outputDir = "tmp/"

func NewChartProvider(log pkg.Logger, repoUrl string, repoToken string) *ChartProvider {
	return &ChartProvider{log: log, repoToken: repoToken, repo: repoUrl,
		restClient: resty.New()}
}

func (c *ChartProvider) DownloadChart(name string, version string) (string, error) {
	var err error
	if version == "" {
		version, err = c.getLatestVersion(name)
		if err != nil {
			return "", err
		}
	}

	chartUrl := c.buildChartUrl(c.repo, c.repoToken, name, version)
	err = os.RemoveAll(outputDir)
	if err != nil {
		c.log.Printf("Error cleaning up output dir: %s", err)
	}
	chartPath := path.Join(outputDir, "chart.tgz")

	resp, err := c.restClient.R().
		SetOutput(chartPath).
		Get(chartUrl)

	if err != nil {
		return "", errors.Wrap(err, "error downloading chart")
	}
	if !resp.IsSuccess() {
		return "", fmt.Errorf("error downloading helm chart. Status code %d, Response: %s", resp.StatusCode(), resp.String())
	}

	return chartPath, err
}

func (c *ChartProvider) RenderDefaultValues(name string, data any) (valPath string, err error) {
	tmplFile, err := c.downloadDefaultValues(name)
	if err != nil {
		return "", err
	}

	valPath = path.Join(outputDir, "values.yaml")

	tmpl, err := template.New("values.gotmpl").Funcs(sprig.TxtFuncMap()).ParseFiles(tmplFile)
	if err != nil {
		return "", errors.Wrap(err, "error parsing values template")
	}
	f, err := os.Create(valPath)
	if err != nil {
		return "", errors.Wrap(err, "error creating values file")
	}
	defer f.Close()
	err = tmpl.Execute(f, map[string]any{
		"Values": data,
	})

	return valPath, err
}

func (c *ChartProvider) downloadDefaultValues(name string) (valPath string, err error) {
	valPath = path.Join(outputDir, "values.gotmpl")
	repoUrl := c.getRepoUrl(c.repo, c.repoToken)
	url := fmt.Sprintf("%s/cli/%s", strings.TrimSuffix(repoUrl, "/"), valueFiles[name])
	c.log.Printf("Downloading values file\n")
	resp, err := c.restClient.R().
		SetOutput(valPath).
		Get(url)

	if err != nil {
		return "", errors.Wrap(err, "error downloading value file")
	}
	if !resp.IsSuccess() {
		return "", fmt.Errorf("error downloading value file. Status code %d, Response: %s", resp.StatusCode(), resp.String())
	}

	return valPath, err
}

func (c *ChartProvider) getLatestVersion(name string) (string, error) {
	s := struct {
		Entries map[string][]struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}
	}{}

	repoUrl := c.getRepoUrl(c.repo, c.repoToken)
	indexUrl := fmt.Sprintf("%s/index.yaml", strings.TrimSuffix(repoUrl, "/"))
	resp, err := c.restClient.R().Get(indexUrl)
	if err != nil {
		return "", errors.Wrap(err, "error getting chart index")
	}
	if resp.IsError() {
		return "", fmt.Errorf("error donwloading chart index. Status code %d, Response: %s", resp.StatusCode(), resp.String())
	}

	err = yaml.Unmarshal(resp.Body(), &s)
	if err != nil {
		return "", errors.Wrap(err, "error parsing chart index")
	}

	ver := "v0.0.0"
	for _, e := range s.Entries[name] {
		curr := "v" + e.Version
		if semver.Compare(curr, ver) > 0 {
			ver = curr
		}
	}

	return strings.TrimPrefix(ver, "v"), nil
}

func (c *ChartProvider) buildChartUrl(repo string, token string, chartName string, version string) string {
	chartUrl := repo

	chartUrl = c.getRepoUrl(chartUrl, token)

	chartUrl = fmt.Sprintf("%s/%s-%s.tgz", strings.TrimSuffix(chartUrl, "/"), chartName, version)
	return chartUrl
}

func (c *ChartProvider) getRepoUrl(chartUrl string, token string) string {
	if token == "" {
		return chartUrl
	}
	// add token to repo url
	if !strings.Contains(chartUrl, "@") {
		if strings.HasPrefix(c.repo, "http://") {
			chartUrl = strings.Replace(c.repo, "http://", "http://"+token+"@", 1)
		}
		if strings.HasPrefix(c.repo, "https://") {
			chartUrl = strings.Replace(c.repo, "https://", "https://"+token+"@", 1)
		}
	}
	return chartUrl
}
