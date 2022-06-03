package server

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/services/cloud/net/pkg"
)

func Test_marshallTargets(t *testing.T) {
	h := HttpServer{
		orgMetricsConf: &pkg.OrgMetricsConf{
			Url: "https://localhost:10251", // full url with port path and http
		},
	}
	output, err := h.marshallTargets(map[string]string{"a": "1", "b": "2"}, map[string]pkg.OrgNet{
		"a": {Org: "org1", Network: "net1"},
	}, 10250)
	assert.NoError(t, err)

	actual := strings.ReplaceAll(string(output), " ", "")
	actual = strings.Replace(actual, `{"targets":["1:10250"],"labels":{"network":"net1","nodeid":"a","org":"org1"}}`, "", 1)
	actual = strings.Replace(actual, `{"targets":["2:10250"],"labels":{"nodeid":"b"}}`, "", 1)

	assert.Equal(t, "[{\"targets\":[\"https://localhost:10251\"],\"labels\":null},,]", actual)
}
