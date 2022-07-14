package helm

import (
	b64 "encoding/base64"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/interfaces/cli/pkg"
)

func TestDownloadChart(t *testing.T) {
	repoToken := "test-token"
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !assert.Contains(t, r.Header.Get("Authorization"), b64.StdEncoding.EncodeToString([]byte(repoToken+":"))) {
				t.FailNow()
			}
			// return index yaml
			if strings.Contains(r.RequestURI, "index.yaml") {
				b, err := ioutil.ReadFile("testdata/index.yaml")
				if err != nil {
					t.Error(err)
				}
				_, err = w.Write(b)
				if err != nil {
					t.Error(err)
				}
			}
			// return chart tgz
			if strings.Contains(r.URL.Path, "ukamax-1.2.3.tgz") {
				b, err := ioutil.ReadFile("testdata/chart.tgz")
				if err != nil {
					t.Error(err)
				}
				_, err = w.Write(b)

				if err != nil {
					t.Error(err)
				}
			}

			if strings.Contains(r.URL.Path, "ukamax-cli-values.gotmpl") {
				b, err := ioutil.ReadFile("testdata/ukamax-cli-values.gotmpl")
				if err != nil {
					t.Error(err)
				}
				_, err = w.Write(b)

				if err != nil {
					t.Error(err)
				}
			}
		},
		))

	cp := NewChartProvider(pkg.NewLogger(os.Stdout, os.Stderr, true), srv.URL, repoToken)

	t.Run("download chart", func(t *testing.T) {
		path, err := cp.DownloadChart("ukamax", "1.2.3")
		if assert.NoError(t, err) {
			assert.FileExistsf(t, path, "File does not exist")
		}
	})

	t.Run("download latest version", func(t *testing.T) {
		path, err := cp.DownloadChart("ukamax", "")
		if assert.NoError(t, err) {
			assert.FileExistsf(t, path, "File does not exist")
		}
	})

	t.Run("download values file", func(t *testing.T) {
		path, err := cp.downloadDefaultValues("ukamax")
		if assert.NoError(t, err) {
			assert.FileExistsf(t, path, "File does not exist")
		}
	})

	t.Run("render value files", func(t *testing.T) {
		path, err := cp.RenderDefaultValues("ukamax", map[string]string{
			"baseDomain": "ukama.com",
			"amqppass":   "test",
		})
		if assert.NoError(t, err) {
			assert.FileExistsf(t, path, "File does not exist")
			b, err := ioutil.ReadFile(path)
			if assert.NoError(t, err) {
				assert.Contains(t, string(b), "domain: \"ukama.com\"")
			}
		}
	})
}

func TestBuildUrl(t *testing.T) {
	cp := NewChartProvider(pkg.NewLogger(os.Stdout, os.Stderr, true), "", "http://example.com")
	tok := cp.buildChartUrl("http://example.com", "", "ukamax", "1.2.3")
	assert.Equal(t, "http://example.com/ukamax-1.2.3.tgz", tok)
}

func TestIntegration(t *testing.T) {
	cp := NewChartProvider(pkg.NewLogger(os.Stdout, os.Stderr, true), "https://raw.githubusercontent.com/ukama/helm-charts/standalone-ukamax", os.Getenv("GH_TOKEN"))
	path, err := cp.downloadDefaultValues("ukamax")
	if assert.NoError(t, err) {
		assert.FileExists(t, path)
	}
}
