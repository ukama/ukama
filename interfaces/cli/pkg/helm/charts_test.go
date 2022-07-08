package helm

import (
	b64 "encoding/base64"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/interfaces/cli/pkg"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
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
		},
		))

	cp := NewChartProvider(pkg.NewLogger(os.Stdout, os.Stderr, true), repoToken, srv.URL)

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
}

func TestBuildUrl(t *testing.T) {
	cp := NewChartProvider(pkg.NewLogger(os.Stdout, os.Stderr, true), "", "http://example.com")
	tok := cp.buildChartUrl("http://example.com", "", "ukamax", "1.2.3")
	assert.Equal(t, "http://example.com/ukamax-1.2.3.tgz", tok)
}
